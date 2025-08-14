package main

// 12-advanced-topics.go - Go语言高级特性详解
// ========================================
//
// 文件功能：
// 本文件深入讲解Go语言的高级特性和最佳实践，包括反射、泛型、上下文管理、
// 网络编程、加密安全、微服务架构等高级主题。
//
// 学习目标：
// 1. 掌握反射机制及其应用场景
// 2. 理解并使用Go 1.18+的泛型编程特性
// 3. 学会使用上下文管理控制程序执行流程
// 4. 掌握网络编程基础
// 5. 了解加密和安全相关的标准库
// 6. 理解微服务架构的基本组件
// 7. 学会配置管理和日志系统的设计
//
// 代码组织：
// 文件按主题分为多个部分，每个部分包含相关的示例代码和说明。

import (
    "context"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "sync"
    "time"
)

// 高级主题
// Go语言的高级特性和最佳实践

// 1. 反射和元编程
// 反射是Go语言中一种强大的机制，允许程序在运行时检查和操作类型信息
// 它使得我们可以编写通用的代码，处理不同类型的数据
// 反射的核心包是reflect，提供了获取类型信息、创建实例、修改变量等功能
import "reflect"

// 2. 泛型编程（Go 1.18+）

// 3. 上下文管理

// 4. 网络编程

// 5. 加密和安全

// 6. 微服务架构

// 1. 反射示例
// Person 是一个用于演示反射操作的示例结构体
// 包含了不同类型的字段和JSON标签
// 注意：反射可以访问私有字段，但不能修改它们
 type Person struct {
    Name    string `json:"name"` // 姓名字段，带JSON标签
    Age     int    `json:"age"`  // 年龄字段，带JSON标签
    Email   string `json:"email"` // 邮箱字段，带JSON标签
    private string // 私有字段，反射可以读取但不能修改
}

// 使用反射分析结构体
// analyzeStruct 函数通过反射分析任意对象的类型信息
// 参数 obj 是任意类型的值
// 该函数会打印出对象的类型、种类、字段数量以及每个字段的详细信息
func analyzeStruct(obj interface{}) {
    val := reflect.ValueOf(obj) // 获取值信息
    typ := reflect.TypeOf(obj)  // 获取类型信息
    
    fmt.Printf("类型: %s\n", typ)     // 打印完整类型名
    fmt.Printf("种类: %s\n", typ.Kind()) // 打印类型种类(如struct, int, string等)
    
    // 如果是结构体类型，则分析其字段
    if typ.Kind() == reflect.Struct {
        fmt.Printf("字段数量: %d\n", typ.NumField())
        
        for i := 0; i < typ.NumField(); i++ {
            field := typ.Field(i)  // 获取字段类型信息
            value := val.Field(i)  // 获取字段值信息
            
            fmt.Printf("  字段 %d: %s (%s) = %v, 标签: %s\n", 
                i, field.Name, field.Type, value, field.Tag)
        }
    }
}

// 动态创建结构体实例
// createInstanceFromType 函数根据给定的类型动态创建一个结构体实例
// 参数 typ 是要创建的类型
// 返回值是创建的实例(interface{}类型)
// 注意：该函数只能创建结构体类型的实例
func createInstanceFromType(typ reflect.Type) interface{} {
    // 检查是否为结构体类型
    if typ.Kind() != reflect.Struct {
        return nil
    }
    
    // 创建结构体指针并获取其元素
    instance := reflect.New(typ).Elem()
    
    // 遍历结构体字段并设置默认值
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        fieldVal := instance.Field(i)
        
        // 根据字段类型设置不同的默认值
        switch field.Type.Kind() {
        case reflect.String:
            fieldVal.SetString("default") // 字符串默认值
        case reflect.Int:
            fieldVal.SetInt(0)             // 整数默认值
        case reflect.Bool:
            fieldVal.SetBool(false)        // 布尔值默认值
        // 可以根据需要添加更多类型的处理
        }
    }
    
    // 返回接口类型的实例
    return instance.Interface()
}

// 2. 泛型编程示例
// 泛型是Go 1.18版本引入的新特性，允许我们编写不依赖于特定类型的代码
// 通过类型参数，我们可以创建适用于多种类型的通用函数和数据结构
// 这提高了代码的复用性并减少了重复代码

// 泛型切片操作
// Filter 函数根据给定的谓词函数过滤切片中的元素
// 类型参数 T 表示切片元素的类型
// 参数 slice 是要过滤的切片
// 参数 predicate 是一个函数，用于判断元素是否应该被保留
// 返回值是包含满足条件的元素的新切片
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range slice {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}

// Map 函数将一个切片的元素通过转换函数映射到另一个切片
// 类型参数 T 表示输入切片元素的类型
// 类型参数 U 表示输出切片元素的类型
// 参数 slice 是要转换的输入切片
// 参数 transform 是一个函数，用于将 T 类型转换为 U 类型
// 返回值是包含转换后元素的新切片
func Map[T any, U any](slice []T, transform func(T) U) []U {
    result := make([]U, len(slice)) // 预分配结果切片
    for i, item := range slice {
        result[i] = transform(item)
    }
    return result
}

// Reduce 函数通过累加器函数将切片元素 reduce 为单个值
// 类型参数 T 表示切片元素的类型
// 类型参数 U 表示结果的类型
// 参数 slice 是要 reduce 的切片
// 参数 initial 是初始累加值
// 参数 reducer 是一个函数，用于将当前累加值和下一个元素结合
// 返回值是最终的累加结果
func Reduce[T any, U any](slice []T, initial U, reducer func(U, T) U) U {
    result := initial
    for _, item := range slice {
        result = reducer(result, item)
    }
    return result
}

// 泛型栈
// Stack 是一个泛型栈实现，可以存储任意类型的元素
// 类型参数 T 表示栈中元素的类型
// items 是用于存储栈元素的切片
// 栈的特点是后进先出(LIFO)
type Stack[T any] struct {
    items []T
}

// NewStack 创建并返回一个新的泛型栈
// 类型参数 T 表示栈中元素的类型
// 返回值是一个指向新创建的 Stack 的指针
func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

// Push 向栈顶添加一个元素
// 参数 item 是要添加的元素
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop 从栈顶移除并返回一个元素
// 返回值1是弹出的元素，如果栈为空则返回T类型的零值
// 返回值2是一个布尔值，表示操作是否成功
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T // T类型的零值
        return zero, false
    }
    
    item := s.items[len(s.items)-1] // 获取栈顶元素
    s.items = s.items[:len(s.items)-1] // 移除栈顶元素
    return item, true
}

// Peek 返回栈顶元素但不移除
// 返回值1是栈顶元素，如果栈为空则返回T类型的零值
// 返回值2是一个布尔值，表示操作是否成功
func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// Size 返回栈中元素的数量
func (s *Stack[T]) Size() int {
    return len(s.items)
}

// 泛型队列
// Queue 是一个泛型队列实现，可以存储任意类型的元素
// 类型参数 T 表示队列中元素的类型
// items 是用于存储队列元素的切片
// 队列的特点是先进先出(FIFO)
type Queue[T any] struct {
    items []T
}

// NewQueue 创建并返回一个新的泛型队列
// 类型参数 T 表示队列中元素的类型
// 返回值是一个指向新创建的 Queue 的指针
func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{items: make([]T, 0)}
}

// Enqueue 向队列尾部添加一个元素
// 参数 item 是要添加的元素
func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

// Dequeue 从队列头部移除并返回一个元素
// 返回值1是出队的元素，如果队列为空则返回T类型的零值
// 返回值2是一个布尔值，表示操作是否成功
func (q *Queue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    
    item := q.items[0]  // 获取队列头部元素
    q.items = q.items[1:] // 移除队列头部元素
    return item, true
}

// Size 返回队列中元素的数量
func (q *Queue[T]) Size() int {
    return len(q.items)
}

// 泛型缓存
// Cache 是一个线程安全的泛型缓存实现
// 类型参数 K 表示缓存键的类型，必须是可比较的
// 类型参数 V 表示缓存值的类型
// items 是用于存储缓存项的映射
// mu 是读写锁，用于保证并发安全
import "sync"

type Cache[K comparable, V any] struct {
    items map[K]V
    mu    sync.RWMutex
}

// NewCache 创建并返回一个新的泛型缓存
// 类型参数 K 表示缓存键的类型，必须是可比较的
// 类型参数 V 表示缓存值的类型
// 返回值是一个指向新创建的 Cache 的指针
func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]V),
    }
}

// Set 设置缓存项
// 参数 key 是缓存键
// 参数 value 是缓存值
// 该方法使用写锁保证并发安全
func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}

// Get 获取缓存项
// 参数 key 是要查找的缓存键
// 返回值1是缓存值，如果键不存在则返回V类型的零值
// 返回值2是一个布尔值，表示键是否存在
// 该方法使用读锁保证并发安全
func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.items[key]
    return value, exists
}

// Delete 删除缓存项
// 参数 key 是要删除的缓存键
// 该方法使用写锁保证并发安全
func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

// 3. 上下文管理示例
// 上下文(Context)是Go语言中用于控制goroutine生命周期和传递请求相关值的机制
// 它可以用于超时控制、取消操作以及在不同goroutine之间传递数据
// 标准库中的context包提供了相关功能

// contextExample 演示了带超时的上下文的使用
func contextExample() {
    fmt.Println("\n=== 上下文管理示例 ===")
    
    // 创建带超时的上下文
    // context.Background() 返回一个空上下文，通常作为根上下文
    // WithTimeout 创建一个会在指定时间后自动取消的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // 确保在函数退出时取消上下文，避免资源泄漏
    
    // 模拟长时间操作
    go func() {
        select {
        case <-time.After(3 * time.Second): // 模拟3秒的操作
            fmt.Println("操作完成")
        case <-ctx.Done(): // 监听上下文取消信号
            fmt.Println("操作超时:", ctx.Err())
        }
    }()
    
    // 等待操作完成或超时
    time.Sleep(3 * time.Second)
}

// 使用上下文取消操作
// longRunningOperation 模拟一个长时间运行的操作
// 参数 ctx 是用于控制操作取消的上下文
// 返回值是错误，如果操作被取消则返回上下文的错误
func longRunningOperation(ctx context.Context) error {
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():// 监听上下文取消信号
            return ctx.Err()
        case <-time.After(500 * time.Millisecond):// 模拟500毫秒的工作
            fmt.Printf("操作进度: %d%%\n", (i+1)*10)
        }
    }
    return nil
}

// 4. 网络编程示例
// 网络编程是Go语言的一个重要应用领域，标准库中的net包提供了丰富的网络编程接口
// 下面的示例演示了如何创建TCP服务器和客户端

// tcpServerExample 演示了TCP服务器和客户端的基本用法
func tcpServerExample() {
    fmt.Println("\n=== TCP服务器示例 ===")
    
    // 创建TCP监听器
    // net.Listen 创建一个TCP监听器，监听指定地址和端口
    listener, err := net.Listen("tcp", ":8081")
    if err != nil {
        fmt.Printf("创建监听器失败: %v\n", err)
        return
    }
    defer listener.Close() // 确保在函数退出时关闭监听器
    
    fmt.Println("TCP服务器启动在 :8081")
    
    // 在后台接受连接
    go func() {
        for {
            // 接受客户端连接
            conn, err := listener.Accept()
            if err != nil {
                fmt.Printf("接受连接失败: %v\n", err)
                continue
            }
            
            // 为每个连接创建一个新的goroutine处理
            go handleTCPConnection(conn)
        }
    }()
    
    // 模拟客户端连接
    time.Sleep(100 * time.Millisecond)
    
    // 建立TCP连接
    conn, err := net.Dial("tcp", "localhost:8081")
    if err != nil {
        fmt.Printf("连接服务器失败: %v\n", err)
        return
    }
    
    // 发送消息给服务器
    fmt.Fprintln(conn, "Hello from client")
    
    // 读取服务器响应
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("读取响应失败: %v\n", err)
        return
    }
    
    fmt.Printf("收到服务器响应: %s\n", string(buffer[:n]))
    conn.Close() // 关闭连接
    
    time.Sleep(100 * time.Millisecond)
}

// handleTCPConnection 处理单个TCP连接
// 参数 conn 是客户端连接
func handleTCPConnection(conn net.Conn) {
    defer conn.Close() // 确保在函数退出时关闭连接
    
    // 读取客户端发送的数据
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("读取数据失败: %v\n", err)
        return
    }
    
    // 处理接收到的数据
    message := string(buffer[:n])
    fmt.Printf("收到消息: %s\n", message)
    
    // 发送响应给客户端
    response := fmt.Sprintf("服务器收到: %s", message)
    conn.Write([]byte(response))
}

// 5. 加密和安全示例
// 加密和安全是网络应用中的重要考虑因素
// Go语言的标准库提供了多种加密算法和安全相关的功能
// 下面的示例演示了如何使用AES-GCM模式进行加密和解密

// encryptionExample 演示了AES-GCM加密和解密的基本用法
import "crypto/aes"
import "crypto/cipher"

func encryptionExample() {
    fmt.Println("\n=== 加密示例 ===")
    
    // 生成随机密钥
    // AES-256需要32字节的密钥
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        fmt.Printf("生成密钥失败: %v\n", err)
        return
    }
    
    // 要加密的数据
    plaintext := []byte("Hello, World! 这是一个需要加密的消息")
    
    // 创建AES加密器
    block, err := aes.NewCipher(key)
    if err != nil {
        fmt.Printf("创建加密器失败: %v\n", err)
        return
    }
    
    // 使用GCM模式
    // GCM(Galois/Counter Mode)是一种认证加密模式，同时提供 confidentiality和完整性
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        fmt.Printf("创建GCM失败: %v\n", err)
        return
    }
    
    // 生成随机nonce
    // nonce是一个只使用一次的随机数，用于确保相同的明文加密后得到不同的密文
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        fmt.Printf("生成nonce失败: %v\n", err)
        return
    }
    
    // 加密
    // Seal方法会将nonce附加到密文前面
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    fmt.Printf("加密后的数据: %x\n", ciphertext)
    
    // 解密
    // Open方法会从密文中提取nonce并使用它进行解密
    decrypted, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
    if err != nil {
        fmt.Printf("解密失败: %v\n", err)
        return
    }
    
    fmt.Printf("解密后的数据: %s\n", string(decrypted))
}

// 6. 微服务架构示例
// 微服务架构是一种将应用程序拆分为独立部署的小型服务的设计模式
// 这些服务通过网络进行通信，各自负责特定的业务功能
// 下面的示例演示了微服务架构中的服务发现和负载均衡组件

// 服务发现接口
// ServiceRegistry 定义了服务注册和发现的接口
// 服务发现是微服务架构中的关键组件，用于管理服务实例的生命周期和位置

type ServiceRegistry interface {
    // Register 注册一个服务实例
    // serviceName 是服务名称
    // address 是服务实例的地址
    Register(serviceName string, address string) error
    
    // Deregister 注销一个服务
    // serviceName 是要注销的服务名称
    Deregister(serviceName string) error
    
    // Discover 发现服务的所有实例
    // serviceName 是要发现的服务名称
    // 返回值是服务实例地址列表
    Discover(serviceName string) ([]string, error)
    
    // HealthCheck 检查服务是否健康
    // serviceName 是要检查的服务名称
    // 返回值是服务是否健康
    HealthCheck(serviceName string) bool
}

// 内存服务注册表
// InMemoryServiceRegistry 是一个基于内存的服务注册表实现
// 它实现了ServiceRegistry接口
// services 存储服务名称到服务实例地址列表的映射
// mu 是读写锁，用于保证并发安全
type InMemoryServiceRegistry struct {
    services map[string][]string
    mu       sync.RWMutex
}

// NewInMemoryServiceRegistry 创建并返回一个新的内存服务注册表
func NewInMemoryServiceRegistry() *InMemoryServiceRegistry {
    return &InMemoryServiceRegistry{
        services: make(map[string][]string),
    }
}

// Register 实现ServiceRegistry接口的Register方法
// 注册一个服务实例到注册表
func (r *InMemoryServiceRegistry) Register(serviceName string, address string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.services[serviceName] = append(r.services[serviceName], address)
    fmt.Printf("服务 %s 注册在 %s\n", serviceName, address)
    return nil
}

// Deregister 实现ServiceRegistry接口的Deregister方法
// 从注册表中注销一个服务
func (r *InMemoryServiceRegistry) Deregister(serviceName string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    delete(r.services, serviceName)
    fmt.Printf("服务 %s 已注销\n", serviceName)
    return nil
}

// Discover 实现ServiceRegistry接口的Discover方法
// 发现服务的所有实例地址
func (r *InMemoryServiceRegistry) Discover(serviceName string) ([]string, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    addresses, exists := r.services[serviceName]
    if !exists {
        return nil, fmt.Errorf("服务 %s 未找到", serviceName)
    }
    
    return addresses, nil
}

// HealthCheck 实现ServiceRegistry接口的HealthCheck方法
// 检查服务是否健康（在本例中，只需检查服务是否存在）
func (r *InMemoryServiceRegistry) HealthCheck(serviceName string) bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    _, exists := r.services[serviceName]
    return exists
}

// 负载均衡器
// LoadBalancer 是一个简单的负载均衡器实现
// registry 是服务注册表
// strategy 是负载均衡策略，支持 "round-robin", "random", "least-connections"
type LoadBalancer struct {
    registry ServiceRegistry
    strategy string // "round-robin", "random", "least-connections"
}

// NewLoadBalancer 创建并返回一个新的负载均衡器
// registry 是服务注册表
func NewLoadBalancer(registry ServiceRegistry) *LoadBalancer {
    return &LoadBalancer{
        registry: registry,
        strategy: "round-robin", // 默认使用轮询策略
    }
}

// GetServiceAddress 根据负载均衡策略获取一个服务实例地址
// serviceName 是服务名称
// 返回值是服务实例地址
func (lb *LoadBalancer) GetServiceAddress(serviceName string) (string, error) {
    addresses, err := lb.registry.Discover(serviceName)
    if err != nil {
        return "", err
    }
    
    if len(addresses) == 0 {
        return "", fmt.Errorf("没有可用的服务实例")
    }
    
    // 简单的轮询策略实现
    // 注意：这是一个简化版实现，真实的轮询策略需要维护当前索引
    return addresses[0], nil
}

// 7. 配置管理
// 配置管理是应用程序开发中的重要组成部分
// 下面的示例演示了如何定义和加载配置

// Config 定义了应用程序的配置结构
// 使用嵌套结构体组织不同部分的配置
// json标签用于JSON序列化和反序列化
type Config struct {
    Server struct {
        Host string `json:"host"` // 服务器主机名
        Port int    `json:"port"` // 服务器端口
    } `json:"server"`
    Database struct {
        Host     string `json:"host"`     // 数据库主机名
        Port     int    `json:"port"`     // 数据库端口
        User     string `json:"user"`     // 数据库用户名
        Password string `json:"password"` // 数据库密码
        Name     string `json:"name"`     // 数据库名称
    } `json:"database"`
    Cache struct {
        Host string `json:"host"` // 缓存主机名
        Port int    `json:"port"` // 缓存端口
    } `json:"cache"`
}

// loadConfig 从JSON数据加载配置
// jsonData 是包含配置的JSON字符串
// 返回值是解析后的配置对象和可能的错误
func loadConfig(jsonData string) (*Config, error) {
    var config Config
    if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
        return nil, fmt.Errorf("解析配置失败: %w", err)
    }
    return &config, nil
}

// 8. 日志系统
// 日志系统是应用程序中记录运行时信息的重要组件
// 下面的示例演示了一个简单的日志接口和实现

// Logger 定义了日志系统的接口
// 支持不同级别的日志记录

type Logger interface {
    // Info 记录信息级别的日志
    Info(msg string, fields ...interface{})
    
    // Error 记录错误级别的日志
    Error(msg string, fields ...interface{})
    
    // Debug 记录调试级别的日志
    Debug(msg string, fields ...interface{})
}

// SimpleLogger 是Logger接口的一个简单实现
// prefix 是日志前缀
type SimpleLogger struct {
    prefix string
}

// NewSimpleLogger 创建并返回一个新的简单日志记录器
// prefix 是日志前缀
func NewSimpleLogger(prefix string) *SimpleLogger {
    return &SimpleLogger{prefix: prefix}
}

// Info 实现Logger接口的Info方法
// 记录信息级别的日志
func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
    fmt.Printf("[INFO] %s: %s %v\n", l.prefix, msg, fields)
}

// Error 实现Logger接口的Error方法
// 记录错误级别的日志
func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
    fmt.Printf("[ERROR] %s: %s %v\n", l.prefix, msg, fields)
}

// Debug 实现Logger接口的Debug方法
// 记录调试级别的日志
func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
    fmt.Printf("[DEBUG] %s: %s %v\n", l.prefix, msg, fields)
}

// 9. 限流器
type RateLimiter struct {
    requests map[string][]time.Time
    limit    int
    window   time.Duration
    mu       sync.RWMutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    
    // 清理过期的请求记录
    if requests, exists := rl.requests[key]; exists {
        validRequests := make([]time.Time, 0)
        for _, reqTime := range requests {
            if now.Sub(reqTime) <= rl.window {
                validRequests = append(validRequests, reqTime)
            }
        }
        rl.requests[key] = validRequests
    }
    
    // 检查是否超过限制
    if len(rl.requests[key]) >= rl.limit {
        return false
    }
    
    // 记录当前请求
    rl.requests[key] = append(rl.requests[key], now)
    return true
}

// 10. 健康检查
type HealthChecker struct {
    services map[string]func() error
}

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        services: make(map[string]func() error),
    }
}

func (hc *HealthChecker) Register(name string, checkFunc func() error) {
    hc.services[name] = checkFunc
}

func (hc *HealthChecker) CheckAll() map[string]string {
    results := make(map[string]string)
    
    for name, checkFunc := range hc.services {
        if err := checkFunc(); err != nil {
            results[name] = "unhealthy: " + err.Error()
        } else {
            results[name] = "healthy"
        }
    }
    
    return results
}

// 主函数演示
func main() {
    fmt.Println("=== Go语言高级主题 ===")
    
    // 1. 反射示例
    fmt.Println("\n1. 反射示例")
    person := Person{Name: "张三", Age: 25, Email: "zhangsan@example.com"}
    analyzeStruct(person)
    
    // 2. 泛型示例
    fmt.Println("\n2. 泛型示例")
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 使用泛型过滤偶数
    evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
    fmt.Printf("偶数: %v\n", evens)
    
    // 使用泛型映射平方
    squares := Map(numbers, func(n int) int { return n * n })
    fmt.Printf("平方: %v\n", squares)
    
    // 使用泛型求和
    sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
    fmt.Printf("总和: %d\n", sum)
    
    // 泛型栈
    stack := NewStack[int]()
    stack.Push(1)
    stack.Push(2)
    stack.Push(3)
    
    if item, ok := stack.Pop(); ok {
        fmt.Printf("栈顶元素: %d\n", item)
    }
    
    // 泛型缓存
    cache := NewCache[string, int]()
    cache.Set("key1", 100)
    cache.Set("key2", 200)
    
    if value, exists := cache.Get("key1"); exists {
        fmt.Printf("缓存值: %d\n", value)
    }
    
    // 3. 上下文示例
    contextExample()
    
    // 测试长时间操作
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    if err := longRunningOperation(ctx); err != nil {
        fmt.Printf("操作结果: %v\n", err)
    }
    
    // 4. TCP网络示例
    tcpServerExample()
    
    // 5. 加密示例
    encryptionExample()
    
    // 6. 微服务架构示例
    fmt.Println("\n6. 微服务架构示例")
    
    registry := NewInMemoryServiceRegistry()
    
    // 注册服务
    registry.Register("user-service", "localhost:8080")
    registry.Register("user-service", "localhost:8081")
    registry.Register("order-service", "localhost:8082")
    
    // 服务发现
    if addresses, err := registry.Discover("user-service"); err == nil {
        fmt.Printf("发现用户服务: %v\n", addresses)
    }
    
    // 负载均衡
    lb := NewLoadBalancer(registry)
    if address, err := lb.GetServiceAddress("user-service"); err == nil {
        fmt.Printf("负载均衡选择: %s\n", address)
    }
    
    // 配置管理
    configJSON := `{
        "server": {"host": "localhost", "port": 8080},
        "database": {"host": "localhost", "port": 3306, "user": "root", "password": "password", "name": "app"},
        "cache": {"host": "localhost", "port": 6379}
    }`
    
    if config, err := loadConfig(configJSON); err == nil {
        fmt.Printf("配置加载成功: %+v\n", config)
    }
    
    // 日志系统
    logger := NewSimpleLogger("app")
    logger.Info("应用启动成功")
    logger.Debug("调试信息", "key", "value")
    
    // 限流器
    rateLimiter := NewRateLimiter(5, 1*time.Minute)
    
    for i := 0; i < 10; i++ {
        if rateLimiter.Allow("user123") {
            fmt.Printf("请求 %d: 允许\n", i+1)
        } else {
            fmt.Printf("请求 %d: 拒绝\n", i+1)
        }
    }
    
    // 健康检查
    healthChecker := NewHealthChecker()
    
    healthChecker.Register("database", func() error {
        // 模拟数据库健康检查
        return nil
    })
    
    healthChecker.Register("cache", func() error {
        // 模拟缓存健康检查
        return nil
    })
    
    healthChecker.Register("external-api", func() error {
        // 模拟外部API健康检查
        return fmt.Errorf("外部API暂时不可用")
    })
    
    healthResults := healthChecker.CheckAll()
    fmt.Printf("健康检查结果: %+v\n", healthResults)
    
    fmt.Println("\n=== 高级主题演示完成 ===")
    fmt.Println("\n练习：")
    fmt.Println("1. 使用反射实现通用的JSON序列化器")
    fmt.Println("2. 创建泛型二叉树数据结构")
    fmt.Println("3. 实现基于上下文的可取消任务队列")
    fmt.Println("4. 构建TCP聊天服务器")
    fmt.Println("5. 实现JWT令牌生成和验证")
    fmt.Println("6. 创建服务网格代理")
    fmt.Println("7. 实现分布式配置中心")
    fmt.Println("8. 构建API网关")
}