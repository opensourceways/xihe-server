package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func wukongOwnerFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func wukongIdFilter(id string) bson.M {
	return bson.M{
		fieldId: id,
	}
}

func NewWuKongPictureMapper(name string) repositories.WuKongPictureMapper {
	return wukongPicture{
		collectionName: name,
	}
}

type wukongPicture struct {
	collectionName string
}

func (col wukongPicture) GetVersion(user string) (version int, err error) {
	v := new(dWuKongPicture)

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
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

func (col wukongPicture) ListLikesByUserName(user string) ([]repositories.WuKongPictureDO, int, error) {
	return col.listFieldNameByUserName(user, fieldLikes)
}

func (col wukongPicture) ListPublicsByUserName(user string) ([]repositories.WuKongPictureDO, int, error) {
	return col.listFieldNameByUserName(user, fieldPublics)
}

func (col wukongPicture) listFieldNameByUserName(user, fieldName string) ([]repositories.WuKongPictureDO, int, error) {
	var v dWuKongPicture

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user),
			nil, &v,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = nil
		}

		return nil, 0, err
	}

	var t []pictureItem
	if fieldName == fieldLikes {
		t = v.Likes
	} else {
		t = v.Publics
	}
	r := make([]repositories.WuKongPictureDO, len(t))

	for i := range t {
		col.toPictureDO(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col wukongPicture) InsertIntoLikes(user string, do *repositories.WuKongPictureDO, version int) (
	identity string, err error,
) {
	return col.insertIntoFieldName(user, do, version, fieldLikes)
}

func (col wukongPicture) InsertIntoPublics(user string, do *repositories.WuKongPictureDO, version int) (
	identity string, err error,
) {
	return col.insertIntoFieldName(user, do, version, fieldPublics)
}

func (col wukongPicture) insertIntoFieldName(
	user string, do *repositories.WuKongPictureDO,
	version int, fieldName string,
) (
	identity string, err error,
) {
	identity, err = col.insert(user, do, version, fieldName)
	if err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert
	if err = col.newDoc(user); err != nil {
		return
	}

	identity, err = col.insert(user, do, version, fieldName)
	if err != nil && isDocNotExists(err) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (col wukongPicture) newDoc(user string) error {
	docFilter := resourceOwnerFilter(user)

	doc := bson.M{
		fieldOwner:   user,
		fieldLikes:   bson.A{},
		fieldPublics: bson.A{},
		fieldVersion: 0,
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && !isDocExists(err) {
		return err
	}

	return nil
}

func (col wukongPicture) insert(
	user string,
	do *repositories.WuKongPictureDO,
	version int,
	filedName string,
) (
	identity string, err error,
) {
	identity = newId()
	do.Id = identity

	doc, err := col.toPictureDoc(do)
	doc[fieldVersion] = do.Version
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user),
			bson.M{filedName: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col wukongPicture) deleteFieldName(user, pid, fieldName string) error {
	f := func(ctx context.Context) error {
		return cli.pullArrayElem(
			ctx, col.collectionName, fieldName,
			resourceOwnerFilter(user),
			resourceIdFilter(pid),
		)
	}

	return withContext(f)
}

func (col wukongPicture) DeleteLike(user string, pid string) error {
	return col.deleteFieldName(user, pid, fieldLikes)
}

func (col wukongPicture) DeletePublic(user string, pid string) error {
	return col.deleteFieldName(user, pid, fieldPublics)
}

func (col wukongPicture) getByUserName(user, pid, field string) (
	do repositories.WuKongPictureDO,
	err error,
) {
	var v []dWuKongPicture

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, field,
			resourceOwnerFilter(user),
			resourceIdFilter(pid),
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
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toPictureDO(&l[0], &do)

	return
}

func (col wukongPicture) GetLikeByUserName(user string, pid string) (
	do repositories.WuKongPictureDO,
	err error,
) {
	return col.getByUserName(user, pid, fieldLikes)
}

func (col wukongPicture) GetPublicByUserName(user string, pid string) (
	do repositories.WuKongPictureDO,
	err error,
) {
	return col.getByUserName(user, pid, fieldPublics)
}

func (col wukongPicture) GetPublicsGlobal() (do []repositories.WuKongPictureDO, err error) {
	var v []dWuKongPicture

	f := func(ctx context.Context) error {
		opts := options.FindOptions{}
		return cli.getDocs(
			ctx, col.collectionName, nil,
			opts.SetProjection(bson.M{
				fieldPublics: 1,
			}), &v,
		)
	}

	if err := withContext(f); err != nil {
		return nil, err
	}

	return col.toArrayPictureDO(v), nil
}

func (col wukongPicture) UpdatePublicPicture(
	user string, pid string, version int,
	do *repositories.WuKongPictureDO,
) (err error) {
	doc, err := col.toPictureDoc(do)
	if err != nil {
		return
	}

	updated := false

	f := func(ctx context.Context) error {
		updated, err = cli.updateArrayElem(
			ctx, col.collectionName, fieldPublics,
			wukongOwnerFilter(user),
			wukongIdFilter(pid),
			doc, version, 0,
		)

		return err
	}

	if err = withContext(f); err != nil {
		return
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(
			errors.New("no update"),
		)
	}

	return
}

func (col wukongPicture) toPictureDO(p *pictureItem, do *repositories.WuKongPictureDO) {
	*do = repositories.WuKongPictureDO{
		Id:        p.Id,
		Owner:     p.Owner,
		OBSPath:   p.OBSPath,
		Level:     p.Level,
		Diggs:     p.Diggs,
		DiggCount: p.DiggCount,
		Version:   p.Version,
		CreatedAt: p.CreatedAt,
	}

	do.Desc = p.Desc
	do.Style = p.Style
}

func (col wukongPicture) toPictureDoc(do *repositories.WuKongPictureDO) (bson.M, error) {
	return genDoc(pictureItem{
		Id:        do.Id,
		Owner:     do.Owner,
		Desc:      do.Desc,
		Style:     do.Style,
		OBSPath:   do.OBSPath,
		Diggs:     do.Diggs,
		DiggCount: do.DiggCount,
		CreatedAt: do.CreatedAt,
	})
}

func (col wukongPicture) toArrayPictureDO(p []dWuKongPicture) []repositories.WuKongPictureDO {
	var length int
	for i := range p {
		length += len(p[i].Publics)
	}

	r := make([]repositories.WuKongPictureDO, length)
	var count = 0
	for i := range p {
		for j := range p[i].Publics {
			col.toPictureDO(&p[i].Publics[j], &r[count])
			count++
		}
	}

	return r
}
