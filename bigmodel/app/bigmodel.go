package app

import (
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	corepo "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type BigModelService interface {
	// taichu
	DescribePicture(types.Account, io.Reader, string, int64) (string, error)
	GenPicture(types.Account, string) (string, string, error)
	GenPictures(types.Account, string) ([]string, string, error)
	Ask(types.Account, domain.Question, string) (string, string, error)
	VQAUploadPicture(io.Reader, types.Account, string) error

	// luojia
	LuoJiaUploadPicture(io.Reader, types.Account) error
	LuoJia(types.Account) (string, error)
	ListLuoJiaRecord(types.Account) ([]LuoJiaRecordDTO, error)

	// pangu
	PanGu(types.Account, string) (string, string, error)

	// codegeex
	CodeGeex(types.Account, *CodeGeexCmd) (CodeGeexDTO, string, error)

	// wukong
	GenWuKongSamples(int) ([]string, error)
	WuKong(types.Account, *WuKongCmd) (map[string]string, string, error)
	WuKongHF(*WuKongHFCmd) (map[string]string, string, error)
	WuKongInferenceAsync(types.Account, *WuKongCmd) (string, error)
	GetWuKongWaitingTaskRank(types.Account) (WuKongRankDTO, error)
	GetWuKongLastTaskResp(types.Account) (WuKongLinksDTO, string, error)
	AddLikeFromTempPicture(*WuKongAddLikeFromTempCmd) (string, string, error)
	AddLikeFromPublicPicture(*WuKongAddLikeFromPublicCmd) (string, string, error)
	AddPublicFromTempPicture(*WuKongAddPublicFromTempCmd) (string, string, error)
	AddPublicFromLikePicture(*WuKongAddPublicFromLikeCmd) (string, string, error)
	CancelPublic(types.Account, string) error
	GetPublicsGlobal(cmd *WuKongListPublicGlobalCmd) (WuKongPublicGlobalDTO, error)
	ListPublics(types.Account) ([]WuKongPublicDTO, error)
	CancelLike(types.Account, string) error
	ListLikes(types.Account) ([]WuKongLikeDTO, error)
	DiggPicture(*WuKongAddDiggCmd) (int, error)
	CancelDiggPicture(*WuKongCancelDiggCmd) (int, error)
	ReGenerateDownloadURL(types.Account, string) (string, string, error)
}

func NewBigModelService(
	fm bigmodel.BigModel,
	user corepo.User,
	luojia repository.LuoJia,
	wukong repository.WuKong,
	wukongPicture repository.WuKongPicture,
	wukongAsyncRepo repository.WuKongAsyncTask,
	sender message.AsyncMessageProducer,
) BigModelService {
	return bigModelService{
		fm:              fm,
		user:            user,
		sender:          sender,
		luojia:          luojia,
		wukong:          wukong,
		wukongPicture:   wukongPicture,
		wukongAsyncRepo: wukongAsyncRepo,
		wukongSampleId:  fm.GetWuKongSampleId(),
	}
}

type bigModelService struct {
	fm bigmodel.BigModel

	sender          message.AsyncMessageProducer
	user            corepo.User
	luojia          repository.LuoJia
	wukong          repository.WuKong
	wukongPicture   repository.WuKongPicture
	wukongAsyncRepo repository.WuKongAsyncTask

	wukongSampleId string
}

func (s bigModelService) setCode(err error) string {
	if err != nil && bigmodel.IsErrorSensitiveInfo(err) {
		return ErrorBigModelSensitiveInfo
	}

	if err != nil && bigmodel.IsErrorBusySource(err) {
		return ErrorBigModelRecourseBusy
	}

	return ""
}

func (s bigModelService) GenWuKongSamples(batchNum int) ([]string, error) {
	num := s.fm.GenWuKongSampleNums(batchNum)
	if len(num) == 0 {
		return nil, nil
	}

	return s.wukong.ListSamples(s.wukongSampleId, num)
}

func (s bigModelService) WuKong(
	user types.Account, cmd *WuKongCmd,
) (links map[string]string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelWuKong)

	links, err = s.fm.GenPicturesByWuKong(user, (*domain.WuKongPictureMeta)(cmd), string(domain.BigmodelWuKong))
	if err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) WuKongHF(cmd *WuKongHFCmd) (
	links map[string]string, code string, err error,
) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelWuKong)

	links, err = s.fm.GenPicturesByWuKong(cmd.User, (*domain.WuKongPictureMeta)(&cmd.WuKongCmd), string(domain.BigmodelWuKongHF))
	if err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) WuKongInferenceAsync(user types.Account, cmd *WuKongCmd) (code string, err error) {
	// content audit
	if err = s.fm.CheckText(cmd.Desc.WuKongPictureDesc()); err != nil {
		code = ErrorBigModelSensitiveInfo

		return
	}

	msg := new(message.MsgTask)
	msg.ToMsgTask(user.Account(), cmd.Desc.WuKongPictureDesc(), cmd.Style)

	return "", s.sender.CreateWuKongTask(msg)
}

func (s bigModelService) GetWuKongWaitingTaskRank(user types.Account) (dto WuKongRankDTO, err error) {
	t, _ := commondomain.NewTime(time.Now().Add(-300 * time.Second).Unix()) // TODO config

	var rank int
	if rank, err = s.wukongAsyncRepo.GetWaitingTaskRank(user, t); err != nil {
		if !commonrepo.IsErrorResourceNotExists(err) {
			return
		}
	}

	dto = WuKongRankDTO{
		Rank: rank,
	}

	return
}

func (s bigModelService) GetWuKongLastTaskResp(user types.Account) (dto WuKongLinksDTO, code string, err error) {
	p, err := s.wukongAsyncRepo.GetLastFinishedTask(user)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			code = ErrorWuKongNoPicture
			err = errors.New("wukong picture record not found")
		}

		return
	}

	if p.Status.IsError() {
		err = errors.New(p.Links.StringLinks())

		if bigmodel.IsErrorSensitiveInfo(err) {
			code = ErrorBigModelSensitiveInfo
		} else {
			code = ErrorCodeSytem
		}

		return
	}

	if p.Status.IsRunning() {
		code = ErrorCodeSytem
		err = errors.New("task is running, please try it later")

		return
	}

	dto = WuKongLinksDTO{
		Pictures: p.Links.Links(),
	}

	return
}

func (s bigModelService) AddLikeFromTempPicture(cmd *WuKongAddLikeFromTempCmd) (
	pid string, code string, err error,
) {
	meta, p, err := s.fm.CheckWuKongPictureTempToLike(cmd.User, cmd.OBSPath)
	if err != nil {
		code = ErrorWuKongInvalidPath

		return
	}

	v, version, err := s.wukongPicture.ListLikesByUserName(cmd.User)
	if err != nil {
		return
	}
	if len(v) >= 10 { // TODO config
		code = ErrorWuKongExccedMaxLikeNum
		err = errors.New("exceed the max num user can add like to pictures")

		return
	}

	for i := range v {
		if v[i].OBSPath == p {
			code = ErrorWuKongDuplicateLike
			err = errors.New("the picture has been saved")

			return
		}
	}

	if err = s.fm.MoveWuKongPictureToDir(p, cmd.OBSPath); err != nil {
		return
	}

	pid, err = s.wukongPicture.SaveLike(
		cmd.User,
		&domain.WuKongPicture{
			Owner:             cmd.User,
			OBSPath:           p,
			CreatedAt:         utils.Date(),
			WuKongPictureMeta: meta,
		},
		version,
	)

	return
}

func (s bigModelService) AddLikeFromPublicPicture(
	cmd *WuKongAddLikeFromPublicCmd,
) (pid string, code string, err error) {
	p, err := s.wukongPicture.GetPublicByUserName(cmd.Owner, cmd.Id)
	if err != nil {
		code = ErrorWuKongInvalidId
		return
	}

	// gen like path
	likePath, _ := s.fm.CheckWuKongPicturePublicToLike(cmd.User, p.OBSPath)

	// check
	v, version, err := s.wukongPicture.ListLikesByUserName(cmd.User)
	if err != nil {
		return
	}
	if len(v) >= 10 { // TODO config
		code = ErrorWuKongExccedMaxLikeNum
		err = errors.New("exceed the max num user can add like to pictures")

		return
	}

	for i := range v {
		if v[i].OBSPath == likePath {
			code = ErrorWuKongDuplicateLike
			err = errors.New("the picture has been liked")

			return
		}
	}

	// copy picture from public dir to like dir on obs
	if err = s.fm.MoveWuKongPictureToDir(likePath, p.OBSPath); err != nil {
		code = ErrorCodeSytem
		return
	}

	// save
	wp := &domain.WuKongPicture{
		Owner:             p.Owner,
		OBSPath:           likePath,
		CreatedAt:         utils.Date(),
		WuKongPictureMeta: p.WuKongPictureMeta,
	}
	pid, err = s.wukongPicture.SaveLike(cmd.User, wp, version)

	return
}

func (s bigModelService) CancelLike(user types.Account, pid string) (err error) {
	v, err := s.wukongPicture.GetLikeByUserName(user, pid)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	if err = s.fm.DeleteWuKongPicture(v.OBSPath); err != nil {
		return
	}

	err = s.wukongPicture.DeleteLike(user, pid)

	return
}

func (s bigModelService) ListLikes(user types.Account) (
	r []WuKongLikeDTO, err error,
) {
	v, _, err := s.wukongPicture.ListLikesByUserName(user)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]WuKongLikeDTO, len(v))
	for i := range v {
		item := &v[i]
		dto := &r[i]

		dto.Link, err = s.fm.GenWuKongPictureLink(item.OBSPath)
		if err != nil {
			return
		}

		dto.IsPublic, err = s.isPublic(item)
		if err != nil {
			return
		}

		avatar, err := s.user.GetUserAvatarId(item.Owner)
		if err != nil {
			return nil, err
		}

		dto.Owner = item.Owner.Account()
		dto.Avatar = avatar.AvatarId()
		dto.Id = item.Id
		dto.Desc = item.Desc.WuKongPictureDesc()
		dto.Style = item.Style
		dto.CreatedAt = item.CreatedAt
	}

	return
}

func (s bigModelService) AddPublicFromTempPicture(cmd *WuKongAddPublicFromTempCmd) (
	pid string, code string, err error,
) {
	// gen meta and public path
	meta, publicPath, err := s.fm.CheckWuKongPictureToPublic(cmd.User, cmd.OBSPath)
	if err != nil {
		code = ErrorWuKongInvalidPath

		return
	}

	// check
	v, version, err := s.wukongPicture.ListPublicsByUserName(cmd.User)
	if err != nil {
		return
	}

	for i := range v {
		if v[i].OBSPath == publicPath {
			code = ErrorWuKongDuplicateLike
			err = errors.New("the picture has been publiced")

			return
		}
	}

	// copy picture from public dir to like dir on obs
	if err = s.fm.MoveWuKongPictureToDir(publicPath, cmd.OBSPath); err != nil {
		code = ErrorCodeSytem

		return
	}

	// save
	p := &domain.WuKongPicture{
		Owner:             cmd.User,
		OBSPath:           publicPath,
		CreatedAt:         utils.Date(),
		WuKongPictureMeta: meta,
	}
	pid, err = s.wukongPicture.SavePublic(p, version)

	return
}

func (s bigModelService) AddPublicFromLikePicture(cmd *WuKongAddPublicFromLikeCmd) (
	pid string, code string, err error,
) {
	// get like infomation
	p, err := s.wukongPicture.GetLikeByUserName(cmd.User, cmd.Id)
	if err != nil {
		code = ErrorWuKongInvalidPath

		return
	}

	// gen public path
	_, publicPath, err := s.fm.CheckWuKongPictureToPublic(cmd.User, p.OBSPath)
	if err != nil {
		code = ErrorWuKongNoAuthorization

		return
	}

	// check
	v, version, err := s.wukongPicture.ListPublicsByUserName(cmd.User)
	if err != nil {
		return
	}

	for i := range v {
		if v[i].OBSPath == publicPath {
			code = ErrorWuKongDuplicateLike
			err = errors.New("the picture has been publiced")

			return
		}
	}

	// copy picture from public dir to like dir on obs
	if err = s.fm.MoveWuKongPictureToDir(publicPath, p.OBSPath); err != nil {
		code = ErrorCodeSytem

		return
	}

	// save
	ps := &domain.WuKongPicture{
		Owner:             p.Owner,
		OBSPath:           publicPath,
		CreatedAt:         p.CreatedAt,
		WuKongPictureMeta: p.WuKongPictureMeta,
	}
	if pid, err = s.wukongPicture.SavePublic(ps, version); err != nil {
		code = ErrorCodeSytem

		return
	}

	return
}

func (s bigModelService) CancelPublic(user types.Account, pid string) (err error) {
	v, err := s.wukongPicture.GetPublicByUserName(user, pid)
	if err != nil {
		return
	}

	if err = s.wukongPicture.DeletePublic(v.Owner, v.Id); err != nil {
		return
	}

	s.fm.DeleteWuKongPicture(v.OBSPath)

	return
}

func (s bigModelService) GetPublicsGlobal(cmd *WuKongListPublicGlobalCmd) (r WuKongPublicGlobalDTO, err error) {
	var v []domain.WuKongPicture
	if cmd.Level != nil && cmd.Level.IsOfficial() {
		v, err = s.wukongPicture.GetOfficialPublicsGlobal()
	} else {
		v, err = s.wukongPicture.GetPublicsGlobal()
	}
	if err != nil {
		return
	}

	var b, e int
	if b = cmd.CountPerPage * (cmd.PageNum - 1); b >= len(v) {
		err = errors.New("paginator error")

		return
	}
	if e = b + cmd.CountPerPage; e > len(v) {
		e = len(v)
	}

	d := make([]WuKongPublicDTO, len(v[b:e]))
	for i := range v[b:e] {
		item := &v[b:e][i]
		link := s.fm.GenWuKongLinkFromOBSPath(item.OBSPath)
		avatarId, _ := s.user.GetUserAvatarId(item.Owner)

		var (
			a       string
			isDigg  bool
			LikeDto WuKongIsLikeDTO
		)

		if cmd.User != nil {
			LikeDto, _ = s.isLike(item, cmd.User)
			isDigg = s.isDigg(cmd.User, item.Diggs)
		} else {
			isDigg = false
		}

		if avatarId != nil {
			a = avatarId.AvatarId()
		} else {
			a = ""
		}

		d[i].toWuKongPublicDTO(item, a, LikeDto.IsLike, LikeDto.LikeID, isDigg, link)
	}

	r = WuKongPublicGlobalDTO{
		Total:    len(v),
		Pictures: d,
	}

	return
}

func (s bigModelService) ListPublics(user types.Account) (
	r []WuKongPublicDTO, err error,
) {
	v, _, err := s.wukongPicture.ListPublicsByUserName(user)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]WuKongPublicDTO, len(v))
	for i := range v {
		item := &v[i]
		dto := &r[i]

		link := s.fm.GenWuKongLinkFromOBSPath(item.OBSPath)
		likeDto, _ := s.isLike(item, user)
		isDigg := s.isDigg(user, item.Diggs)

		dto.toWuKongPublicDTO(item, "", likeDto.IsLike, likeDto.LikeID, isDigg, link)
	}

	return
}

func (s bigModelService) DiggPicture(cmd *WuKongAddDiggCmd) (count int, err error) {
	// get picture info
	p, err := s.wukongPicture.GetPublicByUserName(cmd.Owner, cmd.Id)
	if err != nil {
		return
	}

	// insert digg user and update diggcount
	diggs := p.Diggs
	for _, user := range diggs {
		if user == cmd.User.Account() {
			err = errors.New("the picture had been digged")

			return
		}
	}
	p.Diggs = append(p.Diggs, cmd.User.Account())
	p.DiggCount = len(p.Diggs)

	// save
	if err = s.wukongPicture.UpdatePublicPicture(p.Owner, p.Id, p.Version, &p); err != nil {
		return
	}

	count = p.DiggCount
	return
}

func (s bigModelService) CancelDiggPicture(cmd *WuKongCancelDiggCmd) (count int, err error) {
	// get picture info
	p, err := s.wukongPicture.GetPublicByUserName(cmd.Owner, cmd.Id)
	if err != nil {
		return
	}

	// delete digg user and update diggcount
	f := func(arr []string, s string) []string {
		i := 0
		for _, v := range arr {
			if v != s {
				arr[i] = v
				i++
			}
		}
		return arr[:i]
	}

	l := len(p.Diggs)
	p.Diggs = f(p.Diggs, cmd.User.Account())
	if l == len(p.Diggs) {
		err = errors.New("user not digg this picture")
		return
	}
	p.DiggCount = len(p.Diggs)

	// save
	if err = s.wukongPicture.UpdatePublicPicture(p.Owner, p.Id, p.Version, &p); err != nil {
		return
	}

	count = p.DiggCount
	return
}

func (s bigModelService) isLike(
	p *domain.WuKongPicture,
	user types.Account,
) (WuKongIsLikeDTO, error) {
	pics, _, err := s.wukongPicture.ListLikesByUserName(user)
	if err != nil {
		return WuKongIsLikeDTO{}, err
	}

	for _, pic := range pics {
		likePath, err := s.fm.CheckWuKongPicturePublicToLike(user, p.OBSPath)
		if err != nil {
			return WuKongIsLikeDTO{}, err
		}

		if pic.OBSPath == likePath {
			return WuKongIsLikeDTO{
				IsLike: true,
				LikeID: pic.Id,
			}, nil
		}
	}

	return WuKongIsLikeDTO{}, nil
}

func (s bigModelService) isPublic(
	p *domain.WuKongPicture,
) (bool, error) {
	pics, _, err := s.wukongPicture.ListPublicsByUserName(p.Owner)
	if err != nil {
		return false, err
	}

	for _, pic := range pics {
		_, publicPath, err := s.fm.CheckWuKongPictureToPublic(p.Owner, p.OBSPath)
		if err != nil {
			return false, err
		}

		if pic.OBSPath == publicPath {
			return true, nil
		}
	}

	return false, nil
}

func (s bigModelService) isDigg(
	user types.Account,
	diggs []string,
) bool {
	for _, username := range diggs {
		if user.Account() == username {
			return true
		}
	}

	return false
}

func (s bigModelService) ReGenerateDownloadURL(
	user types.Account, link string,
) (
	newLink string, code string, err error,
) {
	v, err := url.Parse(link)
	if err != nil || !strings.Contains(v.Path, user.Account()) {
		code = ErrorWuKongInvalidLink
		err = errors.New("invalid link")
	} else {
		newLink, err = s.fm.GenWuKongPictureLink(
			strings.TrimPrefix(v.Path, "/"),
		)
	}

	return
}
