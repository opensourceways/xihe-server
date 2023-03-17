package controller

import (
	compapp "github.com/opensourceways/xihe-server/competition/app"
	"github.com/opensourceways/xihe-server/course/app"
)

type homeInfo struct {
	Comp   []compapp.CompetitionSummaryDTO `json:"comp"`
	Course []app.CourseSummaryDTO          `json:"course"`
}
