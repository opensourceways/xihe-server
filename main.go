package main

import (
	"github.com/opensourceways/xihe-server/models"
	"github.com/opensourceways/xihe-server/server"
	"github.com/opensourceways/xihe-server/util"
)

func main() {
	// init config
	err := util.InitConfig("")
	if err != nil {
		util.Log.Errorf("init config failed , err:%v\n", err)
		return
	}
	//init database
	err = util.InitMonogoDB()
	if err != nil {
		util.Log.Errorf("database connect failed , err:%v\n", err)
		return
	}
	//init Authing.cn config
	err = models.InitAuthing("", "")
	if err != nil {
		util.Log.Errorf("init authing failed , err:%v\n", err)
		return
	}
	// start web server
	server.StartWebServer()
}
