package controllers

import (
	"fmt"
)

// Index_index 首页
func (e *Entry) Index_Index_Action() {
	ip := e.GetIp()
	fmt.Fprintf(e.Res, "<br><br><br><h1><center> GAME DankeCQ API INTERFACE </center></h1><h2><center>GOLANG Version 1.0</center></h2><center><p>Your IP:"+ip+"</p></center>")
}
