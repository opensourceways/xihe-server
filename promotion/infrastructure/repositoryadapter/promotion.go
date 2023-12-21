package repositoryadapter

import (
	"context"

	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func PromotionAdapter(cli mongodbClient) repository.Promotion {
	return &promotionAdapter{cli}
}

type promotionAdapter struct {
	cli mongodbClient
}

func (impl *promotionAdapter) Find(promotionid string) (repository.PromotionRepo, error) {
	var do promotionDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, docIdFilter(promotionid),
			nil, &do,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return repository.PromotionRepo{}, err
	}

	return do.toPromotionRepo()
}

func (impl *promotionAdapter) FindAll() (prs []repository.PromotionRepo, err error) {
	var dos []promotionDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, nil,
			nil, &dos,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	prs = make([]repository.PromotionRepo, len(dos))
	for i := range dos {
		if prs[i], err = dos[i].toPromotionRepo(); err != nil {
			return
		}
	}

	return
}

func (impl *promotionAdapter) UserRegister(promotionid string, user types.Account, version int) error {
	regUserDO := RegUserDO{
		User:      user.Account(),
		CreatedAt: utils.Now(),
	}
	doc, err := regUserDO.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.PushElemArrayWithVersion(
			ctx, fieldRegUsers,
			docIdFilter(promotionid),
			doc, version, nil,
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
