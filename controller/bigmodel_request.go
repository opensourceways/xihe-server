package controller

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userd "github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	tempAccountHF       = "wukong_hf"
	tempAccountVQAHF    = "vqa_hf"
	tempAccountLuoJiaHF = "luojia_hf"
)

type describePictureResp struct {
	Desc string `json:"desc"`
}

type pictureGenerateRequest struct {
	Desc string `json:"desc"`
}

func (req *pictureGenerateRequest) toCmd(user types.Account) (cmd app.GenPictureCmd, err error) {
	if cmd.Desc, err = domain.NewDesc(req.Desc); err != nil {
		return
	}

	cmd.User = user

	return
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

	if req.Picture == "" || !utils.IsPictureName(req.Picture) {
		err = errors.New("invalid picture")

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

type questionAskHFReq struct {
	Picture  io.Reader `json:"picture"`
	Question string    `json:"question"`
}

func (req *questionAskHFReq) toCmd() (
	cmd app.VQAHFCmd, err error,
) {
	if cmd.User, err = types.NewAccount(tempAccountVQAHF); err != nil {
		return
	}

	cmd.Picture = req.Picture

	cmd.Ask = req.Question

	err = cmd.Validate()

	return
}

type luojiaHFReq struct {
	Picture io.Reader `json:"picture"`
}

func (req *luojiaHFReq) toCmd() (
	cmd app.LuoJiaHFCmd, err error,
) {
	if cmd.User, err = types.NewAccount(tempAccountLuoJiaHF); err != nil {
		return
	}

	cmd.Picture = req.Picture

	err = cmd.Validate()

	return
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
	Desc        string `json:"desc"`
	Style       string `json:"style"`
	ImgQuantity int    `json:"img_quantity"`
}

func (req *wukongRequest) toCmd() (cmd app.WuKongCmd, err error) {
	cmd.Style = req.Style

	if cmd.Desc, err = domain.NewWuKongPictureDesc(req.Desc); err != nil {
		return
	}

	switch req.ImgQuantity {
	case 4:
		cmd.EsType = string(domain.BigmodelWuKong4Img)
	default:
		cmd.EsType = string(domain.BigmodelWuKong)
	}

	err = cmd.Validate()

	return
}

type wukongHFRequest struct {
	Desc  string `json:"desc"`
	Style string `json:"style"`
}

func (req *wukongHFRequest) toCmd() (cmd app.WuKongHFCmd, err error) {
	if cmd.User, err = types.NewAccount(tempAccountHF); err != nil {
		return
	}

	cmd.Style = req.Style

	if cmd.Desc, err = domain.NewWuKongPictureDesc(req.Desc); err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type wukongApiRequest struct {
	Desc  string `json:"desc"`
	Style string `json:"style"`
}

func (req *wukongApiRequest) toCmd() (cmd app.WuKongApiCmd, err error) {
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

func (req *wukongAddLikeFromTempRequest) toCmd(user types.Account) (cmd app.WuKongAddLikeFromTempCmd, err error) {
	if cmd.OBSPath, err = domain.NewOBSPath(req.OBSPath); err != nil {
		return
	}

	cmd.User = user

	return
}

type wukongAddLikeFromPublicRequest struct {
	Owner string `json:"owner" binding:"required"`
	Id    string `json:"id" binding:"required"`
}

func (req *wukongAddLikeFromPublicRequest) toCmd(user types.Account) (
	cmd app.WuKongAddLikeFromPublicCmd, err error,
) {
	owner, err := types.NewAccount(req.Owner)
	if err != nil {
		return
	}

	cmd = app.WuKongAddLikeFromPublicCmd{
		Owner: owner,
		User:  user,
		Id:    req.Id,
	}

	return
}

type wukongAddPublicFromTempRequest wukongAddLikeFromTempRequest

func (req *wukongAddPublicFromTempRequest) toCmd(user types.Account) (cmd app.WuKongAddPublicFromTempCmd, err error) {
	if cmd.OBSPath, err = domain.NewOBSPath(req.OBSPath); err != nil {
		return
	}

	cmd.User = user

	return
}

type wukongAddPublicFromLikeRequest struct {
	Id string `json:"id" binding:"required"`
}

func (req *wukongAddPublicFromLikeRequest) toCmd(user types.Account) app.WuKongAddPublicFromLikeCmd {
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

func (req *wukongAddDiggPublicRequest) toCmd(user types.Account) (cmd app.WuKongAddDiggCmd, err error) {
	owner, err := types.NewAccount(req.User)
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

func (req *wukongCancelDiggPublicRequest) toCmd(user types.Account) (cmd app.WuKongCancelDiggCmd, err error) {
	owner, err := types.NewAccount(req.User)
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

type aiDetectorReq struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

func (req aiDetectorReq) toCmd(user types.Account) (cmd app.AIDetectorCmd, err error) {
	if cmd.Lang, err = domain.NewLang(req.Lang); err != nil {
		return
	}

	if cmd.Text, err = domain.NewAIDetectorText(req.Text); err != nil {
		return
	}

	cmd.User = user

	err = cmd.Validate()

	return
}

type aiDetectorResp struct {
	IsMachine bool `json:"is_machine"`
}

type applyApiReq struct {
	Name     string            `json:"name"`
	City     string            `json:"city"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Identity string            `json:"identity"`
	Province string            `json:"province"`
	Detail   map[string]string `json:"detail"`
}

func (req *applyApiReq) toCmd(user types.Account) (cmd userapp.UserRegisterInfoCmd, err error) {
	if cmd.Name, err = userd.NewName(req.Name); err != nil {
		return
	}

	if cmd.City, err = userd.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = userd.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = userd.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = userd.NewIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = userd.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}

type isApplyResp struct {
	IsApply bool `json:"is_apply"`
}

// baichuan
type baichuanReq struct {
	Text              string  `json:"text"`
	Sampling          bool    `json:"sampling"`
	TopK              int     `json:"top_k"`
	TopP              float64 `json:"top_p"`
	Temperature       float64 `json:"temperature"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
}

func (req *baichuanReq) toCmd(user types.Account) (cmd app.BaiChuanCmd, err error) {
	if cmd.Text, err = domain.NewBaiChuanText(req.Text); err != nil {
		return
	}

	if req.Sampling {
		if cmd.TopK, err = domain.NewTopK(req.TopK); err != nil {
			return
		}
	
		if cmd.TopP, err = domain.NewTopP(req.TopP); err != nil {
			return
		}
	
		if cmd.Temperature, err = domain.NewTemperature(req.Temperature); err != nil {
			return
		}
	
		if cmd.RepetitionPenalty, err = domain.NewRepetitionPenalty(req.RepetitionPenalty); err != nil {
			return
		}
	} else {
		cmd.SetDefault()
	}

	cmd.User = user
	cmd.Sampling = req.Sampling

	return
}