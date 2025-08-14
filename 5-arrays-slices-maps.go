package main

import "fmt"

/*
5-arrays-slices-maps.go - 复合数据结构详解
========================================

文件功能：
本文件全面介绍Go语言中的复合数据结构：数组(Array)、切片(Slice)和映射(Map)。
这些数据结构是构建复杂应用程序的基础，理解其特性和使用场景至关重要。

学习目标：
1. 理解数组的固定长度特性
2. 掌握切片的动态扩容机制
3. 学会Map的键值对存储
4. 理解切片和数组的关系
5. 掌握多维数据结构的构建
6. 学会实际应用中的数据处理技巧
*/

func main() {
    fmt.Println("=== Go语言数组、切片和Map ===")
    
    /*
    第一节：数组(Array)详解
    
    数组特点：
    1. 固定长度，长度是类型的一部分
    2. 值类型，赋值会复制整个数组
    3. 内存连续分配，访问效率高
    4. 编译时确定大小
    
    适用场景：
    - 长度已知的集合
    - 需要固定大小的缓冲区
    - 性能要求高的场景
    */
    fmt.Println("\n--- 数组 ---")
    
    // 声明和初始化数组的多种方式
    var arr1 [5]int                    // 声明未初始化的数组，元素为零值
    arr2 := [5]int{1, 2, 3, 4, 5}      // 字面量初始化
    arr3 := [...]int{1, 2, 3}  // 自动推断长度
    arr4 := [5]int{1, 2, 3}    // 未赋值的元素为零值
    
    fmt.Printf("arr1: %v, 长度: %d\n", arr1, len(arr1))
    fmt.Printf("arr2: %v, 长度: %d\n", arr2, len(arr2))
    fmt.Printf("arr3: %v, 长度: %d\n", arr3, len(arr3))
    fmt.Printf("arr4: %v\n", arr4)
    
    // 访问和修改数组元素
    arr2[0] = 100
    fmt.Printf("修改后的arr2: %v\n", arr2)
    
    // 遍历数组的两种方式
    fmt.Println("数组遍历:")
    for i := 0; i < len(arr2); i++ {  // 传统索引遍历
        fmt.Printf("arr2[%d] = %d\n", i, arr2[i])
    }
    
    for i, v := range arr2 {            // range遍历，获取索引和值
        fmt.Printf("索引%d: 值%d\n", i, v)
    }
    
    /*
    第二节：切片(Slice)详解
    
    切片特点：
    1. 动态长度，可自动扩容
    2. 引用类型，底层共享数组
    3. 包含指针、长度、容量三个属性
    4. 零值为nil，但长度为0的切片不等同于nil
    
    内部结构：
    - 指针：指向底层数组
    - 长度：当前元素个数
    - 容量：底层数组的最大容量
    */
    fmt.Println("\n--- 切片 ---")
    
    // 创建切片的多种方式
    var slice1 []int                    // nil切片，未分配底层数组
    slice2 := []int{1, 2, 3, 4, 5}      // 字面量创建
    slice3 := make([]int, 5)      // 创建长度为5的切片
    slice4 := make([]int, 3, 10)  // 长度为3，容量为10
    
    fmt.Printf("slice1: %v, 长度: %d, 容量: %d\n", slice1, len(slice1), cap(slice1))
    fmt.Printf("slice2: %v, 长度: %d, 容量: %d\n", slice2, len(slice2), cap(slice2))
    fmt.Printf("slice3: %v, 长度: %d, 容量: %d\n", slice3, len(slice3), cap(slice3))
    fmt.Printf("slice4: %v, 长度: %d, 容量: %d\n", slice4, len(slice4), cap(slice4))
    
    // 切片操作详解
    fmt.Println("\n切片操作:")
    
    // append操作：动态扩容
    slice1 = append(slice1, 10, 20, 30)  // 追加多个元素
    fmt.Printf("append后slice1: %v (len=%d, cap=%d)\n", slice1, len(slice1), cap(slice1))
    
    // append切片：解包操作
    slice1 = append(slice1, slice2...)    // 使用...解包切片
    fmt.Printf("合并切片: %v\n", slice1)
    
    // 切片截取：左闭右开区间
    slice5 := slice2[1:4]    // 从索引1到3（不包含4）
    slice6 := slice2[:3]    // 从开头到索引2
    slice7 := slice2[2:]    // 从索引2到结尾
    slice8 := slice2[:]     // 整个切片
    
    fmt.Printf("slice2[1:4] = %v\n", slice5)
    fmt.Printf("slice2[:3] = %v\n", slice6)
    fmt.Printf("slice2[2:] = %v\n", slice7)
    fmt.Printf("slice2[:] = %v\n", slice8)
    
    // 切片共享底层数组的陷阱
    fmt.Println("\n切片共享底层数组:")
    original := []int{1, 2, 3, 4, 5}
    shared := original              // 共享底层数组
    shared[0] = 100                 // 修改会影响原切片
    fmt.Printf("original: %v (被修改)\n", original)
    fmt.Printf("shared: %v\n", shared)
    
    // 使用copy避免共享：创建独立副本
    independent := make([]int, len(original))
    copy(independent, original)     // 深拷贝
    independent[0] = 999
    fmt.Printf("original保持不变: %v\n", original)
    fmt.Printf("independent: %v\n", independent)
    
    /*
    第三节：Map（映射）详解
    
    Map特点：
    1. 键值对存储，无序集合
    2. 引用类型，零值为nil
    3. 键必须可比较（支持==操作）
    4. 动态扩容，自动处理哈希冲突
    
    适用场景：
    - 快速查找
    - 数据分组
    - 配置存储
    */
    fmt.Println("\n--- Map ---")
    
    // 创建Map的多种方式
    var map1 map[string]int           // nil map，不能直接使用
    map2 := make(map[string]int)       // make创建，可立即使用
    map3 := map[string]int{           // 字面量创建并初始化
        "apple":  5,
        "banana": 3,
        "orange": 8,
    }
    
    fmt.Printf("map1: %v (nil: %t)\n", map1, map1 == nil)
    fmt.Printf("map2: %v (空map)\n", map2)
    fmt.Printf("map3: %v (初始化)\n", map3)
    
    // Map的基本操作
    fmt.Println("\nMap操作:")
    
    // 添加和修改元素
    map2["apple"] = 10                // 添加新键值对
    map2["banana"] = 20               // 修改已有键的值
    map2["grape"] = 15                // 继续添加
    fmt.Printf("操作后map2: %v\n", map2)
    
    // 访问Map元素：两种返回值模式
    value := map3["apple"]            // 直接访问，键不存在返回零值
    fmt.Printf("apple的数量: %d\n", value)
    
    // 检查键是否存在：comma-ok模式
    if value, exists := map3["grape"]; exists {
        fmt.Printf("grape存在，值为: %d\n", value)
    } else {
        fmt.Println("grape不存在，返回零值:", value)
    }
    
    // 删除Map元素
    delete(map3, "banana")            // 删除键值对
    fmt.Printf("删除banana后: %v\n", map3)
    
    // 遍历Map：顺序随机
    fmt.Println("Map遍历:")
    for key, value := range map3 {
        fmt.Printf("%s: %d\n", key, value)
    }
    
    // 获取Map长度
    fmt.Printf("Map长度: %d\n", len(map3))
    
    /*
    第四节：多维数据结构
    
    多维数据结构用于表示复杂的数据关系：
    1. 二维数组：固定大小的矩阵
    2. 切片嵌套：动态大小的多维数据
    3. Map嵌套：层次化数据存储
    */
    fmt.Println("\n--- 多维数据结构 ---")
    
    // 二维数组：矩阵表示
    var matrix [3][3]int              // 3x3零矩阵
    matrix[0][0] = 1                  // 对角线赋值
    matrix[1][1] = 2
    matrix[2][2] = 3
    fmt.Printf("矩阵: %v\n", matrix)
    
    // 遍历二维数组
    fmt.Println("矩阵遍历:")
    for i, row := range matrix {
        for j, val := range row {
            fmt.Printf("matrix[%d][%d] = %d\n", i, j, val)
        }
    }
    
    // 二维切片：动态矩阵
    var slice2D [][]int               // nil二维切片
    slice2D = append(slice2D, []int{1, 2, 3})    // 添加行
    slice2D = append(slice2D, []int{4, 5, 6})    // 添加第二行
    slice2D = append(slice2D, []int{7, 8, 9})  // 添加第三行
    fmt.Printf("二维切片: %v\n", slice2D)
    
    // 访问二维切片
    fmt.Printf("slice2D[1][2] = %d\n", slice2D[1][2])
    
    // Map嵌套：层次化数据结构
    students := map[string]map[string]int{
        "张三": {
            "数学": 90,
            "英语": 85,
            "物理": 88,
        },
        "李四": {
            "数学": 95,
            "英语": 88,
            "物理": 92,
        },
        "王五": {
            "数学": 88,
            "英语": 92,
            "物理": 85,
        },
    }
    
    fmt.Printf("学生成绩: %v\n", students)
    fmt.Printf("张三的数学成绩: %d\n", students["张三"]["数学"])
    
    // 遍历嵌套Map
    fmt.Println("学生成绩详情:")
    for name, scores := range students {
        fmt.Printf("学生: %s\n", name)
        for subject, score := range scores {
            fmt.Printf("  %s: %d分\n", subject, score)
        }
    }
    
    /*
    第五节：实际应用示例
    
    通过实际案例展示复合数据结构的应用：
    1. 数据统计：使用Map进行频率统计
    2. 数据排序：使用切片实现排序算法
    3. 数据转换：不同结构间的转换
    */
    fmt.Println("\n--- 实际应用 ---")
    
    // 示例1：统计单词频率
    fmt.Println("1. 单词频率统计:")
    text := "hello world hello go world hello programming"
    words := []string{"hello", "world", "go", "hello", "world", "hello", "go", "programming"}
    
    wordCount := make(map[string]int)
    for _, word := range words {
        wordCount[word]++               // Map作为计数器
    }
    fmt.Printf("原文: %s\n", text)
    fmt.Printf("单词频率: %v\n", wordCount)
    
    // 示例2：切片排序实现
    fmt.Println("\n2. 切片排序实现:")
    numbers := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
    fmt.Printf("排序前: %v\n", numbers)
    
    // 冒泡排序实现
    bubbleSort := func(nums []int) {
        n := len(nums)
        for i := 0; i < n-1; i++ {
            for j := 0; j < n-i-1; j++ {
                if nums[j] > nums[j+1] {
                    nums[j], nums[j+1] = nums[j+1], nums[j]
                }
            }
        }
    }
    
    bubbleSort(numbers)
    fmt.Printf("排序后: %v\n", numbers)
    
    // 示例3：数据分组统计
    fmt.Println("\n3. 数据分组统计:")
    scores := []int{85, 92, 78, 95, 88, 76, 94, 87, 91, 83}
    
    // 按分数段分组
    scoreGroups := map[string][]int{
        "优秀(90-100)": {},
        "良好(80-89)":  {},
        "中等(70-79)":  {},
        "及格(60-69)":  {},
        "不及格(0-59)": {},
    }
    
    for _, score := range scores {
        switch {
        case score >= 90:
            scoreGroups["优秀(90-100)"] = append(scoreGroups["优秀(90-100)"], score)
        case score >= 80:
            scoreGroups["良好(80-89)"] = append(scoreGroups["良好(80-89)"], score)
        case score >= 70:
            scoreGroups["中等(70-79)"] = append(scoreGroups["中等(70-79)"], score)
        case score >= 60:
            scoreGroups["及格(60-69)"] = append(scoreGroups["及格(60-69)"], score)
        default:
            scoreGroups["不及格(0-59)"] = append(scoreGroups["不及格(0-59)"], score)
        }
    }
    
    // 显示分组结果
    for group, scores := range scoreGroups {
        if len(scores) > 0 {
            fmt.Printf("%s: %v\n", group, scores)
        }
    }
    
    /*
    练习指导：
    
    1. 学生信息Map：
       - 使用map[string]interface{}存储不同类型的值
       - 包含姓名、年龄、成绩、爱好等字段
    
    2. 切片反转函数：
       - 使用双指针交换元素
       - 注意边界条件处理
    
    3. 数组极值查找：
       - 遍历一次找出最大值和最小值
       - 处理空数组的情况
    
    4. Map合并：
       - 处理键冲突的策略
       - 深拷贝vs浅拷贝
    */
    fmt.Println("\n练习指导:")
    fmt.Println("1. 创建一个学生信息Map，包含姓名、年龄、成绩、爱好")
    fmt.Println("2. 实现一个切片的反转函数")
    fmt.Println("3. 找出数组中的最大值、最小值和平均值")
    fmt.Println("4. 合并两个Map，处理键冲突")
    
    // 练习答案示例
    fmt.Println("\n练习答案示例:")
    
    // 1. 学生信息Map
    student := map[string]interface{}{
        "姓名": "张三",
        "年龄": 20,
        "成绩": []int{85, 92, 78},
        "爱好": []string{"编程", "阅读", "运动"},
    }
    fmt.Printf("学生信息: %v\n", student)
    
    // 2. 切片反转函数
    reverseSlice := func(s []int) {
        for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
            s[i], s[j] = s[j], s[i]
        }
    }
    testSlice := []int{1, 2, 3, 4, 5}
    reverseSlice(testSlice)
    fmt.Printf("反转后: %v\n", testSlice)
    
    // 3. 数组极值查找
    findExtremes := func(arr []int) (min, max, avg int) {
        if len(arr) == 0 {
            return 0, 0, 0
        }
        min, max = arr[0], arr[0]
        sum := 0
        for _, v := range arr {
            if v < min {
                min = v
            }
            if v > max {
                max = v
            }
            sum += v
        }
        return min, max, sum / len(arr)
    }
    testArray := []int{3, 7, 2, 9, 1, 5}
    min, max, avg := findExtremes(testArray)
    fmt.Printf("最小值: %d, 最大值: %d, 平均值: %d\n", min, max, avg)
    
    // 4. Map合并
    mergeMaps := func(map1, map2 map[string]int) map[string]int {
        result := make(map[string]int)
        for k, v := range map1 {
            result[k] = v
        }
        for k, v := range map2 {
            result[k] = v  // 键冲突时，map2的值覆盖map1
        }
        return result
    }
    
    mapA := map[string]int{"a": 1, "b": 2}
    mapB := map[string]int{"b": 3, "c": 4}
    merged := mergeMaps(mapA, mapB)
    fmt.Printf("合并结果: %v\n", merged)
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