# Golang并发编程

## Goroutine基础

### 什么是Goroutine

Goroutine是Go语言中的轻量级线程，由Go运行时管理，创建成本极低，可以轻松创建成千上万个goroutine。

### 创建Goroutine

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 5; i++ {
        fmt.Printf("Hello, %s!\n", name)
        time.Sleep(time.Millisecond * 100)
    }
}

func main() {
    // 启动一个goroutine
    go sayHello("Alice")
    
    // 主goroutine继续执行
    sayHello("Bob")
    
    // 等待goroutine完成
    time.Sleep(time.Second * 1)
}
```

### 匿名函数Goroutine

```go
func main() {
    go func(msg string) {
        fmt.Println(msg)
    }("Hello from anonymous goroutine!")
    
    time.Sleep(time.Second)
}
```

## Channel通信

### Channel基础

Channel是goroutine之间的通信机制，遵循CSP（Communicating Sequential Processes）模型。

```go
package main

import "fmt"

func main() {
    // 创建channel
    ch := make(chan int)
    
    // 启动goroutine发送数据
    go func() {
        ch <- 42  // 发送数据到channel
    }()
    
    // 接收数据
    value := <-ch
    fmt.Println("Received:", value)
}
```

### Channel类型

#### 无缓冲Channel

```go
ch := make(chan int)  // 无缓冲，发送和接收必须同时进行
```

#### 有缓冲Channel

```go
ch := make(chan int, 3)  // 缓冲区大小为3
```

### Channel方向

```go
// 只发送
func send(ch chan<- int, value int) {
    ch <- value
}

// 只接收
func receive(ch <-chan int) int {
    return <-ch
}
```

### 关闭Channel

```go
func main() {
    ch := make(chan int, 2)
    
    ch <- 1
    ch <- 2
    close(ch)  // 关闭channel
    
    // 仍然可以接收已发送的数据
    fmt.Println(<-ch)  // 1
    fmt.Println(<-ch)  // 2
    
    // 接收关闭的channel会返回零值
    fmt.Println(<-ch)  // 0
}
```

### range遍历Channel

```go
func main() {
    ch := make(chan int)
    
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch)
    }()
    
    // range自动接收直到channel关闭
    for value := range ch {
        fmt.Println(value)
    }
}
```

## Select语句

### 基本用法

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(time.Second * 1)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(time.Second * 2)
        ch2 <- "from ch2"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }
}
```

### 超时处理

```go
func main() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    case <-time.After(time.Second * 3):
        fmt.Println("Timeout!")
    }
}
```

### 非阻塞通信

```go
func main() {
    ch := make(chan int, 1)
    ch <- 1
    
    select {
    case ch <- 2:
        fmt.Println("Sent 2")
    default:
        fmt.Println("Channel is full!")
    }
}
```

## 并发模式

### Worker Pool模式

```go
package main

import (
    "fmt"
    "time"
)

type Job struct {
    ID   int
    Data int
}

type Result struct {
    JobID  int
    Output int
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        // 模拟处理时间
        time.Sleep(time.Millisecond * 100)
        result := Result{JobID: job.ID, Output: job.Data * 2}
        results <- result
    }
}

func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // 启动3个worker
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // 发送5个任务
    for j := 1; j <= 5; j++ {
        jobs <- Job{ID: j, Data: j}
    }
    close(jobs)
    
    // 收集结果
    for r := 1; r <= 5; r++ {
        result := <-results
        fmt.Printf("Job %d processed with result %d\n", result.JobID, result.Output)
    }
}
```

### Fan-in模式

```go
func producer(name string, ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
        time.Sleep(time.Millisecond * 100)
    }
}

func main() {
    ch := make(chan int)
    
    go producer("A", ch)
    go producer("B", ch)
    
    for i := 0; i < 10; i++ {
        fmt.Println(<-ch)
    }
}
```

### Fan-out模式

```go
func consumer(name string, ch <-chan int) {
    for value := range ch {
        fmt.Printf("Consumer %s got %d\n", name, value)
    }
}

func main() {
    ch := make(chan int)
    
    go consumer("A", ch)
    go consumer("B", ch)
    
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
    
    time.Sleep(time.Second)
}
```

## 并发安全

### Mutex互斥锁

```go
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Final count:", counter.Get())
}
```

### RWMutex读写锁

```go
type SafeMap struct {
    mu    sync.RWMutex
    data  map[string]int
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    val, ok := sm.data[key]
    return val, ok
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.data[key] = value
}
```

### Atomic原子操作

```go
import "sync/atomic"

func main() {
    var counter int64 = 0
    
    for i := 0; i < 1000; i++ {
        go func() {
            atomic.AddInt64(&counter, 1)
        }()
    }
    
    time.Sleep(time.Second)
    fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
```