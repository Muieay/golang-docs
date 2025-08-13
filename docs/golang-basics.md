# Golang基础教程

## 环境安装与配置

### 下载与安装

访问 [Go官网](https://golang.org/dl/) 下载适合你操作系统的安装包。

### 配置环境变量
```bash
# Linux/macOS
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go

# Windows go1.16以后的版本，不需要配置环境变量，直接安装即可。
# 将 C:\Go\bin 添加到 PATH 环境变量
```

### 验证安装

```bash
go version
```

## Hello World

```go
package main

import (
    "fmt"
    "os"
    "sort"
)

func main() {
    fmt.Println("Hello, Golang!")
}
```

## 基础语法
- Go程序由包（package）组成，入口为`main`包的`main()`函数。
- 每行语句以换行结束，无需分号。
- 使用`import`导入标准库或第三方包。

## fmt包
Println、、Printf、Print 比较

- `fmt.Print()`：直接输出内容，不自动换行。
- `fmt.Println()`：输出内容后自动换行。
- `fmt.Printf()`：格式化输出，支持格式化占位符。

### 常用格式化占位符
- `%d`：十进制整数
- `%s`：字符串
- `%f`：浮点数（默认6位小数）
- `%v`：自动匹配变量类型的默认格式
- `%+v`：结构体时，输出字段名和值
- `%#v`：Go语法表示的值
- `%T`：输出变量类型
- `%p`：指针
- `%c`: 字符（rune）
- `%t`：布尔值
- `%x`：十六进制整数
- `%o`：八进制整数
- `%b`：二进制整数

### Print/Println/Printf区别
- `Print`：参数直接输出，无分隔符，不换行。
- `Println`：参数间以空格分隔，输出后自动换行。
- `Printf`：可自定义格式，需指定格式化字符串。

### 示例
```go
package main

import "fmt"

func main() {
    a := 10
    b := "hello"
    c := 3.1415
    fmt.Print(a, b)             // 输出: 10hello
    fmt.Println(a, b)           // 输出: 10 hello（并换行）
    fmt.Printf("%d %s %.2f\n", a, b, c) // 输出: 10 hello 3.14
    fmt.Printf("类型: %T, 值: %v\n", c, c) // 输出: 类型: float64, 值: 3.1415
}
```


## 变量声明

```go
// 方式1：var关键字
var name string = "Golang"
var age int = 10

// 方式2：类型推断
var name = "Golang"

// 方式3：简短声明（只能在函数内部使用）
name := "Golang"
```

## 匿名变量
- 使用`_`作为变量名，表示忽略该值。
- 常用于函数返回多个值时忽略不需要的值。

#### 示例
```go
func foo() (int, string) {
    return 1, "hello"
}

func main() {
    a, _ := foo() // 只接收第一个返回值，忽略第二个
    _, b := foo() // 只接收第二个返回值，忽略第一个
    fmt.Println(a) // 输出: 1
    fmt.Println(b) // 输出: hello
}
```

## 数据类型

#### 基本类型

- 布尔型：`bool`
- 整型：`int`, `int8`, `int16`, `int32`, `int64`
- 无符号整型：`uint`, `uint8`, `uint16`, `uint32`, `uint64`
- 浮点型：`float32`, `float64`
- 复数型：`complex64`, `complex128`
- 字符串：`string`
- 字符：`rune`（Unicode码点）

#### 复合类型

- 数组：`[n]T`
- 切片：`[]T`
- 映射：`map[K]V`
- 结构体：`struct`
- 接口：`interface`

## 控制结构

#### 条件语句

```go
// if语句
if x > 0 {
    return x
} else if x < 0 {
    return -x
} else {
    return 0
}

// switch语句
switch day {
case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
    fmt.Println("工作日")
case "Saturday", "Sunday":
    fmt.Println("周末")
default:
    fmt.Println("无效的日期")
}

// switch不带表达式
switch {
case score >= 90:
    fmt.Println("优秀")
case score >= 80:
    fmt.Println("良好")
case score >= 60:
    fmt.Println("及格")
default:
    fmt.Println("不及格")
}
```

#### 循环语句

```go
// 传统for循环
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// while风格的循环
for i < 100 {
    i += 2
}

// 无限循环
for {
    // 死循环
}

// range循环
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

#### 流程控制语句

##### break语句
`break`用于立即跳出当前循环或switch语句。

```go
// 在循环中使用break
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // 当i等于5时跳出循环
    }
    fmt.Println(i)  // 输出: 0 1 2 3 4
}

// 在嵌套循环中使用break
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 2 {
            break  // 只跳出内层循环
        }
        fmt.Printf("i: %d, j: %d\n", i, j)
    }
}

// 使用标签跳出指定循环
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i*j == 2 {
            break outer  // 跳出外层循环
        }
        fmt.Printf("i: %d, j: %d\n", i, j)
    }
}
```

##### continue语句
`continue`用于跳过当前循环的剩余语句，直接进入下一次循环。

```go
// 跳过偶数
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // 跳过偶数
    }
    fmt.Println(i)  // 输出: 1 3 5 7 9
}

// 在嵌套循环中使用continue
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            continue  // 跳过j等于1的情况
        }
        fmt.Printf("i: %d, j: %d\n", i, j)
    }
}

// 使用标签控制外层循环
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            continue outer  // 跳过外层循环的剩余语句
        }
        fmt.Printf("i: %d, j: %d\n", i, j)
    }
}
```

##### goto语句
`goto`用于跳转到程序中的指定标签位置。虽然功能强大，但应谨慎使用，以免造成代码难以理解和维护。

```go
// 基本使用
func printNumbers() {
    i := 0
loop:
    if i >= 5 {
        goto end
    }
    fmt.Println(i)
    i++
    goto loop
end:
    fmt.Println("循环结束")
}

// 错误处理场景
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    data := make([]byte, 100)
    _, err = file.Read(data)
    if err != nil {
        goto cleanup
    }

    // 处理数据...
    if err != nil {
        goto cleanup
    }

    return nil

cleanup:
    // 清理操作
    fmt.Println("执行清理操作")
    return err
}
```

### 函数

```go
// 基本函数定义
func add(a int, b int) int {
    return a + b
}

// 多返回值
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}
```

## 数据结构详解

### 数组（Array）
数组是值类型，长度固定，声明时需要指定长度。

#### 基本使用
```go
// 声明数组
var arr1 [5]int                           // 声明长度为5的int数组，默认值为0
var arr2 = [5]int{1, 2, 3, 4, 5}         // 声明并初始化
arr3 := [5]int{1, 2, 3, 4, 5}            // 简短声明
arr4 := [...]int{1, 2, 3, 4, 5}          // 自动推断长度
arr5 := [5]int{1: 10, 3: 30}             // 指定索引初始化

// 访问数组元素
fmt.Println(arr2[0])     // 输出: 1
fmt.Println(arr2[len(arr2)-1])  // 输出: 5

// 修改数组元素
arr2[2] = 100

// 遍历数组
for i, v := range arr2 {
    fmt.Printf("索引: %d, 值: %d\n", i, v)
}
```

#### 数组特性
- **值类型**：数组赋值和传参会复制整个数组
- **长度固定**：数组长度是类型的一部分，[5]int和[10]int是不同的类型
- **内存连续**：数组元素在内存中是连续存储的

```go
// 值类型特性演示
func modifyArray(arr [5]int) {
    arr[0] = 100  // 修改的是副本
}

func main() {
    arr := [5]int{1, 2, 3, 4, 5}
    modifyArray(arr)
    fmt.Println(arr)  // 输出: [1 2 3 4 5]，原数组未被修改
}

// 如果想要修改原数组，需要使用指针
func modifyArrayPtr(arr *[5]int) {
    arr[0] = 100
}
```

### 切片（Slice）
切片是引用类型，动态数组，基于数组构建，长度可变。

#### 基本使用
```go
// 声明切片
var slice1 []int                        // 声明nil切片
slice2 := []int{1, 2, 3, 4, 5}          // 声明并初始化
slice3 := make([]int, 5)                // 创建长度为5的切片
slice4 := make([]int, 3, 5)             // 创建长度为3，容量为5的切片

// 从数组创建切片
arr := [5]int{1, 2, 3, 4, 5}
slice5 := arr[1:4]  // 从索引1到4（不包含4），得到[2 3 4]

// 切片操作
fmt.Println(len(slice2))  // 长度: 5
fmt.Println(cap(slice2))  // 容量: 5

// 追加元素
slice2 = append(slice2, 6, 7, 8)  // 追加多个元素

// 切片的切片
subSlice := slice2[1:4]  // 得到[2 3 4]
```

#### 切片底层原理
```go
// 切片结构体（概念理解）
type slice struct {
    ptr *Element  // 指向底层数组的指针
    len int       // 当前长度
    cap int       // 容量（底层数组的最大长度）
}

// 切片扩容机制
func sliceExpansion() {
    s := make([]int, 0, 3)  // 初始容量为3
    fmt.Printf("初始: len=%d, cap=%d\n", len(s), cap(s))
    
    for i := 0; i < 5; i++ {
        s = append(s, i)
        fmt.Printf("添加%d后: len=%d, cap=%d\n", i, len(s), cap(s))
    }
}
```

#### 切片高级操作
```go
// 复制切片
dst := make([]int, len(src))
copy(dst, src)  // 深拷贝

// 删除元素
// 删除索引为2的元素
slice := []int{1, 2, 3, 4, 5}
slice = append(slice[:2], slice[3:]...)  // 得到[1 2 4 5]

// 插入元素
// 在索引2处插入100
slice = append(slice[:2], append([]int{100}, slice[2:]...)...)

// 清空切片
slice = slice[:0]  // 长度变为0，但容量不变
```

#### 切片作为函数参数
```go
// 切片是引用类型，函数内部修改会影响原切片
func modifySlice(s []int) {
    if len(s) > 0 {
        s[0] = 100  // 修改原切片
    }
}

func main() {
    slice := []int{1, 2, 3, 4, 5}
    modifySlice(slice)
    fmt.Println(slice)  // 输出: [100 2 3 4 5]
}
```

### Map（映射）
Map是引用类型，无序的键值对集合。

#### 基本使用
```go
// 声明map
var map1 map[string]int                 // 声明nil map
map2 := make(map[string]int)            // 使用make创建
map3 := map[string]int{                // 字面量创建
    "apple":  5,
    "banana": 3,
    "orange": 8,
}

// 添加/修改元素
map2["key1"] = 100
map2["key2"] = 200

// 访问元素
value := map3["apple"]      // 获取值
value, ok := map3["pear"]  // 检查键是否存在
if ok {
    fmt.Println("存在:", value)
} else {
    fmt.Println("不存在")
}

// 删除元素
delete(map3, "banana")

// 遍历map
for key, value := range map3 {
    fmt.Printf("%s: %d\n", key, value)
}
```

#### Map高级特性
```go
// map的map（嵌套map）
nestedMap := make(map[string]map[string]int)
nestedMap["fruits"] = make(map[string]int)
nestedMap["fruits"]["apple"] = 5

// 值是切片的map
mapWithSlice := make(map[string][]int)
mapWithSlice["scores"] = []int{90, 85, 78, 92}

// 并发安全的map（需要使用sync.Map）
// 或者使用读写锁保护普通map
```

#### Map作为函数参数
```go
// map是引用类型，函数内部修改会影响原map
func modifyMap(m map[string]int) {
    m["newKey"] = 999
}

func main() {
    myMap := map[string]int{"a": 1, "b": 2}
    modifyMap(myMap)
    fmt.Println(myMap)  // 输出包含newKey: 999
}
```

#### Map排序
```go
// map是无序的，需要排序时转换为切片
func sortMapByKey() {
    m := map[string]int{"c": 3, "a": 1, "b": 2}
    
    // 提取键并排序
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    // 按排序后的键访问map
    for _, k := range keys {
        fmt.Printf("%s: %d\n", k, m[k])
    }
}

func sortMapByValue() {
    m := map[string]int{"apple": 5, "banana": 3, "orange": 8}
    
    // 创建键值对切片
    pairs := make([]struct{ Key string; Value int }, 0, len(m))
    for k, v := range m {
        pairs = append(pairs, struct{ Key string; Value int }{k, v})
    }
    
    // 按值排序
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].Value < pairs[j].Value
    })
    
    for _, pair := range pairs {
        fmt.Printf("%s: %d\n", pair.Key, pair.Value)
    }
}
```

## 错误处理

```go
// 错误检查
result, err := someFunction()
if err != nil {
    // 处理错误
    log.Fatal(err)
}

// 自定义错误
type MyError struct {
    Msg string
    Code int
}

func (e *MyError) Error() string {
    return fmt.Sprintf("错误 %d: %s", e.Code, e.Msg)
}

## 综合示例

下面是一个综合使用数组、切片和map的示例：

```go
package main

import "fmt"

func main() {
    // 数组示例
    var scores [5]int = [5]int{90, 85, 78, 92, 88}
    
    // 切片示例 - 从数组创建切片
    topScores := scores[1:4]  // [85 78 92]
    topScores = append(topScores, 95, 97)
    
    // Map示例 - 存储学生成绩
    studentScores := map[string]int{
        "张三": scores[0],
        "李四": scores[1],
        "王五": scores[2],
    }
    
    // 添加新学生
    studentScores["赵六"] = topScores[len(topScores)-1]
    
    // 遍历并显示
    fmt.Println("学生成绩:")
    for name, score := range studentScores {
        fmt.Printf("%s: %d分\n", name, score)
    }
    
    // 使用切片存储所有成绩并计算平均
    allScores := make([]int, 0, len(studentScores))
    for _, score := range studentScores {
        allScores = append(allScores, score)
    }
    
    sum := 0
    for _, score := range allScores {
        sum += score
    }
    
    average := sum / len(allScores)
    fmt.Printf("平均成绩: %d分\n", average)
}
```