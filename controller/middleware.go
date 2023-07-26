package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	cloudAllowedUserName = "MindSpore"
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

func checkUserName(ctl *baseController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pl, _, ok := ctl.checkUserApiTokenNoRefresh(ctx, false)
		if !ok {
			ctx.Abort()

			return
		}

		if !pl.isAllowedUserName(cloudAllowedUserName) {
			ctl.sendCodeMessage(
				ctx, "user_no_permission",
				errors.New("no this interface permission"),
			)

			ctx.Abort()

			return
		}
	}
}
