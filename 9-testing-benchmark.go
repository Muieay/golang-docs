package main

import (
    "fmt"
    "strings"
    "sync"
    "time"
)

/*
Go测试与基准测试最佳实践：
1. 单元测试：验证函数功能正确性
2. 基准测试：测量代码性能
3. 并发测试：验证并发安全性
4. 内存分析：检测内存分配情况
*/

// ==================== 核心功能函数 ====================

/*
Fibonacci - 递归实现斐波那契数列
参数n: 要计算的斐波那契数列位置
返回值: 第n个斐波那契数
特点: 简单但效率低(时间复杂度O(2^n))
*/
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}

/*
FibonacciOptimized - 迭代优化斐波那契数列
参数n: 要计算的斐波那契数列位置
返回值: 第n个斐波那契数
特点: 高效实现(时间复杂度O(n))
实现思路: 使用两个变量交替保存前两个值
*/
func FibonacciOptimized(n int) int {
    if n <= 1 {
        return n
    }
    
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b // 交换变量值
    }
    return b
}

/*
ReverseString - 字符串反转函数
参数s: 要反转的字符串
返回值: 反转后的字符串
实现思路: 将字符串转为rune切片后双指针交换
注意: 使用rune支持Unicode字符
*/
func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i] // 交换字符
    }
    return string(runes)
}

/*
IsPrime - 素数判断函数
参数n: 要判断的数字
返回值: 是否为素数
优化点: 
1. 排除偶数
2. 只检查到sqrt(n)
3. 步进6减少检查次数
*/
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
    // 检查6k±1形式的数
    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}

// ==================== 并发安全数据结构 ====================

/*
Counter - 并发安全计数器
使用sync.Mutex保护共享状态
*/
type Counter struct {
    mu    sync.Mutex // 互斥锁
    value int        // 计数器值
}

func NewCounter() *Counter {
    return &Counter{}
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

/*
Queue - 并发安全队列
使用slice实现FIFO队列
*/
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

// ==================== 测试辅助函数 ====================

/*
generateTestData - 生成测试数据
参数size: 数据量大小
返回值: 包含1-size的整数切片
*/
func generateTestData(size int) []int {
    data := make([]int, size)
    for i := 0; i < size; i++ {
        data[i] = i + 1
    }
    return data
}

/*
benchmarkFibonacci - 斐波那契基准测试
参数n: 计算第n个斐波那契数
返回值: 计算耗时
*/
func benchmarkFibonacci(n int) time.Duration {
    start := time.Now()
    Fibonacci(n)
    return time.Since(start)
}

// ==================== 测试用例 ====================

/*
concurrentTest - 并发安全测试
验证Counter在并发场景下的正确性
使用WaitGroup同步多个goroutine
*/
func concurrentTest() {
    counter := NewCounter()
    var wg sync.WaitGroup
    
    const (
        numGoroutines = 1000
        incrementsPer = 100
    )
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < incrementsPer; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    
    expected := numGoroutines * incrementsPer
    actual := counter.Value()
    
    fmt.Printf("并发测试: 期望 %d, 实际 %d\n", expected, actual)
}

/*
performanceComparison - 性能对比测试
比较递归和迭代斐波那契实现的性能
计算加速比展示优化效果
*/
func performanceComparison() {
    fmt.Println("\n性能比较测试:")
    for n := 10; n <= 40; n += 10 {
        recursiveTime := benchmarkFibonacci(n)
        optimizedTime := benchmarkFibonacciOptimized(n)
        
        speedup := float64(recursiveTime) / float64(optimizedTime)
        fmt.Printf("n=%d: 递归: %v, 优化: %v, 加速比: %.1fx\n", 
            n, recursiveTime, optimizedTime, speedup)
    }
}

/*
testCoverageExample - 测试覆盖率示例
演示如何测试边界条件
包括空字符串、单字符、回文等情况
*/
func testCoverageExample() {
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
        {"回文测试", func() error {
            s := "racecar"
            if ReverseString(s) != s {
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

// ==================== 主函数 ====================

func main() {
    fmt.Println("=== Go测试与基准测试 ===")
    
    // 运行各类测试
    fmt.Println("\n1. 基本功能验证")
    fmt.Printf("Fibonacci(10)=%d\n", Fibonacci(10))
    
    fmt.Println("\n2. 并发安全测试")
    concurrentTest()
    
    fmt.Println("\n3. 性能对比")
    performanceComparison()
    
    fmt.Println("\n4. 测试覆盖率示例")
    testCoverageExample()
    
    fmt.Println("\n=== 测试完成 ===")
}

/*
boolToString - 布尔值转字符串
辅助函数，用于测试输出
*/
func boolToString(b bool) string {
    if b {
        return "true"
    }
    return "false"
}

/*
min - 取最小值
辅助函数，用于字符串截取
*/
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}