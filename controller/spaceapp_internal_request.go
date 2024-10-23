package controller

import (
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/spaceapp/app"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

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

// reqToCreateSpaceApp
type reqToCreateSpaceApp struct {
	SpaceId  string `json:"space_id"`
	CommitId string `json:"commit_id"`
}

func (req *reqToCreateSpaceApp) toCmd() (cmd app.CmdToCreateApp, err error) {
	cmd.SpaceId, err = primitive.NewIdentity(req.SpaceId)
	if err != nil {
		return
	}

	cmd.CommitId = req.CommitId

	return
}
