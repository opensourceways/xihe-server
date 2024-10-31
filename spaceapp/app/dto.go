package app

import (
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

type InferenceDTO struct {
	Error      string `json:"error"`
	AccessURL  string `json:"access_url"`
	InstanceId string `json:"inference_id"`
}

// CmdToNotifyServiceIsStarted is a command to notify that the service has started.
type CmdToNotifyServiceIsStarted struct {
	CmdToNotifyBuildIsStarted

	AppURL domain.AppURL
}

// CmdToNotifyBuildIsStarted is a command to notify that the build has started.
type CmdToNotifyBuildIsStarted struct {
	domain.SpaceAppIndex

	LogURL commondomain.URL
}

// CmdToCreateApp is a command to create an app.
type CmdToCreateApp = domain.SpaceAppIndex

// CmdToNotifyFailedStatus is a command to notify that status has update.
type CmdToNotifyFailedStatus struct {
	domain.SpaceAppIndex

	Reason string
	Status domain.AppStatus

	Logs string
}

// CmdToNotifyStarting is a command to notify that the build has finished.
type CmdToNotifyStarting struct {
	domain.SpaceAppIndex

	Logs string
}

// SpaceAppDTO is a data transfer object for space app.
type SpaceAppDTO struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	AppURL      string `json:"app_url"`
	AppLogURL   string `json:"-"`
	BuildLogURL string `json:"-"`
}

type GetSpaceAppCmd = spacedomain.SpaceIndex

// BuildLogsDTO
type BuildLogsDTO struct {
	Logs string `json:"logs"`
}

func toSpaceDTO(space *spacedomain.Project) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     space.Id,
		Status: space.Exception.Exception(),
		Reason: types.ExceptionMap[space.Exception.Exception()],
	}

	return dto
}

func toSpaceNoCompQuotaDTO(space *spacedomain.Project) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     space.Id,
		Status: types.NoCompQuotaException,
		Reason: types.ExceptionMap[types.NoCompQuotaException],
	}
	return dto
}
