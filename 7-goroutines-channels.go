package main

import (
    "fmt"
    "sync"
    "time"
)

// 并发编程：goroutine和channel
// Go语言的并发是其核心特性

func main() {
    fmt.Println("=== Go语言并发编程 ===")
    
    // 1. goroutine基础
    fmt.Println("\n--- goroutine基础 ---")
    
    // 启动goroutine
    go sayHello("goroutine 1")
    go sayHello("goroutine 2")
    
    // 主goroutine继续执行
    time.Sleep(1 * time.Second)
    fmt.Println("主goroutine结束")
    
    // 2. 使用WaitGroup等待goroutine
    fmt.Println("\n--- WaitGroup示例 ---")
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
    fmt.Println("所有工作完成")
    
    // 3. channel基础
    fmt.Println("\n--- channel基础 ---")
    
    // 创建channel
    ch := make(chan int)
    
    // 启动goroutine发送数据
    go func() {
        ch <- 42  // 发送数据
    }()
    
    // 接收数据
    value := <-ch
    fmt.Printf("接收到: %d\n", value)
    
    // 4. 缓冲channel
    fmt.Println("\n--- 缓冲channel ---")
    bufferedCh := make(chan int, 3)
    
    bufferedCh <- 1
    bufferedCh <- 2
    bufferedCh <- 3
    
    fmt.Printf("从缓冲channel接收: %d\n", <-bufferedCh)
    fmt.Printf("从缓冲channel接收: %d\n", <-bufferedCh)
    
    // 5. 关闭channel
    fmt.Println("\n--- 关闭channel ---")
    closeCh := make(chan int)
    
    go func() {
        for i := 1; i <= 5; i++ {
            closeCh <- i
        }
        close(closeCh)  // 关闭channel
    }()
    
    // 使用range接收数据直到channel关闭
    for num := range closeCh {
        fmt.Printf("接收到: %d\n", num)
    }
    
    // 6. select语句
    fmt.Println("\n--- select语句 ---")
    
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "来自ch1的消息"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "来自ch2的消息"
    }()
    
    // 使用select等待多个channel
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("收到:", msg1)
        case msg2 := <-ch2:
            fmt.Println("收到:", msg2)
        case <-time.After(3 * time.Second):
            fmt.Println("超时")
        }
    }
    
    // 7. 实际应用：并发计算
    fmt.Println("\n--- 并发计算示例 ---")
    
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 并发计算平方和
    result := concurrentSquareSum(numbers)
    fmt.Printf("并发计算平方和: %d\n", result)
    
    // 8. 生产者-消费者模式
    fmt.Println("\n--- 生产者-消费者模式 ---")
    
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // 启动3个worker
    for w := 1; w <= 3; w++ {
        go workerProducerConsumer(w, jobs, results)
    }
    
    // 发送任务
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)
    
    // 收集结果
    for r := 1; r <= 9; r++ {
        result := <-results
        fmt.Printf("结果: %d\n", result)
    }
    
    // 9. 并发安全的数据结构
    fmt.Println("\n--- 并发安全的数据结构 ---")
    
    safeCounter := NewSafeCounter()
    
    for i := 0; i < 1000; i++ {
        go func() {
            safeCounter.Increment()
        }()
    }
    
    time.Sleep(1 * time.Second)
    fmt.Printf("最终计数: %d\n", safeCounter.Value())
    
    // 10. 并发错误处理
    fmt.Println("\n--- 并发错误处理 ---")
    
    errorCh := make(chan error, 1)
    go func() {
        time.Sleep(500 * time.Millisecond)
        errorCh <- fmt.Errorf("模拟错误")
    }()
    
    if err := <-errorCh; err != nil {
        fmt.Printf("收到错误: %v\n", err)
    }
    
    // 练习
    fmt.Println("\n练习：")
    fmt.Println("1. 使用goroutine并发下载多个网页")
    fmt.Println("2. 实现一个并发安全的队列")
    fmt.Println("3. 使用channel实现任务调度器")
}

// 简单的goroutine函数
func sayHello(name string) {
    fmt.Printf("Hello from %s\n", name)
}

// 工作goroutine
func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d 开始工作\n", id)
    time.Sleep(time.Duration(id) * 500 * time.Millisecond)
    fmt.Printf("Worker %d 完成工作\n", id)
}

// 并发计算平方和
func concurrentSquareSum(numbers []int) int {
    ch := make(chan int)
    var wg sync.WaitGroup
    
    for _, num := range numbers {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            ch <- n * n
        }(num)
    }
    
    go func() {
        wg.Wait()
        close(ch)
    }()
    
    sum := 0
    for square := range ch {
        sum += square
    }
    
    return sum
}

// 生产者-消费者模式的worker
func workerProducerConsumer(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("worker %d 处理任务 %d\n", id, j)
        time.Sleep(time.Millisecond * 100)
        results <- j * 2  // 模拟处理结果
    }
}

// 并发安全的计数器
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func NewSafeCounter() *SafeCounter {
    return &SafeCounter{}
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// 并发安全的队列
func demonstrateSafeQueue() {
    fmt.Println("\n--- 并发安全的队列 ---")
    
    queue := NewSafeQueue()
    
    var wg sync.WaitGroup
    
    // 生产者
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(val int) {
            defer wg.Done()
            queue.Enqueue(val)
            fmt.Printf("入队: %d\n", val)
        }(i)
    }
    
    // 消费者
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            time.Sleep(100 * time.Millisecond)
            if val, ok := queue.Dequeue(); ok {
                fmt.Printf("出队: %d\n", val)
            }
        }()
    }
    
    wg.Wait()
}

// 并发安全的队列实现
type SafeQueue struct {
    items []int
    mu    sync.Mutex
}

func NewSafeQueue() *SafeQueue {
    return &SafeQueue{items: make([]int, 0)}
}

func (q *SafeQueue) Enqueue(item int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
}

func (q *SafeQueue) Dequeue() (int, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.items) == 0 {
        return 0, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}