package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "server/bootstrap"
	"server/controllers"
	"server/lib"

	"github.com/Qesy/qesygo"
)

// MyMux 定义数据类型
type MyMux struct{}

var strProt = flag.String("p", "", "启动端口")

func main() {
	flag.Parse()
	if *strProt != "" {
		lib.ConfRs.Conf["Port"] = *strProt
	}
	mux := &MyMux{}
	if Err := http.ListenAndServe(":"+lib.ConfRs.Conf["Port"], mux); Err != nil {
		fmt.Println("Web service start fail : " + Err.Error())
	}
}

func (p *MyMux) ServeHTTP(Res http.ResponseWriter, Req *http.Request) {
	entry := controllers.Entry{Res: Res, Req: Req, Controller: "Index", Method: "Index", Params: []string{}, URL: qesygo.Substr(Req.URL.Path, 1, 0)}
	entry.Run()

}
