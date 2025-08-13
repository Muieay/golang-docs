# Golang数据库操作

## 数据库/sql标准库

### 安装驱动

```bash
go get -u github.com/go-sql-driver/mysql
```

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
    // 连接MySQL
    dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // 测试连接
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Successfully connected!")
}
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
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    user.ID = int(id)
    return nil
}
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
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}
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

### 安装GORM

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

### 定义模型

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "time"
)

type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"not null;size:100"`
    Email     string         `gorm:"unique;not null;size:100"`
    Age       int            `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
    gorm.Model
    Code  string `gorm:"unique"`
    Price uint
}
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
// 创建用户
user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
result := db.Create(&user)
fmt.Println("New record ID:", user.ID)

// 批量创建
users := []User{
    {Name: "Bob", Email: "bob@example.com", Age: 30},
    {Name: "Charlie", Email: "charlie@example.com", Age: 35},
}
db.Create(&users)
```

#### 查询记录

```go
// 根据主键查询
var user User
db.First(&user, 1)  // 查询ID为1的用户
db.First(&user, "email = ?", "alice@example.com")  // 查询email为指定值的用户

// 获取所有记录
var users []User
db.Find(&users)

// 条件查询
db.Where("age > ?", 20).Find(&users)
db.Where("name LIKE ?", "%Alice%").Find(&users)
db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)

// 排序和限制
db.Order("age desc").Limit(10).Find(&users)
db.Offset(10).Limit(5).Find(&users)

// 选择特定字段
db.Select("name", "email").Find(&users)

// 聚合查询
var count int64
db.Model(&User{}).Where("age > ?", 20).Count(&count)

var totalAge int64
db.Model(&User{}).Select("SUM(age)").Scan(&totalAge)
```

#### 更新记录

```go
// 更新单个字段
db.Model(&user).Update("Age", 26)

// 更新多个字段
db.Model(&user).Updates(User{Name: "Alice Smith", Age: 27})
db.Model(&user).Updates(map[string]interface{}{"Name": "Alice Smith", "Age": 27})

// 更新所有记录
db.Model(&User{}).Where("age < ?", 18).Update("age", 18)
```

#### 删除记录

```go
// 删除记录
db.Delete(&user)

// 根据条件删除
db.Where("age < ?", 18).Delete(&User{})

// 批量删除
db.Delete(&User{}, []int{1, 2, 3})
```

### 关联关系

#### 一对多关系

```go
type User struct {
    gorm.Model
    Name  string
    Posts []Post
}

type Post struct {
    gorm.Model
    Title  string
    Body   string
    UserID uint
}

// 查询关联数据
var user User
db.Preload("Posts").First(&user)
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
// 使用事务
db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行操作
    if err := tx.Create(&user).Error; err != nil {
        return err // 回滚事务
    }
    
    if err := tx.Create(&profile).Error; err != nil {
        return err // 回滚事务
    }
    
    return nil // 提交事务
})

// 手动控制事务
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
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

### 安装Redis客户端

```bash
go get -u github.com/redis/go-redis/v9
```

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