package mempool

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

const MaxSizeClass = 32 // 最大 Size Class 数量

// HighConcurrencyPool 高并发内存池
// 对应 C 结构：typedef struct { ... } HighConcurrencyPool;
type HighConcurrencyPool struct {
	baseAddr   unsafe.Pointer // 预分配的内存块起始地址
	totalSize  int            // 总内存大小
	usedSize   int            // 已使用内存大小（非线程安全，仅示例）
	globalLock sync.Mutex     // 全局锁（保护大对象分配）
	tlsCache   *ThreadLocalCache
	mmapData   []byte // 保持 mmap 内存的引用，防止 GC 释放底层数组
}

// ThreadLocalCache 线程本地缓存（TLS）
// 注意：Go 的 goroutine 可能在 OS 线程间迁移，严格意义上需要使用 sync.Pool 或 runtime.LockOSThread
// 此处为代码结构对照，保持与 C 版本一致的字段定义
type ThreadLocalCache struct {
	freeList []unsafe.Pointer // 按 Size Class 分类的空闲链表头，索引对应 size class
}

// HcpInit 初始化高并发内存池（替代 new 的全局分配）
// 对应 C 函数：HighConcurrencyPool* hcp_init(size_t total_size)
func HcpInit(totalSize int) (*HighConcurrencyPool, error) {
	// 1. 预分配连续内存块（对齐到页大小）
	// 对应 C: mmap(NULL, total_size, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0)
	data, err := syscall.Mmap(-1, 0, totalSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		return nil, errors.New("mmap failed: " + err.Error())
	}

	// 2. 初始化线程本地缓存（TLS）
	// 对应 C: tls_cache->free_list = (void**)calloc(MAX_SIZE_CLASS, sizeof(void*))
	tlsCache := &ThreadLocalCache{
		freeList: make([]unsafe.Pointer, MaxSizeClass),
	}

	// 3. 初始化全局内存池
	// 对应 C: malloc + 字段赋值
	pool := &HighConcurrencyPool{
		baseAddr:  unsafe.Pointer(&data[0]),
		totalSize: totalSize,
		usedSize:  0,
		tlsCache:  tlsCache,
		mmapData:  data, // 必须保持引用，否则底层数组可能被 GC 回收
	}

	return pool, nil
}

// alignToSizeClass 对齐到最近的 Size Class
// 简化实现：按 8, 16, 32, 64... 字节对齐（实际实现可能更复杂）
func alignToSizeClass(size uintptr) int {
	if size <= 8 {
		return 0
	}
	sc := 0
	aligned := uintptr(8)
	for size > aligned && sc < MaxSizeClass-1 {
		aligned <<= 1
		sc++
	}
	return sc
}

// allocateFromGlobal 从全局池申请一批对象（简化版）
// 对应 C: allocate_from_global(pool, sc)
func (pool *HighConcurrencyPool) allocateFromGlobal(sc int) unsafe.Pointer {
	// 计算当前 size class 的对象大小（示例：8 << sc）
	objSize := uintptr(8 << sc)

	if pool.usedSize+int(objSize) > pool.totalSize {
		return nil // 内存不足
	}

	// 从预分配内存中切分一块
	addr := unsafe.Pointer(uintptr(pool.baseAddr) + uintptr(pool.usedSize))
	pool.usedSize += int(objSize)

	return addr
}

// HcpAlloc 分配函数（替代 new，无锁线程本地操作）
// 对应 C 函数：void* hcp_alloc(HighConcurrencyPool* pool, size_t size)
func (pool *HighConcurrencyPool) HcpAlloc(size uintptr) unsafe.Pointer {
	tls := pool.tlsCache
	sc := alignToSizeClass(size)

	// 1. 从线程本地缓存分配（O(1)时间）
	// 对应 C: void** bucket = &tls->free_list[sc];
	bucket := &tls.freeList[sc]
	if *bucket != nil {
		obj := *bucket
		// 链表指针前移：*bucket = *(void**)obj
		// 将对象首部的指针值（即下一个节点地址）读出，更新链表头
		nextPtr := (*unsafe.Pointer)(obj)
		*bucket = *nextPtr
		return obj
	}

	// 2. 本地无空闲对象，向全局缓存申请（批量分配）
	pool.globalLock.Lock()
	chunk := pool.allocateFromGlobal(sc)
	pool.globalLock.Unlock()

	if chunk != nil {
		// 将新分配的对象插入本地缓存链表头部
		// 对应 C: *(void**)(chunk) = tls->free_list[sc];
		nextField := (*unsafe.Pointer)(chunk)
		*nextField = tls.freeList[sc]

		// 对应 C: tls->free_list[sc] = chunk;
		tls.freeList[sc] = chunk
		return chunk
	}

	return nil // 全局池无内存，触发扩容或返回 NULL
}

// HcpDestroy 释放内存池（补充函数，对应 munmap）
func (pool *HighConcurrencyPool) HcpDestroy() error {
	pool.globalLock.Lock()
	defer pool.globalLock.Unlock()

	// 解除内存映射
	if pool.mmapData != nil {
		return syscall.Munmap(pool.mmapData)
	}
	return nil
}
