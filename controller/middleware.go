package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/utils"
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

func internalApiCheckMiddleware(ctl *baseController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := checkToken(ctx); err != nil {
			ctl.sendCodeMessage(
				ctx, "invalid_internal_token",
				errors.New("this interface requires the internal token"),
			)

			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func checkToken(ctx *gin.Context) error {
	rawToken := ctx.GetHeader(apiConfig.InternalHeader)

	calcTokenHash, err := utils.EncodeToken(rawToken, apiConfig.InternalTokeSalt)
	if err != nil {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "check token failed", err,
		)
	}

	if calcTokenHash != apiConfig.InternalTokenHash {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "invalid token", errors.New("token mismatch"),
		)
	}

	return nil
}
