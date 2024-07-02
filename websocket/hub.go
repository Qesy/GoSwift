package websocket

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Qesy/qesygo"
	"golang.org/x/net/websocket"
)

// Client 客户端结构
type Client struct {
	Conn   *websocket.Conn
	UserID string
	ConnTs int64
	IP     string
}

// ClientMsg 消息结构体
type ClientMsg struct {
	UserID string
	Msg    string
}

// Hub 结构
type Hub struct {
	Clients    map[string]*Client
	BroadCast  chan string
	Register   chan *Client
	UnRegister chan string
	ClientMsg  chan *ClientMsg
	GetData    chan string //获取数据类型
}

// HubRouter 路由
var HubRouter = newHub()

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		BroadCast:  make(chan string),
		Register:   make(chan *Client),
		UnRegister: make(chan string),
		ClientMsg:  make(chan *ClientMsg),
		GetData:    make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:

			if OldClient, ok := h.Clients[client.UserID]; ok {
				fmt.Println("client.conn", client.Conn, OldClient.Conn)
				if client.Conn != OldClient.Conn {
					fmt.Println("账号被顶掉")
					client.ConnTs = OldClient.ConnTs //让登录时间继续保持之前的登录时间
					HubSend(OldClient.Conn, "Offline")
					OldClient.Conn.Close()
				}
			}
			h.Clients[client.UserID] = client

		case UserID := <-h.UnRegister:
			if client, ok := h.Clients[UserID]; ok {
				delete(h.Clients, UserID)
				client.Conn.Close()
				connectTime := qesygo.Time("Millisecond") - client.ConnTs // 客户端在线时间
				fmt.Printf("User Online %v \n", connectTime)

			}
		case message := <-h.BroadCast: //广播
			for _, usClient := range h.Clients {
				HubSend(usClient.Conn, message)
			}
		case clientMsg := <-h.ClientMsg:
			if Client, ok := h.Clients[clientMsg.UserID]; ok {
				HubSend(Client.Conn, clientMsg.Msg)
			}
		case getData := <-h.GetData: // 获取在线用户
			if getData == "UserCount" {
				Str <- strconv.Itoa(len(h.Clients))
			} else if getData == "UserList" {
				UserIDs := []string{}
				for k := range h.Clients {
					UserIDs = append(UserIDs, k)
				}
				Str <- strings.Join(UserIDs, "|")
			}
		}

	}
}
