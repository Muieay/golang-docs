package main

import "fmt"

/*
3-control-flow.go - 控制流程详解
=================================

文件功能：
本文件全面介绍Go语言的控制流程机制，包括条件判断、循环控制和数据遍历。
Go语言提供了简洁而强大的控制结构，支持多种编程范式。

学习目标：
1. 掌握if-else条件判断的各种形式
2. 理解switch语句的高级特性
3. 学会for循环的多种写法
4. 掌握range遍历的适用场景
5. 理解break、continue和标签的使用
6. 能够编写复杂的嵌套控制结构
*/

func main() {
    fmt.Println("=== Go语言控制流程 ===")
    
    /*
    第一节：if条件语句
    
    Go语言的if语句特点：
    1. 条件表达式不需要括号
    2. 支持在if语句中声明局部变量
    3. 左大括号必须与if在同一行
    4. 支持嵌套和复杂条件判断
    */
    fmt.Println("\n--- if条件语句 ---")
    age := 18
    
    // 基本if语句：单条件判断
    // 适用场景：简单的条件执行
    if age >= 18 {
        fmt.Println("成年人")
    }
    
    // if-else语句：二分支选择
    // 适用场景：非此即彼的条件判断
    if age >= 18 {
        fmt.Println("可以投票")
    } else {
        fmt.Println("未成年，不能投票")
    }
    
    // if-else if-else语句：多分支选择
    // 适用场景：多个条件区间的判断
    if age < 13 {
        fmt.Println("儿童")
    } else if age < 20 {
        fmt.Println("青少年")
    } else {
        fmt.Println("成年人")
    }
    
    // if语句中声明变量：作用域限制在if块内
    // 优点：减少变量污染，提高代码可读性
    if score := 85; score >= 90 {     // 变量score只在if-else块中有效
        fmt.Println("优秀")
    } else if score >= 80 {
        fmt.Println("良好")
    } else {
        fmt.Println("需要努力")
    }
    
    /*
    第二节：switch语句
    
    Go的switch语句特点：
    1. 支持多值匹配（逗号分隔）
    2. 支持无表达式switch（替代复杂if-else）
    3. 默认不穿透（不需要break）
    4. 支持fallthrough强制穿透
    5. case表达式可以是任意类型
    */
    fmt.Println("\n--- switch语句 ---")
    day := "Monday"
    
    // 基本switch：多值匹配
    // 适用场景：离散值的等值判断
    switch day {
    case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
        fmt.Println("工作日")
    case "Saturday", "Sunday":
        fmt.Println("周末")
    default:
        fmt.Println("无效的日期")
    }
    
    // switch不带表达式：替代复杂if-else
    // 适用场景：复杂的条件判断，比if-else更清晰
    temperature := 25
    switch {
    case temperature < 0:
        fmt.Println("冰冻")
    case temperature < 20:
        fmt.Println("凉爽")
    case temperature < 30:
        fmt.Println("温暖")
    default:
        fmt.Println("炎热")
    }
    
    // fallthrough用法：强制执行下一个case
    // 注意：fallthrough会忽略下一个case的条件判断
    switch num := 2; num {
    case 1:
        fmt.Println("一")
        fallthrough  // 继续执行下一个case，无论条件是否满足
    case 2:
        fmt.Println("二")
        fallthrough
    case 3:
        fmt.Println("三")
    default:
        fmt.Println("其他")
    }
    
    /*
    第三节：for循环详解
    
    Go只有for循环，但支持多种写法：
    1. 传统三段式for：初始化;条件;后置
    2. while风格for：只有条件判断
    3. 无限循环for：省略所有表达式
    4. for range：遍历容器类型
    */
    fmt.Println("\n--- for循环 ---")
    
    // 传统for循环：三段式结构
    // 适用场景：已知循环次数的场景
    fmt.Println("传统for循环:")
    for i := 0; i < 5; i++ {           // 初始化;条件判断;后置操作
        fmt.Printf("i = %d\n", i)
    }
    
    // while风格for循环：只有条件判断
    // 适用场景：循环次数不确定，根据条件结束
    fmt.Println("while风格的for循环:")
    j := 0
    for j < 3 {                        // 只有条件判断，相当于while
        fmt.Printf("j = %d\n", j)
        j++
    }
    
    // 无限循环：省略所有表达式
    // 适用场景：服务器监听、事件循环等
    fmt.Println("无限循环（执行3次后break）:")
    counter := 0
    for {                              // 无限循环，需要配合break使用
        if counter >= 3 {
            break                      // 主动退出循环
        }
        fmt.Printf("counter = %d\n", counter)
        counter++
    }
    
    /*
    第四节：range循环详解
    
    range用于遍历容器类型：
    1. 数组/切片：返回索引和值
    2. map：返回key和value
    3. 字符串：返回字节索引和字符（rune）
    4. channel：接收值
    
    注意：如果只想要索引或key，使用_忽略值
    */
    fmt.Println("\n--- range循环 ---")
    
    // 遍历数组：固定长度，值类型
    numbers := [5]int{1, 2, 3, 4, 5}
    fmt.Print("数组元素: ")
    for index, value := range numbers {  // index是索引，value是元素值
        fmt.Printf("索引%d=%d ", index, value)
    }
    fmt.Println()
    
    // 遍历切片：动态长度，引用类型
    fruits := []string{"苹果", "香蕉", "橙子"}
    fmt.Print("切片元素: ")
    for index, fruit := range fruits {
        fmt.Printf("%d:%s ", index, fruit)
    }
    fmt.Println()
    
    // 遍历map：键值对无序集合
    person := map[string]string{
        "name": "张三",
        "city": "北京",
        "job": "工程师",
    }
    fmt.Print("map元素: ")
    for key, value := range person {     // 遍历顺序随机
        fmt.Printf("%s:%s ", key, value)
    }
    fmt.Println()
    
    // 遍历字符串：按字符（rune）遍历
    str := "Hello"
    fmt.Print("字符串字符: ")
    for index, char := range str {       // 返回字节索引和Unicode字符
        fmt.Printf("%d:%c ", index, char)
    }
    fmt.Println()
    
    /*
    第五节：循环控制语句
    
    break和continue的作用：
    - break：立即终止整个循环，跳出循环体
    - continue：跳过本次循环剩余代码，开始下一次循环
    
    标签(label)的作用：
    - 标识循环，配合break/continue控制外层循环
    - 提高代码可读性，避免深层嵌套
    */
    fmt.Println("\n--- break和continue ---")
    
    // break示例：完全终止循环
    fmt.Println("break示例:")
    for i := 0; i < 10; i++ {
        if i == 5 {
            break  // 退出整个循环，不再继续
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()
    
    // continue示例：跳过本次循环
    fmt.Println("continue示例:")
    for i := 0; i < 10; i++ {
        if i%2 == 0 {                   // 偶数跳过
            continue                    // 跳过本次循环剩余代码
        }
        fmt.Printf("%d ", i)           // 只打印奇数
    }
    fmt.Println()
    
    /*
    第六节：嵌套循环和标签
    
    嵌套循环的应用场景：
    1. 多维数据处理（矩阵、表格）
    2. 组合问题求解
    3. 图形打印
    
    标签的语法：LabelName: 放在循环前
    */
    fmt.Println("\n--- 嵌套循环和标签 ---")
    
    // 嵌套循环示例：二维坐标遍历
    fmt.Println("嵌套循环:")
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 3; j++ {
            fmt.Printf("(%d,%d) ", i, j)  // 打印坐标对
        }
        fmt.Println()                      // 每行结束换行
    }
    
    // 使用标签跳出外层循环：避免深层嵌套
    fmt.Println("使用标签跳出外层循环:")
OuterLoop:                              // 定义标签，标识外层循环
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i*j >= 2 {
                break OuterLoop          // 跳出标签标识的外层循环
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    fmt.Println()
    
    /*
    第七节：综合实例
    
    通过实际案例综合运用控制流程：
    1. 乘法表：嵌套循环的典型应用
    2. 素数查找：算法与流程控制的结合
    3. 数据处理：结合条件判断和循环遍历
    */
    fmt.Println("\n--- 综合练习 ---")
    
    // 打印乘法表：嵌套循环的经典应用
    fmt.Println("乘法表:")
    for i := 1; i <= 9; i++ {          // 外层循环：控制行数
        for j := 1; j <= i; j++ {       // 内层循环：控制每行的列数
            fmt.Printf("%d×%d=%2d ", j, i, i*j)
        }
        fmt.Println()                    // 每行结束换行
    }
    
    // 查找素数：算法与流程控制的结合
    fmt.Println("1-20之间的素数:")
    for num := 2; num <= 20; num++ {   // 遍历2-20的所有数字
        isPrime := true                // 假设当前数字是素数
        for i := 2; i < num; i++ {     // 检查是否能被2到num-1整除
            if num%i == 0 {            // 如果能整除，不是素数
                isPrime = false
                break                  // 找到因子即可终止内层循环
            }
        }
        if isPrime {                   // 如果仍然是素数，打印出来
            fmt.Printf("%d ", num)
        }
    }
    fmt.Println()
    
    /*
    练习指导：
    1. 计算1-100的和：使用for循环累加
    2. 成绩等级判断：switch语句比if-else更清晰
    3. 切片最大值：range遍历+条件判断
    
    实现思路：
    - 明确循环条件和终止条件
    - 选择合适的控制结构
    - 注意边界条件处理
    */
    fmt.Println("\n练习：")
    fmt.Println("1. 使用for循环计算1到100的和")
    fmt.Println("2. 使用switch语句判断成绩等级")
    fmt.Println("3. 使用range遍历切片并找出最大值")
    
    // 练习答案示例：
    fmt.Println("\n练习答案示例:")
    
    // 1. 计算1-100的和
    sum := 0
    for i := 1; i <= 100; i++ {
        sum += i
    }
    fmt.Printf("1-100的和: %d\n", sum)
    
    // 2. 成绩等级判断
    score := 85
    switch {
    case score >= 90:
        fmt.Println("成绩等级: 优秀")
    case score >= 80:
        fmt.Println("成绩等级: 良好")
    case score >= 70:
        fmt.Println("成绩等级: 中等")
    case score >= 60:
        fmt.Println("成绩等级: 及格")
    default:
        fmt.Println("成绩等级: 不及格")
    }
    
    // 3. 切片最大值查找
    numbers := []int{23, 45, 12, 67, 34, 89, 56}
    max := numbers[0]
    for _, value := range numbers[1:] {
        if value > max {
            max = value
        }
    }
    fmt.Printf("切片最大值: %d\n", max)
}