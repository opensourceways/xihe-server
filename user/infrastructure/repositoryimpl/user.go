package repositoryimpl

import (
	"context"
	"errors"
	"sort"

	mongo "github.com/opensourceways/xihe-server/common/infrastructure/mongo"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	typesrepo "github.com/opensourceways/xihe-server/infrastructure/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

func NewUserRepo(m mongodbClient) repository.User {
	return &userRepoImpl{m}
}

type userRepoImpl struct {
	cli mongodbClient
}

func (impl *userRepoImpl) Save(u *domain.User) (r domain.User, err error) {
	if u.Id != "" {
		if err = impl.update(u); err != nil {
			err = typesrepo.ConvertError(err)
		} else {
			r = *u
			r.Version += 1
		}

		return
	}

	v, err := impl.insert(u)
	if err != nil {
		err = typesrepo.ConvertError(err)
	} else {
		r = *u
		r.Id = v
	}

	return
}

func (impl *userRepoImpl) GetByAccount(account domain.Account) (r domain.User, err error) {
	if r, _, err = impl.GetByFollower(account, nil); err != nil {
		err = typesrepo.ConvertError(err)

		return
	}

	return
}

func (impl *userRepoImpl) update(u *domain.User) (err error) {
	var user DUser
	toUserDoc(*u, &user)
	doc, err := mongo.GenDoc(user)
	if err != nil {
		return
	}

	filter, err := mongo.ObjectIdFilter(u.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter, doc, mongoCmdSet, u.Version,
		)
	}

	if err = withContext(f); err != nil && impl.cli.IsDocNotExists(err) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (impl *userRepoImpl) insert(u *domain.User) (id string, err error) {
	var user DUser
	toUserDoc(*u, &user)
	doc, err := mongo.GenDoc(user)
	if err != nil {
		return
	}

	doc[fieldVersion] = 0
	doc[fieldFollower] = bson.A{}
	doc[fieldFollowing] = bson.A{}

	f := func(ctx context.Context) error {
		v, err := impl.cli.NewDocIfNotExist(
			ctx, mongo.UserDocFilterByAccount(u.Account.Account()), doc,
		)

		id = v

		return err
	}

	if err = withContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (impl *userRepoImpl) GetByFollower(owner, follower domain.Account) (
	u domain.User, isFollower bool, err error,
) {
	var v []struct {
		DUser `bson:",inline"`

		IsFollower     bool `bson:"is_follower"`
		FollowerCount  int  `bson:"follower_count"`
		FollowingCount int  `bson:"following_count"`
	}

	f := func(ctx context.Context) error {
		fields := bson.M{
			fieldFollowerCount:  bson.M{"$size": "$" + fieldFollower},
			fieldFollowingCount: bson.M{"$size": "$" + fieldFollowing},
		}

		if follower != nil {
			fields[fieldIsFollower] = bson.M{
				"$in": bson.A{follower.Account(), "$" + fieldFollower},
			}
		}

		pipeline := bson.A{
			bson.M{"$match": mongo.UserDocFilterByAccount(owner.Account())},
			bson.M{"$addFields": fields},
			bson.M{"$project": bson.M{
				fieldFollowing: 0,
				fieldFollower:  0,
			}},
		}

		cursor, err := impl.cli.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = repositories.NewErrorDataNotExists(errors.New("no user"))

		return
	}

	item := &v[0]
	if err = toUser(item.DUser, &u); err != nil {
		return
	}

	u.FollowerCount = item.FollowerCount
	u.FollowingCount = item.FollowingCount

	if follower != nil {
		isFollower = item.IsFollower
	}

	return
}

func (impl *userRepoImpl) FindUsersInfo(accounts []domain.Account) (r []domain.UserInfo, err error) {
	var v []DUser

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldName: bson.M{
				"$in": accounts,
			},
		}

		opts := options.FindOptions{}
		return impl.cli.GetDocs(
			ctx, filter, opts.SetProjection(bson.M{
				fieldName:     1,
				fieldAvatarId: 1,
			}), &v,
		)
	}

	if err := withContext(f); err != nil {
		err = typesrepo.ConvertError(err)

		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	r = make([]domain.UserInfo, len(v))
	for i := range v {
		if err := toUserInfo(v[i], &r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl *userRepoImpl) GetUserAvatarId(account domain.Account) (id domain.AvatarId, err error) {

	var v DUser

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, bson.M{fieldName: account},
			bson.M{fieldAvatarId: 1}, &v,
		)
	}

	if err := withContext(f); err != nil {
		err = typesrepo.ConvertError(err)

		return nil, err
	}

	if id, err = domain.NewAvatarId(v.AvatarId); err != nil {
		return
	}

	return
}

func (impl *userRepoImpl) Search(opt *repository.UserSearchOption) (
	r repository.UserSearchResult, err error,
) {
	key := "$" + fieldFollower
	fieldCount := "count"
	fields := bson.M{fieldCount: bson.M{
		"$cond": bson.M{
			"if":   bson.M{"$isArray": key},
			"then": bson.M{"$size": key},
			"else": 0,
		},
	}}

	pipeline := bson.A{
		bson.M{mongo.MongoCmdMatch: bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$" + fieldName,
					"regex":   opt.Name,
					"options": "i",
				},
			},
		}},
		bson.M{"$addFields": fields},
		bson.M{"$project": bson.M{
			fieldName:  1,
			fieldCount: 1,
		}},
	}

	var v []userInfo

	f := func(ctx context.Context) error {
		col := impl.cli.Collection()
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	total := len(v)

	items := make([]*userInfo, total)
	for i := range v {
		items[i] = &v[i]
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count >= items[j].Count
	})

	n := opt.TopNum
	if total < n {
		n = total
	}

	top := make([]domain.Account, n)
	for i := range top {
		if top[i], err = domain.NewAccount(items[i].Name); err != nil {
			return
		}
	}

	r.Top = top
	r.Total = n

	return
}

func (impl *userRepoImpl) AddFollowing(v *domain.FollowerInfo) error {
	err := impl.addFollow(v.Follower.Account(), v.User.Account(), fieldFollowing)
	if err != nil {
		return typesrepo.ConvertError(err)
	}

	return nil
}

func (impl *userRepoImpl) AddFollower(v *domain.FollowerInfo) error {
	err := impl.addFollow(v.User.Account(), v.Follower.Account(), fieldFollower)
	if err != nil {
		return typesrepo.ConvertError(err)
	}

	return nil
}

func (impl *userRepoImpl) addFollow(user, account, field string) error {
	f := func(ctx context.Context) error {
		return impl.cli.AddToSimpleArray(
			ctx, field,
			mongo.UserDocFilterByAccount(user), account,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}

func (impl *userRepoImpl) RemoveFollowing(v *domain.FollowerInfo) error {
	err := impl.removeFollow(v.Follower.Account(), v.User.Account(), fieldFollowing)
	if err != nil {
		return typesrepo.ConvertError(err)
	}

	return nil
}

func (impl *userRepoImpl) RemoveFollower(v *domain.FollowerInfo) error {
	err := impl.removeFollow(v.User.Account(), v.Follower.Account(), fieldFollower)
	if err != nil {
		return typesrepo.ConvertError(err)
	}

	return nil
}

func (impl *userRepoImpl) removeFollow(user, account, field string) error {
	f := func(ctx context.Context) error {
		return impl.cli.RemoveFromSimpleArray(
			ctx, field,
			mongo.UserDocFilterByAccount(user), account,
		)
	}

	return withContext(f)
}

func (impl *userRepoImpl) FindFollowing(owner domain.Account, option *repository.FollowFindOption) (
	info repository.FollowerUserInfos, err error,
) {
	v, err := impl.getFollows(owner.Account(), fieldFollowing)
	if err != nil {
		err = typesrepo.ConvertError(err)
		return
	}

	items := v.Following
	total := len(items)
	if total == 0 {
		return
	}

	items = impl.getPageItems(items, option)
	if len(items) == 0 {
		return
	}

	if option.Follower == nil {
		info.Users, err = impl.listFollowsDirectly(items, false)
	} else if owner == option.Follower {
		info.Users, err = impl.listFollowsDirectly(items, true)
	} else {
		info.Users, err = impl.listFollows(option.Follower.Account(), items)
	}

	info.Total = total

	return
}

func (impl *userRepoImpl) FindFollower(owner domain.Account, option *repository.FollowFindOption) (
	info repository.FollowerUserInfos, err error,
) {
	v, err := impl.getFollows(owner.Account(), fieldFollower)
	if err != nil {
		err = typesrepo.ConvertError(err)
		return
	}

	items := v.Follower
	total := len(items)
	if total == 0 {
		return
	}

	items = impl.getPageItems(items, option)
	if len(items) == 0 {
		return
	}

	if option.Follower == nil {
		info.Users, err = impl.listFollowsDirectly(items, false)
	} else {
		info.Users, err = impl.listFollows(option.Follower.Account(), items)
	}

	info.Total = total

	return
}

func (impl *userRepoImpl) getPageItems(items []string, option *repository.FollowFindOption) []string {
	if option.CountPerPage <= 0 {
		return items
	}

	total := len(items)

	if option.PageNum <= 1 {
		if total > option.CountPerPage {
			return items[:option.CountPerPage]
		}

		return items
	}

	skip := option.CountPerPage * (option.PageNum - 1)
	if skip >= total {
		return nil
	}

	if n := total - skip; n > option.CountPerPage {
		return items[skip : skip+option.CountPerPage]
	}

	return items[skip:]
}

func (impl *userRepoImpl) getFollows(user, field string) (v DUser, err error) {
	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, mongo.UserDocFilterByAccount(user),
			bson.M{field: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}
	}

	return
}

func (impl *userRepoImpl) listFollowsDirectly(accounts []string, isFollower bool) (
	[]domain.FollowerUserInfo, error,
) {
	var v []DUser

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldName: bson.M{
				"$in": accounts,
			},
		}

		opts := options.FindOptions{}
		return impl.cli.GetDocs(
			ctx, filter, opts.SetProjection(bson.M{
				fieldBio:      1,
				fieldName:     1,
				fieldAvatarId: 1,
			}), &v,
		)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]domain.FollowerUserInfo, len(v))
	for i := range v {
		item := &v[i]
		if err := toFollowerUserInfo(*item, &r[i]); err != nil {
			return nil, err
		}
		r[i].IsFollower = isFollower
	}

	return r, nil
}

func (impl *userRepoImpl) listFollows(follower string, accounts []string) (
	[]domain.FollowerUserInfo, error,
) {
	var v []struct {
		DUser `bson:",inline"`

		IsFollower bool `bson:"is_follower"`
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"$expr": bson.M{"$in": bson.A{"$" + fieldName, accounts}},
		}},
		bson.M{"$project": bson.M{
			fieldIsFollower: bson.M{
				"$in": bson.A{follower, "$" + fieldFollower},
			},
			fieldBio:      1,
			fieldName:     1,
			fieldAvatarId: 1,
		}},
	}

	err := withContext(func(ctx context.Context) error {
		col := impl.cli.Collection()
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})

	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]domain.FollowerUserInfo, len(v))
	for i := range v {
		item := &v[i]
		if err = toFollowerUserInfo(item.DUser, &r[i]); err != nil {
			return nil, err
		}
		r[i].IsFollower = item.IsFollower
	}

	return r, nil
}
