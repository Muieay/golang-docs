# Golang并发编程

## Goroutine基础

### 什么是Goroutine

Goroutine是Go语言中的轻量级线程，由Go运行时管理，创建成本极低，可以轻松创建成千上万个goroutine。

### 创建Goroutine

使用 `go` 关键字可以轻松启动一个新的Goroutine。为了等待Goroutine执行完成，推荐使用 `sync.WaitGroup`。

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func sayHello(name string, wg *sync.WaitGroup) {
    defer wg.Done() // Goroutine完成后，计数器减一
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s!\n", name)
        time.Sleep(time.Millisecond * 100)
    }
}

func main() {
    var wg sync.WaitGroup // 创建一个WaitGroup
    
    wg.Add(1) // 计数器加一
    // 启动一个goroutine
    go sayHello("Alice", &wg)
    
    // 主goroutine可以做其他事情
    fmt.Println("Hello from main goroutine.")
    
    // 等待所有goroutine完成
    wg.Wait()
    fmt.Println("All goroutines finished.")
}
```

### 匿名函数Goroutine

```go
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    
    go func(msg string) {
        defer wg.Done()
        fmt.Println(msg)
    }("Hello from anonymous goroutine!")
    
    wg.Wait() // 等待goroutine完成
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

> **注意**：对无缓冲Channel的发送操作会阻塞，直到另一个Goroutine对该Channel进行接收操作。同样，接收操作也会阻塞，直到另一个Goroutine进行发送。

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
    
    // 接收关闭的channel会立即返回零值
    fmt.Println(<-ch)  // 0
    
    // 使用 "comma ok" idiom 判断channel是否关闭
    val, ok := <-ch
    if !ok {
        fmt.Println("Channel is closed!") // Channel is closed!
    }
    fmt.Println(val) // 0
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

当 `select` 等待的所有 `case` 都无法立即执行时，如果设置了超时 `case`，则在指定时间后会执行该 `case`。

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

### 非阻塞操作

`select` 的 `default` 子句使其变为非阻塞的。如果没有任何 `case` 可以立即执行，`default` 子句就会被执行。

```go
func main() {
    ch := make(chan int, 1)
    
    // 非阻塞发送
    select {
    case ch <- 1:
        fmt.Println("Sent 1")
    default:
        fmt.Println("Channel is full!")
    }

    // 非阻塞接收
    select {
    case val := <-ch:
        fmt.Printf("Received %d\n", val)
    default:
        fmt.Println("Channel is empty!")
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

Fan-in（扇入）是一种将多个输入Channel合并到一个输出Channel的并发模式。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(name string, ch chan<- int) {
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Printf("Producer %s sent %d\n", name, i)
		time.Sleep(time.Millisecond * 100)
	}
}

func fanIn(ch1, ch2 <-chan int) <-chan int {
	mergedCh := make(chan int)
	go func() {
		defer close(mergedCh)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for val := range ch1 {
				mergedCh <- val
			}
		}()

		go func() {
			defer wg.Done()
			for val := range ch2 {
				mergedCh <- val
			}
		}()

		wg.Wait()
	}()
	return mergedCh
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		producer("A", ch1)
	}()
	go func() {
		defer close(ch2)
		producer("B", ch2)
	}()

	merged := fanIn(ch1, ch2)

	for val := range merged {
		fmt.Printf("Main received: %d\n", val)
	}
}
```

### Fan-out模式

Fan-out（扇出）是一种将一个输入Channel分发到多个输出Channel或由多个Goroutine处理的模式。

```go
package main

import (
	"fmt"
	"sync"
)

func consumer(name string, ch <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for value := range ch {
        fmt.Printf("Consumer %s got %d\n", name, value)
    }
}

func main() {
    ch := make(chan int)
    var wg sync.WaitGroup
    
    // 启动两个consumer
    wg.Add(2)
    go consumer("A", ch, &wg)
    go consumer("B", ch, &wg)
    
    // 发送数据
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
    
    // 等待所有consumer处理完毕
    wg.Wait()
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

`RWMutex`（读写锁）允许多个读操作同时进行，但写操作是互斥的。适用于读多写少的场景，可以提高性能。

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

`sync/atomic` 包提供了低级的原子内存操作，对于简单的计数器等场景，比使用 `Mutex` 性能更好。

```go
import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
    var counter int64 = 0
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
```

## Context

`context` 包定义了 `Context` 类型，它可以在API边界之间和进程之间传递截止日期、取消信号和其他请求范围的值。这对于管理服务器中的并发请求或长时间运行的操作至关重要。

### 基本用法：取消Goroutine

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done(): // 监听取消信号
            fmt.Printf("%s: worker cancelled\n", name)
            return
        default:
            fmt.Printf("%s: working...\n", name)
            time.Sleep(time.Second)
        }
    }
}

func main() {
    // 创建一个可取消的context
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx, "Worker 1")
    
    // 运行3秒后取消
    time.Sleep(time.Second * 3)
    cancel() // 发出取消信号
    
    // 等待goroutine退出
    time.Sleep(time.Second)
    fmt.Println("Main goroutine finished.")
}
```

### 超时控制

`context.WithTimeout` 和 `context.WithDeadline` 可用于设置超时。

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 创建一个3秒后超时的context
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel() // 及时释放资源

    select {
    case <-time.After(5 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err()) // prints "context deadline exceeded"
    }
}
```

### 传递值

`context.WithValue` 可以将键值对附加到 `Context` 中，实现跨API边界传递请求范围的数据。

> **注意**：应谨慎使用 `WithValue`，只传递请求范围的数据，不要用它来传递可选参数。key的类型应该是自定义类型，以避免冲突。

```go
package main

import (
    "context"
    "fmt"
)

type keyType string

const requestIDKey keyType = "requestID"

func processRequest(ctx context.Context) {
    // 从context中获取值
    requestID := ctx.Value(requestIDKey)
    if requestID != nil {
        fmt.Printf("Processing request with ID: %s\n", requestID)
    } else {
        fmt.Println("No request ID found.")
    }
}

func main() {
    // 创建一个带有值的context
    ctx := context.WithValue(context.Background(), requestIDKey, "12345")
    
    processRequest(ctx)
}
```