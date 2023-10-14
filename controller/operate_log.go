package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	OPERATE_LOG         = "operate_log"
	OPERATE_TYPE_USER   = "user"
	OPERATE_TYPE_SYSTEM = "system"
)

// operate log record each user operation
// which contains userid, clientid, operate time, operate type, operate detail, result
func prepareOperateLog(ctx *gin.Context,
	user string, operateType, detail string) {
	_, time := utils.DateAndTime(utils.Now())

	log := fmt.Sprintf("XIHE_OPERATE_LOG: %s:%s:%s %s at %s",
		user, ctx.ClientIP(), operateType, detail, time)

	ctx.Set(OPERATE_LOG, log)
}

func GetOperateLog(ctx *gin.Context) string {
	var result string

	if ctx.Writer.Status() >= 400 {
		result = "failed"
	} else {
		result = "succeeded"
	}

	l, ok := ctx.Get(OPERATE_LOG)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s %s", l, result)
}
