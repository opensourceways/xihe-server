package repositoryadapter

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	common "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
)

func UserPointsAdapter(cli mongodbClient, cfg *Config) *userPointsAdapter {
	return &userPointsAdapter{
		cli:  cli,
		keep: cfg.Keep,
	}
}

type userPointsAdapter struct {
	cli  mongodbClient
	keep int
}

func (impl *userPointsAdapter) docFilter(account common.Account) bson.M {
	return bson.M{fieldUser: account.Account()}
}

func (impl *userPointsAdapter) docOfDaysFilter(date string) bson.M {
	return bson.M{fieldDate: date}
}

func (impl *userPointsAdapter) SavePointsItem(up *domain.UserPoints, item *domain.PointsItem) error {
	if up.Version == 0 {
		// set version to 1 to avoid running again here.
		up.Version = 1

		return impl.addUserPoints(up, item)
	}

	if up.IsFirstPointsDetailOfDay() {
		return impl.addFirstPointsDetailOfDay(up, item)
	}

	return impl.addPointsDetail(up, item)
}

// insert new user points
func (impl *userPointsAdapter) addUserPoints(up *domain.UserPoints, item *domain.PointsItem) error {
	do := userPointsDO{
		User:    up.User.Account(),
		Total:   up.Total,
		Days:    []pointsDetailsOfDayDO{toPointsItemsOfDayDO(item)},
		Dones:   up.Dones,
		Version: up.Version,
	}

	doc, err := do.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(ctx, impl.docFilter(up.User), doc)

		return err
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}

// push the first points detail of new day
func (impl *userPointsAdapter) addFirstPointsDetailOfDay(up *domain.UserPoints, item *domain.PointsItem) error {
	do := toPointsItemsOfDayDO(item)

	doc, err := do.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.PushElemToLimitedArrayWithVersion(
			ctx, fieldDays, impl.keep,
			impl.docFilter(up.User),
			doc, up.Version,
			bson.M{
				fieldTotal: up.Total,
				fieldDones: up.Dones,
			},
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}

		return err
	}

	return nil
}

// add points detail of day
func (impl *userPointsAdapter) addPointsDetail(up *domain.UserPoints, item *domain.PointsItem) error {
	do := topointsDetailDO(item.TaskId, item.LatestDetail())

	doc, err := do.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.PushNestedArrayElemAndUpdate(
			ctx, fieldDays,
			impl.docFilter(up.User),
			impl.docOfDaysFilter(item.Date),
			bson.M{fieldDetails: doc}, up.Version,
			bson.M{
				fieldTotal: up.Total,
				fieldDones: up.Dones,
			},
		)

		return err
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}

		return err
	}

	return nil
}

func (impl *userPointsAdapter) Find(account common.Account, date string) (domain.UserPoints, error) {
	var dos []userPointsDO

	f := func(ctx context.Context) error {
		project := bson.M{
			fieldUser:    1,
			fieldTotal:   1,
			fieldDones:   1,
			fieldDays:    1,
			fieldVersion: 1,
		}

		return impl.cli.GetArrayElem(
			ctx, fieldDays, impl.docFilter(account),
			impl.docOfDaysFilter(date), project, &dos,
		)
	}

	if err := withContext(f); err != nil {
		return domain.UserPoints{}, err
	}

	if len(dos) == 0 {
		err := repoerr.NewErrorResourceNotExists(errors.New("no data"))

		return domain.UserPoints{}, err
	}

	return dos[0].toUserPoints()
}

func (impl *userPointsAdapter) FindAll(account common.Account) (domain.UserPoints, error) {
	var do userPointsDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, impl.docFilter(account), nil, &do,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.UserPoints{}, err
	}

	return do.toUserPoints()
}
