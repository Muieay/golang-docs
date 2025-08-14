package main

import (
    "fmt"
    "strings"
    "sync"
    "time"
)

// 测试和基准测试
// Go语言的测试框架和最佳实践

// 1. 待测试的函数

// 计算斐波那契数列
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}

// 优化版的斐波那契数列（使用迭代）
func FibonacciOptimized(n int) int {
    if n <= 1 {
        return n
    }
    
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

// 字符串工具函数
func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// 判断素数
func IsPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}

// 并发安全的计数器
type Counter struct {
    mu    sync.Mutex
    value int
}

func NewCounter() *Counter {
    return &Counter{value: 0}
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// 并发安全的队列
type Queue struct {
    items []interface{}
    mu    sync.Mutex
}

func NewQueue() *Queue {
    return &Queue{items: make([]interface{}, 0)}
}

func (q *Queue) Enqueue(item interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (interface{}, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.items) == 0 {
        return nil, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

func (q *Queue) Size() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    return len(q.items)
}

// 缓存系统
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func NewCache() *Cache {
    return &Cache{data: make(map[string]interface{})}
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, exists := c.data[key]
    return val, exists
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.data, key)
}

// 性能测试辅助函数
func generateTestData(size int) []int {
    data := make([]int, size)
    for i := 0; i < size; i++ {
        data[i] = i + 1
    }
    return data
}

// 性能基准测试
func benchmarkFibonacci(n int) time.Duration {
    start := time.Now()
    Fibonacci(n)
    return time.Since(start)
}

func benchmarkFibonacciOptimized(n int) time.Duration {
    start := time.Now()
    FibonacciOptimized(n)
    return time.Since(start)
}

// 内存使用测试
func memoryTest() {
    // 创建大量对象测试内存使用
    objects := make([]*Counter, 10000)
    for i := 0; i < 10000; i++ {
        objects[i] = NewCounter()
    }
    
    // 清理
    objects = nil
}

// 并发性能测试
func concurrentTest() {
    counter := NewCounter()
    var wg sync.WaitGroup
    
    numGoroutines := 1000
    incrementsPerGoroutine := 100
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < incrementsPerGoroutine; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    
    expected := numGoroutines * incrementsPerGoroutine
    actual := counter.Value()
    
    fmt.Printf("并发测试: 期望 %d, 实际 %d\n", expected, actual)
    fmt.Printf("测试状态: %s\n", boolToString(expected == actual))
}

// 字符串处理性能测试
func stringTest() {
    testStrings := []string{
        "hello",
        "世界",
        "Hello, 世界!",
        strings.Repeat("a", 1000),
        strings.Repeat("测试", 100),
    }
    
    fmt.Println("\n字符串处理测试:")
    for _, str := range testStrings {
        start := time.Now()
        reversed := ReverseString(str)
        duration := time.Since(start)
        
        fmt.Printf("字符串 '%s...' (%d字符) 反转耗时: %v\n", 
            str[:min(len(str), 20)], len(str), duration)
        
        // 验证正确性
        doubleReversed := ReverseString(reversed)
        if doubleReversed != str {
            fmt.Printf("错误: 双重反转不匹配\n")
        }
    }
}

// 素数计算性能测试
func primeTest() {
    testNumbers := []int{2, 3, 97, 100, 1000, 10007}
    
    fmt.Println("\n素数测试:")
    for _, num := range testNumbers {
        start := time.Now()
        isPrime := IsPrime(num)
        duration := time.Since(start)
        
        fmt.Printf("数字 %d 是素数: %s (耗时: %v)\n", 
            num, boolToString(isPrime), duration)
    }
    
    // 批量测试
    primes := 0
    start := time.Now()
    for i := 2; i <= 1000; i++ {
        if IsPrime(i) {
            primes++
        }
    }
    duration := time.Since(start)
    fmt.Printf("1-1000中共有 %d 个素数，耗时: %v\n", primes, duration)
}

// 缓存性能测试
func cacheTest() {
    cache := NewCache()
    
    // 写入测试
    start := time.Now()
    for i := 0; i < 10000; i++ {
        cache.Set(fmt.Sprintf("key%d", i), i)
    }
    writeDuration := time.Since(start)
    
    // 读取测试
    start = time.Now()
    hits := 0
    for i := 0; i < 10000; i++ {
        if _, exists := cache.Get(fmt.Sprintf("key%d", i)); exists {
            hits++
        }
    }
    readDuration := time.Since(start)
    
    fmt.Printf("\n缓存测试:")
    fmt.Printf("写入10000条记录耗时: %v\n", writeDuration)
    fmt.Printf("读取10000条记录耗时: %v\n", readDuration)
    fmt.Printf("读取命中率: %.2f%%\n", float64(hits)/10000*100)
}

// 队列性能测试
func queueTest() {
    queue := NewQueue()
    
    // 入队测试
    start := time.Now()
    for i := 0; i < 10000; i++ {
        queue.Enqueue(i)
    }
    enqueueDuration := time.Since(start)
    
    // 出队测试
    start = time.Now()
    dequeued := 0
    for {
        if _, ok := queue.Dequeue(); ok {
            dequeued++
        } else {
            break
        }
    }
    dequeueDuration := time.Since(start)
    
    fmt.Printf("\n队列测试:")
    fmt.Printf("入队10000条记录耗时: %v\n", enqueueDuration)
    fmt.Printf("出队10000条记录耗时: %v\n", dequeueDuration)
    fmt.Printf("实际出队数量: %d\n", dequeued)
}

// 性能比较测试
func performanceComparison() {
    fmt.Println("\n=== 性能比较测试 ===")
    
    // 斐波那契性能比较
    fmt.Println("\n斐波那契性能比较:")
    for n := 10; n <= 40; n += 10 {
        recursiveTime := benchmarkFibonacci(n)
        optimizedTime := benchmarkFibonacciOptimized(n)
        
        speedup := float64(recursiveTime) / float64(optimizedTime)
        fmt.Printf("n=%d: 递归: %v, 优化: %v, 加速比: %.1fx\n", 
            n, recursiveTime, optimizedTime, speedup)
    }
}

// 内存分配测试
func memoryAllocationTest() {
    fmt.Println("\n=== 内存分配测试 ===")
    
    // 测试不同大小的切片分配
    sizes := []int{100, 1000, 10000, 100000}
    
    for _, size := range sizes {
        start := time.Now()
        slice := make([]int, size)
        for i := 0; i < size; i++ {
            slice[i] = i
        }
        duration := time.Since(start)
        
        fmt.Printf("分配 %d 个整数的切片耗时: %v\n", size, duration)
    }
}

// 测试覆盖率示例
func testCoverageExample() {
    fmt.Println("\n=== 测试覆盖率示例 ===")
    
    // 测试各种边界条件
    testCases := []struct {
        name     string
        function func() error
    }{
        {"空字符串反转", func() error {
            if ReverseString("") != "" {
                return fmt.Errorf("空字符串反转失败")
            }
            return nil
        }},
        {"单字符反转", func() error {
            if ReverseString("a") != "a" {
                return fmt.Errorf("单字符反转失败")
            }
            return nil
        }},
        {"回文测试", func() error {
            palindrome := "racecar"
            if ReverseString(palindrome) != palindrome {
                return fmt.Errorf("回文反转失败")
            }
            return nil
        }},
    }
    
    for _, tc := range testCases {
        if err := tc.function(); err != nil {
            fmt.Printf("测试失败: %s - %v\n", tc.name, err)
        } else {
            fmt.Printf("测试通过: %s\n", tc.name)
        }
    }
}

// 并发测试框架
func concurrentTestFramework() {
    fmt.Println("\n=== 并发测试框架 ===")
    
    // 测试不同并发级别
    for concurrency := 10; concurrency <= 1000; concurrency *= 10 {
        counter := NewCounter()
        var wg sync.WaitGroup
        
        start := time.Now()
        
        for i := 0; i < concurrency; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for j := 0; j < 100; j++ {
                    counter.Increment()
                }
            }()
        }
        
        wg.Wait()
        duration := time.Since(start)
        
        expected := concurrency * 100
        actual := counter.Value()
        
        fmt.Printf("并发级别 %d: 耗时 %v, 期望 %d, 实际 %d\n", 
            concurrency, duration, expected, actual)
    }
}

// 辅助函数
func boolToString(b bool) string {
    if b {
        return "true"
    }
    return "false"
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    fmt.Println("=== Go语言测试和基准测试 ===")
    
    // 运行所有测试
    fmt.Println("\n1. 基本功能测试")
    fmt.Printf("Fibonacci(10) = %d\n", Fibonacci(10))
    fmt.Printf("FibonacciOptimized(10) = %d\n", FibonacciOptimized(10))
    fmt.Printf("ReverseString("Hello") = %s\n", ReverseString("Hello"))
    fmt.Printf("IsPrime(97) = %t\n", IsPrime(97))
    
    fmt.Println("\n2. 并发测试")
    concurrentTest()
    
    fmt.Println("\n3. 性能测试")
    performanceComparison()
    
    fmt.Println("\n4. 内存分配测试")
    memoryAllocationTest()
    
    fmt.Println("\n5. 字符串处理测试")
    stringTest()
    
    fmt.Println("\n6. 素数计算测试")
    primeTest()
    
    fmt.Println("\n7. 缓存性能测试")
    cacheTest()
    
    fmt.Println("\n8. 队列性能测试")
    queueTest()
    
    fmt.Println("\n9. 测试覆盖率示例")
    testCoverageExample()
    
    fmt.Println("\n10. 并发测试框架")
    concurrentTestFramework()
    
    fmt.Println("\n=== 测试完成 ===")
    fmt.Println("\n练习：")
    fmt.Println("1. 为上述所有函数编写完整的单元测试")
    fmt.Println("2. 使用go test -bench运行基准测试")
    fmt.Println("3. 使用go test -cover查看测试覆盖率")
    fmt.Println("4. 实现一个性能分析工具")
}