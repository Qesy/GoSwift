package models

import (
	"strconv"

	"github.com/Qesy/qesydb"
	"github.com/Qesy/qesygo"
)

type UserInfo struct {
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}

func WeChatGetCommonApi(Path string, Para map[string]string) map[string]string { //通用调用接口
	Para["access_token"] = WeChatGetAccessToken()
	var c qesygo.CurlModel
	b, _ := c.SetUrl("https://api.weixin.qq.com/" + Path).SetPara(Para).ExecGet()
	Ret := map[string]string{}
	qesygo.JsonDecode(b, &Ret)
	return Ret
}

func WeChatGetOpenId(js_code string) map[string]string { //通过js_code获取openid
	//var m qesydb.Model
	//PlatformRs, _ := m.SetTable("ob_auth").ExecSelectOne()
	SettingRs, _ := SettingGet()
	Path := "sns/jscode2session"
	Para := map[string]string{
		"appid":      SettingRs["WeChatAppID"],
		"secret":     SettingRs["WeChatSecret"],
		"js_code":    js_code,
		"grant_type": "authorization_code",
	}
	var c qesygo.CurlModel
	b, _ := c.SetUrl("https://api.weixin.qq.com/" + Path).SetPara(Para).ExecGet()
	Ret := map[string]string{}
	qesygo.JsonDecode(b, &Ret)
	return Ret
}

func WeChatGetWxacode(Para map[string]string) []byte { // 总共生成的码数量限制为 100,000，请谨慎调用。
	Path := "wxa/getwxacode"
	//Para := map[string]string{"path": "page/index/index"}
	ACCESS_TOKEN := WeChatGetAccessToken()
	var c qesygo.CurlModel
	b, _ := c.SetUrl("https://api.weixin.qq.com/" + Path + "?access_token=" + ACCESS_TOKEN).SetPara(Para).SetIsJson(true).ExecPost()
	return b
}

func WeChatGetAccessToken() string { //获取AccessToken
	var m qesydb.Model
	PlatformRs, _ := m.SetTable("ob_auth").ExecSelectOne()
	TsExpire, _ := strconv.Atoi(PlatformRs["WeChatTsExpire"])
	if PlatformRs["AccessToken"] != "" && TsExpire > int(qesygo.Time("Second")) {
		return PlatformRs["AccessToken"]
	}
	var c qesygo.CurlModel
	b, _ := c.SetUrl("https://api.weixin.qq.com/cgi-bin/token").SetPara(map[string]string{
		"grant_type": "client_credential",
		"appid":      PlatformRs["WeChatAppID"],
		"secret":     PlatformRs["WeChatSecret"],
	}).ExecGet()
	Ret := map[string]string{}
	qesygo.JsonDecode(b, &Ret)

	TsExpire = int(qesygo.Time("Second")) + 7200
	m.SetTable("ob_auth").SetUpdate(map[string]string{"WeChatAccessToken": Ret["access_token"], "WeChatTsExpire": strconv.Itoa(TsExpire)}).ExecUpdate()
	return Ret["access_token"]
}
