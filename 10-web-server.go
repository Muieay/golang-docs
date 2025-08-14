package main

import (
	"encoding/json"   // 提供JSON序列化和反序列化功能，用于处理API的JSON请求和响应
	"fmt"             // 提供格式化输入输出功能
	"log"             // 提供日志记录功能
	"net/http"        // 提供HTTP服务器和客户端功能，是Go Web开发的核心包
	"strconv"         // 提供字符串和基本数据类型之间的转换功能
	"strings"         // 提供字符串操作功能
	"time"            // 提供时间相关的功能，用于处理时间戳和超时等
)

// Web服务器开发
// Go语言的HTTP服务器和RESTful API开发
// 本示例展示了如何使用Go标准库构建一个完整的Web服务器，包括：
// - 数据模型定义
// - 内存存储（模拟数据库）
// - 中间件机制
// - RESTful API实现
// - 静态文件服务
// - 路由管理

// 1. 数据模型
// User：用户数据模型，用于表示系统中的用户信息
// 结构体字段后的`json:"字段名"`是结构体标签（struct tag）
// 作用：在JSON序列化/反序列化时指定字段名称，实现Go字段名与JSON字段名的映射
type User struct {
	ID      int    `json:"id"`       // 用户唯一标识，自增整数
	Name    string `json:"name"`     // 用户名
	Email   string `json:"email"`    // 用户邮箱
	Age     int    `json:"age"`      // 用户年龄
	Created string `json:"created"`  // 账号创建时间，使用RFC3339格式字符串
}

// Post：帖子数据模型，用于表示用户发布的内容
type Post struct {
	ID      int    `json:"id"`       // 帖子唯一标识，自增整数
	Title   string `json:"title"`    // 帖子标题
	Content string `json:"content"`  // 帖子内容
	Author  string `json:"author"`   // 作者名称
	Date    string `json:"date"`     // 发布时间，使用RFC3339格式字符串
}

// 2. 内存存储（模拟数据库）
// 本示例使用map作为内存存储来模拟数据库，适合教学和演示
// 在实际项目中，通常会替换为真实的数据库（如MySQL、PostgreSQL等）
var (
	users      = make(map[int]User)  // 存储用户数据，key为用户ID
	posts      = make(map[int]Post)  // 存储帖子数据，key为帖子ID
	nextUserID = 1                   // 下一个可用的用户ID，用于生成新用户的唯一标识
	nextPostID = 1                   // 下一个可用的帖子ID，用于生成新帖子的唯一标识
)

// 3. 中间件
// 中间件（Middleware）是Go Web开发中的重要概念，用于在请求到达处理器之前或之后添加额外逻辑
// 典型用途：日志记录、身份验证、跨域处理、错误处理等
// 实现原理：接收一个http.HandlerFunc作为参数，返回一个新的http.HandlerFunc

// loggingMiddleware：日志中间件，记录请求处理的相关信息
// 功能：记录请求方法、路径、处理开始时间和耗时
// 参数：next http.HandlerFunc - 下一个要执行的处理器函数
// 返回值：包装后的处理器函数
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// 返回一个匿名函数作为新的处理器
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()                          // 记录请求处理开始时间
		log.Printf("开始处理: %s %s", r.Method, r.URL.Path)  // 记录请求方法和路径
		
		next(w, r)  // 调用下一个处理器，继续处理请求（核心：中间件链的传递）
		
		// 计算并记录请求处理耗时
		log.Printf("完成处理: %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	}
}

// corsMiddleware：跨域资源共享中间件，处理浏览器的跨域请求限制
// 功能：设置CORS相关的HTTP头，允许跨域请求
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置允许的源（*表示允许所有源）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的HTTP方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的请求头
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// 处理预检请求（OPTIONS方法）
		// 浏览器在发送跨域请求前，可能会先发送OPTIONS请求检查服务器是否允许跨域
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)  // 直接返回200表示允许
			return
		}
		
		next(w, r)  // 继续处理请求
	}
}

// 4. 用户管理处理器
// 处理器（Handler）是处理HTTP请求的函数，遵循http.HandlerFunc类型定义
// 签名为：func(w http.ResponseWriter, r *http.Request)
// w：用于构建响应，r：包含请求信息

// handleUsers：用户管理的主处理器，根据HTTP方法分发到不同的处理函数
// 实现RESTful API的核心思想：同一资源路径根据不同HTTP方法执行不同操作
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:    // GET方法：获取资源
		getUsers(w, r)
	case http.MethodPost:   // POST方法：创建资源
		createUser(w, r)
	case http.MethodPut:    // PUT方法：更新资源
		updateUser(w, r)
	case http.MethodDelete: // DELETE方法：删除资源
		deleteUser(w, r)
	default:
		// 处理不支持的HTTP方法，返回405 Method Not Allowed
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// getUsers：处理获取所有用户的请求（GET /users）
// 功能：从内存存储中读取所有用户，以JSON格式返回
func getUsers(w http.ResponseWriter, r *http.Request) {
	// 设置响应头Content-Type为application/json，告诉客户端返回的是JSON数据
	w.Header().Set("Content-Type", "application/json")
	
	// 将map中的用户转换为切片（数组），因为JSON序列化map的顺序不确定
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	
	// 使用json.NewEncoder将用户列表编码为JSON并写入响应
	// Encode方法会自动处理错误，失败时会返回500 Internal Server Error
	json.NewEncoder(w).Encode(userList)
}

// createUser：处理创建新用户的请求（POST /users）
// 功能：解析请求体中的JSON数据，创建新用户并保存到内存存储
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User  // 声明一个User类型变量，用于接收解析后的请求数据
	
	// 解析请求体中的JSON数据到user变量
	// json.NewDecoder(r.Body).Decode(&user)：从请求体读取并反序列化JSON
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// 解析失败时，返回400 Bad Request错误
		http.Error(w, "无效的JSON数据", http.StatusBadRequest)
		return
	}
	
	// 为新用户分配ID和创建时间
	user.ID = nextUserID          // 使用全局变量nextUserID作为新用户ID
	user.Created = time.Now().Format(time.RFC3339)  // 格式化当前时间为RFC3339标准格式
	nextUserID++                  // 更新nextUserID，确保下次创建用户时ID唯一
	
	// 将新用户保存到内存存储
	users[user.ID] = user
	
	// 设置响应头和状态码
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)  // 201 Created表示资源创建成功
	// 返回创建的用户信息
	json.NewEncoder(w).Encode(user)
}

// updateUser：处理更新用户的请求（PUT /users/{id}）
// 功能：根据URL路径中的ID查找用户，更新其信息并保存
func updateUser(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中提取用户ID（例如从"/users/123"中提取"123"）
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	// 将字符串ID转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 转换失败（ID不是数字），返回400 Bad Request
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}
	
	// 解析请求体中的JSON数据到临时user变量
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "无效的JSON数据", http.StatusBadRequest)
		return
	}
	
	// 检查用户是否存在
	if _, exists := users[id]; !exists {
		// 用户不存在，返回404 Not Found
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}
	
	// 确保更新的用户ID与URL中的ID一致（防止ID被篡改）
	user.ID = id
	// 更新内存存储中的用户信息
	users[id] = user
	
	// 返回更新后的用户信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser：处理删除用户的请求（DELETE /users/{id}）
// 功能：根据URL路径中的ID查找并删除用户
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中提取用户ID
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}
	
	// 检查用户是否存在
	if _, exists := users[id]; !exists {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}
	
	// 从内存存储中删除用户
	delete(users, id)
	// 返回204 No Content，表示删除成功且无响应体
	w.WriteHeader(http.StatusNoContent)
}

// 5. 帖子管理处理器
// handlePosts：帖子管理的主处理器，根据HTTP方法分发到不同的处理函数
func handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:    // GET方法：获取帖子列表
		getPosts(w, r)
	case http.MethodPost:   // POST方法：创建新帖子
		createPost(w, r)
	default:
		// 暂不支持PUT和DELETE方法，返回405
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// getPosts：处理获取所有帖子的请求（GET /posts）
// 功能：从内存存储中读取所有帖子，以JSON格式返回
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 将map中的帖子转换为切片
	postList := make([]Post, 0, len(posts))
	for _, post := range posts {
		postList = append(postList, post)
	}
	
	json.NewEncoder(w).Encode(postList)
}

// createPost：处理创建新帖子的请求（POST /posts）
// 功能：解析请求体中的JSON数据，创建新帖子并保存到内存存储
func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	
	// 解析请求体中的JSON数据
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "无效的JSON数据", http.StatusBadRequest)
		return
	}
	
	// 为新帖子分配ID和发布时间
	post.ID = nextPostID
	post.Date = time.Now().Format(time.RFC3339)
	nextPostID++
	
	// 保存新帖子
	posts[post.ID] = post
	
	// 返回创建的帖子信息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// 6. 健康检查处理器
// handleHealth：处理健康检查请求（GET /health）
// 功能：返回服务器的健康状态，常用于监控系统检查服务是否正常运行
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// 构建健康状态数据
	health := map[string]interface{}{
		"status":    "healthy",                // 健康状态
		"timestamp": time.Now().Format(time.RFC3339),  // 当前时间戳
		"version":   "1.0.0",                  // 服务版本
	}
	
	// 返回JSON格式的健康状态
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// 7. 统计信息处理器
// handleStats：处理统计信息请求（GET /stats）
// 功能：返回服务器的统计数据，如用户数量、帖子数量、运行时间等
func handleStats(w http.ResponseWriter, r *http.Request) {
	// 构建统计数据
	stats := map[string]interface{}{
		"users":  len(users),                 // 用户数量
		"posts":  len(posts),                 // 帖子数量
		"uptime": time.Since(startTime).String(),  // 服务器运行时间
	}
	
	// 返回JSON格式的统计信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// 8. 路由处理器
// handleAPI：API路由分发处理器，根据URL路径分发到不同的功能处理器
// 注意：本示例中此函数未被直接使用，主路由设置在main函数中
func handleAPI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path  // 获取请求的URL路径
	
	// 根据路径匹配不同的处理器
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
		// 路径未匹配，返回404 Not Found
		http.NotFound(w, r)
	}
}

// 9. 初始化数据
// initData：服务器启动时初始化示例数据
// 功能：添加一些默认用户和帖子，方便测试API功能
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
	
	// 更新下一个可用ID，确保新创建的资源ID不会冲突
	nextUserID = 3
	nextPostID = 3
}

// 全局变量
var startTime = time.Now()  // 记录服务器启动时间，用于计算运行时间

// 10. 静态文件服务
// handleStatic：处理静态文件请求和主页请求
// 功能：提供静态资源访问和默认主页
func handleStatic(w http.ResponseWriter, r *http.Request) {
	// 根路径请求（/）返回主页
	if r.URL.Path == "/" {
		serveHomePage(w, r)
		return
	}
	
	// 处理静态文件请求（/static/前缀的路径）
	if strings.HasPrefix(r.URL.Path, "/static/") {
		// http.ServeFile：从本地文件系统读取文件并返回给客户端
		// r.URL.Path[1:]：去掉路径开头的斜杠，获取正确的文件路径
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}
	
	// 未匹配的路径返回404
	http.NotFound(w, r)
}

// 11. 主页服务
// serveHomePage：生成并返回HTML格式的主页
// 功能：展示服务器的API端点和使用示例，方便用户了解如何使用API
func serveHomePage(w http.ResponseWriter, r *http.Request) {
	// 定义HTML内容，包含API端点说明和使用示例
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
	
	// 设置响应头Content-Type为text/html，告诉客户端返回的是HTML内容
	w.Header().Set("Content-Type", "text/html")
	// 将HTML内容写入响应
	fmt.Fprint(w, html)
}

// 12. 中间件链
// withMiddleware：组合多个中间件，形成中间件链
// 功能：将CORS中间件和日志中间件组合，应用到处理器上
// 注意：中间件的顺序很重要，这里先执行corsMiddleware，再执行loggingMiddleware
func withMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return corsMiddleware(loggingMiddleware(next))
}

// 主函数：程序入口点
func main() {
	fmt.Println("=== Go语言Web服务器开发 ===")
	
	// 初始化示例数据
	initData()
	
	// 设置路由规则
	// http.HandleFunc：将URL模式与处理器函数关联
	// 第一个参数是URL模式，第二个参数是处理器函数（可以是经过中间件包装的）
	http.HandleFunc("/", withMiddleware(handleStatic))
	http.HandleFunc("/health", withMiddleware(handleHealth))
	http.HandleFunc("/stats", withMiddleware(handleStats))
	http.HandleFunc("/users", withMiddleware(handleUsers))
	http.HandleFunc("/users/", withMiddleware(handleUsers))  // 处理带ID的用户路径
	http.HandleFunc("/posts", withMiddleware(handlePosts))
	http.HandleFunc("/posts/", withMiddleware(handlePosts))  // 处理带ID的帖子路径
	
	// 启动HTTP服务器
	fmt.Println("服务器启动在 http://localhost:8080")
	fmt.Println("按 Ctrl+C 停止服务器")
	
	// http.ListenAndServe：启动服务器，监听指定地址和端口
	// 第一个参数是地址（格式为"host:port"），第二个参数是处理器（nil表示使用默认的DefaultServeMux）
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// 服务器启动失败时，记录错误并退出
		log.Fatal("服务器启动失败:", err)
	}
}

// 13. 客户端测试函数（在另一个文件中）
// testClient：用于测试服务器API的客户端函数
// 功能：发送测试请求到服务器，验证API是否正常工作
func testClient() {
	// 这个函数可以在另一个文件中用于测试服务器
	fmt.Println("\n=== 客户端测试 ===")
	
	// 测试健康检查API
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()  // defer确保响应体在函数退出时关闭，避免资源泄露
	
	fmt.Printf("健康检查响应状态: %s\n", resp.Status)
	
	// 测试获取用户列表API
	resp, err = http.Get("http://localhost:8080/users")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()  // 关闭响应体
	
	fmt.Printf("获取用户响应状态: %s\n", resp.Status)
}
