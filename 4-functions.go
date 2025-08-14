package main

import (
	"errors"
	"fmt"
)

// 函数详解
// Go语言的函数是一等公民，支持多返回值、命名返回值、可变参数等特性
// 本示例涵盖了Go函数的核心特性，适合初学者理解函数的各种用法

func main() {
	fmt.Println("=== Go语言函数 ===")
	
	// 1. 基本函数调用
	// 调用无状态的基础函数，传入参数并接收返回值
	result := add(3, 5)
	fmt.Printf("3 + 5 = %d\n", result)
	
	// 2. 多返回值函数
	// Go支持函数返回多个值，无需像其他语言那样使用数组或对象包装
	quotient, remainder := divide(10, 3)
	fmt.Printf("10 ÷ 3 = 商: %d, 余数: %d\n", quotient, remainder)
	
	// 3. 命名返回值
	// 函数定义时指定返回值名称，可简化return语句并提高代码可读性
	area, perimeter := rectangleInfo(4, 6)
	fmt.Printf("矩形 - 面积: %.1f, 周长: %.1f\n", area, perimeter)
	
	// 4. 错误处理
	// Go通过多返回值实现错误处理，通常最后一个返回值为error类型
	// 调用时先判断错误是否为nil，再处理结果
	if result, err := safeDivide(10, 0); err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("安全除法结果: %.2f\n", result)
	}
	
	// 5. 可变参数函数
	// 可接收任意数量的同类型参数，适合不确定参数数量的场景
	fmt.Printf("平均值: %.2f\n", average(1, 2, 3, 4, 5))
	fmt.Printf("平均值: %.2f\n", average(10, 20, 30))  // 支持不同数量的参数
	
	// 6. 匿名函数
	// 没有名称的函数，可直接定义并赋值给变量，适合临时使用的简单逻辑
	greet := func(name string) {
		fmt.Printf("你好, %s!\n", name)
	}
	greet("张三")  // 像调用普通函数一样使用
	
	// 7. 闭包示例
	// 闭包是引用了外部变量的匿名函数，能"捕获"并维护外部状态
	counter := makeCounter()
	fmt.Printf("计数器: %d\n", counter())  // 每次调用都能保持状态
	fmt.Printf("计数器: %d\n", counter())
	fmt.Printf("计数器: %d\n", counter())
	
	// 8. 函数作为参数
	// Go中函数是一等公民，可作为参数传递，实现灵活的逻辑注入
	numbers := []int{1, 2, 3, 4, 5}
	doubled := mapNumbers(numbers, func(n int) int {  // 匿名函数作为参数
		return n * 2
	})
	fmt.Printf("原数组: %v, 翻倍后: %v\n", numbers, doubled)
	
	// 9. 递归函数
	// 函数调用自身的实现方式，适合分治思想的问题（如阶乘、斐波那契数列）
	fmt.Printf("5! = %d\n", factorial(5))
	fmt.Printf("斐波那契第10项: %d\n", fibonacci(10))
	
	// 10. 延迟执行（defer）
	// 延迟执行语句，在函数退出前执行，常用于资源释放（文件、数据库连接等）
	deferExample()
	
	// 11. 高阶函数示例
	// 返回函数的函数，可动态生成具有特定行为的函数
	resultFunc := createMultiplier(3)  // 生成一个"乘以3"的函数
	fmt.Printf("3 × 4 = %d\n", resultFunc(4))  // 使用生成的函数
	
	// 练习
	fmt.Println("\n练习：")
	fmt.Println("1. 编写一个计算圆面积和周长的函数")
	fmt.Println("2. 编写一个判断素数的函数")
	fmt.Println("3. 编写一个字符串反转函数")
}

// 1. 基本函数定义
// 功能：计算两个整数的和
// 参数：
//   a: 第一个加数（int类型）
//   b: 第二个加数（int类型）
// 返回值：两个数的和（int类型）
func add(a int, b int) int {
	return a + b
}

// 参数类型相同的简写形式
// 当多个参数类型相同时，可只在最后一个参数后指定类型
func add2(a, b int) int {  // 等价于add(a int, b int)，语法糖简化代码
	return a + b
}

// 2. 多返回值函数
// 功能：计算两个整数的商和余数
// 参数：
//   a: 被除数（int类型）
//   b: 除数（int类型，注意未处理b=0的情况，仅作示例）
// 返回值：
//   第一个返回值：商（int类型）
//   第二个返回值：余数（int类型）
// 说明：Go的多返回值避免了使用结构体或数组包装多个返回结果的麻烦
func divide(a, b int) (int, int) {
	quotient := a / b   // 整数除法求商
	remainder := a % b  // 取余运算求余数
	return quotient, remainder  // 用逗号分隔多个返回值
}

// 3. 命名返回值
// 功能：计算矩形的面积和周长
// 参数：
//   width: 矩形宽度（float64类型）
//   height: 矩形高度（float64类型）
// 返回值：
//   area: 面积（width * height）
//   perimeter: 周长（2 * (width + height)）
// 说明：在函数签名中指定返回值名称，函数内部可直接使用这些变量
//      return语句可省略返回值列表，自动返回命名变量的值
func rectangleInfo(width, height float64) (area float64, perimeter float64) {
	area = width * height  // 直接给命名返回值赋值
	perimeter = 2 * (width + height)
	return  // 无需指定返回值，自动返回area和perimeter
}

// 4. 错误处理
// 功能：安全除法，处理除数为0的情况
// 参数：
//   a: 被除数（float64类型）
//   b: 除数（float64类型）
// 返回值：
//   第一个返回值：除法结果（float64类型）
//   第二个返回值：错误信息（error类型，nil表示无错误）
// 说明：Go采用"显式错误处理"模式，通过多返回值返回错误，
//      而非其他语言的异常机制，更强调错误处理的显式性
func safeDivide(a, b float64) (float64, error) {
	if b == 0 {
		// 当除数为0时，返回错误信息（使用errors.New创建错误）
		return 0, errors.New("除数不能为零")
	}
	// 正常情况下，返回结果和nil（表示无错误）
	return a / b, nil
}

// 5. 可变参数函数
// 功能：计算任意数量浮点数的平均值
// 参数：
//   numbers: 可变参数（...float64表示接收0个或多个float64参数）
// 返回值：平均值（float64类型）
// 说明：可变参数在函数内部会被转换为对应类型的切片（[]float64）
//      适合参数数量不确定的场景（如求和、平均值计算等）
func average(numbers ...float64) float64 {
	if len(numbers) == 0 {  // 处理空参数情况
		return 0
	}
	
	sum := 0.0
	// 遍历可变参数（内部作为切片处理）
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))  // 计算并返回平均值
}

// 6. 闭包示例
// 功能：创建一个计数器函数
// 返回值：一个无参数、返回int的函数
// 说明：闭包是能够捕获外部变量的函数，即使外部函数执行完毕，
//      捕获的变量仍能被内部函数访问和修改，实现状态保持
func makeCounter() func() int {
	count := 0  // 被闭包捕获的外部变量，生命周期被延长
	
	// 返回一个匿名函数，该函数引用了外部的count变量
	return func() int {
		count++  // 每次调用计数器加1
		return count
	}
}

// 7. 函数作为参数
// 功能：对切片中的每个元素应用映射函数，返回新切片
// 参数：
//   numbers: 输入的整数切片（[]int）
//   mapper: 映射函数（接收int，返回int），用于转换每个元素
// 返回值：转换后的新切片（[]int）
// 说明：函数作为参数传递，实现了"策略模式"，可灵活替换处理逻辑
//      这是函数式编程的核心特性之一
func mapNumbers(numbers []int, mapper func(int) int) []int {
	result := make([]int, len(numbers))  // 创建同长度的结果切片
	// 遍历每个元素，应用映射函数
	for i, num := range numbers {
		result[i] = mapper(num)  // 调用传入的函数处理元素
	}
	return result
}

// 8. 递归函数：阶乘
// 功能：计算n的阶乘（n! = n × (n-1) × ... × 1）
// 参数：n: 非负整数（int类型）
// 返回值：n的阶乘（int类型）
// 说明：递归的核心是"自己调用自己"，需满足两个条件：
//      1. 基准 case（终止条件）：n <= 1 时返回1
//      2. 递归 case：将问题分解为更小的子问题（n * factorial(n-1)）
func factorial(n int) int {
	if n <= 1 {  // 基准 case：终止递归的条件
		return 1
	}
	// 递归 case：调用自身处理更小的问题
	return n * factorial(n-1)
}

// 8. 递归函数：斐波那契数列
// 功能：计算斐波那契数列的第n项（f(n) = f(n-1) + f(n-2)）
// 参数：n: 正整数（int类型）
// 返回值：第n项的数值（int类型）
// 说明：斐波那契数列定义：f(0)=0, f(1)=1, f(n)=f(n-1)+f(n-2)
//      注意：此实现为演示用，性能较差，实际应用需优化（如备忘录法）
func fibonacci(n int) int {
	if n <= 1 {  // 基准 case：前两项直接返回
		return n
	}
	// 递归 case：分解为两个更小的子问题
	return fibonacci(n-1) + fibonacci(n-2)
}

// 9. defer示例
// 功能：演示defer语句的执行时机和顺序
// 说明：defer语句会将其后的函数调用延迟到当前函数退出前执行
//      多个defer按"后进先出"（LIFO）顺序执行，适合资源释放操作
func deferExample() {
	fmt.Println("开始执行")
	
	// 第一个defer：会在函数结束前最后执行
	defer fmt.Println("延迟执行1")
	
	// 第二个defer：会在函数结束前先于"延迟执行1"执行
	defer fmt.Println("延迟执行2")  // 后注册的defer先执行
	
	fmt.Println("正常执行")  // 此语句先于所有defer执行
}

// 10. 高阶函数
// 功能：创建一个"乘以指定因子"的函数
// 参数：factor: 乘数因子（int类型）
// 返回值：一个接收int参数并返回int的函数
// 说明：高阶函数是指能接收函数作为参数或返回函数的函数
//      此例通过闭包捕获factor变量，动态生成特定功能的函数
func createMultiplier(factor int) func(int) int {
	// 返回一个匿名函数，该函数使用外部的factor变量
	return func(x int) int {
		return x * factor  // 用捕获的factor作为乘数
	}
}
