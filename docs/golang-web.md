# Golang Web开发

## HTTP基础

### 简单的HTTP服务器

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", helloHandler)
    fmt.Println("Server starting on port 8080...")
    http.ListenAndServe(":8080", nil)
}
```

### 路由处理

```go
package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to Home Page</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>About Us</h1>")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Path[len("/user/"):]
    fmt.Fprintf(w, "<h1>Hello, %s!</h1>", name)
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/user/", userHandler)
    
    fmt.Println("Server starting on port 8080...")
    http.ListenAndServe(":8080", nil)
}
```

## 使用Gin框架

### 安装Gin

```bash
go get -u github.com/gin-gonic/gin
```

### 基本使用

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })
    
    r.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.JSON(200, gin.H{
            "message": "Hello " + name,
        })
    })
    
    r.GET("/welcome", func(c *gin.Context) {
        firstname := c.DefaultQuery("firstname", "Guest")
        lastname := c.Query("lastname")
        
        c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
    })
    
    r.Run(":8080")
}
```

### POST请求处理

```go
type LoginForm struct {
    User     string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func main() {
    r := gin.Default()
    
    r.POST("/login", func(c *gin.Context) {
        var form LoginForm
        if err := c.ShouldBind(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if form.User == "admin" && form.Password == "password" {
            c.JSON(http.StatusOK, gin.H{"status": "logged in"})
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
        }
    })
    
    r.Run(":8080")
}
```

### JSON数据处理

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    r := gin.Default()
    
    r.POST("/user", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": "User created",
            "user":    user,
        })
    })
    
    r.Run(":8080")
}
```

### 路由组

```go
func main() {
    r := gin.Default()
    
    // API路由组
    api := r.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.GET("/users/:id", getUser)
        api.POST("/users", createUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }
    
    r.Run(":8080")
}

func getUsers(c *gin.Context) {
    c.JSON(200, gin.H{"users": "list of users"})
}

func getUser(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user": "user " + id})
}

func createUser(c *gin.Context) {
    c.JSON(201, gin.H{"message": "user created"})
}

func updateUser(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"message": "user " + id + " updated"})
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"message": "user " + id + " deleted"})
}
```

## 中间件

### 自定义中间件

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        
        // 设置变量
        c.Set("example", "12345")
        
        // 请求前
        
        c.Next()
        
        // 请求后
        latency := time.Since(t)
        log.Print(latency)
        
        // 访问状态
        status := c.Writer.Status()
        log.Println(status)
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())
    
    r.GET("/test", func(c *gin.Context) {
        example := c.MustGet("example").(string)
        log.Println(example)
    })
    
    r.Run(":8080")
}
```

### 认证中间件

```go
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
            c.Abort()
            return
        }
        
        // 验证token逻辑
        if token != "valid-token" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func main() {
    r := gin.Default()
    
    // 公开路由
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "public"})
    })
    
    // 需要认证的路由
    authorized := r.Group("/admin")
    authorized.Use(AuthRequired())
    {
        authorized.GET("/dashboard", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "admin dashboard"})
        })
    }
    
    r.Run(":8080")
}
```

## 静态文件服务

```go
func main() {
    r := gin.Default()
    
    // 单个文件
    r.StaticFile("/favicon.ico", "./resources/favicon.ico")
    
    // 整个目录
    r.Static("/static", "./static")
    
    // 使用前缀
    r.StaticFS("/more_static", http.Dir("my_file_system"))
    
    r.Run(":8080")
}
```

## 模板渲染

### HTML模板

```go
func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*")
    
    router.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "title": "Main website",
        })
    })
    
    router.Run(":8080")
}
```

### 模板文件示例 (templates/index.tmpl)

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
</body>
</html>
```

## 错误处理

```go
func main() {
    r := gin.Default()
    
    r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if err, ok := recovered.(string); ok {
            c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
        }
        c.AbortWithStatus(http.StatusInternalServerError)
    }))
    
    r.GET("/panic", func(c *gin.Context) {
        panic("foo")
    })
    
    r.Run(":8080")
}
```