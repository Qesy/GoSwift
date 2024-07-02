package controllers

import (
	"fmt"
	"log"
	"server/lib"
	"server/models"
	"sync"
	"time"

	"github.com/Qesy/qesygo"
)

type PingStr struct {
	Times int
	L     sync.Mutex
}

var Ping PingStr
var SettingKv map[string]string

func Crontab_Run() { //业务定时器
	fmt.Println("Crontab Run Success !")
	Key := lib.RedisKey_Crontab_SetNx()
	for {
		SettingKv, _ = models.SettingGet()

		if qesygo.Date(qesygo.Time("Second"), "15:04:05") == "00:00:00" { // 记录每日日志
			qesygo.LogSave(lib.LogFilePath) //每日Log保存
		}

		if Result, Err := lib.RedisCr.Exists(Key); Result || Err != nil { // 分布式锁已经存在
			time.Sleep(1 * time.Second)
			continue
		}
		if Result, Err := lib.RedisCr.SetEx(Key, 60, "Crontab"+lib.ConfRs.Conf["Port"]); Result != "OK" || Err != nil {
			log.Printf("Crontab_Run Err : Result %v, Err %v", Result, Err)
			time.Sleep(1 * time.Second)
			continue
		}

		time.Sleep(1 * time.Second)
		lib.RedisCr.Del(Key) //执行完成，解锁

	}
}
