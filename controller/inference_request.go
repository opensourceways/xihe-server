/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/spaceapp/app"
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
