package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func checkUserEmailMiddleware(ctl *baseController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pl, _, _ := ctl.checkUserApiToken(ctx, false)

		if pl.hasEmail() {
			ctl.sendCodeMessage(
				ctx, "no email",
				errors.New("this api need email of user"),
			)

			ctx.Abort()

			return
		}

		ctx.Next()

	}
}
