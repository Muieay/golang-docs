# Golang进阶教程

## 指针与内存管理

### 指针基础

```go
package main

import "fmt"

func main() {
    x := 42
    p := &x  // p是指向x的指针
    
    fmt.Println("x的值:", x)    // 42
    fmt.Println("x的地址:", &x)  // 内存地址
    fmt.Println("p的值:", p)     // 与&x相同
    fmt.Println("p指向的值:", *p)  // 42
    
    *p = 100  // 通过指针修改值
    fmt.Println("新的x值:", x)     // 100
}
```

### 指针与结构体

```go
type Person struct {
    Name string
    Age  int
}

func main() {
    p := &Person{"Alice", 30}
    p.Age = 31  // 自动解引用
    fmt.Println(p.Name)
}
```

## 结构体与方法

### 结构体定义

```go
type Rectangle struct {
    Width, Height float64
}

// 值接收者方法
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 指针接收者方法
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

### 嵌套结构体

```go
type Address struct {
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address Address  // 嵌套结构体
}

func main() {
    p := Person{
        Name: "Bob",
        Age:  25,
        Address: Address{
            City:    "Beijing",
            Country: "China",
        },
    }
    
    fmt.Println(p.Address.City)
}
```

## 接口与多态

### 接口定义

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Rectangle struct {
    Width, Height float64
}

type Circle struct {
    Radius float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

func printShapeInfo(s Shape) {
    fmt.Printf("面积: %.2f\n", s.Area())
    fmt.Printf("周长: %.2f\n", s.Perimeter())
}
```

### 空接口

```go
func printAnything(v interface{}) {
    fmt.Printf("类型: %T, 值: %v\n", v, v)
}

func main() {
    printAnything(42)
    printAnything("hello")
    printAnything([]int{1, 2, 3})
}
```

## 反射

```go
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    u := User{"Alice", 30}
    t := reflect.TypeOf(u)
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("字段名: %s, 类型: %s, JSON标签: %s\n", 
            field.Name, field.Type, field.Tag.Get("json"))
    }
}
```

## 泛型

```go
// 泛型函数
func Min[T comparable](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 泛型结构体
type Stack[T any] struct {
    elements []T
}

func (s *Stack[T]) Push(element T) {
    s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.elements) == 0 {
        var zero T
        return zero, false
    }
    element := s.elements[len(s.elements)-1]
    s.elements = s.elements[:len(s.elements)-1]
    return element, true
}
```

## 测试

### 单元测试

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### 基准测试

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

### 表驱动测试

```go
func TestAddTableDriven(t *testing.T) {
    tests := []struct {
        name     string
        a, b, expected int
    }{
        {"2+3", 2, 3, 5},
        {"0+0", 0, 0, 0},
        {"-1+1", -1, 1, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```