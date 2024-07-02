package models

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/Qesy/qesygo"
)

var Static = struct {
	sync.RWMutex
	Data map[string]map[string]map[string]string //文件|ID|map[string]string
}{Data: make(map[string]map[string]map[string]string)}

var Filter = struct { //屏蔽字库
	sync.RWMutex
	Data []string
}{Data: []string{}}

func FilterCache() {
	if Str, err := qesygo.ReadFile("static/txt/filter.txt"); err != nil {
		log.Panic("Static File filter.txt Open Err !")
	} else {
		defer Filter.Unlock()
		Filter.Lock()
		Filter.Data = strings.Split(string(Str), "|")
		fmt.Println("Import File filter.txt Sueccss !")
	}
}

func StaticCache(FileName string) { //仅用于初始化和刷新静态表
	if json, err := qesygo.ReadFile("static/json/" + FileName + ".json"); err != nil {
		panic("Static File " + FileName + " Open Err !")
	} else {
		Static.Lock()
		defer Static.Unlock()
		list := make(map[string]map[string]string)
		if err := qesygo.JsonDecode(json, &list); err != nil {
			panic("Static File " + FileName + ".json JsonDecode Err !")
		}
		Static.Data[FileName] = list
		fmt.Println("Import File " + FileName + ".json Sueccss !")
	}
}

func StaticGet(FileName string) map[string]map[string]string {
	Static.RLock()
	defer Static.RUnlock()
	if sData, ok := Static.Data[FileName]; !ok {
		panic("Static Data Get Err : " + FileName)
	} else {
		return sData
	}
}

func StaticKeys(FileName string) []string {
	Arr := StaticGet(FileName)
	Keys := []string{}
	for k := range Arr {
		Keys = append(Keys, k)
	}
	return Keys
}

func StaticGetKeyIsHaver(FileName string, Id string) bool {
	list := StaticGet(FileName)
	_, ok := list[Id]
	return ok
}

func StaticGetByKey(FileName string, Id string) map[string]string {
	list := StaticGet(FileName)
	if rs, ok := list[Id]; !ok {
		log.Panicf("StaticGetByKey Err : FileName %v Id %v ", FileName, Id)
	} else {
		return rs
	}
	return map[string]string{}
}

func StaticGetByField(FileName string, Id string, Field string) string {
	rs := StaticGetByKey(FileName, Id)
	if str, ok := rs[Field]; !ok {
		log.Panicf("StaticGetByField Err : FileName %v Id %v Field %v", FileName, Id, Field)
	} else {
		return str
	}
	return ""
}

func StaticGetListByField(FileName string, Key string) map[string]map[string]string { //把某字段作为Key生成新的MAP
	list := StaticGet(FileName)
	NewMap := map[string]map[string]string{}
	for _, v := range list {
		NewMap[v[Key]] = v
	}
	return NewMap
}

func StaticCommonKV() map[string]string { //获取通用配置
	Arr := StaticGet("common")
	ArrMap := map[string]string{}
	for _, v := range Arr {
		ArrMap[v["key"]] = v["value"]
	}
	return ArrMap
}

func StaticRolePropKV() map[string]string { //获取角色配置
	Arr := StaticGet("RoleProp")
	ArrMap := map[string]string{}
	for _, v := range Arr {
		ArrMap[v["YName"]] = v["Coefficient"]
	}
	return ArrMap
}

func IsAllowNickName(Str string) bool {
	if utf8.RuneCountInString(Str) < 2 || utf8.RuneCountInString(Str) > 12 {
		log.Printf("IsAllowNickName Err : MinLen %v MaxLen %v", utf8.RuneCountInString(Str), utf8.RuneCountInString(Str))
		return false
	}
	if IsSpecialCharacters(Str) { //是否包含特殊字符
		log.Printf("IsAllowNickName Err : IsSpecialCharacters %v", Str)
		return false
	}
	if IsFilter(Str) {
		log.Printf("IsAllowNickName Err : IsFilter %v", Str)
		return false
	}
	return true
}

func IsFilter(Str string) bool { // 是否触发屏蔽字
	for _, v := range Filter.Data {
		if strings.Contains(Str, v) {
			log.Printf("IsFilter Str %v v %v", Str, v)
			return true
		}
	}
	return false
}

func IsSpecialCharacters(str string) bool { //是否包含特殊字符
	reg := regexp.MustCompile("^[a-zA-Z0-9\u4e00-\u9fa5]{1,12}")
	result := reg.FindStringSubmatch(str)
	if len(result) == 0 || str != result[0] {
		return true
	}
	return false
}

func IsSpecialCharacters2(str string) bool { //是否包含特殊字符
	reg := regexp.MustCompile("^[a-zA-Z0-9\u4e00-\u9fa5]{1,50}")
	result := reg.FindStringSubmatch(str)
	if len(result) == 0 || str != result[0] {
		return true
	}
	return false
}
