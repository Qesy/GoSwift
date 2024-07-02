package websocket

import (
	"context"

	"fmt"
	"server/models"
	"sync"
	"time"

	"github.com/Qesy/qesygo"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMsgResult struct { //消息结构体
	Act      string
	Data     map[string]string
	ServerID string
}

type RabbitConnection struct {
	URL        string
	Conn       *amqp.Connection
	Ch         *amqp.Channel
	Queue      amqp.Queue
	L          sync.Locker
	Messages   <-chan amqp.Delivery
	ConnErrNum int //链接失败数字大于3，休息3秒钟
	ServerID   string
}

var RabbitStr = &RabbitConnection{
	URL:        "",
	Conn:       &amqp.Connection{},
	Ch:         &amqp.Channel{},
	Queue:      amqp.Queue{},
	L:          &sync.Mutex{},
	Messages:   make(<-chan amqp.Delivery),
	ConnErrNum: 0,
	ServerID:   "",
}

func RabbitInit(Host string, Port string, User string, Password string, ServerID string) {
	RabbitStr = &RabbitConnection{
		URL:   "amqp://" + User + ":" + Password + "@" + Host + ":" + Port + "/",
		Conn:  &amqp.Connection{},
		Ch:    &amqp.Channel{},
		Queue: amqp.Queue{},
		L:     &sync.Mutex{},
		//Messages:   <-chan amqp.Delivery,
		ConnErrNum: 0,
		ServerID:   ServerID,
	}
	RabbitStr.ConnMq()
}

func (R *RabbitConnection) ConnMq() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if conn, err := amqp.Dial(R.URL); err != nil {
		qesygo.FailOnError(err, "Failed to Dial")
	} else {
		fmt.Println("Rabbit Connection Success !")
		R.Conn = conn
	}
	//defer R.Conn.Close() //不关闭了
	if ch, err := R.Conn.Channel(); err != nil {
		qesygo.FailOnError(err, "Failed to Channel")
	} else {
		R.Ch = ch
	}
	//defer R.Ch.Close() //不关闭了
	if err := R.Ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		qesygo.FailOnError(err, "Failed to ExchangeDeclare")
	}

	if q, err := R.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		qesygo.FailOnError(err, "Failed to declare a queue")
	} else {
		R.Queue = q
	}

	if err := R.Ch.QueueBind(
		R.Queue.Name, // queue name
		"",           // routing key
		"logs",       // exchange
		false,
		nil,
	); err != nil {
		qesygo.FailOnError(err, "Failed to bind a queue")
	}

	if msgs, err := R.Ch.Consume(
		R.Queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	); err != nil {
		qesygo.FailOnError(err, "Failed to register a consumer")
	} else {
		R.Messages = msgs
	}
	go R.Receive() //开启接收

}

func (R *RabbitConnection) reConn() { //重连
	defer func() {
		R.L.Unlock()
	}()
	R.L.Lock()

	if R.Conn.IsClosed() {
		R.ConnErrNum++
		if R.ConnErrNum >= 3 {
			time.Sleep(3 * time.Second) //休息1秒
			R.ConnErrNum = 0
		}

		R.ConnMq() //重连
	}

}

func (R *RabbitConnection) PublishMsg(Msg RabbitMsgResult) error {
	defer func() {
		if err := recover(); err != nil {
			if R.Conn.IsClosed() {
				R.reConn()
			}

		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	MsgStr, _ := qesygo.JsonEncode(Msg)
	err := R.Ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        MsgStr,
		})
	if err != nil {
		qesygo.FailOnError(err, "Failed to PublishWithContext")
	}
	if Msg.Act != "Ping" { //心跳太烦了，屏蔽
		fmt.Printf(" [x] Sent %s \n", Msg)
	}

	return err
}

func (R *RabbitConnection) Receive() {
	for d := range R.Messages {
		Result := &RabbitMsgResult{}
		qesygo.JsonDecode(d.Body, Result)
		if Result.Act != "Ping" { //心跳太烦了，屏蔽
			fmt.Printf(" [x] Receive %s \n", *Result)
		}

		if Result.Act == "RepeatLogin" { //重复登录
			rabbitRepeatLogin(Result)
		} else if Result.Act == "Broadcast" { //全平台广播
			rabbitBroadcast(Result)
		} else if Result.Act == "heartbeat" { //心跳
			fmt.Println("RabbitMQ " + Result.Act)
		} else if Result.Act == "FlushItem" { //刷新装备，人物属性，道具
			rabbitFlushItem(Result)
		}

	}
	fmt.Println("关闭接收")
}

func rabbitRepeatLogin(Result *RabbitMsgResult) { //重复登陆了处理(跨服被顶掉，本服务的不处理)
	if !qesygo.VeriPara(Result.Data, []string{"UserID"}) { // 缺少参数不处理
		fmt.Println("rabbitRepeatLogin ：缺少UserID")
		return
	}
	if Result.ServerID != RabbitStr.ServerID {
		Send(Result.Data["UserID"], "Offline")
		Unregister(Result.Data["UserID"])
	}
}

func rabbitBroadcast(Result *RabbitMsgResult) { //全服广播
	if !qesygo.VeriPara(Result.Data, []string{"Content"}) {
		fmt.Println("rabbitBroadcast ：缺少Content")
		return
	}
	Broadcast(Result.Data["Content"])
}

func rabbitFlushItem(Result *RabbitMsgResult) { // 刷新用户所有信息
	if !qesygo.VeriPara(Result.Data, []string{"UserID"}) {
		fmt.Println("rabbitFlushItem ：缺少UserID")
		return
	}
	Rs, Err := models.CacheGetOne("dk_user", "UserID", Result.Data["UserID"])
	if len(Rs) == 0 || Err != nil { // 操作失败
		fmt.Println("rabbitFlushItem ：没有此用户")
		return
	}
	// 刷新用户信息
}

func RabbitMQHeartBeat() { // RabbitMQ 心跳
	RabbitMsgResult := RabbitMsgResult{
		Act:      "Ping",
		Data:     map[string]string{},
		ServerID: RabbitStr.ServerID,
	}
	for { //30秒发一次心跳给RabbitMQ
		RabbitStr.PublishMsg(RabbitMsgResult)
		time.Sleep(30 * time.Second) // 改为30秒，断线后重连周期也依赖这个时间
	}
}
