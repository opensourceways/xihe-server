package controller

import (
	"github.com/opensourceways/xihe-server/app"
	compapp "github.com/opensourceways/xihe-server/competition/app"
	courseapp "github.com/opensourceways/xihe-server/course/app"
)

type homeInfo struct {
	Comp   []compapp.CompetitionSummaryDTO `json:"comp"`
	Course []courseapp.CourseSummaryDTO    `json:"course"`
}

type homeElectricityInfo struct {
	Comp    []compapp.CompetitionSummaryDTO `json:"comp"`
	Course  []courseapp.CourseSummaryDTO    `json:"course"`
	Peoject app.GlobalProjectsDTO           `json:"project"`
	Dataset app.GlobalDatasetsDTO           `json:"dataset"`
	Model   app.GlobalModelsDTO             `json:"model"`
}
