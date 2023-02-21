package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type describePictureResp struct {
	Desc string `json:"desc"`
}

type pictureGenerateRequest struct {
	Desc string `json:"desc"`
}

func (req *pictureGenerateRequest) validate() error {
	if req.Desc == "" {
		return errors.New("missing desc")
	}

	return nil
}

type pictureGenerateResp struct {
	Picture string `json:"picture"`
}

type multiplePicturesGenerateResp struct {
	Pictures []string `json:"pictures"`
}

type questionAskRequest struct {
	Question string `json:"question"`
	Picture  string `json:"picture"`
}

func (req *questionAskRequest) toCmd() (
	q domain.Question, p string, err error,
) {
	if q, err = domain.NewQuestion(req.Question); err != nil {
		return
	}

	if req.Picture == "" {
		err = errors.New("missing picture")

		// TODO doesnot allow chinese name for picture

		return
	}

	p = req.Picture

	return
}

type questionAskResp struct {
	Answer string `json:"answer"`
}

type pictureUploadResp struct {
	Path string `json:"path"`
}

type panguRequest struct {
	Question string `json:"question"`
}

type panguResp struct {
	Answer string `json:"answer"`
}

type luojiaResp struct {
	Answer string `json:"answer"`
}

type CodeGeexRequest struct {
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

func (req *CodeGeexRequest) toCmd() (
	cmd app.CodeGeexCmd, err error,
) {
	cmd.Lang = req.Lang
	cmd.Content = req.Content

	err = cmd.Validate()

	return
}

type wukongRequest struct {
	Desc  string `json:"desc"`
	Style string `json:"style"`
}

func (req *wukongRequest) toCmd() (cmd app.WuKongCmd, err error) {
	cmd.Style = req.Style

	if cmd.Desc, err = domain.NewWuKongPictureDesc(req.Desc); err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type wukongPicturesGenerateResp struct {
	Pictures map[string]string `json:"pictures"`
}

type wukongAddLikeFromTempRequest struct {
	OBSPath string `json:"obspath" binding:"required"`
}

func (req *wukongAddLikeFromTempRequest) toCmd(user domain.Account) app.WuKongAddLikeFromTempCmd {
	return app.WuKongAddLikeFromTempCmd{
		User:    user,
		OBSPath: req.OBSPath,
	}
}

type wukongAddLikeFromPublicRequest struct {
	Owner string `json:"owner" binding:"required"`
	Id    string `json:"id" binding:"required"`
}

func (req *wukongAddLikeFromPublicRequest) toCmd(user domain.Account) app.WuKongAddLikeFromPublicCmd {
	owner, _ := domain.NewAccount(req.Owner)
	return app.WuKongAddLikeFromPublicCmd{
		Owner: owner,
		User:  user,
		Id:    req.Id,
	}
}

type wukongAddPublicFromTempRequest wukongAddLikeFromTempRequest

func (req *wukongAddPublicFromTempRequest) toCmd(user domain.Account) app.WuKongAddPublicFromTempCmd {
	return app.WuKongAddPublicFromTempCmd{
		User:    user,
		OBSPath: req.OBSPath,
	}
}

type wukongAddPublicFromLikeRequest struct {
	Id string `json:"id" biding:"required"`
}

func (req *wukongAddPublicFromLikeRequest) toCmd(user domain.Account) app.WuKongAddPublicFromLikeCmd {
	return app.WuKongAddPublicFromLikeCmd{
		User: user,
		Id:   req.Id,
	}
}

type wukongAddDiggPublicRequest struct {
	User string `json:"user"`
	Id   string `json:"id"`
}

type wukongCancelDiggPublicRequest wukongAddDiggPublicRequest

func (req *wukongAddDiggPublicRequest) toCmd(user domain.Account) (cmd app.WuKongAddDiggCmd, err error) {
	owner, err := domain.NewAccount(req.User)
	if err != nil {
		return
	}
	cmd = app.WuKongAddDiggCmd{
		User:  user,
		Owner: owner,
		Id:    req.Id,
	}

	return
}

func (req *wukongCancelDiggPublicRequest) toCmd(user domain.Account) (cmd app.WuKongCancelDiggCmd, err error) {
	owner, err := domain.NewAccount(req.User)
	if err != nil {
		return
	}
	cmd = app.WuKongCancelDiggCmd{
		User:  user,
		Owner: owner,
		Id:    req.Id,
	}

	return
}

type wukongAddLikeResp struct {
	Id string `json:"id"`
}

type wukongAddPublicResp struct {
	Id string `json:"id"`
}

type wukongPictureLink struct {
	Link string `json:"link"`
}

type wukongDiggResp struct {
	DiggCount int `json:"digg_count"`
}
