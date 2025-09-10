package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	room string
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Room string          `json:"room"`
}

type Room struct {
	clients []*Client
	mu      sync.Mutex
}

var rooms = make(map[string]*Room)
var roomsMu sync.Mutex

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("服务器启动在 :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("读取消息失败:", err)
			removeClient(client)
			break
		}

		switch msg.Type {
		case "join":
			handleJoin(client, msg)
		case "offer", "answer", "ice-candidate":
			broadcastToRoom(client, msg)
		case "leave":
			removeClient(client)
		default:
			log.Println("未知的消息类型:", msg.Type)
		}
	}
}

func handleJoin(client *Client, msg Message) {
	var data struct {
		Room string `json:"room"`
	}

	if err := json.Unmarshal(msg.Data, &data); err != nil {
		log.Println("解析join消息失败:", err)
		return
	}

	client.room = data.Room
	addClientToRoom(client, data.Room)

	response := map[string]interface{}{
		"type": "room-joined",
		"room": data.Room,
	}

	client.conn.WriteJSON(response)
}

func addClientToRoom(client *Client, roomID string) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		room = &Room{}
		rooms[roomID] = room
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for _, c := range room.clients {
		if c == client {
			return
		}
	}

	room.clients = append(room.clients, client)
	log.Printf("客户端加入房间: %s, 当前人数: %d", roomID, len(room.clients))
}

func removeClient(client *Client) {
	if client.room == "" {
		return
	}

	roomsMu.Lock()
	room, exists := rooms[client.room]
	roomsMu.Unlock()

	if !exists {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for i, c := range room.clients {
		if c == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			log.Printf("客户端离开房间: %s, 剩余人数: %d", client.room, len(room.clients))
			break
		}
	}

	// 如果房间为空，删除房间
	if len(room.clients) == 0 {
		roomsMu.Lock()
		delete(rooms, client.room)
		roomsMu.Unlock()
	}
}

func broadcastToRoom(sender *Client, msg Message) {
	if sender.room == "" {
		return
	}

	roomsMu.Lock()
	room, exists := rooms[sender.room]
	roomsMu.Unlock()

	if !exists {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for _, client := range room.clients {
		if client != sender {
			// 设置写入超时
			client.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Println("发送消息失败:", err)
				// 如果发送失败，从房间中移除客户端
				go removeClient(client)
			}
		}
	}
}
