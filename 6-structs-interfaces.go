package main

import (
    "fmt"
    "math"
)

// 结构体和接口
// Go语言的面向对象编程

// 1. 结构体定义
// 定义一个Person结构体

/*
结构体是Go语言中的复合数据类型，可以包含多个不同类型的字段
*/

type Person struct {
    Name    string
    Age     int
    Email   string
    Address string
}

// 2. 结构体方法
// 给Person结构体添加方法

// 值接收者方法（不会修改原对象）
func (p Person) GetInfo() string {
    return fmt.Sprintf("姓名: %s, 年龄: %d, 邮箱: %s", p.Name, p.Age, p.Email)
}

// 指针接收者方法（会修改原对象）
func (p *Person) UpdateAge(newAge int) {
    p.Age = newAge
}

// 3. 嵌套结构体
type Address struct {
    City     string
    Street   string
    ZipCode  string
}

type Employee struct {
    Person  // 匿名嵌套，继承Person的所有字段
    ID      int
    Salary  float64
    Address Address  // 嵌套结构体
}

// 4. 接口定义
// 定义几何图形接口

type Shape interface {
    Area() float64
    Perimeter() float64
}

// 5. 实现接口的结构体
type Rectangle struct {
    Width  float64
    Height float64
}

type Circle struct {
    Radius float64
}

type Triangle struct {
    Base   float64
    Height float64
    SideA  float64
    SideB  float64
    SideC  float64
}

// 实现Shape接口
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func (t Triangle) Area() float64 {
    return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

// 6. 空接口和类型断言
func describe(i interface{}) {
    fmt.Printf("类型: %T, 值: %v\n", i, i)
    
    // 类型断言
    switch v := i.(type) {
    case int:
        fmt.Printf("这是整数: %d\n", v)
    case string:
        fmt.Printf("这是字符串: %s\n", v)
    case Person:
        fmt.Printf("这是Person: %s\n", v.Name)
    default:
        fmt.Printf("未知类型\n")
    }
}

func main() {
    fmt.Println("=== Go语言结构体和接口 ===")
    
    // 1. 结构体使用
    fmt.Println("\n--- 结构体基础 ---")
    
    // 创建结构体实例
    person1 := Person{
        Name:  "张三",
        Age:   25,
        Email: "zhangsan@example.com",
    }
    
    // 另一种创建方式
    person2 := Person{"李四", 30, "lisi@example.com", "北京市"}
    
    // 零值结构体
    var person3 Person
    
    fmt.Printf("person1: %+v\n", person1)
    fmt.Printf("person2: %+v\n", person2)
    fmt.Printf("person3: %+v\n", person3)
    
    // 访问和修改字段
    person1.Age = 26
    person1.Address = "上海市浦东新区"
    fmt.Printf("更新后的person1: %+v\n", person1)
    
    // 调用结构体方法
    fmt.Printf("person1信息: %s\n", person1.GetInfo())
    person1.UpdateAge(27)
    fmt.Printf("更新年龄后: %d\n", person1.Age)
    
    // 2. 嵌套结构体
    fmt.Println("\n--- 嵌套结构体 ---")
    
    employee := Employee{
        Person: Person{
            Name:  "王五",
            Age:   35,
            Email: "wangwu@company.com",
        },
        ID:     1001,
        Salary: 8000.0,
        Address: Address{
            City:    "深圳",
            Street:  "科技园路",
            ZipCode: "518000",
        },
    }
    
    fmt.Printf("员工信息: %+v\n", employee)
    fmt.Printf("员工姓名: %s\n", employee.Name)  // 直接访问嵌套结构体的字段
    fmt.Printf("员工城市: %s\n", employee.Address.City)
    
    // 3. 接口使用
    fmt.Println("\n--- 接口和多态 ---")
    
    shapes := []Shape{
        Rectangle{Width: 4, Height: 5},
        Circle{Radius: 3},
        Triangle{Base: 4, Height: 3, SideA: 3, SideB: 4, SideC: 5},
    }
    
    for _, shape := range shapes {
        fmt.Printf("图形: %T\n", shape)
        fmt.Printf("面积: %.2f\n", shape.Area())
        fmt.Printf("周长: %.2f\n", shape.Perimeter())
        fmt.Println()
    }
    
    // 4. 接口作为函数参数
    fmt.Println("\n--- 接口作为参数 ---")
    
    rect := Rectangle{Width: 5, Height: 3}
    circle := Circle{Radius: 2}
    
    printShapeInfo(rect)
    printShapeInfo(circle)
    
    // 5. 空接口和类型断言
    fmt.Println("\n--- 空接口和类型断言 ---")
    
    describe(42)
    describe("hello")
    describe(person1)
    describe(3.14)
    
    // 类型断言示例
    var i interface{} = "这是一个字符串"
    if str, ok := i.(string); ok {
        fmt.Printf("类型断言成功，字符串内容: %s\n", str)
    } else {
        fmt.Println("类型断言失败")
    }
    
    // 6. 结构体标签（用于JSON序列化等）
    fmt.Println("\n--- 结构体标签 ---")
    
    type User struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
        Email    string `json:"email"`
        Age      int    `json:"age,omitempty"`  // omitempty表示空值时省略
    }
    
    user := User{
        ID:       1,
        Username: "alice",
        Email:    "alice@example.com",
        Age:      25,
    }
    
    fmt.Printf("用户结构体: %+v\n", user)
    
    // 7. 结构体比较
    fmt.Println("\n--- 结构体比较 ---")
    
    p1 := Person{Name: "张三", Age: 25}
    p2 := Person{Name: "张三", Age: 25}
    p3 := Person{Name: "李四", Age: 30}
    
    fmt.Printf("p1 == p2: %t\n", p1 == p2)
    fmt.Printf("p1 == p3: %t\n", p1 == p3)
    
    // 8. 结构体指针
    fmt.Println("\n--- 结构体指针 ---")
    
    personPtr := &Person{Name: "赵六", Age: 28}
    fmt.Printf("personPtr: %p\n", personPtr)
    fmt.Printf("personPtr指向的值: %+v\n", *personPtr)
    
    personPtr.UpdateAge(29)  // 指针方法可以通过值或指针调用
    fmt.Printf("更新后的年龄: %d\n", personPtr.Age)
    
    // 练习
    fmt.Println("\n练习：")
    fmt.Println("1. 创建一个Book结构体，包含标题、作者、价格字段")
    fmt.Println("2. 定义一个Stringer接口，实现String()方法")
    fmt.Println("3. 创建一个图书馆管理系统，使用结构体切片存储图书")
}

// 打印图形信息的函数，使用接口作为参数
func printShapeInfo(s Shape) {
    fmt.Printf("图形类型: %T\n", s)
    fmt.Printf("面积: %.2f\n", s.Area())
    fmt.Printf("周长: %.2f\n", s.Perimeter())
}

// 额外的接口示例
type Stringer interface {
    String() string
}

// 实现Stringer接口
type Book struct {
    Title  string
    Author string
    Price  float64
}

func (b Book) String() string {
    return fmt.Sprintf("《%s》- %s (¥%.2f)", b.Title, b.Author, b.Price)
}