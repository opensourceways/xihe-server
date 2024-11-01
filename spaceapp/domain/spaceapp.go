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

// StartBuilding starts the building process for the space app and sets the build log URL.
func (app *SpaceApp) StartBuilding(logURL domain.URL) error {
	if !app.Status.IsInit() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = AppStatusBuilding
	app.BuildLogURL = logURL

	return nil
}

// SetBuildFailed set app status is build failed.
func (app *SpaceApp) SetBuildFailed(status AppStatus, reason string) error {
	if !app.Status.IsBuilding() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
}

// SetStarting sets the starting status of the space app based on the success parameter.
func (app *SpaceApp) SetStarting() error {
	if !app.Status.IsBuilding() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = AppStatusServeStarting

	return nil
}

// SetStartFailed set app status is start failed.
func (app *SpaceApp) SetStartFailed(status AppStatus, reason string) error {
	if !app.Status.IsStarting() {
		e := fmt.Errorf("old status is %s, can not set", app.Status.AppStatus())
		return allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
	}

	app.Status = status
	app.Reason = reason

	return nil
}

// SpaceAppBuildLog is the value object of log
type SpaceAppBuildLog struct {
	AppId domain.Identity
	Logs  string
}

// IsAppNotAllowToInit app can be init if return false
func (app *SpaceApp) IsAppNotAllowToInit() bool {
	if app.Status.IsPaused() || app.Status.IsResuming() || app.Status.IsResumeFailed() || app.Status.IsRestarting() {
		return true
	}

	return false
}
