/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package privacycheck provides functionality for checking user privacy agreement.
package privacycheck

// import (
// 	"github.com/gin-gonic/gin"

// 	commonctl "github.com/opensourceways/xihe-server/common/controller"
// 	"github.com/opensourceways/xihe-server/common/controller/middleware"
// 	"github.com/opensourceways/xihe-server/common/domain/primitive"
// 	userapp "github.com/opensourceways/xihe-server/user/app"
// )

// // PrivacyCheck is a function that creates a new instance of the privacyCheck.
// func PrivacyCheck(u middleware.UserMiddleWare, ua userapp.UserService,
// ) *privacyCheck {
// 	return &privacyCheck{
// 		user:    u,
// 		userApp: ua,
// 	}
// }

// type privacyCheck struct {
// 	user    middleware.UserMiddleWare
// 	userApp userapp.UserService
// }

// // Check is used to determine whether the user agrees to the privacy agreement
// func (c *privacyCheck) Check(ctx *gin.Context) {
// 	user := c.user.GetUserAndExitIfFailed(ctx)
// 	if user == nil {
// 		return
// 	}

// 	c.check(ctx, user)
// }

// func (c *privacyCheck) CheckOwner(ctx *gin.Context) {
// 	owner, err := primitive.NewAccount(ctx.Param("owner"))
// 	if err != nil {
// 		commonctl.SendBadRequestParam(ctx, err)

// 		ctx.Abort()
// 	}

// 	c.check(ctx, owner)
// }

// func (c *privacyCheck) CheckName(ctx *gin.Context) {
// 	name, err := primitive.NewAccount(ctx.Param("name"))
// 	if err != nil {
// 		commonctl.SendBadRequestParam(ctx, err)

// 		return
// 	}

// 	c.check(ctx, name)
// }

// func (c *privacyCheck) check(ctx *gin.Context, user primitive.Account) {
// 	isAgree, err := c.userApp.IsAgreePrivacy(ctx, user)
// 	if err != nil {
// 		commonctl.SendError(ctx, allerror.New(allerror.ErrorFailedToGetUserInfo,
// 			"failed to get user info when checking privacy agreement", err))

// 		ctx.Abort()
// 	}

// 	if !isAgree {
// 		err = allerror.New(allerror.ErrorCodeDisAgreedPrivacy,
// 			"disagreed privacy", fmt.Errorf("disagreed privacy"))
// 		commonctl.SendError(ctx, err)

// 		ctx.Abort()
// 	}

// 	ctx.Next()
// }
