package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
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
