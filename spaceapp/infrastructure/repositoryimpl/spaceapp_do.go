package repositoryimpl

import (
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/domain"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
)

const tableSpaceApp = "space_app"

func toSpaceAppDO(m *spaceappdomain.SpaceApp) spaceappDO {
	do := spaceappDO{
		Status:      m.Status.AppStatus(),
		SpaceId:     m.SpaceId.Integer(),
		Version:     m.Version,
		CommitId:    m.CommitId,
		Reason:      m.Reason,
		RestartedAt: m.RestartedAt,
		ResumedAt:   m.ResumedAt,
	}

	if m.Id != nil {
		do.Id = m.Id.Integer()
	}

	if m.AppURL != nil {
		do.AppURL = m.AppURL.AppURL()
	}

	if m.AppLogURL != nil {
		do.AppLogURL = m.AppLogURL.URL()
	}

	if m.BuildLogURL != nil {
		do.BuildLogURL = m.BuildLogURL.URL()
	}

	return do
}

// spaceappDO
type spaceappDO struct {
	Id       int64  `gorm:"primarykey"`
	SpaceId  int64  `gorm:"column:space_id;index:,unique"`
	CommitId string `gorm:"column:commit_id"`

	Status string `gorm:"column:status"`
	Reason string `gorm:"column:reason"`

	RestartedAt int64 `gorm:"column:restarted_at"`
	ResumedAt   int64 `gorm:"column:resumed_at"`

	AppURL    string `gorm:"column:app_url"`
	AppLogURL string `gorm:"column:app_log_url"`

	AllBuildLog string `gorm:"column:all_build_log;type:text;"`
	BuildLogURL string `gorm:"column:build_log_url"`

	Version int `gorm:"column:version"`
}

func (spaceappDO) TableName() string {
	return tableSpaceApp
}

func (do *spaceappDO) toSpaceApp() spaceappdomain.SpaceApp {
	v := spaceappdomain.SpaceApp{
		Id: domain.CreateIdentity(do.Id),
		SpaceAppIndex: spaceappdomain.SpaceAppIndex{
			SpaceId:  domain.CreateIdentity(do.SpaceId),
			CommitId: do.CommitId,
		},
		Status:      spaceappdomain.CreateAppStatus(do.Status),
		Reason:      do.Reason,
		RestartedAt: do.RestartedAt,
		ResumedAt:   do.ResumedAt,
		Version:     do.Version,
	}

	if do.AppURL != "" {
		v.AppURL = spaceappdomain.CreateAppURL(do.AppURL)
	}

	if do.AppLogURL != "" {
		v.AppLogURL = commondomain.CreateURL(do.AppLogURL)
	}

	if do.BuildLogURL != "" {
		v.BuildLogURL = commondomain.CreateURL(do.BuildLogURL)
	}

	return v
}
