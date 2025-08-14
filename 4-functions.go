package main

import (
    "fmt"
    "errors"
)

// 函数详解
// Go语言的函数是一等公民，支持多返回值、命名返回值、可变参数等特性

func main() {
    fmt.Println("=== Go语言函数 ===")
    
    // 1. 基本函数调用
    result := add(3, 5)
    fmt.Printf("3 + 5 = %d\n", result)
    
    // 2. 多返回值函数
    quotient, remainder := divide(10, 3)
    fmt.Printf("10 ÷ 3 = 商: %d, 余数: %d\n", quotient, remainder)
    
    // 3. 命名返回值
    area, perimeter := rectangleInfo(4, 6)
    fmt.Printf("矩形 - 面积: %.1f, 周长: %.1f\n", area, perimeter)
    
    // 4. 错误处理
    if result, err := safeDivide(10, 0); err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("安全除法结果: %.2f\n", result)
    }
    
    // 5. 可变参数函数
    fmt.Printf("平均值: %.2f\n", average(1, 2, 3, 4, 5))
    fmt.Printf("平均值: %.2f\n", average(10, 20, 30))
    
    // 6. 匿名函数
    greet := func(name string) {
        fmt.Printf("你好, %s!\n", name)
    }
    greet("张三")
    
    // 7. 闭包示例
    counter := makeCounter()
    fmt.Printf("计数器: %d\n", counter())
    fmt.Printf("计数器: %d\n", counter())
    fmt.Printf("计数器: %d\n", counter())
    
    // 8. 函数作为参数
    numbers := []int{1, 2, 3, 4, 5}
    doubled := mapNumbers(numbers, func(n int) int {
        return n * 2
    })
    fmt.Printf("原数组: %v, 翻倍后: %v\n", numbers, doubled)
    
    // 9. 递归函数
    fmt.Printf("5! = %d\n", factorial(5))
    fmt.Printf("斐波那契第10项: %d\n", fibonacci(10))
    
    // 10. 延迟执行（defer）
    deferExample()
    
    // 11. 高阶函数示例
    resultFunc := createMultiplier(3)
    fmt.Printf("3 × 4 = %d\n", resultFunc(4))
    
    // 练习
    fmt.Println("\n练习：")
    fmt.Println("1. 编写一个计算圆面积和周长的函数")
    fmt.Println("2. 编写一个判断素数的函数")
    fmt.Println("3. 编写一个字符串反转函数")
}

// 1. 基本函数定义
func add(a int, b int) int {
    return a + b
}

// 参数类型相同的可以简写
func add2(a, b int) int {
    return a + b
}

// 2. 多返回值函数
func divide(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

// 3. 命名返回值
func rectangleInfo(width, height float64) (area float64, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // 可以直接return，返回值已经在函数签名中定义
}

// 4. 错误处理
func safeDivide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 5. 可变参数函数
func average(numbers ...float64) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    sum := 0.0
    for _, num := range numbers {
        sum += num
    }
    return sum / float64(len(numbers))
}

// 6. 闭包示例
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 7. 函数作为参数
func mapNumbers(numbers []int, mapper func(int) int) []int {
    result := make([]int, len(numbers))
    for i, num := range numbers {
        result[i] = mapper(num)
    }
    return result
}

// 8. 递归函数
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// 9. defer示例
deferExample() {
    fmt.Println("开始执行")
    defer fmt.Println("延迟执行1")
    defer fmt.Println("延迟执行2")  // 后注册的defer先执行
    fmt.Println("正常执行")
}

// 10. 高阶函数
func createMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}