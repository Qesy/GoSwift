package bootstrap

import (
	"fmt"
	"server/controllers"
	"server/lib"
	"server/models"
	"server/websocket"
	"strconv"

	"github.com/Qesy/qesydb"
	"github.com/Qesy/qesygo"
)

func init() {
	defer func() {
		ErrStr := recover()
		if ErrStr != nil {
			qesygo.Die(ErrStr)
		}
	}()
	GetConf()

	lib.RedisCr.Conninfo = lib.ConfRs.Cache["Host"] + ":" + lib.ConfRs.Cache["Port"]
	lib.RedisCr.Auth = lib.ConfRs.Cache["Auth"]
	if Err := lib.RedisCr.Connect(); Err != nil {
		qesygo.Die("Redis Connect Err : " + Err.Error())
	}
	fmt.Println("Redis Connection Success !")
	qesydb.OpenLog = 1 //打开MYSQL错误日志
	qesydb.MaxOpenConns = 600
	qesydb.MaxIdleConns = 600
	if Err := qesydb.Connect(lib.ConfRs.Db["User"] + ":" + lib.ConfRs.Db["Password"] + "@tcp(" + lib.ConfRs.Db["Host"] + ":" + lib.ConfRs.Db["Port"] + ")/" + lib.ConfRs.Db["Name"] + "?charset=" + lib.ConfRs.Db["Charset"]); Err != nil {
		qesygo.Die("Mysql Connect Err : " + Err.Error())
	}
	fmt.Println("Mysql Connection Success !")

	websocket.RabbitInit(lib.ConfRs.Amqp["Host"], lib.ConfRs.Amqp["Port"], lib.ConfRs.Amqp["User"], lib.ConfRs.Amqp["Password"], lib.ConfRs.Amqp["ServerID"])

	// 导入静态文件，
	for _, v := range lib.StaticFiles {
		models.StaticCache(v)
	}
	// 导入屏蔽字库
	models.FilterCache()

	// go websocket.RabbitMQHeartBeat() // RabbitMQ心跳

	lib.ServerID, _ = strconv.ParseInt(lib.ConfRs.Amqp["ServerID"], 10, 64)
	lib.SnowWorker, _ = qesygo.NewWorker(lib.ServerID)

	qesygo.Log(lib.LogFilePath)  // 开启日志记录
	go controllers.Crontab_Run() // 开启定时器
	go websocket.HubRouter.Run() //开启chan接收

	fmt.Println("Service Start Success !")

}

// getConf 获取配置文件
func GetConf() {
	if str, err := qesygo.ReadFile("conf.ini"); err == nil {
		if err := qesygo.JsonDecode(str, &lib.ConfRs); err != nil {
			qesygo.Die("config reload error: " + err.Error())
		}
	} else {
		qesygo.Die(err)
	}

}
