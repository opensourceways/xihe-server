package interfaces

import (
	"github.com/opensourceways/xihe-server/app"
)

type Project struct {
	Name      string `json:"name" required:"true"`
	Desc      string `json:"desc" required:"true"`
	Type      string `json:"type" required:"true"`
	CoverId   string `json:"cover_id" required:"true"`
	Protocol  string `json:"protocol" required:"true"`
	Training  string `json:"training" required:"true"`
	Inference string `json:"inference" required:"true"`
}

func (p *Project) GenCreateProjectCmd() app.CreateProjectCmd {
	return app.CreateProjectCmd{
		Name:      p.Name,
		Desc:      p.Desc,
		Type:      p.Type,
		CoverId:   p.CoverId,
		Protocol:  p.Protocol,
		Training:  p.Training,
		Inference: p.Inference,
	}
}
