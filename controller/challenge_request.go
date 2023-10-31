package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type competitorApplyRequest struct {
	Name      string            `json:"name"`
	City      string            `json:"city"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Identity  string            `json:"identity"`
	Province  string            `json:"province"`
	Detail    map[string]string `json:"detail"`
	Agreement bool              `json:"agreement"`
}

func (req *competitorApplyRequest) toCmd(user domain.Account) (cmd app.CompetitorApplyCmd, err error) {
	if cmd.Name, err = domain.NewCompetitorName(req.Name); err != nil {
		return
	}

	if cmd.City, err = domain.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = domain.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = domain.NewcompetitionIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = domain.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}

type aiQuestionAnswerSubmitRequest struct {
	Times  int      `json:"times"`
	Result []string `json:"result"`
	Answer string   `json:"answer"`
}

func (req *aiQuestionAnswerSubmitRequest) toCmd() (app.AIQuestionAnswerSubmitCmd, error) {
	if len(req.Result) == 0 || req.Answer == "" {
		return app.AIQuestionAnswerSubmitCmd{}, errors.New("invalid cmd")
	}

	return app.AIQuestionAnswerSubmitCmd{
		Times:  req.Times,
		Result: req.Result,
		Answer: req.Answer,
	}, nil
}

type aiQuestionAnswerSubmitResp struct {
	Score int `json:"score"`
}
