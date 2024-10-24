/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/spaceapp/app"
	appdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
)

// reqToCreateSpaceApp
type reqToCreateSpaceApp struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

func (req *reqToCreateSpaceApp) toCmd() (cmd app.CmdToCreateApp, err error) {
	cmd.SpaceId, err = domain.NewIdentity(req.SpaceId)
	if err != nil {
		return
	}

	cmd.CommitId = req.CommitId

	return
}

// reqToUpdateServiceInfo
type reqToUpdateServiceInfo struct {
	reqToUpdateBuildInfo

	AppURL string `json:"app_url"`
}

func (req *reqToUpdateServiceInfo) toCmd() (cmd app.CmdToNotifyServiceIsStarted, err error) {
	if cmd.CmdToNotifyBuildIsStarted, err = req.reqToUpdateBuildInfo.toCmd(); err != nil {
		return
	}

	cmd.AppURL, err = appdomain.NewAppURL(req.AppURL)

	return
}

// reqToUpdateBuildInfo
type reqToUpdateBuildInfo struct {
	reqToCreateSpaceApp

	LogURL string `json:"log_url"`
}

func (req *reqToUpdateBuildInfo) toCmd() (cmd app.CmdToNotifyBuildIsStarted, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.LogURL, err = commondomain.NewURL(req.LogURL)

	return
}
