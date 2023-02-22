package app

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type BigModelService interface {
	DescribePicture(domain.Account, io.Reader, string, int64) (string, error)
	GenPicture(domain.Account, string) (string, string, error)
	GenPictures(domain.Account, string) ([]string, string, error)
	Ask(domain.Account, domain.Question, string) (string, string, error)
	VQAUploadPicture(io.Reader, domain.Account, string) error
	LuoJiaUploadPicture(io.Reader, domain.Account) error
	PanGu(domain.Account, string) (string, string, error)
	CodeGeex(domain.Account, *CodeGeexCmd) (CodeGeexDTO, string, error)
	LuoJia(domain.Account) (string, error)
	ListLuoJiaRecord(domain.Account) ([]LuoJiaRecordDTO, error)
	GenWuKongSamples(int) ([]string, error)
	WuKong(domain.Account, *WuKongCmd) (map[string]string, string, error)
	AddLikeFromTempPicture(*WuKongAddLikeFromTempCmd) (string, string, error)
	AddLikeFromPublicPicture(*WuKongAddLikeFromPublicCmd) (string, string, error)
	AddPublicFromTempPicture(*WuKongAddPublicFromTempCmd) (string, string, error)
	AddPublicFromLikePicture(*WuKongAddPublicFromLikeCmd) (string, string, error)
	CancelPublic(domain.Account, string) error
	GetPublicsGlobal(cmd *WuKongListPublicGlobalCmd) (WuKongPublicGlobalDTO, error)
	ListPublics(domain.Account) ([]WuKongPublicDTO, error)
	CancelLike(domain.Account, string) error
	ListLikes(domain.Account) ([]WuKongLikeDTO, error)
	DiggPicture(*WuKongAddDiggCmd) (int, error)
	CancelDiggPicture(*WuKongCancelDiggCmd) (int, error)
	ReGenerateDownloadURL(domain.Account, string) (string, string, error)
}

func NewBigModelService(
	fm bigmodel.BigModel,
	user repository.User,
	luojia repository.LuoJia,
	wukong repository.WuKong,
	wukongPicture repository.WuKongPicture,
	sender message.Sender,
) BigModelService {
	return bigModelService{
		fm:             fm,
		user:           user,
		sender:         sender,
		luojia:         luojia,
		wukong:         wukong,
		wukongPicture:  wukongPicture,
		wukongSampleId: fm.GetWuKongSampleId(),
	}
}

type bigModelService struct {
	fm bigmodel.BigModel

	sender        message.Sender
	user          repository.User
	luojia        repository.LuoJia
	wukong        repository.WuKong
	wukongPicture repository.WuKongPicture

	wukongSampleId string
}

func (s bigModelService) DescribePicture(
	user domain.Account, picture io.Reader, name string, length int64,
) (string, error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelDescPicture)

	return s.fm.DescribePicture(picture, name, length)
}

func (s bigModelService) GenPicture(
	user domain.Account, desc string,
) (link string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelGenPicture)

	if link, err = s.fm.GenPicture(user, desc); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) GenPictures(
	user domain.Account, desc string,
) (links []string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelGenPicture)

	if links, err = s.fm.GenPictures(user, desc); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) Ask(
	u domain.Account, q domain.Question, f string,
) (v string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(u, domain.BigmodelVQA)

	if v, err = s.fm.Ask(q, f); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) VQAUploadPicture(f io.Reader, user domain.Account, fileName string) error {
	return s.fm.VQAUploadPicture(f, user, fileName)
}

func (s bigModelService) LuoJiaUploadPicture(f io.Reader, user domain.Account) error {
	return s.fm.LuoJiaUploadPicture(f, user)
}

func (s bigModelService) PanGu(u domain.Account, q string) (v string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(u, domain.BigmodelPanGu)

	if v, err = s.fm.PanGu(q); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) LuoJia(user domain.Account) (v string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelLuoJia)

	if v, err = s.fm.LuoJia(user.Account()); err != nil {
		return
	}

	record := domain.UserLuoJiaRecord{User: user}
	record.CreatedAt = utils.Now()

	s.luojia.Save(&record)

	return
}

func (s bigModelService) ListLuoJiaRecord(user domain.Account) (
	dtos []LuoJiaRecordDTO, err error,
) {
	v, err := s.luojia.List(user)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = append(dtos, LuoJiaRecordDTO{
		CreatedAt: utils.ToDate(v[0].CreatedAt),
		Id:        v[0].Id,
	})

	return
}

func (s bigModelService) CodeGeex(user domain.Account, cmd *CodeGeexCmd) (
	dto CodeGeexDTO, code string, err error,
) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelCodeGeex)

	if dto, err = s.fm.CodeGeex((*bigmodel.CodeGeexReq)(cmd)); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) setCode(err error) string {
	if err != nil && bigmodel.IsErrorSensitiveInfo(err) {
		return ErrorBigModelSensitiveInfo
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
	user domain.Account, cmd *WuKongCmd,
) (links map[string]string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelWuKong)

	links, err = s.fm.GenPicturesByWuKong(user, (*domain.WuKongPictureMeta)(cmd))
	if err != nil {
		code = s.setCode(err)
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
	if len(v) >= appConfig.WuKongMaxLikeNum {
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
	if len(v) >= appConfig.WuKongMaxLikeNum {
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

func (s bigModelService) CancelLike(user domain.Account, pid string) (err error) {
	v, err := s.wukongPicture.GetLikeByUserName(user, pid)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
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

func (s bigModelService) ListLikes(user domain.Account) (
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

		dto.Owner = item.Owner.Account()
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
	pid, err = s.wukongPicture.SavePublic(ps, version)

	return
}

func (s bigModelService) CancelPublic(user domain.Account, pid string) (err error) {
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
	if cmd.Level != nil && cmd.Level.WuKongPictureLevel() == "official" {
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
	v = v[b:e]

	d := make([]WuKongPublicDTO, len(v))
	for i := range v {
		item := &v[i]
		link := s.fm.GenWuKongLinkFromOBSPath(item.OBSPath)
		avatarId, _ := s.user.GetUserAvatarId(item.Owner)

		var (
			isDigg  bool
			LikeDto WuKongIsLikeDTO
		)
		if cmd.User != nil {
			LikeDto, _ = s.isLike(item, cmd.User)
			isDigg = s.isDigg(cmd.User, item.Diggs)
		} else {
			isDigg = false
		}

		d[i].toWuKongPublicDTO(item, avatarId.AvatarId(), LikeDto.IsLike, LikeDto.LikeID, isDigg, link)
	}

	r = WuKongPublicGlobalDTO{
		Total:    len(d),
		Pictures: d,
	}

	return
}

func (s bigModelService) ListPublics(user domain.Account) (
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
	user domain.Account,
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
	user domain.Account,
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
	user domain.Account, link string,
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
