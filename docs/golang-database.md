# Golang数据库操作

## 数据库/sql标准库

Go 语言通过 `database/sql` 包提供了通用的数据库接口，允许开发者以统一的方式操作不同的数据库系统。

### 安装驱动

```bash
go get -u github.com/go-sql-driver/mysql
```

> **注意**：导入驱动时使用下划线 `_` 前缀，表示只初始化驱动而不直接使用其导出的函数。

### 连接数据库

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // 连接MySQL - 注意：dsn中的密码如果包含特殊字符需要进行URL编码
    dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("无法打开数据库连接: %v", err)
    }
    defer db.Close()
    
    // 测试连接 - 确保能真正连接到数据库
    err = db.Ping()
    if err != nil {
        log.Fatalf("无法连接到数据库: %v", err)
    }
    
    fmt.Println("Successfully connected!")
}
```

> **最佳实践**：
> - 不要在程序退出前关闭数据库连接，应使用 `defer` 语句确保在函数返回时关闭
> - 始终检查连接错误，`sql.Open` 可能返回无效连接但不报错
> - 使用 `db.Ping()` 验证连接是否可用
> - 对于生产环境，考虑使用配置文件或环境变量存储数据库凭证，而不是硬编码

### 连接池配置

```go
// 设置连接池参数
sqlDB, err := db.DB()
if err != nil {
    log.Fatalf("获取数据库连接池失败: %v", err)
}

sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数，根据并发量调整
sqlDB.SetMaxOpenConns(100)       // 最大打开连接数，避免连接耗尽
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期，防止连接过期
```

> **注意**：连接池参数应根据应用的并发量和数据库服务器性能进行调整。过高的连接数可能导致数据库服务器资源耗尽。
```

### 创建表

```go
func createTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        age INT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
    
    _, err := db.Exec(query)
    return err
}
```

### 插入数据

```go
type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
    CreatedAt time.Time
}

func insertUser(db *sql.DB, user *User) error {
    query := "INSERT INTO users (name, email, age) VALUES (?, ?, ?)"
    result, err := db.Exec(query, user.Name, user.Email, user.Age)
    if err != nil {
        return fmt.Errorf("插入用户失败: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("获取插入ID失败: %w", err)
    }
    
    user.ID = int(id)
    return nil
}
```

### 预处理语句

预处理语句可以提高性能并防止 SQL 注入攻击，特别适合重复执行的SQL语句。

```go
func batchInsertUsers(db *sql.DB, users []*User) error {
    // 准备预处理语句
    stmt, err := db.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return fmt.Errorf("准备预处理语句失败: %w", err)
    }
    defer stmt.Close()
    
    // 开始事务提高性能
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("开始事务失败: %w", err)
    }
    
    for _, user := range users {
        result, err := tx.Stmt(stmt).Exec(user.Name, user.Email, user.Age)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("插入用户 %s 失败: %w", user.Name, err)
        }
        
        id, err := result.LastInsertId()
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("获取插入ID失败: %w", err)
        }
        
        user.ID = int(id)
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("提交事务失败: %w", err)
    }
    
    return nil
}
```

> **安全提示**：始终使用参数化查询（? 占位符）而不是字符串拼接来构建 SQL 语句，以防止 SQL 注入攻击。
```

### 查询数据

#### 查询单行

```go
func getUserByID(db *sql.DB, id int) (*User, error) {
    query := "SELECT id, name, email, age, created_at FROM users WHERE id = ?"
    
    var user User
    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("用户不存在 (ID: %d)", id)
        }
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }
    
    return &user, nil
}
```

#### 查询多行与分页

```go
func getUsersWithPagination(db *sql.DB, page, pageSize int) ([]User, error) {
    // 计算偏移量
    offset := (page - 1) * pageSize
    
    query := "SELECT id, name, email, age, created_at FROM users LIMIT ? OFFSET ?"
    
    rows, err := db.Query(query, pageSize, offset)
    if err != nil {
        return nil, fmt.Errorf("查询用户列表失败: %w", err)
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("扫描用户数据失败: %w", err)
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("行迭代错误: %w", err)
    }
    
    return users, nil
}
```

> **性能提示**：对于大型结果集，始终使用分页查询并限制返回的行数，以避免内存溢出和提高查询效率。
```

#### 查询多行

```go
func getAllUsers(db *sql.DB) ([]User, error) {
    query := "SELECT id, name, email, age, created_at FROM users"
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

### 更新数据

```go
func updateUser(db *sql.DB, user *User) error {
    query := "UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?"
    _, err := db.Exec(query, user.Name, user.Email, user.Age, user.ID)
    return err
}
```

### 删除数据

```go
func deleteUser(db *sql.DB, id int) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := db.Exec(query, id)
    return err
}
```

## 使用GORM

GORM 是 Go 语言中最流行的 ORM 框架，它提供了简洁、强大的 API 来操作数据库。

### 安装GORM

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

### 连接数据库与配置

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "time"
)

func main() {
    dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    
    // 创建数据库连接
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        // 配置日志级别
        Logger: logger.Default.LogMode(logger.Info),
        // 禁用外键约束
        DisableForeignKeyConstraintWhenMigrating: true,
        // 禁用默认事务
        SkipDefaultTransaction: true,
    })
    if err != nil {
        log.Fatalf("连接数据库失败: %v", err)
    }
    
    // 配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("获取连接池失败: %v", err)
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

> **配置提示**：
> - 根据实际需求调整日志级别（Silent, Error, Warn, Info）
> - 对于性能要求高的应用，可考虑禁用默认事务 `SkipDefaultTransaction: true`
> - 禁用外键约束可以提高迁移和查询性能

### 定义模型

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"not null;size:100;comment:'用户名称'"`
    Email     string         `gorm:"unique;not null;size:100;comment:'邮箱地址'"`
    Age       int            `gorm:"default:0;comment:'年龄'"`
    Role      string         `gorm:"size:20;default:'user';comment:'用户角色'"`
    Active    bool           `gorm:"default:true;comment:'是否激活'"`
    CreatedAt time.Time      `gorm:"comment:'创建时间'"`
    UpdatedAt time.Time      `gorm:"comment:'更新时间'"`
    DeletedAt gorm.DeletedAt `gorm:"index;comment:'删除时间'"`
}

type Product struct {
    gorm.Model
    Code  string  `gorm:"unique;size:50;comment:'产品编码'"`
    Name  string  `gorm:"not null;size:100;comment:'产品名称'"`
    Price float64 `gorm:"not null;comment:'产品价格'"`
    Stock int     `gorm:"default:0;comment:'库存数量'"`
    CategoryID uint `gorm:"index;comment:'分类ID'"`
}
```

> **模型设计提示**：
> - 使用 `comment` 标签为字段添加注释，提高可维护性
> - 为常用查询字段添加索引 `index`，提高查询效率
> - 合理设置默认值，减少应用层处理逻辑
> - 对于金额等精确数值，考虑使用 `decimal` 类型而不是 `float`
```

### 连接数据库

```go
dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
    panic("failed to connect database")
}
```

### 自动迁移

```go
// 自动迁移表结构
db.AutoMigrate(&User{})
db.AutoMigrate(&Product{})
```

### CRUD操作

#### 创建记录

```go
// 创建单个用户
user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
result := db.Create(&user)
if result.Error != nil {
    log.Printf("创建用户失败: %v", result.Error)
} else {
    fmt.Println("New record ID:", user.ID)
}

// 批量创建 - 适用于大量数据插入
users := []User{
    {Name: "Bob", Email: "bob@example.com", Age: 30},
    {Name: "Charlie", Email: "charlie@example.com", Age: 35},
}
// 使用CreateInBatches提高批量插入性能
err := db.CreateInBatches(users, 100).Error
if err != nil {
    log.Printf("批量创建用户失败: %v", err)
}
```

> **性能提示**：批量插入时使用 `CreateInBatches` 而不是普通的 `Create`，可以显著提高性能。第二个参数指定每批插入的记录数。
```

#### 查询记录

```go
// 根据主键查询\var user User
// 查找第一个匹配记录
db.First(&user, 1)  // 查询ID为1的用户
// 条件查询
db.First(&user, "email = ?", "alice@example.com")

// 安全的错误处理
if err := db.First(&user, 1).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        log.Println("用户不存在")
    } else {
        log.Printf("查询用户失败: %v", err)
    }
}

// 高级条件查询
db.Where(&User{Name: "Alice", Age: 25}).Find(&users)  // 结构体条件
db.Where(map[string]interface{}{"name": "Alice", "age": 25}).Find(&users)  // Map条件

db.Where("age > ? AND name LIKE ?", 20, "%A%").Find(&users)  // 多条件

db.Where("age IN ?", []int{20, 25, 30}).Find(&users)  // IN条件

// 子查询
db.Where("age > (?) ", db.Table("users").Select("AVG(age)as avg_age")).Find(&users)

// 高级排序
db.Order("age desc, name asc").Find(&users)  // 多字段排序

db.Clauses(clause.OrderBy{  // 复杂排序
    Expression: clause.Expr{SQL: "FIELD(status, 'active', 'pending', 'inactive')"},
}).Find(&users)

// 分页查询
page := 1
pageSize := 10
db.Offset((page-1)*pageSize).Limit(pageSize).Find(&users)

// 关联预加载
db.Preload("Posts").First(&user)  // 预加载关联

db.Preload("Posts", "status = ?", "published").First(&user)  // 带条件的预加载
```

> **查询优化**：
> - 只选择需要的字段，使用 `Select` 方法减少数据传输
> - 合理使用索引，加速查询
> - 对于复杂查询，考虑使用原生 SQL
> - 避免在循环中执行查询，尽量批量查询
```

#### 更新记录

```go
// 更新单个字段
db.Model(&user).Update("Age", 26)

// 更新多个字段 - 注意：零值字段不会被更新
db.Model(&user).Updates(User{Name: "Alice Smith", Age: 27})

// 使用Map更新，可以更新零值字段
db.Model(&user).Updates(map[string]interface{}{
    "Name": "Alice Smith", 
    "Age": 27, 
    "Active": false,  // 可以更新布尔值false
})

// 更新选定字段，忽略零值
db.Model(&user).Select("Name", "Age").Updates(User{Name: "Alice Smith", Age: 0})

// 批量更新
db.Model(&User{}).Where("age < ?", 18).Update("age", 18)

// 表达式更新
db.Model(&Product{}).Update("price", gorm.Expr("price * 1.1"))
```

> **更新提示**：
> - 使用 `Select` 明确指定要更新的字段，避免意外更新
> - 对于批量更新，始终添加 `Where` 条件，避免全表更新
> - 使用 `gorm.Expr` 执行数据库端计算
```

#### 删除记录

```go
// 软删除 - 只标记删除时间，不实际删除数据
db.Delete(&user)

// 物理删除 - 实际从数据库中删除
// 注意：物理删除不可逆，请谨慎使用
db.Unscoped().Delete(&user)

// 根据条件删除
db.Where("age < ?", 18).Delete(&User{})

// 批量删除
db.Delete(&User{}, []int{1, 2, 3})  // 根据ID批量删除

db.Where("created_at < ?", time.Now().AddDate(0, -3, 0)).Delete(&User{})  // 删除3个月前的数据
```

> **删除提示**：
> - 默认情况下，GORM 执行软删除，仅设置 `DeletedAt` 字段
> - 使用 `Unscoped()` 进行物理删除
> - 批量删除前最好先查询确认，避免误删数据
```

### 关联关系

#### 一对多关系

```go
type User struct {
    gorm.Model
    Name  string
    Posts []Post `gorm:"foreignKey:UserID"`  // 显式指定外键
}

type Post struct {
    gorm.Model
    Title  string
    Body   string
    UserID uint `gorm:"index"`  // 外键，添加索引提高查询效率
    User   User  `gorm:"foreignKey:UserID"`  // 反向引用
}

// 创建关联
user := User{Name: "Alice"}
db.Create(&user)

db.Create(&Post{Title: "Hello", Body: "World", UserID: user.ID})

// 查询关联数据
var user User
db.Preload("Posts").First(&user)  // 预加载所有帖子

db.Preload("Posts", "title LIKE ?", "%Go%").First(&user)  // 预加载标题包含Go的帖子
```

#### 多对多关系

```go
type User struct {
    gorm.Model
    Name      string
    Languages []Language `gorm:"many2many:user_languages;foreignKey:ID;references:ID"`
}

type Language struct {
    gorm.Model
    Name  string
    Users []User `gorm:"many2many:user_languages;foreignKey:ID;references:ID"`
}

// 创建关联
user := User{Name: "Alice"}
language := Language{Name: "Go"}

db.Create(&user)
db.Create(&language)

// 添加关联
db.Model(&user).Association("Languages").Append(&language)

// 查询关联
db.Preload("Languages").First(&user)
```

> **关联提示**：
> - 显式指定外键和关联表，提高代码可读性
> - 为外键字段添加索引，提高查询性能
> - 使用 `Preload` 预加载关联数据，避免N+1查询问题
```

#### 多对多关系

```go
type User struct {
    gorm.Model
    Name   string
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name string
    Users []User `gorm:"many2many:user_languages;"`
}
```

### 事务处理

```go
// 自动事务处理 - 推荐使用
db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行操作
    if err := tx.Create(&user).Error; err != nil {
        return fmt.Errorf("创建用户失败: %w", err) // 回滚事务
    }
    
    profile := Profile{UserID: user.ID, Bio: "Golang Developer"}
    if err := tx.Create(&profile).Error; err != nil {
        return fmt.Errorf("创建用户资料失败: %w", err) // 回滚事务
    }
    
    return nil // 提交事务
})

// 手动控制事务
tx := db.Begin()
if tx.Error != nil {
    return fmt.Errorf("开始事务失败: %w", tx.Error)
}

defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return fmt.Errorf("创建用户失败: %w", err)
}

if err := tx.Commit().Error; err != nil {
    return fmt.Errorf("提交事务失败: %w", err)
}
```

> **事务最佳实践**：
> - 保持事务简短，减少锁定时间
> - 仅在必要时使用事务
> - 明确处理事务中的错误
> - 使用 `defer` 确保在 panic 时回滚事务
```

## 连接池配置

```go
sqlDB, err := db.DB()
if err != nil {
    panic(err)
}

// 设置连接池参数
sqlDB.SetMaxIdleConns(10)       // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)      // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最大生命周期
```

## Redis操作

Redis 是一个高性能的键值存储数据库，常用于缓存、会话管理和消息队列等场景。

### 安装Redis客户端

```bash
go get -u github.com/redis/go-redis/v9
```

### 连接Redis

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
    // 连接单机Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",  // 无密码
        DB:       0,   // 使用默认数据库
        PoolSize: 10,  // 连接池大小
    })
    
    // 测试连接
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("连接Redis失败: %v", err)
    }
    fmt.Println("Redis连接成功:", pong)
    
    // 连接Redis集群
    /*
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
        },
    })
    */
}
```

### 基本数据结构操作

```go
// 字符串操作
err := rdb.Set(ctx, "username", "alice", 10*time.Minute).Err()
if err != nil {
    log.Printf("设置键失败: %v", err)
}

val, err := rdb.Get(ctx, "username").Result()
if err != nil {
    if err == redis.Nil {
        log.Println("键不存在")
    } else {
        log.Printf("获取键失败: %v", err)
    }
}

// 列表操作
rdb.RPush(ctx, "users", "alice", "bob", "charlie")
users, err := rdb.LRange(ctx, "users", 0, -1).Result()

// 哈希操作
rdb.HSet(ctx, "user:1001", map[string]interface{}{
    "name": "alice",
    "age": 25,
    "city": "new york",
})

name, err := rdb.HGet(ctx, "user:1001", "name").Result()

// 集合操作
rdb.SAdd(ctx, "tags", "go", "python", "java")
members, err := rdb.SMembers(ctx, "tags").Result()

// 有序集合操作
rdb.ZAdd(ctx, "scores", redis.Z{Score: 95, Member: "alice"})
rdb.ZAdd(ctx, "scores", redis.Z{Score: 85, Member: "bob"})
rdb.ZAdd(ctx, "scores", redis.Z{Score: 90, Member: "charlie"})

// 获取前三名
rankings, err := rdb.ZRevRangeWithScores(ctx, "scores", 0, 2).Result()
```

### Redis高级功能

#### 事务

```go
// Redis事务
pipe := rdb.TxPipeline()

pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)

// 执行事务
_, err := pipe.Exec(ctx)
if err != nil {
    log.Printf("事务执行失败: %v", err)
}
```

#### 发布/订阅

```go
// 订阅频道
sub := rdb.Subscribe(ctx, "news")

// 接收消息
go func() {
    for {
        msg, err := sub.ReceiveMessage(ctx)
        if err != nil {
            log.Printf("接收消息失败: %v", err)
            return
        }
        fmt.Printf("接收到消息: %s 来自频道: %s\n", msg.Payload, msg.Channel)
    }
}()

// 发布消息
rdb.Publish(ctx, "news", "Hello Redis Pub/Sub")
```

#### 分布式锁

```go
// 简单的分布式锁实现
func acquireLock(rdb *redis.Client, lockKey string, expiration time.Duration) (string, error) {
    // 生成唯一标识符
    lockValue := fmt.Sprintf("lock:%d", time.Now().UnixNano())
    
    // 尝试获取锁
    ok, err := rdb.SetNX(ctx, lockKey, lockValue, expiration).Result()
    if err != nil {
        return "", err
    }
    
    if !ok {
        return "", fmt.Errorf("无法获取锁")
    }
    
    return lockValue, nil
}

func releaseLock(rdb *redis.Client, lockKey, lockValue string) error {
    // 使用Lua脚本确保原子性
    script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
    
    _, err := rdb.Eval(ctx, script, []string{lockKey}, lockValue).Result()
    return err
}
```

> **Redis最佳实践**：
> - 合理设置键的过期时间，避免内存泄漏
> - 使用连接池提高性能
> - 对于分布式系统，考虑使用RedLock算法实现更安全的分布式锁
> - 避免在Redis中存储大量大键，影响性能
> - 对敏感数据进行加密后再存储

## 数据库操作最佳实践总结

1. **连接管理**：
   - 使用连接池并合理配置参数
   - 避免频繁创建和关闭连接
   - 始终检查连接错误

2. **查询优化**：
   - 只选择必要的字段
   - 为常用查询添加索引
   - 避免在循环中执行查询
   - 使用分页查询处理大数据集

3. **安全措施**：
   - 使用参数化查询防止SQL注入
   - 加密敏感数据
   - 避免在代码中硬编码凭证
   - 限制数据库用户权限

4. **事务处理**：
   - 保持事务简短
   - 明确处理事务错误
   - 仅在必要时使用事务

5. **错误处理**：
   - 详细记录错误信息
   - 区分不同类型的错误（如不存在、权限不足等）
   - 提供有意义的错误信息给调用方

6. **性能监控**：
   - 监控慢查询
   - 分析查询执行计划
   - 定期优化数据库结构和索引

通过遵循这些最佳实践，可以编写更安全、更高效、更可维护的数据库操作代码。

### 基本操作

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // 设置值
    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    // 获取值
    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    // 设置过期时间
    err = rdb.Set(ctx, "expire_key", "expire_value", 10*time.Second).Err()
    if err != nil {
        panic(err)
    }

    // 列表操作
    rdb.RPush(ctx, "mylist", "1", "2", "3")
    length := rdb.LLen(ctx, "mylist").Val()
    fmt.Println("List length:", length)

    // Hash操作
    rdb.HSet(ctx, "user:1000", "name", "Alice", "age", 25)
    name := rdb.HGet(ctx, "user:1000", "name").Val()
    fmt.Println("User name:", name)
}
```