package main

import "fmt"

// 数组、切片和Map
// Go语言的核心数据结构

func main() {
    fmt.Println("=== Go语言数据结构 ===")
    
    // 1. 数组（固定长度）
    fmt.Println("\n--- 数组 ---")
    
    // 声明和初始化数组
    var arr1 [5]int
    arr2 := [5]int{1, 2, 3, 4, 5}
    arr3 := [...]int{1, 2, 3}  // 自动推断长度
    arr4 := [5]int{1, 2, 3}    // 未赋值的元素为零值
    
    fmt.Printf("arr1: %v, 长度: %d\n", arr1, len(arr1))
    fmt.Printf("arr2: %v, 长度: %d\n", arr2, len(arr2))
    fmt.Printf("arr3: %v, 长度: %d\n", arr3, len(arr3))
    fmt.Printf("arr4: %v\n", arr4)
    
    // 访问和修改数组元素
    arr2[0] = 100
    fmt.Printf("修改后的arr2: %v\n", arr2)
    
    // 遍历数组
    for i, v := range arr2 {
        fmt.Printf("索引: %d, 值: %d\n", i, v)
    }
    
    // 2. 切片（动态数组）
    fmt.Println("\n--- 切片 ---")
    
    // 创建切片
    var slice1 []int
    slice2 := []int{1, 2, 3, 4, 5}
    slice3 := make([]int, 5)      // 创建长度为5的切片
    slice4 := make([]int, 3, 10)  // 长度为3，容量为10
    
    fmt.Printf("slice1: %v, 长度: %d, 容量: %d\n", slice1, len(slice1), cap(slice1))
    fmt.Printf("slice2: %v, 长度: %d, 容量: %d\n", slice2, len(slice2), cap(slice2))
    fmt.Printf("slice3: %v, 长度: %d, 容量: %d\n", slice3, len(slice3), cap(slice3))
    fmt.Printf("slice4: %v, 长度: %d, 容量: %d\n", slice4, len(slice4), cap(slice4))
    
    // 切片操作
    fmt.Println("\n--- 切片操作 ---")
    
    // 追加元素
    slice1 = append(slice1, 1, 2, 3)
    fmt.Printf("追加后slice1: %v\n", slice1)
    
    // 追加切片
    slice1 = append(slice1, slice2...)
    fmt.Printf("追加切片后slice1: %v\n", slice1)
    
    // 切片切割
    slice5 := slice2[1:4]   // 从索引1到4（不包括4）
    slice6 := slice2[:3]    // 从开头到索引3
    slice7 := slice2[2:]    // 从索引2到结尾
    fmt.Printf("slice2: %v\n", slice2)
    fmt.Printf("slice2[1:4]: %v\n", slice5)
    fmt.Printf("slice2[:3]: %v\n", slice6)
    fmt.Printf("slice2[2:]: %v\n", slice7)
    
    // 复制切片
    slice8 := make([]int, len(slice2))
    copy(slice8, slice2)
    fmt.Printf("复制后的slice8: %v\n", slice8)
    
    // 3. 多维数组和切片
    fmt.Println("\n--- 多维数组 ---")
    
    // 二维数组
    var matrix1 [3][3]int
    matrix2 := [2][3]int{{1, 2, 3}, {4, 5, 6}}
    fmt.Printf("matrix1: %v\n", matrix1)
    fmt.Printf("matrix2: %v\n", matrix2)
    
    // 二维切片
    var matrix3 [][]int
    matrix3 = append(matrix3, []int{1, 2, 3})
    matrix3 = append(matrix3, []int{4, 5, 6})
    fmt.Printf("matrix3: %v\n", matrix3)
    
    // 4. Map（映射）
    fmt.Println("\n--- Map ---")
    
    // 创建map
    var map1 map[string]int
    map1 = make(map[string]int)
    
    map2 := map[string]int{
        "apple":  5,
        "banana": 3,
        "orange": 8,
    }
    
    map3 := make(map[string]int)
    
    fmt.Printf("map1: %v\n", map1)
    fmt.Printf("map2: %v\n", map2)
    fmt.Printf("map3: %v\n", map3)
    
    // map操作
    fmt.Println("\n--- Map操作 ---")
    
    // 添加和修改
    map3["apple"] = 10
    map3["grape"] = 15
    fmt.Printf("添加元素后map3: %v\n", map3)
    
    // 访问元素
    value := map2["apple"]
    fmt.Printf("apple的数量: %d\n", value)
    
    // 检查键是否存在
    if value, exists := map2["pear"]; exists {
        fmt.Printf("pear的数量: %d\n", value)
    } else {
        fmt.Println("pear不存在")
    }
    
    // 删除元素
    delete(map2, "banana")
    fmt.Printf("删除banana后map2: %v\n", map2)
    
    // 遍历map
    fmt.Println("遍历map:")
    for key, value := range map2 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 5. 复杂数据结构
    fmt.Println("\n--- 复杂数据结构 ---")
    
    // 切片中包含map
    students := []map[string]interface{}{
        {"name": "张三", "age": 20, "scores": []int{90, 85, 88}},
        {"name": "李四", "age": 22, "scores": []int{78, 92, 85}},
    }
    
    fmt.Printf("学生信息: %v\n", students)
    
    // map中包含切片
    class := map[string][]string{
        "math":    {"张三", "李四"},
        "english": {"王五", "赵六", "钱七"},
    }
    fmt.Printf("班级: %v\n", class)
    
    // 6. 实际应用示例
    fmt.Println("\n--- 实际应用示例 ---")
    
    // 统计单词频率
    text := "hello world hello go world"
    wordCount := wordFrequency(text)
    fmt.Printf("文本: %s\n", text)
    fmt.Printf("单词频率: %v\n", wordCount)
    
    // 练习
    fmt.Println("\n练习：")
    fmt.Println("1. 创建一个学生成绩表，使用map[string]float64")
    fmt.Println("2. 实现一个函数，找出切片中的最大值和最小值")
    fmt.Println("3. 创建一个二维切片表示棋盘，初始化8x8的棋盘")
}

// 统计单词频率的函数
func wordFrequency(text string) map[string]int {
    words := []string{}
    word := ""
    
    for _, char := range text {
        if char == ' ' {
            if word != "" {
                words = append(words, word)
                word = ""
            }
        } else {
            word += string(char)
        }
    }
    
    if word != "" {
        words = append(words, word)
    }
    
    frequency := make(map[string]int)
    for _, w := range words {
        frequency[w]++
    }
    
    return frequency
}