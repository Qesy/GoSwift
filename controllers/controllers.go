package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"server/lib"
	"server/models"
	wsController "server/websocket"
	"sort"
	"strconv"
	"strings"

	"github.com/Qesy/qesygo"
	"golang.org/x/net/websocket"
)

// Entry 结构
type Entry struct {
	Res        http.ResponseWriter
	Req        *http.Request
	Controller string
	Method     string
	URL        string
	Params     []string
	PostJson   map[string]string
	Body       []byte
}

type Api struct {
	Method     string
	Permission []string
}

// 登录平台对接
const PLATFORM_WECHAT = "wechat"

// 角色
const ROLE_USER string = "user"
const ROLE_EVERYONE string = "everyone"

var ApiArr []Api = []Api{
	/********** 用户接口 **********/
	{Method: "ApiSystem_Get", Permission: []string{ROLE_EVERYONE}},    //获取通用配置
	{Method: "ApiSystem_Notice", Permission: []string{ROLE_EVERYONE}}, //获取公告列表
	/********** 测试专用 ***********/
	{Method: "ApiTest_Login", Permission: []string{ROLE_EVERYONE}},    //用户一键登录
	{Method: "ApiGm_ItemSend", Permission: []string{ROLE_EVERYONE}},   //用户一键登录
	{Method: "ApiGm_Upgrade", Permission: []string{ROLE_EVERYONE}},    //用户一键登录
	{Method: "ApiGm_SelectGate", Permission: []string{ROLE_EVERYONE}}, //用户一键登录
	{Method: "ApiGm_Quarantine", Permission: []string{ROLE_EVERYONE}}, //用户一键登录
	{Method: "ApiGm_Mail", Permission: []string{ROLE_EVERYONE}},       //用户一键登录
	{Method: "ApiGm_OneKeySend", Permission: []string{ROLE_EVERYONE}}, //用户一键登录
	/********** Xy平台 ***********/
	{Method: "ApiXy_Login", Permission: []string{ROLE_EVERYONE}}, //用户一键登录
	/********** CP后台专用 ***********/
	{Method: "Cp_roles", Permission: []string{ROLE_EVERYONE}},      //玩家信息
	{Method: "Cp_role_items", Permission: []string{ROLE_EVERYONE}}, //玩家道具
	{Method: "Cp_role_equip", Permission: []string{ROLE_EVERYONE}}, //玩家道具

}

// Fetch 路由匹配
func (e *Entry) fetch() {
	e.PostJson = map[string]string{}
	if e.URL == "" {
		return
	}
	urlArr := strings.Split(e.URL, "/")
	if len(urlArr) <= 1 {
		e.Controller = urlArr[0]
	} else {
		e.Controller = urlArr[0]
		if urlArr[1] != "" {
			e.Method = urlArr[1]
		}
		e.Params = urlArr[2:]
	}
}

// Run 启动所有服务
func (e *Entry) Run() {
	e.fetch()
	if e.Controller == "static" { // 静态文件服务器
		e.Res.Header().Add("Access-Control-Allow-Origin", "*")
		e.Res.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, Cookie,X-CSRF-TOKEN, Accept,Authorization")
		e.Res.Header().Add("Access-Control-Expose-Headers", "Authorization,authenticated")
		e.Res.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS")
		e.Res.Header().Add("Access-Control-Allow-Credentials", "true")
		fmt.Println("FILERECEIVE:" + e.URL)
		http.ServeFile(e.Res, e.Req, e.URL)
		return
	}
	if e.Controller == "ws" { // Websocket
		e.Req.Header.Set("Origin", "file://")
		handle := websocket.Handler(wsController.Echo)
		handle.ServeHTTP(e.Res, e.Req)
		return
	}
	t := reflect.TypeOf(e)
	Method := e.Controller + "_" + e.Method + "_Action"
	m, ok := t.MethodByName(Method)
	if !ok {
		fmt.Println(ok, e.URL, Method)
		e.ErrorStatus(404)
		return
	}
	fmt.Println("HTTPRECEIVE:" + Method) //打印请求接口
	if strings.Contains(Method, "Upload") {
		e.Req.ParseMultipartForm(5 * 1024 * 1024 * 1024) //接收上传文件用
		for k, v := range e.Req.PostForm {
			e.PostJson[k] = v[0]
		}
		fmt.Println("PostJsonUpload", e.PostJson)
	} else {
		e.Body, _ = io.ReadAll(e.Req.Body)
		qesygo.JsonDecode(e.Body, &e.PostJson)
		fmt.Println("PostJson", string(e.Body), e.PostJson)
	}
	if strings.Contains(Method, "Api") { //调用接口，需验证签名
		e.Res.Header().Add("Access-Control-Allow-Origin", "*")
		e.Res.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, Cookie,X-CSRF-TOKEN, Accept,Authorization")
		e.Res.Header().Add("Access-Control-Expose-Headers", "Authorization,authenticated")
		e.Res.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS")
		e.Res.Header().Add("Access-Control-Allow-Credentials", "true")
		e.Res.Header().Add("Content-Type", "application/json;charset=utf-8")
		if e.Req.Method == "OPTIONS" { //跨域 预检测（跳过）
			return
		}

		if !e.VeriSign() { //验证签名
			e.Error(4001)
			return
		}
		if !e.Permission() { //验证权限
			e.Error(4004)
			return
		}

	}
	m.Func.Call([]reflect.Value{reflect.ValueOf(e)})
}

func (e *Entry) Permission() bool { //验证权限
	Method := e.Controller + "_" + e.Method
	Permission := []string{}
	MethodArr := []string{}
	for _, val := range ApiArr {
		MethodArr = append(MethodArr, val.Method)
		if Method == val.Method {
			Permission = val.Permission
		}
	}
	if !qesygo.InArray(MethodArr, Method) { //没定义的网址不允许访问
		return false
	}
	if qesygo.InArray(Permission, ROLE_EVERYONE) { //可匿名访问
		return true
	}
	if !qesygo.VeriPara(e.PostJson, []string{"PostUserID"}) { //不可匿名访问接口必须传UserID
		return false
	}
	return true

}

// ErrorStatus 错误状态
func (e *Entry) ErrorStatus(code int) {
	e.Res.WriteHeader(code)
	fmt.Fprintf(e.Res, "%d page not found", code)
}

// Error 返回数据
func (e *Entry) Error(code int) {
	e.ErrorData(code, map[string]string{})
}

// Error 返回数据
func (e *Entry) ErrorData(code int, Data map[string]string) {
	retArr := make(map[string]interface{})
	retArr["Act"] = e.Controller + "/" + e.Method
	retArr["Code"] = code
	retArr["Data"] = Data
	retArr["Msg"] = models.StaticGetByField("codeError", strconv.Itoa(code), "des")
	jsonByte, _ := qesygo.JsonEncode(retArr)
	log.Println("Error:", retArr, e.PostJson, e.GetIp())
	fmt.Fprint(e.Res, string(jsonByte))
}

// Error 返回数据
func (e *Entry) ErrorDesc(code int, Msg string) {
	retArr := make(map[string]interface{})
	retArr["Act"] = e.Controller + "/" + e.Method
	retArr["Code"] = code
	retArr["Data"] = make(map[string]string)
	retArr["Msg"] = models.StaticGetByField("codeError", strconv.Itoa(code), "des") + ":" + Msg
	jsonByte, _ := qesygo.JsonEncode(retArr)
	log.Println("Error:", retArr, e.PostJson)
	fmt.Fprint(e.Res, string(jsonByte))
}

// Success
func (e *Entry) Success(ret interface{}) {
	retArr := make(map[string]interface{})
	retArr["Act"] = e.Controller + "/" + e.Method
	retArr["Code"] = 0
	retArr["Data"] = ret
	retArr["Msg"] = ""
	jsonByte, _ := qesygo.JsonEncode(retArr)
	fmt.Println(string(jsonByte))
	fmt.Fprint(e.Res, string(jsonByte))
}

// Show 展示数据
func (e *Entry) Show(ret interface{}) {
	var str string
	switch ret.(type) {
	case []byte:
		str = fmt.Sprintf("%s", ret)
	case string:
		str = fmt.Sprintf("%s", ret)
	case map[string]interface{}:
		json, err := qesygo.JsonEncode(ret)
		if err != nil {
			return
		}
		str = string(json)
	default:
		return
	}
	fmt.Println(str)
	fmt.Fprint(e.Res, str)
}

// 验证签名
func (e *Entry) VeriSign() bool {
	if _, ok := e.PostJson["Time"]; !ok { //强制提交Time
		return false
	}
	PostData := map[string]string{}
	for k, v := range e.PostJson {
		PostData[k] = v
	}
	delete(PostData, "Sign")
	keys := []string{}
	for k := range PostData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	NewData := []string{}
	for _, v := range keys {
		NewData = append(NewData, v+"="+PostData[v])
	}
	Sign := qesygo.Md5(qesygo.Implode(NewData, "&") + "&Secret=" + lib.ConfRs.Conf["Secret"])
	fmt.Println("Veri", Sign, e.PostJson["Sign"], qesygo.Implode(NewData, "&")+"&Secret="+lib.ConfRs.Conf["Secret"])
	return Sign == e.PostJson["Sign"]
}

func (e *Entry) GetPage() (int, int) { //分页必要
	Page := 1
	if PagePost, ok := e.PostJson["Page"]; ok {
		Page, _ = strconv.Atoi(PagePost)
	}
	Num := 20
	if NumPost, ok := e.PostJson["Num"]; ok && NumPost != "" {
		Num, _ = strconv.Atoi(NumPost)
		if Num > 100 || Num < 1 {
			Num = 100
		}
	}
	return Page, Num
}

func (e *Entry) GetIp() string {
	Ip := e.Req.Header.Get("X-Forwarded-For")
	if strings.Contains(Ip, "127.0.0.1") || Ip == "" {
		Ip = e.Req.Header.Get("X-real-ip")
	}
	if Ip == "" {
		Ip = "127.0.0.1"
	}
	return Ip
}
