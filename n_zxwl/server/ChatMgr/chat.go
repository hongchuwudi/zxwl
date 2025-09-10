package chat

// 导入所需的包
import (
	"database/sql"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"mymod/new_zxwl/config"
	"net/http"
	"time"
)

// Dbd 是一个指向 sql.DB 的全局变量，用于数据库操作
var Dbd *sql.DB

// InitDB 初始化数据库连接，设置数据源并测试连接
func InitDB() {
	var err error
	// 定义数据库连接字符串
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	//dsn := "root:hongchu@tcp(localhost:3306)/chatroom?charset=utf8mb4&parseTime=True"
	// 打开数据库连接
	Dbd, err = sql.Open("mysql", dsn)
	if err != nil {
		// 若连接失败，记录错误日志并终止程序
		log.Fatal("数据库连接失败:", err)
	}

	// 测试数据库连接
	err = Dbd.Ping()
	if err != nil {
		// 若测试失败，记录错误日志并终止程序
		log.Fatal("数据库连接测试失败:", err)
	}
}

// upgrader 用于将 HTTP 连接升级为 WebSocket 连接，允许所有来源的请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client 表示一个 WebSocket 客户端，包含连接、学校 ID 和邮箱信息
type Client struct {
	conn     *websocket.Conn // WebSocket 连接
	schoolID int             // 学校 ID
	email    string          // 用户邮箱
}

// clients 是一个映射，存储每个学校 ID 对应的客户端集合
var clients = make(map[int]map[*Client]bool)

// Broadcast 是一个消息通道，用于广播消息
var Broadcast = make(chan Message)

// WsMessage 表示从 WebSocket 接收到的消息结构
type WsMessage struct {
	Action   string `json:"action"`    // 消息动作，如 "join" 或 "message"
	SchoolID int    `json:"school_id"` // 学校 ID
	Content  string `json:"content"`   // 消息内容
	Email    string `json:"email"`     // 用户邮箱
}

// ChatHandler 处理 WebSocket 连接请求，处理客户端的加入和消息发送
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// 若升级失败，记录错误日志并返回
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	// 函数结束时关闭 WebSocket 连接
	defer conn.Close()

	var client *Client

	for {
		var msg WsMessage
		// 从 WebSocket 连接读取 JSON 消息
		err := conn.ReadJSON(&msg)
		if err != nil {
			// 若读取失败，跳出循环
			break
		}

		switch msg.Action {
		case "join":
			// 客户端加入处理
			client = &Client{
				conn:     conn,
				schoolID: msg.SchoolID,
				email:    msg.Email,
			}

			if clients[msg.SchoolID] == nil {
				// 若该学校 ID 对应的客户端集合不存在，则创建一个新的集合
				clients[msg.SchoolID] = make(map[*Client]bool)
			}
			// 将客户端添加到对应的学校 ID 集合中
			clients[msg.SchoolID][client] = true

			// 异步发送历史消息给客户端
			go sendHistory(conn, msg.SchoolID)

		case "message":
			if client == nil {
				// 若客户端未加入，跳过处理
				continue
			}

			// 将消息插入到数据库中
			_, err := Dbd.Exec(
				"INSERT INTO chat_messages (school_id, email, content) VALUES (?, ?, ?)",
				client.schoolID, client.email, msg.Content,
			)
			if err != nil {
				// 若插入失败，记录错误日志并跳过
				log.Println("Insert message error:", err)
				continue
			}

			var name, picture string
			// 从数据库中查询用户的姓名和头像
			Dbd.QueryRow(
				"SELECT name, picture FROM profile WHERE email = ?",
				client.email,
			).Scan(&name, &picture)

			// 将消息发送到广播通道
			Broadcast <- Message{
				SchoolID:  client.schoolID,
				Email:     client.email,
				Content:   msg.Content,
				CreatedAt: time.Now(),
			}
		}
	}

	if client != nil {
		// 客户端断开连接时，从集合中移除该客户端
		delete(clients[client.schoolID], client)
	}
}

// sendHistory 向客户端发送指定学校 ID 的历史聊天消息
func sendHistory(conn *websocket.Conn, schoolID int) {
	// 存储历史消息的切片
	messages := make([]map[string]interface{}, 0)

	// 从数据库中查询历史消息
	rows, err := Dbd.Query(`
        SELECT 
            m.content, 
            m.created_at,
            COALESCE(NULLIF(p.name, ''), '匿名用户') as name,
            COALESCE(NULLIF(p.picture, ''), '/default-avatar.png') as picture
        FROM chat_messages m
        LEFT JOIN profile p ON m.email = p.email
        WHERE m.school_id = ?
        ORDER BY m.created_at DESC
        LIMIT 50`, schoolID)

	if err != nil {
		// 若查询失败，记录错误日志并发送空消息给客户端
		log.Printf("历史消息查询失败: %v", err)
		conn.WriteJSON(map[string]interface{}{
			"action":   "history",
			"messages": messages,
		})
		return
	}
	// 函数结束时关闭查询结果集
	defer rows.Close()

	for rows.Next() {
		var content string
		var createdAt time.Time
		var name, picture string

		// 扫描查询结果的每一行
		if err := rows.Scan(&content, &createdAt, &name, &picture); err != nil {
			// 若扫描失败，记录错误日志并跳过
			log.Printf("行扫描失败: %v", err)
			continue
		}

		// 将消息添加到切片中
		messages = append(messages, map[string]interface{}{
			"content":    content,
			"created_at": createdAt.Format("2006-01-02 15:04:05"),
			"name":       name,
			"picture":    picture,
		})
	}

	if err = rows.Err(); err != nil {
		// 若结果集遍历出错，记录错误日志
		log.Printf("结果集遍历错误: %v", err)
	}

	// 将历史消息发送给客户端
	conn.WriteJSON(map[string]interface{}{
		"action":   "history",
		"messages": messages,
	})
}

// HandleMessages 处理广播通道中的消息，并将消息发送给所有相关客户端
func HandleMessages() {
	for msg := range Broadcast {
		// 记录准备广播的消息信息
		log.Printf("准备广播消息: SchoolID=%d, Content=%s", msg.SchoolID, msg.Content)
		schoolClients := clients[msg.SchoolID]
		if schoolClients == nil {
			// 若该学校 ID 没有客户端监听，记录日志并跳过
			log.Printf("当前没有客户端监听学校ID: %d", msg.SchoolID)
			continue
		}

		// 格式化消息创建时间
		formattedTime := msg.CreatedAt.Format("2006-01-02 15:04")
		// 构建要发送的消息数据
		messageData := map[string]interface{}{
			"action":     "message",
			"schoolID":   msg.SchoolID,
			"email":      msg.Email,
			"content":    msg.Content,
			"name":       msg.Name,
			"picture":    msg.Picture,
			"time":       formattedTime,
			"created_at": formattedTime,
		}

		if messageData["name"] == "" {
			// 若用户姓名为空，设置为匿名用户
			messageData["name"] = "匿名用户"
		}
		if messageData["picture"] == "" {
			// 若用户头像为空，设置为默认头像
			messageData["picture"] = "/default-avatar.png"
		}

		for client := range schoolClients {
			go func(c *Client) {
				if err := c.conn.WriteJSON(messageData); err != nil {
					log.Printf("消息发送失败: %v", err)
					c.conn.Close()
					delete(schoolClients, c)
				}
			}(client)
		}
	}
}
