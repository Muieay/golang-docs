package main

import (
    "fmt"
    "sync"
    "time"
)

/*
Go并发编程核心概念：
1. goroutine - 轻量级线程，由Go运行时管理
2. channel - goroutine间的通信管道
3. sync包 - 提供同步原语如WaitGroup和Mutex
*/

func main() {
    fmt.Println("=== Go语言并发编程 ===")
    
    /* 
    1. goroutine基础
    go关键字启动goroutine，非阻塞执行
    注意：主goroutine退出会导致所有子goroutine终止
    */
    fmt.Println("\n--- goroutine基础 ---")
    go sayHello("goroutine 1") // 启动goroutine1
    go sayHello("goroutine 2") // 启动goroutine2
    time.Sleep(1 * time.Second) // 等待goroutine执行
    
    /*
    2. WaitGroup使用
    用于等待一组goroutine完成
    Add()增加计数器，Done()减少计数器，Wait()阻塞直到计数器为0
    */
    fmt.Println("\n--- WaitGroup示例 ---")
    var wg sync.WaitGroup
    for i := 1; i <= 5; i++ {
        wg.Add(1) // 计数器+1
        go worker(i, &wg) // 传递指针避免复制
    }
    wg.Wait() // 等待所有worker完成
    
    /*
    3. channel基础
    make(chan T)创建无缓冲channel
    发送和接收会阻塞直到另一端准备好
    */
    fmt.Println("\n--- channel基础 ---")
    ch := make(chan int) // 创建int类型channel
    go func() { ch <- 42 }() // 发送数据
    value := <-ch // 接收数据
    fmt.Printf("接收到: %d\n", value)
    
    /*
    4. 缓冲channel
    make(chan T, size)创建带缓冲区的channel
    缓冲区满时发送才会阻塞，空时接收才会阻塞
    */
    fmt.Println("\n--- 缓冲channel ---")
    bufferedCh := make(chan int, 3) // 缓冲区大小3
    bufferedCh <- 1; bufferedCh <- 2; bufferedCh <- 3 // 不会阻塞
    
    /*
    5. 关闭channel
    close(ch)关闭channel
    range会自动检测channel关闭
    */
    fmt.Println("\n--- 关闭channel ---")
    closeCh := make(chan int)
    go func() {
        for i := 1; i <= 5; i++ { closeCh <- i }
        close(closeCh) // 关闭channel
    }()
    for num := range closeCh { // 自动检测关闭
        fmt.Printf("接收到: %d\n", num)
    }
    
    /*
    6. select语句
    监听多个channel操作
    case处理就绪的channel操作
    time.After实现超时机制
    */
    fmt.Println("\n--- select语句 ---")
    ch1, ch2 := make(chan string), make(chan string)
    go func() { time.Sleep(1 * time.Second); ch1 <- "消息1" }()
    go func() { time.Sleep(2 * time.Second); ch2 <- "消息2" }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1: fmt.Println("收到:", msg1)
        case msg2 := <-ch2: fmt.Println("收到:", msg2)
        case <-time.After(3 * time.Second): fmt.Println("超时")
        }
    }
    
    /*
    7. 并发计算示例
    使用goroutine并发计算，通过channel收集结果
    WaitGroup确保所有计算完成
    */
    fmt.Println("\n--- 并发计算示例 ---")
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    fmt.Printf("平方和: %d\n", concurrentSquareSum(numbers))
    
    /*
    8. 生产者-消费者模式
    jobs channel分发任务
    results channel收集结果
    多个worker并发处理
    */
    fmt.Println("\n--- 生产者-消费者模式 ---")
    jobs, results := make(chan int, 100), make(chan int, 100)
    for w := 1; w <= 3; w++ { go workerProducerConsumer(w, jobs, results) }
    for j := 1; j <= 9; j++ { jobs <- j } // 发送任务
    close(jobs) // 关闭任务channel
    for r := 1; r <= 9; r++ { fmt.Printf("结果: %d\n", <-results) }
    
    /*
    9. 并发安全数据结构
    使用sync.Mutex保护共享数据
    Lock()/Unlock()确保原子操作
    */
    fmt.Println("\n--- 并发安全的数据结构 ---")
    counter := NewSafeCounter()
    for i := 0; i < 1000; i++ { go counter.Increment() }
    time.Sleep(1 * time.Second) // 等待所有goroutine完成
    fmt.Printf("计数器值: %d\n", counter.Value())
    
    /*
    10. 并发错误处理
    专用channel传递错误
    主goroutine统一处理错误
    */
    fmt.Println("\n--- 并发错误处理 ---")
    errCh := make(chan error, 1)
    go func() { errCh <- fmt.Errorf("模拟错误") }()
    if err := <-errCh; err != nil { fmt.Printf("错误: %v\n", err) }
}

/* 
sayHello - 演示goroutine的基本使用
参数name标识不同的goroutine
*/
func sayHello(name string) {
    fmt.Printf("Hello from %s\n", name)
}

/*
worker - 演示WaitGroup的使用
参数:
  id - worker标识
  wg - WaitGroup指针(必须传指针)
使用defer确保Done()调用
*/
func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // 确保计数器减1
    fmt.Printf("Worker %d 开始\n", id)
    time.Sleep(time.Duration(id) * 500 * time.Millisecond)
    fmt.Printf("Worker %d 完成\n", id)
}

/*
concurrentSquareSum - 并发计算平方和
1. 为每个数字启动goroutine计算平方
2. 通过channel收集结果
3. 使用WaitGroup等待所有计算完成
4. 关闭channel后汇总结果
*/
func concurrentSquareSum(numbers []int) int {
    ch := make(chan int)
    var wg sync.WaitGroup
    
    // 启动计算goroutine
    for _, num := range numbers {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            ch <- n * n // 发送结果
        }(num)
    }
    
    // 关闭channel的goroutine
    go func() {
        wg.Wait()
        close(ch) // 所有计算完成后关闭
    }()
    
    // 汇总结果
    sum := 0
    for square := range ch { // 自动结束于channel关闭
        sum += square
    }
    return sum
}

/*
workerProducerConsumer - 生产者消费者worker
参数:
  id - worker标识
  jobs - 只接收channel (<-chan)
  results - 只发送channel (chan<-)
从jobs接收任务，处理后将结果发送到results
*/
func workerProducerConsumer(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs { // 自动结束于jobs关闭
        fmt.Printf("worker %d 处理任务 %d\n", id, j)
        time.Sleep(100 * time.Millisecond)
        results <- j * 2 // 发送处理结果
    }
}

/*
SafeCounter - 并发安全计数器
使用Mutex保护value字段
*/
type SafeCounter struct {
    mu    sync.Mutex // 互斥锁
    value int        // 受保护的值
}

func NewSafeCounter() *SafeCounter {
    return &SafeCounter{}
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()         // 获取锁
    defer c.mu.Unlock() // 确保释放锁
    c.value++           // 安全操作
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

/*
SafeQueue - 并发安全队列
使用Mutex保护slice操作
*/
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