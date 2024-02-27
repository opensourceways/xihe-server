package controller

import (
	"github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userd "github.com/opensourceways/xihe-server/user/domain"
)

type pictureUploadResp struct {
	Path string `json:"path"`
}

type luojiaResp struct {
	Answer string `json:"answer"`
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
	Name      string            `json:"name"`
	City      string            `json:"city"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Identity  string            `json:"identity"`
	Province  string            `json:"province"`
	Detail    map[string]string `json:"detail"`
	Agreement bool              `json:"agreement"`
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

type newApiTokenResp struct {
	Token string `json:"token"`
	Date  string `json:"date"`
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

// glm2
type glm2Request struct {
	Text              string      `json:"text"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

func (req *glm2Request) toCmd(ch chan string, user types.Account) (cmd app.GLM2Cmd, err error) {
	if cmd.Text, err = domain.NewGLM2Text(req.Text); err != nil {
		return
	}

	history := make([]domain.History, len(req.History))
	for i := range req.History {
		if history[i], err = domain.NewHistory(req.History[i][0], req.History[i][1]); err != nil {
			return
		}
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

	cmd.CH = ch
	cmd.Sampling = req.Sampling
	cmd.User = user

	return
}

// llama2
type llama2Request struct {
	Text              string      `json:"text"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

func (req *llama2Request) toCmd(ch chan string, user types.Account) (cmd app.LLAMA2Cmd, err error) {
	if cmd.Text, err = domain.NewLLAMA2Text(req.Text); err != nil {
		return
	}

	history := make([]domain.History, len(req.History))
	for i := range req.History {
		if history[i], err = domain.NewHistory(req.History[i][0], req.History[i][1]); err != nil {
			return
		}
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

	cmd.CH = ch
	cmd.Sampling = req.Sampling
	cmd.User = user

	return
}

// skywork 13b
type skyWorkRequest struct {
	Text              string      `json:"text"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

func (req *skyWorkRequest) toCmd(ch chan string, user types.Account) (cmd app.SkyWorkCmd, err error) {
	if cmd.Text, err = domain.NewSkyWorkText(req.Text); err != nil {
		return
	}

	history := make([]domain.History, len(req.History))
	for i := range req.History {
		if history[i], err = domain.NewHistory(req.History[i][0], req.History[i][1]); err != nil {
			return
		}
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

	cmd.CH = ch
	cmd.Sampling = req.Sampling
	cmd.User = user

	return
}

// iflytekspark
type iflyteksparkRequest struct {
	Text              string  `json:"text"`
	Sampling          bool    `json:"sampling"`
	TopK              int     `json:"top_k"`
	Temperature       float64 `json:"temperature"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
}

func (req *iflyteksparkRequest) toCmd(ch chan string, user types.Account) (cmd app.IFlytekSparkCmd, err error) {
	if cmd.Text, err = domain.NewIFlytekSparkText(req.Text); err != nil {
		return
	}

	if req.Sampling {
		if cmd.TopK, err = domain.NewTopK(req.TopK); err != nil {
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

	cmd.CH = ch
	cmd.Sampling = req.Sampling
	cmd.User = user

	return
}
