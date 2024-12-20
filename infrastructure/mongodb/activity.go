package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewActivityMapper(name string, keep int) repositories.ActivityMapper {
	return activity{
		collectionName: name,
		keepNum:        -keep,
	}
}

type activity struct {
	collectionName string
	keepNum        int
}

func (col activity) Insert(owner string, do repositories.ActivityDO) (err error) {
	if err = col.insert(owner, do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = newResourceDoc(col.collectionName, owner); err == nil {

		if err = col.insert(owner, do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

// Two same kind of activities can happen at different time,
// for example, like, unlike, and like again.
func (col activity) insert(owner string, do repositories.ActivityDO) error {
	doc, err := col.toActivityDoc(&do)
	if err != nil {
		return err
	}

	docFilter := resourceOwnerFilter(owner)
	appendElemMatchToFilter(fieldItems, false, doc, docFilter)

	f := func(ctx context.Context) error {
		return cli.pushElemToLimitedArray(
			ctx, col.collectionName, fieldItems, col.keepNum,
			docFilter, doc,
		)
	}

	return withContext(f)
}

func (col activity) List(owner string, opt repositories.ActivityListDO) (
	r []repositories.ActivityDO, err error,
) {
	var v dActivity

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, resourceOwnerFilter(owner), nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	items := v.Items
	r = make([]repositories.ActivityDO, len(items))
	for i := range items {
		col.toActivityDO(&items[i], &r[i])
	}
	fmt.Printf("=======================r: %+v\n", r)
	return
}

func (col activity) toActivityDoc(do *repositories.ActivityDO) (bson.M, error) {
	v := activityItem{
		Type:           do.Type,
		Time:           do.Time,
		RepoType:       do.RepoType,
		ResourceObject: toResourceObject(&do.ResourceObjectDO),
	}

	return genDoc(v)
}

func (col activity) toActivityDO(item *activityItem, do *repositories.ActivityDO) {
	*do = repositories.ActivityDO{
		Type:             item.Type,
		Time:             item.Time,
		RepoType:         item.RepoType,
		ResourceObjectDO: toResourceObjectDO(&item.ResourceObject),
	}
}
