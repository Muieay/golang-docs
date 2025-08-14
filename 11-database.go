package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	// 导入SQLite驱动，下划线表示只使用其初始化函数而不直接调用其方法
	// go-sqlite3是SQLite的Go语言驱动，实现了database/sql接口
	_ "github.com/mattn/go-sqlite3"
)

// 数据库操作
// Go语言的数据库编程和ORM概念（本示例使用原生SQL，未使用ORM框架）
// 本示例展示了如何使用Go标准库database/sql操作SQLite数据库，包括：
// - 数据模型定义（与数据库表对应）
// - 数据库连接管理
// - 表结构初始化
// - CRUD（创建、读取、更新、删除）操作
// - 事务处理（保证数据一致性）
// - 关联查询（多表连接）
// - 数据统计与分析

// 1. 数据模型
// 数据模型（结构体）与数据库表结构一一对应，字段名和类型保持一致
// 结构体标签`json:"字段名"`用于JSON序列化/反序列化时的字段映射

// Product：产品数据模型，对应products表
type Product struct {
	ID          int       `json:"id"`           // 产品唯一标识，自增主键
	Name        string    `json:"name"`         // 产品名称
	Price       float64   `json:"price"`        // 产品价格
	Category    string    `json:"category"`     // 产品分类
	Description string    `json:"description"`  // 产品描述
	CreatedAt   time.Time `json:"created_at"`   // 创建时间，数据库自动生成
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间，数据库自动更新
}

// Order：订单数据模型，对应orders表
type Order struct {
	ID         int         `json:"id"`          // 订单唯一标识，自增主键
	UserID     int         `json:"user_id"`     // 关联的用户ID（外键逻辑）
	Total      float64     `json:"total"`       // 订单总金额
	Status     string      `json:"status"`      // 订单状态（pending/paid/shipped等）
	CreatedAt  time.Time   `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time   `json:"updated_at"`  // 更新时间
	Items      []OrderItem `json:"items"`       // 订单项列表（关联数据）
}

// OrderItem：订单项数据模型，对应order_items表
// 用于关联订单和产品，记录购买数量和单价
type OrderItem struct {
	ID        int     `json:"id"`         // 订单项唯一标识
	OrderID   int     `json:"order_id"`   // 关联的订单ID（外键）
	ProductID int     `json:"product_id"` // 关联的产品ID（外键）
	Quantity  int     `json:"quantity"`   // 购买数量
	Price     float64 `json:"price"`      // 购买时的单价
}

// 2. 数据库管理器
// DatabaseManager：数据库管理器，封装数据库连接和操作方法
// 采用面向对象风格设计，通过结构体方法实现数据库操作的封装
type DatabaseManager struct {
	db *sql.DB // 数据库连接对象，*sql.DB是线程安全的，可在多个goroutine中共享
}

// NewDatabaseManager：创建新的数据库管理器
// 参数：dbPath - 数据库文件路径（SQLite使用文件存储数据库）
// 返回值：*DatabaseManager - 数据库管理器实例；error - 可能的错误
func NewDatabaseManager(dbPath string) (*DatabaseManager, error) {
	// sql.Open：打开数据库连接
	// 第一个参数是驱动名称（"sqlite3"对应导入的驱动）
	// 第二个参数是数据源名称（SQLite为文件路径）
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		// 使用fmt.Errorf包装错误，保留原始错误信息（%w格式符）
		return nil, fmt.Errorf("无法打开数据库: %w", err)
	}
	
	// db.Ping()：验证数据库连接是否有效
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}
	
	// 返回数据库管理器实例，包含有效的数据库连接
	return &DatabaseManager{db: db}, nil
}

// Close：关闭数据库连接
// 实现资源释放，应在程序退出前调用（通常使用defer）
func (dm *DatabaseManager) Close() error {
	return dm.db.Close()
}

// 3. 初始化数据库表
// InitializeSchema：创建数据库表结构和索引
// 功能：如果表不存在则创建，确保数据库结构正确
func (dm *DatabaseManager) InitializeSchema() error {
	// 创建产品表SQL语句
	// IF NOT EXISTS：确保表不存在时才创建，避免重复创建错误
	productTable := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增主键
        name TEXT NOT NULL,                    -- 产品名称，非空
        price REAL NOT NULL,                   -- 产品价格，非空
        category TEXT NOT NULL,                -- 产品分类，非空
        description TEXT,                      -- 产品描述，可空
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认当前时间
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP   -- 更新时间，默认当前时间
    );`
	
	// 创建订单表SQL语句
	orderTable := `
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,              -- 用户ID，非空
        total REAL NOT NULL,                   -- 订单总金额，非空
        status TEXT NOT NULL DEFAULT 'pending', -- 订单状态，默认pending
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	
	// 创建订单项表SQL语句
	// 包含外键约束（FOREIGN KEY），确保数据引用完整性
	orderItemTable := `
    CREATE TABLE IF NOT EXISTS order_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        order_id INTEGER NOT NULL,             -- 关联订单ID，非空
        product_id INTEGER NOT NULL,           -- 关联产品ID，非空
        quantity INTEGER NOT NULL,             -- 数量，非空
        price REAL NOT NULL,                   -- 单价，非空
        -- 外键约束：order_id引用orders表的id
        FOREIGN KEY (order_id) REFERENCES orders(id),
        -- 外键约束：product_id引用products表的id
        FOREIGN KEY (product_id) REFERENCES products(id)
    );`
	
	// 创建索引SQL语句，提高查询效率
	// 索引通常创建在频繁查询的字段上
	createIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);", // 按分类查询产品
		"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);",       // 按用户查询订单
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);",         // 按状态查询订单
		"CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);", // 按订单查询订单项
	}
	
	// 执行创建表的SQL语句
	// db.Exec()：执行不返回结果的SQL语句（CREATE/INSERT/UPDATE/DELETE等）
	if _, err := dm.db.Exec(productTable); err != nil {
		return fmt.Errorf("创建产品表失败: %w", err)
	}
	
	if _, err := dm.db.Exec(orderTable); err != nil {
		return fmt.Errorf("创建订单表失败: %w", err)
	}
	
	if _, err := dm.db.Exec(orderItemTable); err != nil {
		return fmt.Errorf("创建订单项表失败: %w", err)
	}
	
	// 创建索引
	for _, indexSQL := range createIndexes {
		if _, err := dm.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("创建索引失败: %w", err)
		}
	}
	
	return nil
}

// 4. 产品相关操作
// 以下方法实现产品的CRUD（创建、读取、更新、删除）操作

// CreateProduct：创建新产品
// 参数：product - 指向Product结构体的指针，包含产品信息（ID会被自动生成）
// 返回值：error - 可能的错误
func (dm *DatabaseManager) CreateProduct(product *Product) error {
	// SQL插入语句，使用?作为占位符（防止SQL注入）
	query := `INSERT INTO products (name, price, category, description) VALUES (?, ?, ?, ?)`
	
	// 执行插入操作，参数按顺序对应占位符
	// db.Exec()返回sql.Result，包含插入的ID和受影响的行数
	result, err := dm.db.Exec(query, product.Name, product.Price, product.Category, product.Description)
	if err != nil {
		return fmt.Errorf("创建产品失败: %w", err)
	}
	
	// 获取插入记录的自增ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取产品ID失败: %w", err)
	}
	
	// 将生成的ID赋值给product结构体
	product.ID = int(id)
	return nil
}

// GetProduct：根据ID查询产品
// 参数：id - 产品ID
// 返回值：*Product - 产品信息；error - 可能的错误
func (dm *DatabaseManager) GetProduct(id int) (*Product, error) {
	// SQL查询语句，根据ID查询产品
	query := `SELECT id, name, price, category, description, created_at, updated_at FROM products WHERE id = ?`
	
	var product Product
	// db.QueryRow()：执行查询并返回单行结果
	// Scan()：将查询结果映射到结构体字段（顺序必须与SELECT一致）
	err := dm.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Price, &product.Category,
		&product.Description, &product.CreatedAt, &product.UpdatedAt,
	)
	
	// 处理查询结果为空的情况
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("产品不存在")
	}
	// 处理其他查询错误
	if err != nil {
		return nil, fmt.Errorf("查询产品失败: %w", err)
	}
	
	return &product, nil
}

// GetProductsByCategory：按分类查询产品
// 参数：category - 产品分类
// 返回值：[]Product - 产品列表；error - 可能的错误
func (dm *DatabaseManager) GetProductsByCategory(category string) ([]Product, error) {
	// SQL查询语句，按分类查询产品
	query := `SELECT id, name, price, category, description, created_at, updated_at FROM products WHERE category = ?`
	
	// db.Query()：执行查询并返回多行结果（*sql.Rows）
	rows, err := dm.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("查询产品失败: %w", err)
	}
	// defer rows.Close()：确保查询结果集在函数退出时关闭，释放资源
	defer rows.Close()
	
	var products []Product
	// rows.Next()：迭代结果集中的行
	for rows.Next() {
		var product Product
		// 将当前行数据映射到结构体
		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.Category,
			&product.Description, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描产品失败: %w", err)
		}
		products = append(products, product)
	}
	
	// 检查迭代过程中是否发生错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("行迭代错误: %w", err)
	}
	
	return products, nil
}

// UpdateProduct：更新产品信息
// 参数：product - 包含更新信息的Product结构体（必须包含ID）
// 返回值：error - 可能的错误
func (dm *DatabaseManager) UpdateProduct(product *Product) error {
	// SQL更新语句，更新产品信息并刷新updated_at
	query := `UPDATE products SET name = ?, price = ?, category = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	
	// 执行更新操作
	result, err := dm.db.Exec(query, product.Name, product.Price, product.Category, product.Description, product.ID)
	if err != nil {
		return fmt.Errorf("更新产品失败: %w", err)
	}
	
	// 获取受影响的行数，判断是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	// 受影响行数为0表示没有找到对应ID的产品
	if rowsAffected == 0 {
		return fmt.Errorf("产品不存在")
	}
	
	return nil
}

// DeleteProduct：删除产品
// 参数：id - 产品ID
// 返回值：error - 可能的错误
func (dm *DatabaseManager) DeleteProduct(id int) error {
	// SQL删除语句
	query := `DELETE FROM products WHERE id = ?`
	
	// 执行删除操作
	result, err := dm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除产品失败: %w", err)
	}
	
	// 检查是否有记录被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("产品不存在")
	}
	
	return nil
}

// 5. 订单相关操作
// 以下方法实现订单的创建和查询操作，包含事务处理

// CreateOrder：创建订单（包含订单项）
// 参数：order - 订单信息，包含订单项列表
// 返回值：error - 可能的错误
// 关键：使用事务确保订单和订单项要么同时创建，要么都不创建
func (dm *DatabaseManager) CreateOrder(order *Order) error {
	// db.Begin()：开始事务，返回*sql.Tx（事务对象）
	tx, err := dm.db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	
	// 事务处理逻辑：
	// 1. 创建订单
	// 2. 创建订单项
	// 3. 提交事务
	// 任何一步失败都需要回滚事务
	
	// 创建订单
	orderQuery := `INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`
	result, err := tx.Exec(orderQuery, order.UserID, order.Total, order.Status)
	if err != nil {
		tx.Rollback() // 失败时回滚事务
		return fmt.Errorf("创建订单失败: %w", err)
	}
	
	// 获取订单ID
	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // 失败时回滚事务
		return fmt.Errorf("获取订单ID失败: %w", err)
	}
	
	order.ID = int(orderID)
	
	// 创建订单项
	itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	for _, item := range order.Items {
		_, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback() // 失败时回滚事务
			return fmt.Errorf("创建订单项失败: %w", err)
		}
	}
	
	// tx.Commit()：提交事务，所有操作生效
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	
	return nil
}

// GetOrder：查询订单详情（包含订单项）
// 参数：id - 订单ID
// 返回值：*Order - 订单详情；error - 可能的错误
func (dm *DatabaseManager) GetOrder(id int) (*Order, error) {
	// 获取订单基本信息
	orderQuery := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE id = ?`
	
	var order Order
	// 查询订单基本信息
	err := dm.db.QueryRow(orderQuery, id).Scan(
		&order.ID, &order.UserID, &order.Total, &order.Status,
		&order.CreatedAt, &order.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("订单不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	
	// 获取订单项
	itemQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ?`
	
	rows, err := dm.db.Query(itemQuery, id)
	if err != nil {
		return nil, fmt.Errorf("查询订单项失败: %w", err)
	}
	defer rows.Close()
	
	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("扫描订单项失败: %w", err)
		}
		items = append(items, item)
	}
	
	// 关联订单项到订单
	order.Items = items
	return &order, nil
}

// 6. 高级查询
// GetOrdersByUser：查询用户的所有订单（包含订单项）
// 参数：userID - 用户ID
// 返回值：[]Order - 订单列表；error - 可能的错误
// 关键：使用LEFT JOIN关联订单和订单项表，一次性获取关联数据
func (dm *DatabaseManager) GetOrdersByUser(userID int) ([]Order, error) {
	// SQL查询：关联订单表和订单项表，按用户ID查询
	query := `
    SELECT o.id, o.user_id, o.total, o.status, o.created_at, o.updated_at,
           oi.id, oi.product_id, oi.quantity, oi.price
    FROM orders o
    LEFT JOIN order_items oi ON o.id = oi.order_id  -- 左连接，确保没有订单项的订单也能被查询到
    WHERE o.user_id = ?
    ORDER BY o.created_at DESC`  -- 按创建时间降序排列（最新的在前）
	
	rows, err := dm.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户订单失败: %w", err)
	}
	defer rows.Close()
	
	// 使用map临时存储订单，避免重复（因为一条订单可能对应多条订单项）
	orders := make(map[int]*Order)
	
	for rows.Next() {
		var orderID int
		var order Order
		// 使用sql.NullXXX类型处理可能为NULL的字段（当订单没有订单项时）
		var itemID sql.NullInt64
		var productID sql.NullInt64
		var quantity sql.NullInt64
		var price sql.NullFloat64
		
		// 扫描查询结果到变量
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Total, &order.Status,
			&order.CreatedAt, &order.UpdatedAt,
			&itemID, &productID, &quantity, &price,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描订单失败: %w", err)
		}
		
		// 如果订单不在map中，添加到map
		if _, exists := orders[order.ID]; !exists {
			orders[order.ID] = &order
			orders[order.ID].Items = []OrderItem{}  // 初始化订单项切片
		}
		
		// 如果订单项ID有效（非NULL），添加到订单的订单项列表
		if itemID.Valid {
			orders[order.ID].Items = append(orders[order.ID].Items, OrderItem{
				ID:        int(itemID.Int64),
				OrderID:   order.ID,
				ProductID: int(productID.Int64),
				Quantity:  int(quantity.Int64),
				Price:     price.Float64,
			})
		}
	}
	
	// 检查行迭代错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("行迭代错误: %w", err)
	}
	
	// 将map中的订单转换为切片返回
	result := make([]Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, *order)
	}
	
	return result, nil
}

// 7. 事务示例
// ProcessOrder：处理订单流程（包含业务逻辑的事务示例）
// 参数：userID - 用户ID；items - 订单项列表
// 返回值：*Order - 创建的订单；error - 可能的错误
// 功能：查询产品价格 -> 计算总价 -> 创建订单 -> 创建订单项（全流程事务保证）
func (dm *DatabaseManager) ProcessOrder(userID int, items []OrderItem) (*Order, error) {
	// 开始事务
	tx, err := dm.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("开始事务失败: %w", err)
	}
	
	// 计算订单总价（业务逻辑）
	var total float64
	for i := range items {
		var price float64
		// 查询产品当前价格（确保使用最新价格）
		query := `SELECT price FROM products WHERE id = ?`
		err := tx.QueryRow(query, items[i].ProductID).Scan(&price)
		if err != nil {
			tx.Rollback() // 失败回滚
			return nil, fmt.Errorf("获取产品价格失败: %w", err)
		}
		
		// 累加总价
		total += price * float64(items[i].Quantity)
		// 记录购买时的单价
		items[i].Price = price
	}
	
	// 创建订单
	order := &Order{
		UserID: userID,
		Total:  total,
		Status: "pending",
		Items:  items,
	}
	
	orderQuery := `INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`
	result, err := tx.Exec(orderQuery, order.UserID, order.Total, order.Status)
	if err != nil {
		tx.Rollback() // 失败回滚
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}
	
	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // 失败回滚
		return nil, fmt.Errorf("获取订单ID失败: %w", err)
	}
	
	order.ID = int(orderID)
	
	// 创建订单项
	itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	for _, item := range items {
		_, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback() // 失败回滚
			return nil, fmt.Errorf("创建订单项失败: %w", err)
		}
	}
	
	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}
	
	return order, nil
}

// 8. 数据库工具函数
// GetProductStats：获取产品统计信息
// 返回值：map[string]interface{} - 统计数据；error - 可能的错误
// 功能：产品总数、平均价格、分类分布等统计
func (dm *DatabaseManager) GetProductStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 产品总数
	var totalProducts int
	err := dm.db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&totalProducts)
	if err != nil {
		return nil, fmt.Errorf("查询产品总数失败: %w", err)
	}
	stats["total_products"] = totalProducts
	
	// 平均价格
	var avgPrice float64
	err = dm.db.QueryRow(`SELECT AVG(price) FROM products`).Scan(&avgPrice)
	if err != nil {
		return nil, fmt.Errorf("查询平均价格失败: %w", err)
	}
	stats["average_price"] = avgPrice
	
	// 分类统计（按分类分组计数）
	categoryQuery := `SELECT category, COUNT(*) FROM products GROUP BY category`
	rows, err := dm.db.Query(categoryQuery)
	if err != nil {
		return nil, fmt.Errorf("查询分类统计失败: %w", err)
	}
	defer rows.Close()
	
	categories := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("扫描分类统计失败: %w", err)
		}
		categories[category] = count
	}
	stats["categories"] = categories
	
	return stats, nil
}

// 9. 数据清理
// Cleanup：清空所有表数据（用于测试或重置）
// 注意：实际生产环境慎用！
func (dm *DatabaseManager) Cleanup() error {
	// 清理顺序：先清理子表（有外键关联的表），再清理主表
	tables := []string{"order_items", "orders", "products"}
	
	for _, table := range tables {
		_, err := dm.db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("清理表 %s 失败: %w", table, err)
		}
	}
	
	return nil
}

// 主函数：程序入口，演示数据库操作流程
func main() {
	fmt.Println("=== Go语言数据库操作 ===")
	
	// 创建数据库管理器，连接到ecommerce.db文件
	dm, err := NewDatabaseManager("ecommerce.db")
	if err != nil {
		log.Fatal("创建数据库管理器失败:", err)
	}
	// 延迟关闭数据库连接，确保程序退出时释放资源
	defer dm.Close()
	
	// 初始化数据库表结构
	if err := dm.InitializeSchema(); err != nil {
		log.Fatal("初始化数据库表失败:", err)
	}
	
	fmt.Println("数据库初始化完成")
	
	// 添加示例产品
	products := []Product{
		{Name: "iPhone 15", Price: 6999.99, Category: "Electronics", Description: "最新款苹果手机"},
		{Name: "MacBook Pro", Price: 12999.99, Category: "Electronics", Description: "专业级笔记本电脑"},
		{Name: "Nike Air Max", Price: 899.99, Category: "Shoes", Description: "舒适运动鞋"},
		{Name: "Coffee Maker", Price: 299.99, Category: "Home", Description: "全自动咖啡机"},
	}
	
	for _, product := range products {
		if err := dm.CreateProduct(&product); err != nil {
			log.Printf("创建产品失败: %v", err)
		} else {
			fmt.Printf("创建产品: %s (ID: %d)\n", product.Name, product.ID)
		}
	}
	
	// 查询单个产品
	if product, err := dm.GetProduct(1); err != nil {
		log.Printf("查询产品失败: %v", err)
	} else {
		fmt.Printf("查询到产品: %+v\n", product)
	}
	
	// 按分类查询产品
	if electronics, err := dm.GetProductsByCategory("Electronics"); err != nil {
		log.Printf("查询电子产品失败: %v", err)
	} else {
		fmt.Printf("电子产品数量: %d\n", len(electronics))
		for _, product := range electronics {
			fmt.Printf("  - %s: $%.2f\n", product.Name, product.Price)
		}
	}
	
	// 创建订单（通过ProcessOrder处理完整流程）
	orderItems := []OrderItem{
		{ProductID: 1, Quantity: 2},  // 购买2个ID=1的产品（iPhone 15）
		{ProductID: 3, Quantity: 1},  // 购买1个ID=3的产品（Nike Air Max）
	}
	
	order, err := dm.ProcessOrder(1, orderItems)  // 用户ID=1
	if err != nil {
		log.Printf("创建订单失败: %v", err)
	} else {
		fmt.Printf("创建订单成功: ID=%d, 总价=%.2f\n", order.ID, order.Total)
	}
	
	// 查询订单详情
	if order, err := dm.GetOrder(order.ID); err != nil {
		log.Printf("查询订单失败: %v", err)
	} else {
		fmt.Printf("订单详情:\n")
		fmt.Printf("  订单ID: %d\n", order.ID)
		fmt.Printf("  用户ID: %d\n", order.UserID)
		fmt.Printf("  总价: %.2f\n", order.Total)
		fmt.Printf("  状态: %s\n", order.Status)
		fmt.Printf("  订单项:\n")
		for _, item := range order.Items {
			fmt.Printf("    - 产品ID: %d, 数量: %d, 价格: %.2f\n", 
				item.ProductID, item.Quantity, item.Price)
		}
	}
	
	// 获取产品统计信息
	if stats, err := dm.GetProductStats(); err != nil {
		log.Printf("获取统计信息失败: %v", err)
	} else {
		fmt.Printf("\n产品统计信息:\n")
		// 将统计信息转换为格式化的JSON输出
		statsJSON, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(statsJSON))
	}
	
	fmt.Println("\n数据库操作演示完成")
	fmt.Println("\n练习：")
	fmt.Println("1. 为用户表添加更多字段（如地址、电话等）")
	fmt.Println("2. 实现用户注册和登录功能")
	fmt.Println("3. 添加产品搜索功能（按名称模糊搜索）")
	fmt.Println("4. 实现订单状态更新功能")
	fmt.Println("5. 添加产品库存管理")
	fmt.Println("6. 实现分页查询功能")
	fmt.Println("7. 添加数据库连接池配置")
	fmt.Println("8. 实现数据迁移脚本")
}
    