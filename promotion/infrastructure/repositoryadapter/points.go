package repositoryadapter

import (
	"context"
	"errors"

	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func PointsAdapter(cli mongodbClient) repository.Points {
	return &pointsAdapter{cli}
}

type pointsAdapter struct {
	cli mongodbClient
}

func (impl *pointsAdapter) docFilter(username string) bson.M {
	return bson.M{fieldUser: username}
}

func (impl *pointsAdapter) docOfTotal(points int) bson.M {
	return bson.M{fieldTotal: points}
}

func (impl *pointsAdapter) docOfPromotionId(promotionid string) bson.M {
	return bson.M{fieldPromotionId: promotionid}
}

func (impl *pointsAdapter) Find(user types.Account) (domain.UserPoints, error) {
	var do pointsDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(ctx, impl.docFilter(user.Account()), nil, &do)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.UserPoints{}, err
	}

	return do.toUserPoints()
}

func (impl *pointsAdapter) FindAll(promotionid string) (ups []domain.UserPoints, err error) {
	var dos []pointsDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(ctx, impl.docOfPromotionId(promotionid), nil, &dos)
	}

	if err := withContext(f); err != nil || len(dos) == 0 {
		return nil, err
	}

	ups = make([]domain.UserPoints, len(dos))
	for i := range dos {
		if ups[i], err = dos[i].toUserPoints(); err != nil {
			return
		}
	}

	return
}

func (impl *pointsAdapter) Update(user types.Account, item domain.Item, version int) (err error) {
	// version = 0, it means userpoints not created
	if version <= 0 {
		return impl.save(user, item)
	}

	return impl.update(user, item, version)
}

func (impl *pointsAdapter) update(user types.Account, item domain.Item, version int) error {
	if version <= 0 {
		return errors.New("cannot update object version = 0")
	}

	itemDO := toItemDO(&item)
	doc, err := itemDO.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.PushArrayElemAndInc(
			ctx, fieldItems, impl.docFilter(user.Account()),
			doc, impl.docOfTotal(item.Points), version,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}

func (impl *pointsAdapter) save(user types.Account, item domain.Item) error {
	ups := &domain.UserPoints{
		User:    user,
		Total:   item.Points,
		Items:   []domain.Item{item},
		Version: 1,
	}
	do := toPointsDO(ups)

	doc, err := do.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, bson.M{
				fieldUser: user.Account(),
			}, doc,
		)

		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}
