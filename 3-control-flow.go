package main

import "fmt"

// 控制流程：条件语句和循环
// Go语言的控制结构简洁明了

func main() {
    fmt.Println("=== Go语言控制流程 ===")
    
    // 1. if条件语句
    age := 18
    
    // 基本if语句
    if age >= 18 {
        fmt.Println("你已成年")
    }
    
    // if-else语句
    if age < 18 {
        fmt.Println("未成年")
    } else {
        fmt.Println("已成年")
    }
    
    // if-else if-else语句
    if age < 12 {
        fmt.Println("儿童")
    } else if age < 18 {
        fmt.Println("青少年")
    } else if age < 60 {
        fmt.Println("成年人")
    } else {
        fmt.Println("老年人")
    }
    
    // if语句中可以声明变量（作用域仅限于if块）
    if score := 85; score >= 90 {
        fmt.Println("优秀")
    } else if score >= 80 {
        fmt.Println("良好")
    } else if score >= 60 {
        fmt.Println("及格")
    } else {
        fmt.Println("不及格")
    }
    
    // 2. switch语句
    day := "Monday"
    
    switch day {
    case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
        fmt.Println("工作日")
    case "Saturday", "Sunday":
        fmt.Println("周末")
    default:
        fmt.Println("无效的日期")
    }
    
    // switch的表达式可以省略（相当于if-else链）
    temperature := 25
    switch {
    case temperature < 0:
        fmt.Println("寒冷")
    case temperature < 20:
        fmt.Println("凉爽")
    case temperature < 30:
        fmt.Println("温暖")
    default:
        fmt.Println("炎热")
    }
    
    // 3. for循环（Go只有for循环，没有其他循环关键字）
    
    // 传统for循环
    fmt.Println("\n--- 传统for循环 ---")
    for i := 0; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }
    
    // 省略初始化和后置语句
    fmt.Println("\n--- 省略初始化和后置 ---")
    j := 0
    for ; j < 3; j++ {
        fmt.Printf("j = %d\n", j)
    }
    
    // 省略所有语句（相当于while循环）
    fmt.Println("\n--- 相当于while循环 ---")
    k := 0
    for k < 3 {
        fmt.Printf("k = %d\n", k)
        k++
    }
    
    // 无限循环
    fmt.Println("\n--- 无限循环（使用break退出） ---")
    counter := 0
    for {
        if counter >= 3 {
            break
        }
        fmt.Printf("counter = %d\n", counter)
        counter++
    }
    
    // 4. range循环（遍历数组、切片、map、字符串）
    fmt.Println("\n--- range循环 ---")
    
    // 遍历字符串
    str := "Hello"
    for index, char := range str {
        fmt.Printf("索引: %d, 字符: %c\n", index, char)
    }
    
    // 遍历数组
    numbers := [5]int{1, 2, 3, 4, 5}
    for index, value := range numbers {
        fmt.Printf("索引: %d, 值: %d\n", index, value)
    }
    
    // 只获取值（使用下划线忽略索引）
    for _, value := range numbers {
        fmt.Printf("值: %d\n", value)
    }
    
    // 5. break和continue
    fmt.Println("\n--- break和continue ---")
    
    // break示例
    fmt.Println("break示例：")
    for i := 0; i < 10; i++ {
        if i == 5 {
            break  // 完全退出循环
        }
        fmt.Printf("i = %d\n", i)
    }
    
    // continue示例
    fmt.Println("continue示例：")
    for i := 0; i < 5; i++ {
        if i == 2 {
            continue  // 跳过本次循环
        }
        fmt.Printf("i = %d\n", i)
    }
    
    // 6. 嵌套循环和标签
    fmt.Println("\n--- 嵌套循环和标签 ---")
    
    // 使用标签跳出多层循环
    OuterLoop:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i*j >= 2 {
                break OuterLoop  // 跳出外层循环
            }
            fmt.Printf("i = %d, j = %d\n", i, j)
        }
    }
    
    // 练习：
    fmt.Println("\n练习：")
    fmt.Println("1. 使用for循环打印1到100的偶数")
    fmt.Println("2. 使用switch判断成绩等级（A:90-100, B:80-89, C:70-79, D:60-69, F:0-59）")
    fmt.Println("3. 使用range计算数组[1,2,3,4,5]的和")
}