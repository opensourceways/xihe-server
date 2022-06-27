package main

import (
	"fmt"

	"github.com/opensourceways/xihe-server/models"
	"github.com/opensourceways/xihe-server/routers"
	"github.com/opensourceways/xihe-server/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// init config
	util.InitConfig("")
	if util.GetConfig().AppModel == "dev" || util.GetConfig().AppModel == "debug" {
		util.Log.SetLevel(logrus.DebugLevel)
		util.GetConfig().AppModel = gin.DebugMode
	} else {
		util.Log.SetLevel(logrus.InfoLevel)
		util.GetConfig().AppModel = gin.ReleaseMode
	}
	gin.SetMode(util.GetConfig().AppModel)
	//init database
	err := util.InitMonogoDB()
	if err != nil {
		util.Log.Errorf("database connect failed , err:%v\n", err)
		return
	}
	//init Authing.cn config
	models.InitAuthing("", "")
	// init router
	r := routers.InitRouter()
	address := fmt.Sprintf(":%d", util.GetConfig().AppPort)
	util.Log.Infof(" startup meta http service at port %s .and %s mode \n", address, util.GetConfig().AppModel)
	if err := r.Run(address); err != nil {
		util.Log.Infof("startup meta  http service failed, err:%v\n", err)
	}
}
