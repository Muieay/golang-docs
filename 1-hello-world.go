package main

import "fmt"

// 第一个Go程序：Hello World
// 这是每个程序员学习新语言的第一步

/*
Go程序的基本结构：
1. package main - 定义包名，main包是程序的入口
2. import "fmt" - 导入格式化I/O的包
3. func main() - 主函数，程序执行的入口点
*/

func main() {
    // Println打印字符串并换行
    fmt.Println("Hello, World!")
    
    // 打印中文
    fmt.Println("你好，世界！")
    
    // 格式化打印
    name := "Go语言"
    version := "1.21"
    fmt.Printf("欢迎使用%s，版本%s\n", name, version)
    
    // 练习：尝试修改上面的字符串，看看输出结果
    fmt.Println("尝试修改代码，重新运行看看效果！")
}