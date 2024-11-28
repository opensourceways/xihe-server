package controller

import (
	"github.com/opensourceways/xihe-server/app"
	compapp "github.com/opensourceways/xihe-server/competition/app"
	courseapp "github.com/opensourceways/xihe-server/course/app"
	promapp "github.com/opensourceways/xihe-server/promotion/app"
	spaceapp "github.com/opensourceways/xihe-server/space/app"
)

type homeInfo struct {
	Comp   []compapp.CompetitionSummaryDTO `json:"comp"`
	Course []courseapp.CourseSummaryDTO    `json:"course"`
}

type IndustryDTO struct {
	Comp       []compapp.CompetitionSummaryDTO `json:"comp"`
	Course     []courseapp.CourseSummaryDTO    `json:"course"`
	Promotions []promapp.PromotionDTO          `json:"promotions"`
	Dataset    app.GlobalDatasetsDTO           `json:"dataset"`
	Model      app.GlobalModelsDTO             `json:"model"`
	Peoject    spaceapp.GlobalProjectsDTO      `json:"project"`
}
