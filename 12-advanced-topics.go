package main

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
import "reflect"

// 2. 泛型编程（Go 1.18+）

// 3. 上下文管理

// 4. 网络编程

// 5. 加密和安全

// 6. 微服务架构

// 1. 反射示例
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
    private string // 私有字段
}

// 使用反射分析结构体
func analyzeStruct(obj interface{}) {
    val := reflect.ValueOf(obj)
    typ := reflect.TypeOf(obj)
    
    fmt.Printf("类型: %s\n", typ)
    fmt.Printf("种类: %s\n", typ.Kind())
    
    if typ.Kind() == reflect.Struct {
        fmt.Printf("字段数量: %d\n", typ.NumField())
        
        for i := 0; i < typ.NumField(); i++ {
            field := typ.Field(i)
            value := val.Field(i)
            
            fmt.Printf("  字段 %d: %s (%s) = %v, 标签: %s\n", 
                i, field.Name, field.Type, value, field.Tag)
        }
    }
}

// 动态创建结构体实例
func createInstanceFromType(typ reflect.Type) interface{} {
    if typ.Kind() != reflect.Struct {
        return nil
    }
    
    instance := reflect.New(typ).Elem()
    
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        fieldVal := instance.Field(i)
        
        switch field.Type.Kind() {
        case reflect.String:
            fieldVal.SetString("default")
        case reflect.Int:
            fieldVal.SetInt(0)
        case reflect.Bool:
            fieldVal.SetBool(false)
        }
    }
    
    return instance.Interface()
}

// 2. 泛型编程示例

// 泛型切片操作
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range slice {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}

func Map[T any, U any](slice []T, transform func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = transform(item)
    }
    return result
}

func Reduce[T any, U any](slice []T, initial U, reducer func(U, T) U) U {
    result := initial
    for _, item := range slice {
        result = reducer(result, item)
    }
    return result
}

// 泛型栈
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
    return len(s.items)
}

// 泛型队列
type Queue[T any] struct {
    items []T
}

func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

func (q *Queue[T]) Size() int {
    return len(q.items)
}

// 泛型缓存
import "sync"

type Cache[K comparable, V any] struct {
    items map[K]V
    mu    sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]V),
    }
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.items[key]
    return value, exists
}

func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

// 3. 上下文管理示例
func contextExample() {
    fmt.Println("\n=== 上下文管理示例 ===")
    
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // 模拟长时间操作
    go func() {
        select {
        case <-time.After(3 * time.Second):
            fmt.Println("操作完成")
        case <-ctx.Done():
            fmt.Println("操作超时:", ctx.Err())
        }
    }()
    
    // 等待操作完成或超时
    time.Sleep(3 * time.Second)
}

// 使用上下文取消操作
func longRunningOperation(ctx context.Context) error {
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(500 * time.Millisecond):
            fmt.Printf("操作进度: %d%%\n", (i+1)*10)
        }
    }
    return nil
}

// 4. 网络编程示例
func tcpServerExample() {
    fmt.Println("\n=== TCP服务器示例 ===")
    
    // 创建TCP监听器
    listener, err := net.Listen("tcp", ":8081")
    if err != nil {
        fmt.Printf("创建监听器失败: %v\n", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("TCP服务器启动在 :8081")
    
    // 在后台接受连接
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Printf("接受连接失败: %v\n", err)
                continue
            }
            
            go handleTCPConnection(conn)
        }
    }()
    
    // 模拟客户端连接
    time.Sleep(100 * time.Millisecond)
    
    conn, err := net.Dial("tcp", "localhost:8081")
    if err != nil {
        fmt.Printf("连接服务器失败: %v\n", err)
        return
    }
    
    fmt.Fprintln(conn, "Hello from client")
    
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("读取响应失败: %v\n", err)
        return
    }
    
    fmt.Printf("收到服务器响应: %s\n", string(buffer[:n]))
    conn.Close()
    
    time.Sleep(100 * time.Millisecond)
}

func handleTCPConnection(conn net.Conn) {
    defer conn.Close()
    
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("读取数据失败: %v\n", err)
        return
    }
    
    message := string(buffer[:n])
    fmt.Printf("收到消息: %s\n", message)
    
    response := fmt.Sprintf("服务器收到: %s", message)
    conn.Write([]byte(response))
}

// 5. 加密和安全示例
import "crypto/aes"
import "crypto/cipher"

func encryptionExample() {
    fmt.Println("\n=== 加密示例 ===")
    
    // 生成随机密钥
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
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        fmt.Printf("创建GCM失败: %v\n", err)
        return
    }
    
    // 生成随机nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        fmt.Printf("生成nonce失败: %v\n", err)
        return
    }
    
    // 加密
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    fmt.Printf("加密后的数据: %x\n", ciphertext)
    
    // 解密
    decrypted, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
    if err != nil {
        fmt.Printf("解密失败: %v\n", err)
        return
    }
    
    fmt.Printf("解密后的数据: %s\n", string(decrypted))
}

// 6. 微服务架构示例
// 服务发现接口
type ServiceRegistry interface {
    Register(serviceName string, address string) error
    Deregister(serviceName string) error
    Discover(serviceName string) ([]string, error)
    HealthCheck(serviceName string) bool
}

// 内存服务注册表
type InMemoryServiceRegistry struct {
    services map[string][]string
    mu       sync.RWMutex
}

func NewInMemoryServiceRegistry() *InMemoryServiceRegistry {
    return &InMemoryServiceRegistry{
        services: make(map[string][]string),
    }
}

func (r *InMemoryServiceRegistry) Register(serviceName string, address string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.services[serviceName] = append(r.services[serviceName], address)
    fmt.Printf("服务 %s 注册在 %s\n", serviceName, address)
    return nil
}

func (r *InMemoryServiceRegistry) Deregister(serviceName string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    delete(r.services, serviceName)
    fmt.Printf("服务 %s 已注销\n", serviceName)
    return nil
}

func (r *InMemoryServiceRegistry) Discover(serviceName string) ([]string, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    addresses, exists := r.services[serviceName]
    if !exists {
        return nil, fmt.Errorf("服务 %s 未找到", serviceName)
    }
    
    return addresses, nil
}

func (r *InMemoryServiceRegistry) HealthCheck(serviceName string) bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    _, exists := r.services[serviceName]
    return exists
}

// 负载均衡器
type LoadBalancer struct {
    registry ServiceRegistry
    strategy string // "round-robin", "random", "least-connections"
}

func NewLoadBalancer(registry ServiceRegistry) *LoadBalancer {
    return &LoadBalancer{
        registry: registry,
        strategy: "round-robin",
    }
}

func (lb *LoadBalancer) GetServiceAddress(serviceName string) (string, error) {
    addresses, err := lb.registry.Discover(serviceName)
    if err != nil {
        return "", err
    }
    
    if len(addresses) == 0 {
        return "", fmt.Errorf("没有可用的服务实例")
    }
    
    // 简单的轮询策略
    return addresses[0], nil
}

// 7. 配置管理
type Config struct {
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    Database struct {
        Host     string `json:"host"`
        Port     int    `json:"port"`
        User     string `json:"user"`
        Password string `json:"password"`
        Name     string `json:"name"`
    } `json:"database"`
    Cache struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"cache"`
}

func loadConfig(jsonData string) (*Config, error) {
    var config Config
    if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
        return nil, fmt.Errorf("解析配置失败: %w", err)
    }
    return &config, nil
}

// 8. 日志系统
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Debug(msg string, fields ...interface{})
}

type SimpleLogger struct {
    prefix string
}

func NewSimpleLogger(prefix string) *SimpleLogger {
    return &SimpleLogger{prefix: prefix}
}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
    fmt.Printf("[INFO] %s: %s %v\n", l.prefix, msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
    fmt.Printf("[ERROR] %s: %s %v\n", l.prefix, msg, fields)
}

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