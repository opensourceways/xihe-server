/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/spaceapp/app"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

// reqToCreateSpaceApp
type reqToCreateSpaceApp struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

func (req *reqToCreateSpaceApp) toCmd() (cmd app.CmdToCreateApp, err error) {
	cmd.SpaceId, err = types.NewIdentity(req.SpaceId)
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

	cmd.AppURL, err = domain.NewAppURL(req.AppURL)

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

// reqToNotifyStarting
type reqToNotifyStarting struct {
	reqToCreateSpaceApp

	AllBuildLog string `json:"all_build_log"`
}

func (req *reqToNotifyStarting) toCmd() (cmd app.CmdToNotifyStarting, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.Logs = req.AllBuildLog

	return
}

// reqToFailedStatus
type reqToFailedStatus struct {
	reqToCreateSpaceApp

	Status string `json:"status"`
	Reason string `json:"reason"`

	AllBuildLog string `json:"all_build_log"`
}

func (req *reqToFailedStatus) toCmd() (cmd app.CmdToNotifyFailedStatus, err error) {
	if cmd.SpaceAppIndex, err = req.reqToCreateSpaceApp.toCmd(); err != nil {
		return
	}

	cmd.Reason = req.Reason
	cmd.Status, err = domain.NewAppStatus(req.Status)
	if cmd.Status == domain.AppStatusBuildFailed {
		cmd.Logs = req.AllBuildLog
	}

	return
}
