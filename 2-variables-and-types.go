package main

import (
    "fmt"
    "reflect"
)

// 变量和数据类型详解
// Go是静态类型语言，每个变量都有明确的类型

func main() {
    fmt.Println("=== Go语言变量和数据类型 ===")
    
    // 1. 变量声明方式
    
    // 方式1：var关键字声明
    var age int = 25
    var name string = "张三"
    var isStudent bool = true
    
    fmt.Printf("年龄: %d, 姓名: %s, 是否学生: %t\n", age, name, isStudent)
    
    // 方式2：类型推断
    var score = 95.5  // 自动推断为float64
    var city = "北京"   // 自动推断为string
    fmt.Printf("分数类型: %T, 城市类型: %T\n", score, city)
    
    // 方式3：简短声明（最常用）
    country := "中国"
    year := 2024
    fmt.Printf("国家: %s, 年份: %d\n", country, year)
    
    // 2. 基本数据类型
    
    // 整数类型
    var intNum int = 100
    var int8Num int8 = 127        // -128 到 127
    var int16Num int16 = 32767    // -32768 到 32767
    var uintNum uint = 200        // 无符号整数
    var byteNum byte = 255        // byte是uint8的别名
    
    fmt.Printf("int: %d, int8: %d, uint: %d, byte: %d\n", 
        intNum, int8Num, uintNum, byteNum)
    
    // 浮点类型
    var float32Num float32 = 3.14
    var float64Num float64 = 3.141592653589793
    fmt.Printf("float32: %.2f, float64: %.15f\n", float32Num, float64Num)
    
    // 字符串类型
    str1 := "Hello"
    str2 := `Raw string
    可以换行
    支持"引号"`
    fmt.Printf("str1: %s\nstr2: %s\n", str1, str2)
    
    // 布尔类型
    var bool1 bool = true
    var bool2 bool = false
    var bool3 bool // 默认为false
    fmt.Printf("bool1: %t, bool2: %t, bool3: %t\n", bool1, bool2, bool3)
    
    // 3. 零值（默认值）
    var zeroInt int
    var zeroString string
    var zeroBool bool
    var zeroFloat float64
    fmt.Printf("零值 - int: %d, string: %q, bool: %t, float: %f\n", 
        zeroInt, zeroString, zeroBool, zeroFloat)
    
    // 4. 类型转换
    var a int = 100
    var b float64 = float64(a)  // 显式类型转换
    var c int32 = int32(a)
    fmt.Printf("类型转换: %d -> %.1f -> %d\n", a, b, c)
    
    // 5. 使用reflect查看类型
    fmt.Printf("reflect - age的类型: %v\n", reflect.TypeOf(age))
    
    // 练习：尝试声明不同类型的变量并打印
    fmt.Println("\n练习：声明以下变量并打印")
    fmt.Println("1. 一个int32类型的变量，值为1000")
    fmt.Println("2. 一个complex64类型的复数，值为1+2i")
    fmt.Println("3. 一个rune类型的字符，值为'中'")
}