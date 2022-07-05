package main

import (
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/server"
	"github.com/opensourceways/xihe-server/util"
)

func main() {

	// init config
	err := util.InitConfig("")
	if err != nil {
		util.Log.Errorf("init config failed , err:%v\n", err.Error())
		return
	}
	//init Authing.cn config
	err = infrastructure.InitAuthing()
	if err != nil {
		util.Log.Errorf("init authing failed , err:%v\n", err.Error())
		return
	}

	// start web server
	server.StartWebServer()
}
