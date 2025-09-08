下面给出 3 组“能直接拷贝到生产”的写法，每组都包含

1. 限流（令牌桶）
2. 熔断（错误率触发 + 半开恢复）
3. 降级（返回缓存/默认值）

全部基于官方/社区最活跃依赖，无过度封装，2025-07-28 之后仍维护

## 一、最小依赖版（仅用标准库 + sony/gobreaker + golang.org/x/time/rate）

go.mod

```go
require (
    github.com/sony/gobreaker v0.5.0
    golang.org/x/time v0.3.0
)
```

main.go

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/sony/gobreaker"
    "golang.org/x/time/rate"
)

// ---------- 1. 限流：令牌桶 ----------
var limiter = rate.NewLimiter(100, 200) // 每秒 100 个，突发 200

func limitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        next(w, r)
    }
}

// ---------- 2. 熔断：gobreaker ----------
var cb *gobreaker.CircuitBreaker
var once sync.Once

func initCB() {
    once.Do(func() {
        cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
            Name:        "svc",
            MaxRequests: 5,               // 半开时最多放 5 个探测请求
            Interval:    10 * time.Second, // 统计滑动窗口
            Timeout:     3 * time.Second,  // 熔断后多久转半开
            ReadyToTrip: func(count gobreaker.Counts) bool {
                return count.Requests >= 10 &&
                    float64(count.TotalFailures)/float64(count.Requests) > 0.6
            },
            OnStateChange: func(name string, from, to gobreaker.State) {
                fmt.Printf("CB %s: %s -> %s\n", name, from, to)
            },
        })
    })
}

// ---------- 3. 降级：缓存 or 默认值 ----------
var fakeCache = map[string]string{"1": "cached-1"}

func fallbackHandler(w http.ResponseWriter, _ *http.Request) {
    w.Write([]byte("fallback: service is busy, try later"))
}

// ---------- 4. 业务 Handler ----------
func queryHandler(w http.ResponseWriter, r *http.Request) {
    initCB()
    id := r.URL.Query().Get("id")
    if v, ok := fakeCache[id]; ok { // 先读缓存
        w.Write([]byte("cache:" + v))
        return
    }

    _, err := cb.Execute(func() (interface{}, error) {
        // 模拟 RPC/DB 调用
        if id == "error" {
            return nil, fmt.Errorf("backend error")
        }
        time.Sleep(50 * time.Millisecond)
        return "real-data", nil
    })

    if err != nil { // 触发熔断 or 后端异常 → 降级
        fallbackHandler(w, r)
        return
    }
    w.Write([]byte("real-data"))
}

func main() {
    http.HandleFunc("/query", limitMiddleware(queryHandler))
    fmt.Println("listen :8080")
    http.ListenAndServe(":8080", nil)
}
```

## 二、框架版：Go-Zero 内置「限流+熔断」一行配置即可

api/hello.api

```proto
syntax = "v1"

info(
    title: "demo"
    desc: "limit & breaker"
)

type (
    Req  { Id string `path:"id"` }
    Resp { Msg string `json:"msg"` }
)

service hello {
    @handler GetById
    get /query/:id (Req) returns (Resp)
}
```

etc/hello.yaml

```yaml
Name: hello
Host: 0.0.0.0
Port: 8888
# ---------- 1. 限流 ----------
MaxConns: 100          # 最大并发
MaxQps: 200            # 每秒最大请求数
# ---------- 2. 熔断 ----------
Breaker: true          # 开启熔断
BreakerWindow: 10s     # 统计窗口
BreakerThreshold: 0.5  # 失败率阈值
BreakerSuccess: 3      # 半开成功个数即恢复
```

internal/logic/getbyid_logic.go

```go
func (l *GetByIdLogic) GetById(req *types.Req) (*types.Resp, error) {
    // 3. 降级：返回本地缓存
    if v := l.svcCtx.Cache.Get(req.Id); v != "" {
        return &types.Resp{Msg: "cache:" + v}, nil
    }
    // 模拟 RPC
    if req.Id == "error" {
        return nil, errors.New("backend error")
    }
    return &types.Resp{Msg: "real:" + req.Id}, nil
}
```

启动即自带「令牌桶 + 滑动窗口 + 熔断器 + 半开探测」，无需手写中间件。

## 三、纯手撸 100 行「令牌桶 + 熔断器」实现（面试/学习向）

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// ---------- 1. 令牌桶 ----------
type TokenBucket struct {
    rate   int64 // 每秒放令牌数
    cap    int64 // 桶容量
    tokens int64 // 剩余令牌
    last   int64 // 上次放令牌时间戳，纳秒
}

func NewTokenBucket(rate, cap int64) *TokenBucket {
    return &TokenBucket{rate: rate, cap: cap, tokens: cap, last: time.Now().UnixNano()}
}

func (t *TokenBucket) Allow() bool {
    now := time.Now().UnixNano()
    // 计算这段时间应放多少令牌
    elapse := now - atomic.LoadInt64(&t.last)
    add := elapse * t.rate / 1e9
    if add > 0 {
        atomic.StoreInt64(&t.last, now)
        newTokens := atomic.AddInt64(&t.tokens, add)
        if newTokens > t.cap {
            atomic.StoreInt64(&t.tokens, t.cap)
        }
    }
    return atomic.AddInt64(&t.tokens, -1) >= 0
}

// ---------- 2. 熔断器 ----------
type Breaker struct {
    window       int           // 窗口请求数
    failRate     float64       // 失败率阈值
    reqCnt       int64
    failCnt      int64
    state        int32 // 0=close 1=open 2=half
    lastOpenTime int64
    mu           sync.Mutex
}

func NewBreaker(window int, failRate float64) *Breaker {
    return &Breaker{window: window, failRate: failRate}
}

func (b *Breaker) Call(fn func() error) error {
    if !b.allow() {
        return fmt.Errorf("breaker open")
    }
    err := fn()
    b.after(err != nil)
    return err
}

func (b *Breaker) allow() bool {
    st := atomic.LoadInt32(&b.state)
    if st == 0 {
        return true
    }
    if st == 2 {
        return true // 半开允许探测
    }
    // open 状态等待 timeout
    if time.Now().Unix()-atomic.LoadInt64(&b.lastOpenTime) > 3 {
        atomic.StoreInt32(&b.state, 2) // 转半开
        return true
    }
    return false
}

func (b *Breaker) after(isFail bool) {
    b.mu.Lock()
    defer b.mu.Unlock()
    atomic.AddInt64(&b.reqCnt, 1)
    if isFail {
        atomic.AddInt64(&b.failCnt, 1)
    }
    if atomic.LoadInt64(&b.reqCnt) < int64(b.window) {
        return
    }
    rate := float64(atomic.LoadInt64(&b.failCnt)) / float64(atomic.LoadInt64(&b.reqCnt))
    if rate > b.failRate {
        atomic.StoreInt32(&b.state, 1)
        atomic.StoreInt64(&b.lastOpenTime, time.Now().Unix())
    } else {
        atomic.StoreInt32(&b.state, 0)
    }
    atomic.StoreInt64(&b.reqCnt, 0)
    atomic.StoreInt64(&b.failCnt, 0)
}

// ---------- 3. 使用示例 ----------
func main() {
    tb := NewTokenBucket(10, 20)
    cb := NewBreaker(20, 0.5)

    for i := 0; i < 30; i++ {
        if !tb.Allow() {
            fmt.Println("req", i, "reject by limit")
            continue
        }
        err := cb.Call(func() error {
            if i%3 == 0 {
                return fmt.Errorf("backend error")
            }
            return nil
        })
        if err != nil {
            fmt.Println("req", i, "err:", err)
        } else {
            fmt.Println("req", i, "success")
        }
    }
}
```
