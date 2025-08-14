package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time" // 注意：原代码遗漏了time包导入，此处补充以支持时间相关操作
)

// 1. 自定义错误类型：验证错误
// 用途：专门用于字段验证场景，包含具体字段名和错误信息
// 实现思路：通过结构体存储错误详情，实现error接口的Error()方法使其成为合法错误类型
type ValidationError struct {
	Field   string // 出错的字段名
	Message string // 错误描述信息
}

// Error()方法：实现error接口，返回格式化的错误信息
// 作用：使ValidationError能够被当作error类型使用，便于统一错误处理
func (e *ValidationError) Error() string {
	return fmt.Sprintf("字段 %s: %s", e.Field, e.Message)
}

// 2. 网络错误类型
// 用途：专门处理网络请求相关错误，包含状态码、描述和请求URL
// 设计思路：聚合网络错误的关键信息，方便后续根据状态码等进行针对性处理
type NetworkError struct {
	Code    int    // 网络错误状态码（如404、500等）
	Message string // 错误描述
	URL     string // 发生错误的URL地址
}

// Error()方法：实现error接口，格式化网络错误信息
func (e *NetworkError) Error() string {
	return fmt.Sprintf("网络错误 %d: %s (URL: %s)", e.Code, e.Message, e.URL)
}

// 3. 业务错误类型
// 用途：处理业务逻辑相关错误，包含错误类型和详细信息
// 设计思路：使用map存储灵活的错误详情，适应不同业务场景的错误信息需求
type BusinessError struct {
	Type    string                 // 业务错误类型（如"insufficient_balance"）
	Details map[string]interface{} // 错误详情，键值对形式存储额外信息
}

// Error()方法：实现error接口，格式化业务错误信息
func (e *BusinessError) Error() string {
	return fmt.Sprintf("业务错误: %s, 详情: %v", e.Type, e.Details)
}

// main函数：程序入口，展示各种错误处理机制的使用示例
// 实现思路：通过分步骤演示不同错误处理场景，从基础到复杂，循序渐进展示Go错误处理特性
func main() {
	fmt.Println("=== Go语言错误处理 ===")
	
	// 1. 基本错误处理：演示最简单的错误返回与判断
	// 核心思想：函数通过返回值返回错误，调用方检查错误是否为nil来判断执行结果
	fmt.Println("\n--- 基本错误处理 ---")
	
	// 正常情况：除数不为0，无错误
	if result, err := divide(10, 2); err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}
	
	// 错误情况：除数为0，返回错误
	if result, err := divide(10, 0); err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}
	
	// 2. 自定义错误类型：展示如何使用自定义错误传递更丰富的信息
	// 优势：相比标准错误，能携带更多上下文信息，便于精确处理
	fmt.Println("\n--- 自定义错误类型 ---")
	
	// 验证错误示例：验证用户信息不符合要求时返回ValidationError
	if err := validateUser("", 15); err != nil {
		fmt.Printf("验证错误: %v\n", err)
		
		// 类型断言：判断错误具体类型，进行针对性处理
		// 关键字说明：err.(*ValidationError)用于将error类型转换为具体的自定义错误类型
		// ok为true表示转换成功，确保类型安全避免panic
		if validationErr, ok := err.(*ValidationError); ok {
			fmt.Printf("字段错误: %s, 消息: %s\n", validationErr.Field, validationErr.Message)
		}
	}
	
	// 网络错误示例：模拟网络请求出错返回NetworkError
	if err := fetchData("https://example.com"); err != nil {
		fmt.Printf("网络错误: %v\n", err)
	}
	
	// 3. 错误包装和链式错误：使用fmt.Errorf的%w格式符包装原始错误
	// 作用：保留错误链，既添加上下文信息，又不丢失原始错误，便于追溯
	fmt.Println("\n--- 错误包装 ---")
	
	if err := processFile("nonexistent.txt"); err != nil {
		fmt.Printf("文件处理错误: %v\n", err)
		
		// errors.Is：检查错误链中是否包含目标错误
		// 优势：无需手动遍历错误链，简化错误类型判断
		if errors.Is(err, ErrFileNotFound) {
			fmt.Println("这是文件未找到错误")
		}
	}
	
	// 4. 错误处理最佳实践：使用专门的错误处理函数统一处理不同类型错误
	// 好处：集中处理逻辑，避免代码重复，便于维护
	fmt.Println("\n--- 错误处理最佳实践 ---")
	
	// 使用错误处理函数处理用户数据验证错误
	if err := processUserData(UserData{Name: "张三", Age: 25}); err != nil {
		handleError(err)
	}
	
	// 5. 错误聚合：当有多个错误需要同时返回时（如表单多字段验证）
	// 优势：一次性返回所有错误，避免用户多次提交才能发现所有问题
	fmt.Println("\n--- 错误聚合 ---")
	
	if err := validateMultipleFields(); err != nil {
		fmt.Printf("多个验证错误: %v\n", err)
	}
	
	// 6. 重试机制：针对可能暂时失败的操作（如网络请求），重试提高成功率
	// 适用场景：临时性错误（网络波动、资源暂时不可用等）
	fmt.Println("\n--- 重试机制 ---")
	
	// 调用重试函数，最多重试3次
	if result, err := retryOperation(func() (int, error) {
		return riskyOperation() // 传入可能失败的操作
	}, 3); err != nil {
		fmt.Printf("操作失败: %v\n", err)
	} else {
		fmt.Printf("操作成功: %d\n", result)
	}
	
	// 7. 错误恢复：使用recover捕获panic，防止程序崩溃
	// 注意：Go中不推荐用panic处理预期错误，主要用于处理不可恢复的严重错误
	fmt.Println("\n--- 错误恢复 ---")
	
	safeDivideWithRecovery(10, 0) // 故意传入0作为除数，测试panic捕获
	
	// 8. 错误分类和处理：根据错误类型执行不同处理逻辑
	// 实现思路：使用类型断言区分错误类型，针对性处理
	fmt.Println("\n--- 错误分类 ---")
	
	errors := []error{
		&ValidationError{Field: "email", Message: "格式不正确"},
		&NetworkError{Code: 404, Message: "页面未找到", URL: "/api/users"},
		&BusinessError{Type: "insufficient_balance", Details: map[string]interface{}{"balance": 100.0}},
	}
	
	for _, err := range errors {
		handleCategorizedError(err)
	}
	
	// 9. 错误日志记录：收集和管理错误信息，便于调试和监控
	// 作用：持久化错误记录，为问题排查提供依据
	fmt.Println("\n--- 错误日志记录 ---")
	
	logger := NewErrorLogger() // 创建错误日志记录器
	
	if err := parseAndValidate("123"); err != nil {
		logger.Log(err) // 记录错误
	}
	
	if err := parseAndValidate("abc"); err != nil {
		logger.Log(err) // 记录错误
	}
	
	// 10. 错误转换和映射：将底层错误转换为上层业务可理解的错误
	// 用途：隐藏实现细节，提供符合业务域的错误信息
	fmt.Println("\n--- 错误转换 ---")
	
	if err := convertErrorExample(); err != nil {
		fmt.Printf("转换后的错误: %v\n", err)
	}
	
	// 练习提示：引导扩展错误处理的实际应用场景
	fmt.Println("\n练习：")
	fmt.Println("1. 实现一个文件解析器，处理各种文件格式错误")
	fmt.Println("2. 创建一个HTTP客户端，实现重试和错误恢复机制")
	fmt.Println("3. 设计一个验证框架，支持字段级错误聚合")
}

// divide：基础除法函数，演示简单错误返回
// 参数：a（被除数）、b（除数）
// 返回值：商（int）和错误（error）
// 实现思路：检查除数为0的非法情况，返回相应错误；否则返回计算结果
func divide(a, b int) (int, error) {
	if b == 0 {
		// errors.New：创建一个简单的错误实例
		// 适用场景：不需要额外上下文信息的简单错误
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil // 无错误时返回结果和nil
}

// validateUser：用户信息验证函数，演示自定义错误的使用
// 参数：name（用户名）、age（年龄）
// 返回值：错误（error），验证通过返回nil
// 验证逻辑：检查用户名非空和年龄不小于18岁，不符合则返回对应ValidationError
func validateUser(name string, age int) error {
	if name == "" {
		// 返回自定义验证错误，指定出错字段和原因
		return &ValidationError{Field: "name", Message: "姓名不能为空"}
	}
	if age < 18 {
		return &ValidationError{Field: "age", Message: "年龄必须大于等于18岁"}
	}
	return nil // 验证通过
}

// fetchData：模拟网络请求函数，演示网络错误处理
// 参数：url（请求的URL地址）
// 返回值：错误（error），模拟请求失败场景
// 实现思路：检查URL协议是否合法，模拟返回不同网络错误
func fetchData(url string) error {
	// 检查URL是否以https://开头
	if !strings.HasPrefix(url, "https://") {
		return &NetworkError{Code: 400, Message: "无效的URL协议", URL: url}
	}
	// 模拟页面未找到错误
	return &NetworkError{Code: 404, Message: "页面未找到", URL: url}
}

// 预定义错误变量：使用errors.New创建可复用的错误实例
// 优势：便于使用errors.Is进行错误判断，确保错误比较的准确性
// 注意：变量名通常以Err开头，符合Go命名规范
var (
	ErrFileNotFound    = errors.New("文件未找到")    // 表示文件不存在错误
	ErrPermissionDenied = errors.New("权限被拒绝")   // 表示权限不足错误
	ErrInvalidFormat    = errors.New("格式无效")     // 表示数据格式错误
)

// processFile：文件处理函数，演示错误包装（错误链）
// 参数：filename（文件名）
// 返回值：错误（error），包含上下文的包装错误
// 实现思路：使用fmt.Errorf的%w格式符包装预定义错误，添加上下文信息
func processFile(filename string) error {
	if filename == "nonexistent.txt" {
		// %w：将原始错误包装到新错误中，形成错误链
		// 优势：既保留原始错误信息，又添加了当前操作的上下文（文件名）
		return fmt.Errorf("处理文件 %s: %w", filename, ErrFileNotFound)
	}
	return nil
}

// UserData：用户数据结构体，存储用户相关信息
// 用途：作为processUserData函数的输入，演示复杂数据的验证
type UserData struct {
	Name  string // 姓名
	Age   int    // 年龄
	Email string // 邮箱
}

// processUserData：用户数据处理函数，演示多字段验证
// 参数：data（UserData结构体实例）
// 返回值：错误（error），验证失败返回ValidationError
// 实现思路：对用户的姓名、年龄、邮箱字段分别进行验证，返回对应错误
func processUserData(data UserData) error {
	if data.Name == "" {
		return &ValidationError{Field: "name", Message: "姓名不能为空"}
	}
	
	if data.Age <= 0 || data.Age > 150 {
		return &ValidationError{Field: "age", Message: "年龄必须在1-150之间"}
	}
	
	if !strings.Contains(data.Email, "@") {
		return &ValidationError{Field: "email", Message: "邮箱格式不正确"}
	}
	
	return nil // 所有字段验证通过
}

// handleError：错误处理函数，演示根据错误类型进行不同处理
// 参数：err（需要处理的错误）
// 实现思路：使用类型断言区分错误类型，执行对应的处理逻辑
func handleError(err error) {
	fmt.Printf("处理错误: %v\n", err)
	
	// switch类型断言：高效区分不同错误类型
	// 优势：比多个if-else断言更简洁，处理多种错误类型时更清晰
	switch e := err.(type) {
	case *ValidationError:
		// 验证错误：通常需要提示用户修正输入
		fmt.Printf("显示验证错误给用户: %s\n", e.Message)
	case *NetworkError:
		// 网络错误：根据状态码给出不同提示
		if e.Code >= 500 {
			fmt.Println("服务器错误，请稍后重试")
		} else if e.Code >= 400 {
			fmt.Println("客户端错误，请检查输入")
		}
	default:
		// 未知错误：通用处理逻辑
		fmt.Println("未知错误，请联系技术支持")
	}
}

// validateMultipleFields：多字段验证函数，演示错误聚合
// 返回值：错误（error），包含所有验证错误的聚合错误
// 实现思路：收集所有字段的验证错误，统一返回一个包含多个错误的聚合错误
func validateMultipleFields() error {
	var validationErrors []error // 用于存储多个验证错误
	
	// 模拟多个字段验证失败的场景
	validationErrors = append(validationErrors, 
		&ValidationError{Field: "email", Message: "格式不正确"})
	validationErrors = append(validationErrors,
		&ValidationError{Field: "password", Message: "长度不足"})
	validationErrors = append(validationErrors,
		&ValidationError{Field: "phone", Message: "格式无效"})
	
	// 如果有错误，返回错误聚合实例
	if len(validationErrors) > 0 {
		return &ValidationErrors{Errors: validationErrors}
	}
	
	return nil // 无错误
}

// ValidationErrors：验证错误聚合类型，实现error接口
// 用途：一次性返回多个验证错误，便于批量处理
type ValidationErrors struct {
	Errors []error // 存储多个错误的切片
}

// Error()方法：实现error接口，将所有错误信息合并为一个字符串
func (e *ValidationErrors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	// 用分号分隔多个错误信息
	return fmt.Sprintf("验证失败: %s", strings.Join(messages, "; "))
}

// retryOperation：重试机制实现函数
// 参数：
//   operation：需要重试的操作（函数类型，返回结果和错误）
//   maxRetries：最大重试次数
// 返回值：操作结果和最终错误
// 实现思路：循环执行操作，成功则返回；失败则重试，直到达到最大次数
func retryOperation(operation func() (int, error), maxRetries int) (int, error) {
	for i := 0; i < maxRetries; i++ {
		// 执行操作
		if result, err := operation(); err == nil {
			return result, nil // 操作成功，返回结果
		}
		
		// 不是最后一次重试，则等待后继续
		if i < maxRetries-1 {
			fmt.Printf("重试 %d/%d\n", i+1, maxRetries)
			time.Sleep(time.Millisecond * 500) // 等待500毫秒后重试
		}
	}
	
	// 达到最大重试次数仍失败
	return 0, fmt.Errorf("操作失败，已重试%d次", maxRetries)
}

// riskyOperation：模拟有风险的业务操作，随机失败
// 返回值：结果和错误
// 实现思路：使用当前时间的纳秒数判断，模拟偶发失败的场景
func riskyOperation() (int, error) {
	// 随机决定成功或失败（根据当前时间的纳秒数奇偶性）
	if time.Now().UnixNano()%2 == 0 {
		return 42, nil // 成功返回结果
	}
	return 0, errors.New("操作失败") // 失败返回错误
}

// safeDivideWithRecovery：带错误恢复的安全除法函数
// 参数：a（被除数）、b（除数）
// 实现思路：使用defer和recover捕获可能的panic，防止程序崩溃
func safeDivideWithRecovery(a, b int) {
	// defer语句：确保在函数退出前执行匿名函数
	// 作用：即使发生panic，也能执行恢复逻辑
	defer func() {
		// recover()：捕获当前goroutine中的panic，返回panic的值
		// 注意：仅在defer语句中调用才有效
		if r := recover(); r != nil {
			fmt.Printf("捕获到panic: %v\n", r) // 处理panic，避免程序退出
		}
	}()
	
	if b == 0 {
		// panic：触发程序异常，通常用于不可恢复的错误
		// 注意：Go中应避免用panic处理预期错误，这里仅作演示
		panic("除数不能为零")
	}
	
	result := a / b
	fmt.Printf("除法结果: %d\n", result)
}

// handleCategorizedError：错误分类处理函数
// 参数：err（需要处理的错误）
// 实现思路：根据错误的具体类型，提取相关信息并进行处理
func handleCategorizedError(err error) {
	fmt.Printf("处理错误: %v\n", err)
	
	// 类型断言区分不同错误类型
	switch e := err.(type) {
	case *ValidationError:
		fmt.Printf("验证错误 - 字段: %s, 消息: %s\n", e.Field, e.Message)
	case *NetworkError:
		fmt.Printf("网络错误 - 代码: %d, 消息: %s\n", e.Code, e.Message)
	case *BusinessError:
		fmt.Printf("业务错误 - 类型: %s, 详情: %v\n", e.Type, e.Details)
	}
}

// ErrorLogger：错误日志记录器，用于收集和管理错误信息
// 用途：集中记录程序运行过程中发生的错误，便于后续分析
type ErrorLogger struct {
	errors []error // 存储错误的切片
}

// NewErrorLogger：创建新的错误日志记录器实例
// 返回值：*ErrorLogger
func NewErrorLogger() *ErrorLogger {
	return &ErrorLogger{errors: make([]error, 0)}
}

// Log：记录错误到日志器中
// 参数：err（需要记录的错误）
func (l *ErrorLogger) Log(err error) {
	l.errors = append(l.errors, err)
	fmt.Printf("记录错误: %v\n", err)
}

// GetErrors：获取所有记录的错误
// 返回值：错误切片
func (l *ErrorLogger) GetErrors() []error {
	return l.errors
}

// parseAndValidate：解析并验证输入字符串
// 参数：input（需要解析的字符串）
// 返回值：错误（error）
// 实现思路：先将字符串转换为整数，再验证范围，演示错误包装的使用
func parseAndValidate(input string) error {
	// 尝试将字符串转换为整数
	num, err := strconv.Atoi(input)
	if err != nil {
		// 包装原始错误，添加上下文信息
		return fmt.Errorf("解析数字失败: %w", err)
	}
	
	// 验证数字是否在有效范围内
	if num < 0 || num > 100 {
		return &ValidationError{Field: "number", Message: "数字必须在0-100之间"}
	}
	
	return nil // 解析和验证都通过
}

// convertErrorExample：错误转换示例函数
// 返回值：转换后的错误
// 实现思路：将底层操作错误转换为上层业务错误，隐藏实现细节
func convertErrorExample() error {
	// 执行可能出错的操作
	err := processFile("test.txt")
	if err != nil {
		// 将文件处理错误转换为业务错误
		return &BusinessError{
			Type:    "file_processing_failed",
			Details: map[string]interface{}{"original_error": err.Error()},
		}
	}
	return nil
}
