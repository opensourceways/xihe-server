package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

func checkUserEmailMiddleware(ctl *baseController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pl, _, ok := ctl.checkUserApiTokenNoRefresh(ctx, false)
		if !ok {
			ctx.Abort()

			return
		}

		if !pl.hasEmail() {
			ctl.sendCodeMessage(
				ctx, "user_no_email",
				errors.New("this interface requires the users email"),
			)

			ctx.Abort()

			return
		}

		ctx.Next()

	}
}

func ClearSensitiveInfoMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pl, exist := ctx.Get(PayLoad)
		if !exist {
			logrus.Debugf("cannot found payload")

			return
		}

		payload, ok := pl.(*oldUserTokenPayload)
		if !ok {
			logrus.Debugf("payload assert error")

			return
		}

		utils.ClearStringMemory(payload.PlatformToken)

		ctx.Next()
	}
}
