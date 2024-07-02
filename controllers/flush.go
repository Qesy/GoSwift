package controllers

import (
	"fmt"
	"server/lib"
	"server/models"

	"github.com/Qesy/qesygo"
)

func (e *Entry) Flush_Static_Action() { //刷新配置文件
	for _, v := range lib.StaticFiles {
		models.StaticCache(v)
	}
	fmt.Fprintf(e.Res, "刷新静态文件成功")
}

func (e *Entry) Flush_Get_Action() {
	file := e.Req.FormValue("file")
	Arr := map[string]map[string]string{}
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Fprintf(e.Res, "获取静态文件失败:"+rec.(string))
			return
		}
		Json, _ := qesygo.JsonEncode(Arr)
		fmt.Fprint(e.Res, string(Json))

	}()
	Arr = models.StaticGet(file)
}

func (e *Entry) Flush_Redis_Action() {
	lib.RedisCr.FlushAll()
	fmt.Fprint(e.Res, "刷新成功")
}
