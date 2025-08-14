package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// 数据库操作
// Go语言的数据库编程和ORM概念

// 1. 数据模型
type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    Category    string    `json:"category"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
    ID         int       `json:"id"`
    UserID     int       `json:"user_id"`
    Total      float64   `json:"total"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    Items      []OrderItem `json:"items"`
}

type OrderItem struct {
    ID        int     `json:"id"`
    OrderID   int     `json:"order_id"`
    ProductID int     `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// 2. 数据库管理器
type DatabaseManager struct {
    db *sql.DB
}

// 创建新的数据库管理器
func NewDatabaseManager(dbPath string) (*DatabaseManager, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("无法打开数据库: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("无法连接到数据库: %w", err)
    }
    
    return &DatabaseManager{db: db}, nil
}

// 关闭数据库连接
func (dm *DatabaseManager) Close() error {
    return dm.db.Close()
}

// 3. 初始化数据库表
func (dm *DatabaseManager) InitializeSchema() error {
    // 创建产品表
    productTable := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        category TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    // 创建订单表
    orderTable := `
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        total REAL NOT NULL,
        status TEXT NOT NULL DEFAULT 'pending',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    // 创建订单项表
    orderItemTable := `
    CREATE TABLE IF NOT EXISTS order_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        order_id INTEGER NOT NULL,
        product_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL,
        price REAL NOT NULL,
        FOREIGN KEY (order_id) REFERENCES orders(id),
        FOREIGN KEY (product_id) REFERENCES products(id)
    );`
    
    // 创建索引
    createIndexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);",
        "CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);",
        "CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);",
        "CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);",
    }
    
    // 执行创建表的SQL
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
func (dm *DatabaseManager) CreateProduct(product *Product) error {
    query := `INSERT INTO products (name, price, category, description) VALUES (?, ?, ?, ?)`
    
    result, err := dm.db.Exec(query, product.Name, product.Price, product.Category, product.Description)
    if err != nil {
        return fmt.Errorf("创建产品失败: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("获取产品ID失败: %w", err)
    }
    
    product.ID = int(id)
    return nil
}

func (dm *DatabaseManager) GetProduct(id int) (*Product, error) {
    query := `SELECT id, name, price, category, description, created_at, updated_at FROM products WHERE id = ?`
    
    var product Product
    err := dm.db.QueryRow(query, id).Scan(
        &product.ID, &product.Name, &product.Price, &product.Category,
        &product.Description, &product.CreatedAt, &product.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("产品不存在")
    }
    if err != nil {
        return nil, fmt.Errorf("查询产品失败: %w", err)
    }
    
    return &product, nil
}

func (dm *DatabaseManager) GetProductsByCategory(category string) ([]Product, error) {
    query := `SELECT id, name, price, category, description, created_at, updated_at FROM products WHERE category = ?`
    
    rows, err := dm.db.Query(query, category)
    if err != nil {
        return nil, fmt.Errorf("查询产品失败: %w", err)
    }
    defer rows.Close()
    
    var products []Product
    for rows.Next() {
        var product Product
        err := rows.Scan(
            &product.ID, &product.Name, &product.Price, &product.Category,
            &product.Description, &product.CreatedAt, &product.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("扫描产品失败: %w", err)
        }
        products = append(products, product)
    }
    
    return products, nil
}

func (dm *DatabaseManager) UpdateProduct(product *Product) error {
    query := `UPDATE products SET name = ?, price = ?, category = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
    
    result, err := dm.db.Exec(query, product.Name, product.Price, product.Category, product.Description, product.ID)
    if err != nil {
        return fmt.Errorf("更新产品失败: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("获取影响行数失败: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("产品不存在")
    }
    
    return nil
}

func (dm *DatabaseManager) DeleteProduct(id int) error {
    query := `DELETE FROM products WHERE id = ?`
    
    result, err := dm.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("删除产品失败: %w", err)
    }
    
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
func (dm *DatabaseManager) CreateOrder(order *Order) error {
    tx, err := dm.db.Begin()
    if err != nil {
        return fmt.Errorf("开始事务失败: %w", err)
    }
    
    // 创建订单
    orderQuery := `INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`
    result, err := tx.Exec(orderQuery, order.UserID, order.Total, order.Status)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("创建订单失败: %w", err)
    }
    
    orderID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("获取订单ID失败: %w", err)
    }
    
    order.ID = int(orderID)
    
    // 创建订单项
    itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
    for _, item := range order.Items {
        _, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("创建订单项失败: %w", err)
        }
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("提交事务失败: %w", err)
    }
    
    return nil
}

func (dm *DatabaseManager) GetOrder(id int) (*Order, error) {
    // 获取订单基本信息
    orderQuery := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE id = ?`
    
    var order Order
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
    
    order.Items = items
    return &order, nil
}

// 6. 高级查询
func (dm *DatabaseManager) GetOrdersByUser(userID int) ([]Order, error) {
    query := `
    SELECT o.id, o.user_id, o.total, o.status, o.created_at, o.updated_at,
           oi.id, oi.product_id, oi.quantity, oi.price
    FROM orders o
    LEFT JOIN order_items oi ON o.id = oi.order_id
    WHERE o.user_id = ?
    ORDER BY o.created_at DESC`
    
    rows, err := dm.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("查询用户订单失败: %w", err)
    }
    defer rows.Close()
    
    orders := make(map[int]*Order)
    
    for rows.Next() {
        var orderID int
        var order Order
        var itemID sql.NullInt64
        var productID sql.NullInt64
        var quantity sql.NullInt64
        var price sql.NullFloat64
        
        err := rows.Scan(
            &order.ID, &order.UserID, &order.Total, &order.Status,
            &order.CreatedAt, &order.UpdatedAt,
            &itemID, &productID, &quantity, &price,
        )
        if err != nil {
            return nil, fmt.Errorf("扫描订单失败: %w", err)
        }
        
        if _, exists := orders[order.ID]; !exists {
            orders[order.ID] = &order
            orders[order.ID].Items = []OrderItem{}
        }
        
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
    
    result := make([]Order, 0, len(orders))
    for _, order := range orders {
        result = append(result, *order)
    }
    
    return result, nil
}

// 7. 事务示例
func (dm *DatabaseManager) ProcessOrder(userID int, items []OrderItem) (*Order, error) {
    tx, err := dm.db.Begin()
    if err != nil {
        return nil, fmt.Errorf("开始事务失败: %w", err)
    }
    
    // 计算订单总价
    var total float64
    for _, item := range items {
        var price float64
        query := `SELECT price FROM products WHERE id = ?`
        err := tx.QueryRow(query, item.ProductID).Scan(&price)
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("获取产品价格失败: %w", err)
        }
        
        total += price * float64(item.Quantity)
        item.Price = price
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
        tx.Rollback()
        return nil, fmt.Errorf("创建订单失败: %w", err)
    }
    
    orderID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("获取订单ID失败: %w", err)
    }
    
    order.ID = int(orderID)
    
    // 创建订单项
    itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
    for _, item := range items {
        _, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
        if err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("创建订单项失败: %w", err)
        }
    }
    
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("提交事务失败: %w", err)
    }
    
    return order, nil
}

// 8. 数据库工具函数
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
    
    // 分类统计
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
func (dm *DatabaseManager) Cleanup() error {
    tables := []string{"order_items", "orders", "products"}
    
    for _, table := range tables {
        _, err := dm.db.Exec(fmt.Sprintf("DELETE FROM %s", table))
        if err != nil {
            return fmt.Errorf("清理表 %s 失败: %w", table, err)
        }
    }
    
    return nil
}

func main() {
    fmt.Println("=== Go语言数据库操作 ===")
    
    // 创建数据库管理器
    dm, err := NewDatabaseManager("ecommerce.db")
    if err != nil {
        log.Fatal("创建数据库管理器失败:", err)
    }
    defer dm.Close()
    
    // 初始化数据库表
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
    
    // 查询产品
    if product, err := dm.GetProduct(1); err != nil {
        log.Printf("查询产品失败: %v", err)
    } else {
        fmt.Printf("查询到产品: %+v\n", product)
    }
    
    // 按分类查询
    if electronics, err := dm.GetProductsByCategory("Electronics"); err != nil {
        log.Printf("查询电子产品失败: %v", err)
    } else {
        fmt.Printf("电子产品数量: %d\n", len(electronics))
        for _, product := range electronics {
            fmt.Printf("  - %s: $%.2f\n", product.Name, product.Price)
        }
    }
    
    // 创建订单
    orderItems := []OrderItem{
        {ProductID: 1, Quantity: 2, Price: 6999.99},
        {ProductID: 3, Quantity: 1, Price: 899.99},
    }
    
    order, err := dm.ProcessOrder(1, orderItems)
    if err != nil {
        log.Printf("创建订单失败: %v", err)
    } else {
        fmt.Printf("创建订单成功: ID=%d, 总价=%.2f\n", order.ID, order.Total)
    }
    
    // 查询订单
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
    
    // 获取统计信息
    if stats, err := dm.GetProductStats(); err != nil {
        log.Printf("获取统计信息失败: %v", err)
    } else {
        fmt.Printf("\n产品统计信息:\n")
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