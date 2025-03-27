# Gin源码阅读与分析

> 很典型的一个web框架

```go
package main

import "github.com/gin-gonic/gin"

func main() {
        r := gin.Default()
        r.GET("/ping", func(c *gin.Context) {
                c.JSON(200, gin.H{
                        "message": "pong",
                })
        })
        r.Run() // listen and serve on 0.0.0.0:8080
}
```

- 先看 `gin.Default`:

```go
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
        debugPrintWARNINGDefault()
        engine := New()
        engine.Use(Logger(), Recovery())
        return engine
}
```

- 看 `engine := New()` 所返回的结构体：

```go
func New() *Engine {
        debugPrintWARNINGNew()
        engine := &Engine{
                RouterGroup: RouterGroup{
                        Handlers: nil,
                        basePath: "/",
                        root:     true,
                },
                FuncMap:                template.FuncMap{},
                RedirectTrailingSlash:  true,
                RedirectFixedPath:      false,
                HandleMethodNotAllowed: false,
                ForwardedByClientIP:    true,
                AppEngine:              defaultAppEngine,
                UseRawPath:             false,
                UnescapePathValues:     true,
                MaxMultipartMemory:     defaultMultipartMemory,
                trees:                  make(methodTrees, 0, 9),
                delims:                 render.Delims{Left: "{{", Right: "}}"},
                secureJsonPrefix:       "while(1);",
        }
        engine.RouterGroup.engine = engine
        engine.pool.New = func() interface{} {
                return engine.allocateContext()
        }
        return engine
}
```
