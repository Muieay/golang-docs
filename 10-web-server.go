package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// Web服务器开发
// Go语言的HTTP服务器和RESTful API开发

// 1. 数据模型
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    Created  string `json:"created"`
}

type Post struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Author  string `json:"author"`
    Date    string `json:"date"`
}

// 2. 内存存储（模拟数据库）
var (
    users = make(map[int]User)
    posts = make(map[int]Post)
    nextUserID = 1
    nextPostID = 1
)

// 3. 中间件
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("开始处理: %s %s", r.Method, r.URL.Path)
        next(w, r)
        log.Printf("完成处理: %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
    }
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

// 4. 用户管理处理器
func handleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUsers(w, r)
    case http.MethodPost:
        createUser(w, r)
    case http.MethodPut:
        updateUser(w, r)
    case http.MethodDelete:
        deleteUser(w, r)
    default:
        http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
    }
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    userList := make([]User, 0, len(users))
    for _, user := range users {
        userList = append(userList, user)
    }
    
    json.NewEncoder(w).Encode(userList)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "无效的JSON数据", http.StatusBadRequest)
        return
    }
    
    user.ID = nextUserID
    user.Created = time.Now().Format(time.RFC3339)
    nextUserID++
    
    users[user.ID] = user
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "无效的用户ID", http.StatusBadRequest)
        return
    }
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "无效的JSON数据", http.StatusBadRequest)
        return
    }
    
    if _, exists := users[id]; !exists {
        http.Error(w, "用户不存在", http.StatusNotFound)
        return
    }
    
    user.ID = id
    users[id] = user
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "无效的用户ID", http.StatusBadRequest)
        return
    }
    
    if _, exists := users[id]; !exists {
        http.Error(w, "用户不存在", http.StatusNotFound)
        return
    }
    
    delete(users, id)
    w.WriteHeader(http.StatusNoContent)
}

// 5. 帖子管理处理器
func handlePosts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getPosts(w, r)
    case http.MethodPost:
        createPost(w, r)
    default:
        http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
    }
}

func getPosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    postList := make([]Post, 0, len(posts))
    for _, post := range posts {
        postList = append(postList, post)
    }
    
    json.NewEncoder(w).Encode(postList)
}

func createPost(w http.ResponseWriter, r *http.Request) {
    var post Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "无效的JSON数据", http.StatusBadRequest)
        return
    }
    
    post.ID = nextPostID
    post.Date = time.Now().Format(time.RFC3339)
    nextPostID++
    
    posts[post.ID] = post
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

// 6. 健康检查处理器
func handleHealth(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
        "version":   "1.0.0",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

// 7. 统计信息处理器
func handleStats(w http.ResponseWriter, r *http.Request) {
    stats := map[string]interface{}{
        "users": len(users),
        "posts": len(posts),
        "uptime": time.Since(startTime).String(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// 8. 路由处理器
func handleAPI(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    
    switch {
    case path == "/health":
        handleHealth(w, r)
    case path == "/stats":
        handleStats(w, r)
    case path == "/users" || strings.HasPrefix(path, "/users/"):
        handleUsers(w, r)
    case path == "/posts" || strings.HasPrefix(path, "/posts/"):
        handlePosts(w, r)
    default:
        http.NotFound(w, r)
    }
}

// 9. 初始化数据
func initData() {
    // 添加示例用户
    users[1] = User{
        ID:      1,
        Name:    "张三",
        Email:   "zhangsan@example.com",
        Age:     25,
        Created: time.Now().Format(time.RFC3339),
    }
    
    users[2] = User{
        ID:      2,
        Name:    "李四",
        Email:   "lisi@example.com",
        Age:     30,
        Created: time.Now().Format(time.RFC3339),
    }
    
    // 添加示例帖子
    posts[1] = Post{
        ID:      1,
        Title:   "Go语言入门",
        Content: "Go语言是一门现代化的编程语言，具有并发支持...",
        Author:  "张三",
        Date:    time.Now().Format(time.RFC3339),
    }
    
    posts[2] = Post{
        ID:      2,
        Title:   "Web开发基础",
        Content: "使用Go语言构建Web应用非常简单...",
        Author:  "李四",
        Date:    time.Now().Format(time.RFC3339),
    }
    
    nextUserID = 3
    nextPostID = 3
}

// 全局变量
var startTime = time.Now()

// 10. 静态文件服务
func handleStatic(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        serveHomePage(w, r)
        return
    }
    
    // 简单的静态文件服务
    if strings.HasPrefix(r.URL.Path, "/static/") {
        http.ServeFile(w, r, r.URL.Path[1:])
        return
    }
    
    http.NotFound(w, r)
}

// 11. 主页服务
func serveHomePage(w http.ResponseWriter, r *http.Request) {
    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go Web服务器示例</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { margin: 20px 0; padding: 10px; background: #f5f5f5; border-radius: 5px; }
        .method { color: #007cba; font-weight: bold; }
        .path { color: #d73a49; font-family: monospace; }
    </style>
</head>
<body>
    <h1>Go语言Web服务器示例</h1>
    <p>这是一个使用Go语言构建的简单Web服务器示例</p>
    
    <h2>可用API端点</h2>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/health</span> - 健康检查
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/stats</span> - 统计信息
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/users</span> - 获取所有用户
    </div>
    
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/users</span> - 创建新用户
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/posts</span> - 获取所有帖子
    </div>
    
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/posts</span> - 创建新帖子
    </div>
    
    <h2>使用示例</h2>
    <pre>
# 健康检查
curl http://localhost:8080/health

# 获取用户列表
curl http://localhost:8080/users

# 创建新用户
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"王五","email":"wangwu@example.com","age":28}'
    </pre>
</body>
</html>
`
    
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

// 12. 中间件链
func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return corsMiddleware(loggingMiddleware(next))
}

// 主函数
func main() {
    fmt.Println("=== Go语言Web服务器开发 ===")
    
    // 初始化数据
    initData()
    
    // 设置路由
    http.HandleFunc("/", withMiddleware(handleStatic))
    http.HandleFunc("/health", withMiddleware(handleHealth))
    http.HandleFunc("/stats", withMiddleware(handleStats))
    http.HandleFunc("/users", withMiddleware(handleUsers))
    http.HandleFunc("/users/", withMiddleware(handleUsers))
    http.HandleFunc("/posts", withMiddleware(handlePosts))
    http.HandleFunc("/posts/", withMiddleware(handlePosts))
    
    // 启动服务器
    fmt.Println("服务器启动在 http://localhost:8080")
    fmt.Println("按 Ctrl+C 停止服务器")
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("服务器启动失败:", err)
    }
}

// 13. 客户端测试函数（在另一个文件中）
func testClient() {
    // 这个函数可以在另一个文件中用于测试服务器
    fmt.Println("\n=== 客户端测试 ===")
    
    // 测试健康检查
    resp, err := http.Get("http://localhost:8080/health")
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Printf("健康检查响应状态: %s\n", resp.Status)
    
    // 测试获取用户
    resp, err = http.Get("http://localhost:8080/users")
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Printf("获取用户响应状态: %s\n", resp.Status)
}