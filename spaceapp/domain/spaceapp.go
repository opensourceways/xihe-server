package domain

import (
	"fmt"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/domain"
)

// SpaceAppIndex represents the index for a space app.
type SpaceAppIndex struct {
	SpaceId  domain.Identity
	CommitId string
}

// SpaceApp represents a space app.
type SpaceApp struct {
	Id domain.Identity

	SpaceAppIndex

	Status AppStatus
	Reason string

	ResumedAt   int64
	RestartedAt int64

	AppURL      AppURL
	AppLogURL   domain.URL
	BuildLogURL domain.URL

	Version int
}

// StartServing starts the service for the space app with the specified app URL and log URL.
func (app *SpaceApp) StartServing(appURL AppURL, logURL domain.URL) error {
	if app.Status.IsStarting() || app.Status.IsRestarting() || app.Status.IsResuming() {

		app.Status = AppStatusServing
		app.AppURL = appURL
		app.AppLogURL = logURL

		return nil
	}

	e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
	return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
}

// GetFailedReason app only return failed reason
func (app *SpaceApp) GetFailedReason() string {
	if !app.Status.IsUpdateStatusAccept() {
		return ""
	}
	return app.Reason
}
