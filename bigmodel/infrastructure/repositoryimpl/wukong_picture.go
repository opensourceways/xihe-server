package repositoryimpl

import (
	"context"
	"errors"
	"sort"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	commoninfra "github.com/opensourceways/xihe-server/common/infrastructure"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func NewWuKongPictureRepo(m mongodbClient) repository.WuKongPicture {
	return &wukongPictureRepoImpl{m}
}

type wukongPictureRepoImpl struct {
	cli mongodbClient
}

func (impl *wukongPictureRepoImpl) GetVersion(user types.Account) (version int, err error) {
	v := new(dWuKongPicture)

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx,
			bson.M{fieldOwner: user},
			bson.M{fieldVersion: 1},
			v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	version = v.Version
	return
}

func (impl *wukongPictureRepoImpl) ListLikesByUserName(user types.Account) (
	[]domain.WuKongPicture, int, error,
) {
	v, version, err := impl.listFieldNameByUserName(user.Account(), fieldLikes)
	if err != nil {
		return nil, 0, err
	}

	return v, version, nil
}

func (impl *wukongPictureRepoImpl) ListPublicsByUserName(user types.Account) (
	[]domain.WuKongPicture, int, error,
) {
	v, version, err := impl.listFieldNameByUserName(user.Account(), fieldPublics)
	if err != nil {
		return nil, 0, err
	}

	sortWuKongPictureByTime(v)

	return v, version, nil
}

func (impl *wukongPictureRepoImpl) listFieldNameByUserName(user, fieldName string) ([]domain.WuKongPicture, int, error) {
	var v dWuKongPicture

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx,
			wukongOwnerFilter(user),
			nil, &v,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = nil
		}

		return nil, 0, err
	}

	var t []pictureItem

	switch fieldName {
	case fieldLikes:
		t = v.Likes
	case fieldPublics:
		t = v.Publics
	}

	r := make([]domain.WuKongPicture, len(t))

	for i := range t {
		t[i].toWuKongPicture(&r[i])
	}

	return r, v.Version, nil
}

func (impl *wukongPictureRepoImpl) SaveLike(user types.Account, p *domain.WuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	p.SetDefaultDiggs()

	v, err := impl.insertIntoFieldName(user.Account(), p, version, fieldLikes)
	if err != nil {
		return "", commoninfra.ConvertError(err)
	}

	return v, nil
}

func (impl *wukongPictureRepoImpl) SavePublic(p *domain.WuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	p.SetDefaultDiggs()

	v, err := impl.insertIntoFieldName(p.Owner.Account(), p, version, fieldPublics)
	if err != nil {
		return "", commoninfra.ConvertError(err)
	}

	return v, nil
}

func (impl *wukongPictureRepoImpl) insertIntoFieldName(
	user string, d *domain.WuKongPicture,
	version int, fieldName string,
) (
	identity string, err error,
) {

	f := func(ctx context.Context) (err error) {
		// 1. try to create user in wukong picture
		doc, err := toWuKongPictureEmptyLikesPublicsDoc(user, 0)
		if err != nil {
			return
		}

		doc[fieldVersion] = 0
		doc[fieldLikes] = bson.A{}
		doc[fieldPublics] = bson.A{}

		filter := bson.M{
			fieldOwner: user,
		}

		if _, err = impl.cli.NewDocIfNotExist(ctx, filter, doc); err != nil {
			if !impl.cli.IsDocExists(err) {

				return
			}
		}

		// 2. insert picture under order user
		if identity, err = impl.insert(user, d, version, fieldName); err != nil {
			return
		}

		return
	}

	if err = withContext(f); err != nil {

		return
	}

	return
}

func (impl *wukongPictureRepoImpl) insert(
	user string, d *domain.WuKongPicture,
	version int, filedName string,
) (
	identity string, err error,
) {
	identity = newId()
	d.Id = identity

	doc, err := toPictureItemDoc(d)
	if err != nil {
		return
	}

	doc[fieldVersion] = d.Version

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx,
			wukongOwnerFilter(user),
			bson.M{filedName: doc}, mongoCmdPush, version,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	return
}

func (impl *wukongPictureRepoImpl) DeleteLike(user types.Account, pid string) error {
	if err := impl.deleteFieldName(user.Account(), pid, fieldLikes); err != nil {
		return commoninfra.ConvertError(err)
	}

	return nil
}

func (impl *wukongPictureRepoImpl) DeletePublic(user types.Account, pid string) error {
	if err := impl.deleteFieldName(user.Account(), pid, fieldPublics); err != nil {
		return commoninfra.ConvertError(err)
	}

	return nil
}

func (impl *wukongPictureRepoImpl) deleteFieldName(user, pid, fieldName string) error {
	f := func(ctx context.Context) error {
		return impl.cli.PullArrayElem(
			ctx, fieldName,
			wukongOwnerFilter(user),
			wukongIdFilter(pid),
		)
	}

	return withContext(f)
}

func (impl *wukongPictureRepoImpl) GetLikeByUserName(user types.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if p, err = impl.getByUserName(user.Account(), pid, fieldLikes); err != nil {
		err = commoninfra.ConvertError(err)

		return
	}

	return
}

func (impl *wukongPictureRepoImpl) GetPublicByUserName(user types.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if p, err = impl.getByUserName(user.Account(), pid, fieldPublics); err != nil {
		err = commoninfra.ConvertError(err)

		return
	}

	return
}

func (impl *wukongPictureRepoImpl) getByUserName(user, pid, field string) (
	p domain.WuKongPicture,
	err error,
) {
	var v []dWuKongPicture

	f := func(ctx context.Context) error {
		return impl.cli.GetArrayElem(
			ctx, field,
			wukongOwnerFilter(user),
			wukongIdFilter(pid),
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	var l []pictureItem
	if field == fieldLikes {
		l = v[0].Likes
	} else {
		l = v[0].Publics
	}

	if len(v) == 0 || len(l) == 0 {
		err = commoninfra.NewErrorDataNotExists(errDocNotExists)

		return
	}

	l[0].toWuKongPicture(&p)

	return
}

func (impl *wukongPictureRepoImpl) GetPublicsGlobal() (r []domain.WuKongPicture, err error) {

	var v []dWuKongPicture

	f := func(ctx context.Context) error {
		project := bson.M{
			fieldPublics: 1,
		}

		return impl.cli.GetDocs(
			ctx, nil, project, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	// convert
	// count amount of publics picture
	f2 := func(v []dWuKongPicture) (c int) {
		for i := range v {
			c += len(v[i].Publics)
		}

		return
	}

	r = make([]domain.WuKongPicture, f2(v))
	var c int
	for i := range v {
		for j := range v[i].Publics {
			v[i].Publics[j].toWuKongPicture(&r[c])
			c++
		}
	}

	sortWuKongPictureByTime(r)

	return
}

func (impl *wukongPictureRepoImpl) GetOfficialPublicsGlobal() (r []domain.WuKongPicture, err error) {
	var d []domain.WuKongPicture
	if d, err = impl.GetPublicsGlobal(); err != nil {
		return
	}

	// filter: official picture
	for i := range d {
		if d[i].IsOfficial() {
			r = append(r, d[i])
		}
	}

	sortWuKongPictureByTime(d)

	return
}

func (impl *wukongPictureRepoImpl) UpdatePublicPicture(
	user types.Account, pid string, version int,
	d *domain.WuKongPicture,
) (err error) {

	doc, err := toPictureItemDoc(d)
	if err != nil {
		return
	}

	var updated bool
	f := func(ctx context.Context) error {
		updated, err = impl.cli.UpdateArrayElem(
			ctx, fieldPublics,
			wukongOwnerFilter(user.Account()),
			wukongIdFilter(pid),
			doc, version, 0,
		)

		return err
	}

	if err = withContext(f); err != nil {
		return
	}

	if !updated {
		return commoninfra.NewErrorConcurrentUpdating(errDocNoUpdate)
	}

	return
}

func (r *pictureItem) toWuKongPicture(d *domain.WuKongPicture) (err error) {
	if err = r.toWuKongPictureMeta(&d.WuKongPictureMeta); err != nil {
		return
	}

	if d.Owner, err = types.NewAccount(r.Owner); err != nil {
		return
	}

	if d.OBSPath, err = domain.NewOBSPath(r.OBSPath); err != nil {
		return
	}

	d.Level = domain.NewWuKongPictureLevelByNum(r.Level)
	d.Id = r.Id
	d.Diggs = r.Diggs
	d.DiggCount = r.DiggCount
	d.Version = r.Version
	d.CreatedAt = r.CreatedAt

	return
}

func (r *pictureItem) toWuKongPictureMeta(d *domain.WuKongPictureMeta) (err error) {
	if d.Desc, err = domain.NewWuKongPictureDesc(r.Desc); err != nil {
		return
	}

	d.Style = r.Style

	return
}

func toWuKongPictureEmptyLikesPublicsDoc(owner string, version int) (bson.M, error) {
	return genDoc(dWuKongPicture{
		Owner:   owner,
		Version: version,
		Likes:   []pictureItem{},
		Publics: []pictureItem{},
	})
}

func toPictureItemDoc(d *domain.WuKongPicture) (bson.M, error) {
	p := pictureItem{
		Id:        d.Id,
		Style:     d.Style,
		Diggs:     d.Diggs,
		DiggCount: d.DiggCount,
		Version:   d.Version,
		CreatedAt: d.CreatedAt,
	}

	if d.Owner != nil {
		p.Owner = d.Owner.Account()
	}

	if d.Desc != nil {
		p.Desc = d.Desc.WuKongPictureDesc()
	}

	if d.Level != nil {
		p.Level = d.Level.Int()
	}

	if d.OBSPath != nil {
		p.OBSPath = d.OBSPath.OBSPath()
	}

	return genDoc(p)
}

func sortWuKongPictureByTime(p []domain.WuKongPicture) {
	sort.Slice(p, func(i, j int) bool {
		ti, _ := utils.ToUnixTime(p[i].CreatedAt)
		tj, _ := utils.ToUnixTime(p[j].CreatedAt)

		return ti.After(tj)
	})
}
