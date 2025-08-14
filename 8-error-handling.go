package main

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)

// 错误处理
// Go语言的错误处理机制

// 1. 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("字段 %s: %s", e.Field, e.Message)
}

// 2. 网络错误类型
type NetworkError struct {
    Code    int
    Message string
    URL     string
}

func (e *NetworkError) Error() string {
    return fmt.Sprintf("网络错误 %d: %s (URL: %s)", e.Code, e.Message, e.URL)
}

// 3. 业务错误类型
type BusinessError struct {
    Type    string
    Details map[string]interface{}
}

func (e *BusinessError) Error() string {
    return fmt.Sprintf("业务错误: %s, 详情: %v", e.Type, e.Details)
}

func main() {
    fmt.Println("=== Go语言错误处理 ===")
    
    // 1. 基本错误处理
    fmt.Println("\n--- 基本错误处理 ---")
    
    if result, err := divide(10, 2); err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("结果: %d\n", result)
    }
    
    if result, err := divide(10, 0); err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("结果: %d\n", result)
    }
    
    // 2. 自定义错误类型
    fmt.Println("\n--- 自定义错误类型 ---")
    
    // 验证错误
    if err := validateUser("", 15); err != nil {
        fmt.Printf("验证错误: %v\n", err)
        
        // 类型断言检查错误类型
        if validationErr, ok := err.(*ValidationError); ok {
            fmt.Printf("字段错误: %s, 消息: %s\n", validationErr.Field, validationErr.Message)
        }
    }
    
    // 网络错误
    if err := fetchData("https://example.com"); err != nil {
        fmt.Printf("网络错误: %v\n", err)
    }
    
    // 3. 错误包装和链式错误
    fmt.Println("\n--- 错误包装 ---")
    
    if err := processFile("nonexistent.txt"); err != nil {
        fmt.Printf("文件处理错误: %v\n", err)
        
        // 使用errors.Is和errors.As检查错误
        if errors.Is(err, ErrFileNotFound) {
            fmt.Println("这是文件未找到错误")
        }
    }
    
    // 4. 错误处理最佳实践
    fmt.Println("\n--- 错误处理最佳实践 ---")
    
    // 使用错误处理函数
    if err := processUserData(UserData{Name: "张三", Age: 25}); err != nil {
        handleError(err)
    }
    
    // 5. 错误聚合
    fmt.Println("\n--- 错误聚合 ---")
    
    if err := validateMultipleFields(); err != nil {
        fmt.Printf("多个验证错误: %v\n", err)
    }
    
    // 6. 重试机制
    fmt.Println("\n--- 重试机制 ---")
    
    if result, err := retryOperation(func() (int, error) {
        return riskyOperation()
    }, 3); err != nil {
        fmt.Printf("操作失败: %v\n", err)
    } else {
        fmt.Printf("操作成功: %d\n", result)
    }
    
    // 7. 错误恢复
    fmt.Println("\n--- 错误恢复 ---")
    
    safeDivideWithRecovery(10, 0)
    
    // 8. 错误分类和处理
    fmt.Println("\n--- 错误分类 ---")
    
    errors := []error{
        &ValidationError{Field: "email", Message: "格式不正确"},
        &NetworkError{Code: 404, Message: "页面未找到", URL: "/api/users"},
        &BusinessError{Type: "insufficient_balance", Details: map[string]interface{}{"balance": 100.0}},
    }
    
    for _, err := range errors {
        handleCategorizedError(err)
    }
    
    // 9. 错误日志记录
    fmt.Println("\n--- 错误日志记录 ---")
    
    logger := NewErrorLogger()
    
    if err := parseAndValidate("123"); err != nil {
        logger.Log(err)
    }
    
    if err := parseAndValidate("abc"); err != nil {
        logger.Log(err)
    }
    
    // 10. 错误转换和映射
    fmt.Println("\n--- 错误转换 ---")
    
    if err := convertErrorExample(); err != nil {
        fmt.Printf("转换后的错误: %v\n", err)
    }
    
    // 练习
    fmt.Println("\n练习：")
    fmt.Println("1. 实现一个文件解析器，处理各种文件格式错误")
    fmt.Println("2. 创建一个HTTP客户端，实现重试和错误恢复机制")
    fmt.Println("3. 设计一个验证框架，支持字段级错误聚合")
}

// 基本除法函数
divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 用户验证函数
validateUser(name string, age int) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "姓名不能为空"}
    }
    if age < 18 {
        return &ValidationError{Field: "age", Message: "年龄必须大于等于18岁"}
    }
    return nil
}

// 模拟网络请求
fetchData(url string) error {
    if !strings.HasPrefix(url, "https://") {
        return &NetworkError{Code: 400, Message: "无效的URL协议", URL: url}
    }
    return &NetworkError{Code: 404, Message: "页面未找到", URL: url}
}

// 预定义错误
var (
    ErrFileNotFound    = errors.New("文件未找到")
    ErrPermissionDenied = errors.New("权限被拒绝")
    ErrInvalidFormat    = errors.New("格式无效")
)

// 文件处理函数
processFile(filename string) error {
    if filename == "nonexistent.txt" {
        return fmt.Errorf("处理文件 %s: %w", filename, ErrFileNotFound)
    }
    return nil
}

// 用户数据结构
type UserData struct {
    Name  string
    Age   int
    Email string
}

// 用户数据处理
processUserData(data UserData) error {
    if data.Name == "" {
        return &ValidationError{Field: "name", Message: "姓名不能为空"}
    }
    
    if data.Age <= 0 || data.Age > 150 {
        return &ValidationError{Field: "age", Message: "年龄必须在1-150之间"}
    }
    
    if !strings.Contains(data.Email, "@") {
        return &ValidationError{Field: "email", Message: "邮箱格式不正确"}
    }
    
    return nil
}

// 错误处理函数
handleError(err error) {
    fmt.Printf("处理错误: %v\n", err)
    
    // 根据错误类型进行不同处理
    switch e := err.(type) {
    case *ValidationError:
        fmt.Printf("显示验证错误给用户: %s\n", e.Message)
    case *NetworkError:
        if e.Code >= 500 {
            fmt.Println("服务器错误，请稍后重试")
        } else if e.Code >= 400 {
            fmt.Println("客户端错误，请检查输入")
        }
    default:
        fmt.Println("未知错误，请联系技术支持")
    }
}

// 验证多个字段
validateMultipleFields() error {
    var validationErrors []error
    
    // 模拟多个验证错误
    validationErrors = append(validationErrors, 
        &ValidationError{Field: "email", Message: "格式不正确"})
    validationErrors = append(validationErrors,
        &ValidationError{Field: "password", Message: "长度不足"})
    validationErrors = append(validationErrors,
        &ValidationError{Field: "phone", Message: "格式无效"})
    
    if len(validationErrors) > 0 {
        return &ValidationErrors{Errors: validationErrors}
    }
    
    return nil
}

// 验证错误聚合类型
type ValidationErrors struct {
    Errors []error
}

func (e *ValidationErrors) Error() string {
    var messages []string
    for _, err := range e.Errors {
        messages = append(messages, err.Error())
    }
    return fmt.Sprintf("验证失败: %s", strings.Join(messages, "; "))
}

// 重试机制
retryOperation(operation func() (int, error), maxRetries int) (int, error) {
    for i := 0; i < maxRetries; i++ {
        if result, err := operation(); err == nil {
            return result, nil
        }
        
        if i < maxRetries-1 {
            fmt.Printf("重试 %d/%d\n", i+1, maxRetries)
            time.Sleep(time.Millisecond * 500)
        }
    }
    
    return 0, fmt.Errorf("操作失败，已重试%d次", maxRetries)
}

// 模拟有风险的业务操作
riskyOperation() (int, error) {
    if time.Now().UnixNano()%2 == 0 {
        return 42, nil
    }
    return 0, errors.New("操作失败")
}

// 安全除法（带恢复）
safeDivideWithRecovery(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("捕获到panic: %v\n", r)
        }
    }()
    
    if b == 0 {
        panic("除数不能为零")
    }
    
    result := a / b
    fmt.Printf("除法结果: %d\n", result)
}

// 错误分类处理
handleCategorizedError(err error) {
    fmt.Printf("处理错误: %v\n", err)
    
    switch e := err.(type) {
    case *ValidationError:
        fmt.Printf("验证错误 - 字段: %s, 消息: %s\n", e.Field, e.Message)
    case *NetworkError:
        fmt.Printf("网络错误 - 代码: %d, 消息: %s\n", e.Code, e.Message)
    case *BusinessError:
        fmt.Printf("业务错误 - 类型: %s, 详情: %v\n", e.Type, e.Details)
    }
}

// 错误日志记录器
type ErrorLogger struct {
    errors []error
}

func NewErrorLogger() *ErrorLogger {
    return &ErrorLogger{errors: make([]error, 0)}
}

func (l *ErrorLogger) Log(err error) {
    l.errors = append(l.errors, err)
    fmt.Printf("记录错误: %v\n", err)
}

func (l *ErrorLogger) GetErrors() []error {
    return l.errors
}

// 解析和验证
parseAndValidate(input string) error {
    num, err := strconv.Atoi(input)
    if err != nil {
        return fmt.Errorf("解析数字失败: %w", err)
    }
    
    if num < 0 || num > 100 {
        return &ValidationError{Field: "number", Message: "数字必须在0-100之间"}
    }
    
    return nil
}

// 错误转换示例
convertErrorExample() error {
    err := processFile("test.txt")
    if err != nil {
        return &BusinessError{
            Type:    "file_processing_failed",
            Details: map[string]interface{}{"original_error": err.Error()},
        }
    }
    return nil
}