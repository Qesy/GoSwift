package websocket

import (
	"fmt"
	"log"
	"reflect"
	"server/models"
	"strings"
	"sync"
	"time"

	"github.com/Qesy/qesydb"
	"github.com/Qesy/qesygo"
	"golang.org/x/net/websocket"
)

// Entry 结构
type Entry struct {
	Conn    *websocket.Conn
	Act     string
	UID     string
	ProtoId int
	Data    map[string]string
}

// MsgStr 消息结构体
type MsgStr struct {
	Act  string
	Data map[string]string
}

// 结构权限结构体
type Api struct {
	Method     string
	Permission []string
}

// Deadline Ping值
var Deadline = 30
var WsLock sync.Locker

// 登录平台对接
const PLATFORM_WECHAT = "wechat"
const PLATFORM_BYTE = "byte"

var Str chan string = make(chan string)

// Echo 运行
func Echo(conn *websocket.Conn) {
	e := &Entry{Conn: conn}
	e.SetDeadline()
	i := 0
	for {
		i++
		reply := []byte{}
		if err := websocket.Message.Receive(e.Conn, &reply); err != nil {
			fmt.Println("Receive Err: ", err)
			e.Unregister()
			break
		}

		var replyArr MsgStr
		//fmt.Println("reply", reply)
		if err := qesygo.JsonDecode([]byte(reply), &replyArr); err != nil {
			fmt.Println("Receive Err : ", err)
			e.SendError(4000)
			continue
		}

		// buf := buffer.NewBuffer(reply)
		// ProtoId := buf.ReadNInt16()
		// MethodName, ok := packets.MsgSend_value[ProtoId]
		// if ProtoId != 1001 { //心跳包打印隐藏
		// 	fmt.Println("reply", ProtoId, MethodName, string(reply))
		// }
		// if !ok { //协议未定义
		// 	log.Printf("Interface Not Defined, ProtoId : %v MethodName %v", ProtoId, MethodName)
		// 	e.SendError(4007)
		// 	continue
		// }
		e.Act, e.Data = replyArr.Act, replyArr.Data
		DbStats := qesydb.Db.Stats()
		if DbStats.OpenConnections > DbStats.MaxOpenConnections/2 { // Mysql链接数大于等于最大连接数-10，报错（服务器繁忙）
			e.SendError(5000)
			continue
		}
		if m, ok := reflect.TypeOf(e).MethodByName(e.Act); ok {
			if !e.isLogin() {
				log.Printf("Is Not Login , Act : %v Data %v", e.Act, e.Data)
				e.SendError(4008)
			} else {
				m.Func.Call([]reflect.Value{reflect.ValueOf(e)})
			}
			continue
		} else { //接口未定义
			log.Printf("MethodName Not Defined, Act : %v Data %v", e.Act, e.Data)
			e.SendError(4006)
			continue
		}

	}
}

func (e *Entry) isLogin() bool { //是否登录（未登录支允许访问登录接口）
	// if e.UID == "" && e.ProtoId != packets.CLogin_Login {
	// 	return false
	// }
	return true
}

// func (e *Entry) Permission() bool {
// 	Api := Api{}
// 	for _, v := range ApiArr {
// 		if v.Method == e.Act {
// 			Api = v
// 		}
// 	}
// 	if Api.Method == "" { //无定义此接口
// 		return false
// 	}
// 	if qesygo.InArray(Api.Permission, ROLE_EVERYONE) { // 所有人都可以访问
// 		return true
// 	}
// 	if _, ok := HubRouter.clients[e.UID]; ok {
// 		return true
// 	}
// 	return false
// }

// 获取IP
func (e *Entry) GetIp() string {
	Ip := e.Conn.Request().Header.Get("X-Forwarded-For")
	if strings.Contains(Ip, "127.0.0.1") || Ip == "" {
		Ip = e.Conn.Request().Header.Get("X-real-ip")
	}
	if Ip == "" {
		Ip = "127.0.0.1"
	}
	return Ip
}

// SetDeadline 设置PING值
func (e *Entry) SetDeadline() *Entry {
	t := time.Now().Add(time.Duration(Deadline) * time.Second)
	e.Conn.SetDeadline(t)
	return e
}

// Register 用户唯一登录口
func (e *Entry) Register(UserID string) {
	ConnTs := qesygo.Time("Millisecond")
	e.UID = UserID
	Client := Client{UserID: e.UID, Conn: e.Conn, ConnTs: ConnTs, IP: e.GetIp()}
	fmt.Println("Login : Register", Client)
	HubRouter.Register <- &Client

}

// Unregister 用户退出
func (e *Entry) Unregister() {
	fmt.Println("LogOut:", e.UID)
	if e.UID != "" {
		HubRouter.UnRegister <- e.UID
	} else { //未登录直接关闭链接
		e.Conn.Close()
	}
}

func ErrorDescGet(ErrCode int32) string {
	return models.StaticGetByField("codeError", fmt.Sprint(ErrCode), "des")
}

// SendError 发送错误
func (e *Entry) SendError(ErrCode int32) *Entry {
	Content := models.StaticGetByField("codeError", fmt.Sprint(ErrCode), "Content")
	e.Send(Content)
	return e
}

// SendError 发送错误
func (e *Entry) SendErrorDesc(ErrCode int32, Msg string) *Entry {
	e.Send(Msg)
	return e
}

// Send 发送消息
func (e *Entry) Send(Msg string) *Entry {
	Print("Single", Msg, []string{e.UID})
	if e.UID == "" { //未登录直接发
		Result := HubSend(e.Conn, Msg)
		fmt.Println("SendResult", Result)
	} else {
		HubRouter.ClientMsg <- &ClientMsg{UserID: e.UID, Msg: Msg}
	}

	return e
}

// Send 发送消息 指定UID
func SendError(UserID string, ErrCode int32) {
	Msg := ""
	HubRouter.ClientMsg <- &ClientMsg{UserID: UserID, Msg: Msg}
}

// Send 发送消息 指定UID
func Send(UserID string, Msg string) {
	Print("Single", Msg, []string{UserID})
	HubRouter.ClientMsg <- &ClientMsg{UserID: UserID, Msg: Msg}
}

// SendMultiple 指定多发
func SendMultiple(IDArr []string, Msg string) {
	for _, v := range IDArr {
		HubRouter.ClientMsg <- &ClientMsg{UserID: v, Msg: Msg}
	}
}

// HubSend 底层发送
func HubSend(conn *websocket.Conn, msg string) error {
	err := websocket.Message.Send(conn, msg)
	if err != nil {
		log.Println("HubSend", err, msg)
	}
	return err
}

// Broadcast 广播
func Broadcast(Msg string) {
	//msgByte, _ := qesygo.JsonEncode(msgArr)
	Print("Broadcast", Msg, []string{})
	HubRouter.BroadCast <- Msg
}

func Unregister(UserId string) { //踢出
	fmt.Println("LogOut:", UserId)
	HubRouter.UnRegister <- UserId

}

// Print 发送打印
func Print(SendType string, Msg string, userIDArr []string) {
	UIDStr := qesygo.Implode(userIDArr, ",")
	switch SendType {
	case "Single", "Multiple":
		fmt.Println("SOCKETSEND:(", "Type:"+SendType, ", UidArr:"+UIDStr, ")"+" \nData:", Msg)
	case "Broadcast":
		fmt.Println("SOCKETSEND:(", "Type:"+SendType, ")"+" \nData:", Msg)
	}
}

// RecPrint 接收打印
func RecPrint(str string, str2 string, str3 map[string]string, str4 string) {
	fmt.Println(str, str2, str3, str4)
}

// UserCount 统计用户
func UserCount() int {
	return len(HubRouter.Clients)

}

func UserOnline() []string {
	UserIDs := []string{}
	for k := range HubRouter.Clients {
		UserIDs = append(UserIDs, k)
	}
	return UserIDs

}
