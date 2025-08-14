package main

import (
    "fmt"
    "reflect"
)

/*
2-variables-and-types.go - 变量与数据类型详解
==========================================

文件功能：
本文件全面介绍Go语言的变量声明方式、基本数据类型、类型转换机制以及零值概念。
Go是静态类型语言，每个变量在编译时都必须有明确的类型。

学习目标：
1. 掌握Go语言的三种变量声明方式
2. 理解基本数据类型的范围和用途
3. 学会类型转换的正确方法
4. 了解零值的概念和重要性
5. 使用反射查看变量类型信息
*/

func main() {
    fmt.Println("=== Go语言变量和数据类型 ===")
    
    /*
    第一节：变量声明方式
    
    Go语言提供了三种变量声明方式，各有适用场景：
    1. var关键字声明：最传统的声明方式，适用于包级变量和需要明确类型的场景
    2. 类型推断：使用var但不指定类型，编译器根据初始值自动推断
    3. 简短声明(:=)：最简洁的声明方式，只能在函数内部使用
    */
    
    // 方式1：var关键字声明（完整形式）
    // 适用场景：需要明确指定类型、包级变量、零值初始化
    var age int = 25                    // 显式声明int类型变量
    var name string = "张三"            // 显式声明string类型变量
    var isStudent bool = true           // 显式声明bool类型变量
    fmt.Printf("年龄: %d, 姓名: %s, 是否学生: %t\n", age, name, isStudent)
    
    // 方式2：类型推断（省略类型）
    // 适用场景：初始值类型明确，希望代码更简洁
    var score = 95.5                    // 自动推断为float64类型（默认浮点类型）
    var city = "北京"                    // 自动推断为string类型
    fmt.Printf("分数类型: %T, 城市类型: %T\n", score, city)
    
    // 方式3：简短声明（最常用）
    // 适用场景：函数内部变量声明，代码简洁明了
    country := "中国"                   // 简短声明并初始化字符串变量
    year := 2024                        // 简短声明并初始化整型变量
    fmt.Printf("国家: %s, 年份: %d\n", country, year)
    
    /*
    第二节：基本数据类型详解
    
    Go语言提供了丰富的基本数据类型，每种类型都有明确的取值范围和内存占用
    */
    
    // 整数类型家族
    // 特点：不同位数表示不同范围的整数，可根据需求选择合适类型以节省内存
    var intNum int = 100                // int类型：32位系统为int32，64位系统为int64
    var int8Num int8 = 127              // int8类型：-128 到 127（8位有符号整数）
    var int16Num int16 = 32767          // int16类型：-32768 到 32767（16位有符号整数）
    var uintNum uint = 200              // uint类型：无符号整数，只能表示正数
    var byteNum byte = 255              // byte类型：uint8的别名，常用于处理二进制数据
    
    fmt.Printf("int: %d, int8: %d, uint: %d, byte: %d\n", 
        intNum, int8Num, uintNum, byteNum)
    
    // 浮点类型
    // 特点：表示小数，float64精度更高，是默认的浮点类型
    var float32Num float32 = 3.14       // float32：单精度浮点数，约6-7位有效数字
    var float64Num float64 = 3.141592653589793  // float64：双精度浮点数，约15-17位有效数字
    fmt.Printf("float32: %.2f, float64: %.15f\n", float32Num, float64Num)
    
    // 字符串类型
    // 特点：UTF-8编码，支持两种字面量形式
    str1 := "Hello"                     // 双引号字符串：支持转义字符
    str2 := `Raw string                 // 反引号字符串：原始字符串，支持换行
    可以换行
    支持"引号"`                        // 反引号内内容原样保留，包括换行和引号
    fmt.Printf("str1: %s\nstr2: %s\n", str1, str2)
    
    // 布尔类型
    // 特点：只有true和false两个值，零值为false
    var bool1 bool = true               // 显式赋值为true
    var bool2 bool = false              // 显式赋值为false
    var bool3 bool                      // 未赋值的bool变量，默认为false（零值）
    fmt.Printf("bool1: %t, bool2: %t, bool3: %t\n", bool1, bool2, bool3)
    
    /*
    第三节：零值（默认值）概念
    
    零值是Go语言的重要特性：
    - 所有变量声明后都有默认值，不存在未初始化的状态
    - 数值类型零值为0
    - 字符串类型零值为空字符串""
    - 布尔类型零值为false
    - 引用类型零值为nil
    */
    var zeroInt int                     // 整型零值：0
    var zeroString string               // 字符串零值：""
    var zeroBool bool                   // 布尔零值：false
    var zeroFloat float64               // 浮点零值：0.0
    fmt.Printf("零值 - int: %d, string: %q, bool: %t, float: %f\n", 
        zeroInt, zeroString, zeroBool, zeroFloat)
    
    /*
    第四节：类型转换机制
    
    Go是强类型语言，不同类型之间必须显式转换
    转换语法：目标类型(表达式)
    注意：高精度到低精度转换可能丢失数据
    */
    var a int = 100                     // 原始int值
    var b float64 = float64(a)          // 显式转换为float64，可能增加小数位
    var c int32 = int32(a)              // 显式转换为int32，范围可能缩小
    fmt.Printf("类型转换: %d -> %.1f -> %d\n", a, b, c)
    
    // 类型转换示例：不同整数类型之间的转换
    var big int64 = 1000000
    var small int8 = int8(big)          // 注意：可能溢出，1000000超出int8范围
    fmt.Printf("大数转小数: %d -> %d (可能溢出)\n", big, small)
    
    /*
    第五节：使用reflect包查看类型信息
    
    reflect包提供了运行时反射能力，可以动态获取变量类型信息
    常用于：
    - 调试和日志记录
    - 通用函数编写
    - 序列化和反序列化
    */
    fmt.Printf("reflect - age的类型: %v\n", reflect.TypeOf(age))
    fmt.Printf("reflect - name的类型: %v, 值: %v\n", reflect.TypeOf(name), reflect.ValueOf(name))
    
    /*
    第六节：其他重要数据类型
    
    补充介绍几个常用的特殊类型
    */
    
    // rune类型：int32的别名，用于表示Unicode字符
    var chineseChar rune = '中'          // rune类型，表示Unicode码点
    fmt.Printf("rune字符: %c, Unicode码点: %d\n", chineseChar, chineseChar)
    
    // complex64/complex128：复数类型
    var complexNum complex64 = 1 + 2i   // complex64：32位实部+32位虚部
    var complexNum128 complex128 = complex(3, 4)  // complex128：64位实部+64位虚部
    fmt.Printf("复数: %v, 模: %.2f\n", complexNum128, real(complexNum128))
    
    /*
    练习指导：
    1. int32类型变量：适用于需要明确32位整数的场景
    2. complex64类型：适用于复数计算，如信号处理
    3. rune类型：适用于处理Unicode字符和文本处理
    */
    fmt.Println("\n练习：声明以下变量并打印")
    fmt.Println("1. 一个int32类型的变量，值为1000")
    fmt.Println("2. 一个complex64类型的复数，值为1+2i")
    fmt.Println("3. 一个rune类型的字符，值为'中'")
    
    // 练习答案示例：
    var int32Var int32 = 1000
    var complexVar complex64 = 1 + 2i
    var runeVar rune = '中'
    
    fmt.Printf("\n练习答案:\n")
    fmt.Printf("int32变量: %d (类型: %T)\n", int32Var, int32Var)
    fmt.Printf("复数变量: %v (类型: %T)\n", complexVar, complexVar)
    fmt.Printf("rune字符: %c (码点: %d, 类型: %T)\n", runeVar, runeVar, runeVar)
}