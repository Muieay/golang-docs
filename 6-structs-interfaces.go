package main

import (
    "fmt"
    "math"
)

// 结构体和接口
// Go语言的面向对象编程

/*
1. 结构体定义
type 关键字用于定义新的类型
struct 关键字定义结构体，可以包含多个不同类型的字段
结构体是Go语言中实现面向对象编程的基础
*/
type Person struct {
    Name    string // 姓名字段
    Age     int    // 年龄字段
    Email   string // 邮箱字段
    Address string // 地址字段
}

/*
2. 结构体方法
方法接收者分为值接收者和指针接收者两种
值接收者(p Person)不会修改原对象，适用于只读操作
指针接收者(p *Person)会修改原对象，适用于写操作
*/

// GetInfo 方法返回Person的格式化信息
func (p Person) GetInfo() string {
    return fmt.Sprintf("姓名: %s, 年龄: %d, 邮箱: %s", p.Name, p.Age, p.Email)
}

// UpdateAge 方法更新Person的年龄
// 使用指针接收者确保修改能作用到原对象
func (p *Person) UpdateAge(newAge int) {
    p.Age = newAge
}

/*
3. 嵌套结构体
Go支持结构体嵌套实现组合关系
匿名嵌套(嵌入)可以直接访问嵌套结构体的字段
命名嵌套需要通过字段名访问
*/
type Address struct {
    City     string // 城市
    Street   string // 街道
    ZipCode  string // 邮编
}

type Employee struct {
    Person  // 匿名嵌套Person，可以直接访问Person的字段
    ID      int
    Salary  float64
    Address Address // 命名嵌套，需要通过Address字段访问
}

/*
4. 接口定义
interface 关键字定义接口
接口是一组方法的集合，任何实现了这些方法的类型都实现了该接口
Go的接口是隐式实现的，不需要显式声明
*/
type Shape interface {
    Area() float64      // 计算面积
    Perimeter() float64 // 计算周长
}

/*
5. 实现接口的结构体
Rectangle、Circle、Triangle都实现了Shape接口
通过为每个类型定义Area()和Perimeter()方法实现
*/
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

// Rectangle的Area实现
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Rectangle的Perimeter实现
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle的Area实现
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Circle的Perimeter实现
func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Triangle的Area实现
func (t Triangle) Area() float64 {
    return 0.5 * t.Base * t.Height
}

// Triangle的Perimeter实现
func (t Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

/*
6. 空接口和类型断言
interface{} 是空接口，可以表示任何类型
类型断言用于判断接口值的具体类型
switch type语句可以方便地进行类型判断
*/
func describe(i interface{}) {
    fmt.Printf("类型: %T, 值: %v\n", i, i)
    
    // 类型断言示例
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
    
    // 1. 结构体使用示例
    fmt.Println("\n--- 结构体基础 ---")
    
    // 结构体初始化方式1：指定字段名
    person1 := Person{
        Name:  "张三",
        Age:   25,
        Email: "zhangsan@example.com",
    }
    
    // 结构体初始化方式2：按字段顺序
    person2 := Person{"李四", 30, "lisi@example.com", "北京市"}
    
    // 零值结构体
    var person3 Person
    
    // 打印结构体内容，%+v会显示字段名
    fmt.Printf("person1: %+v\n", person1)
    fmt.Printf("person2: %+v\n", person2)
    fmt.Printf("person3: %+v\n", person3)
    
    // 修改结构体字段
    person1.Age = 26
    person1.Address = "上海市浦东新区"
    
    // 调用结构体方法
    fmt.Printf("person1信息: %s\n", person1.GetInfo())
    person1.UpdateAge(27) // 调用指针接收者方法
    
    // 2. 嵌套结构体示例
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
    
    // 访问嵌套字段
    fmt.Printf("员工姓名: %s\n", employee.Name) // 匿名嵌套直接访问
    fmt.Printf("员工城市: %s\n", employee.Address.City) // 命名嵌套通过字段访问
    
    // 3. 接口和多态示例
    fmt.Println("\n--- 接口和多态 ---")
    
    // 创建Shape接口切片，实现多态
    shapes := []Shape{
        Rectangle{Width: 4, Height: 5},
        Circle{Radius: 3},
        Triangle{Base: 4, Height: 3, SideA: 3, SideB: 4, SideC: 5},
    }
    
    // 遍历调用接口方法
    for _, shape := range shapes {
        fmt.Printf("面积: %.2f\n", shape.Area())
        fmt.Printf("周长: %.2f\n", shape.Perimeter())
    }
    
    // 4. 接口作为函数参数
    fmt.Println("\n--- 接口作为参数 ---")
    printShapeInfo(Rectangle{Width: 5, Height: 3})
    printShapeInfo(Circle{Radius: 2})
    
    // 5. 空接口和类型断言
    fmt.Println("\n--- 空接口和类型断言 ---")
    describe(42)
    describe("hello")
    
    // 6. 结构体标签示例
    fmt.Println("\n--- 结构体标签 ---")
    type User struct {
        ID       int    `json:"id"` // 结构体标签，用于JSON序列化
        Username string `json:"username"`
    }
    
    // 7. 结构体比较
    fmt.Println("\n--- 结构体比较 ---")
    p1 := Person{Name: "张三", Age: 25}
    p2 := Person{Name: "张三", Age: 25}
    fmt.Printf("p1 == p2: %t\n", p1 == p2) // 结构体可比较
    
    // 8. 结构体指针
    fmt.Println("\n--- 结构体指针 ---")
    personPtr := &Person{Name: "赵六", Age: 28}
    personPtr.UpdateAge(29) // 指针调用方法
}

// printShapeInfo 函数展示接口作为参数
func printShapeInfo(s Shape) {
    fmt.Printf("面积: %.2f\n", s.Area())
}

// Stringer 接口示例
type Stringer interface {
    String() string
}

// Book 实现Stringer接口
type Book struct {
    Title  string
    Author string
    Price  float64
}

func (b Book) String() string {
    return fmt.Sprintf("《%s》- %s (¥%.2f)", b.Title, b.Author, b.Price)
}