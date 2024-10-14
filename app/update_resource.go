package app

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
)

type ResourceTagsUpdateCmd struct {
	ToAdd    []string
	ToRemove []string
	All      []domain.DomainTags
}

func (cmd *ResourceTagsUpdateCmd) GenTagKinds(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	r := make([]string, 0, len(cmd.All))

	for i := range cmd.All {
		if v := cmd.All[i].GetKindsOfTags(tags); len(v) > 0 {
			r = append(r, v...)
		}
	}

	return r
}

func (cmd *ResourceTagsUpdateCmd) ToTags(old []string) ([]string, bool) {
	tags := sets.NewString(old...)

	if len(cmd.ToAdd) > 0 {
		tags.Insert(cmd.ToAdd...)
	}

	if len(cmd.ToRemove) > 0 {
		tags.Delete(cmd.ToRemove...)
	}

	if tags.Equal(sets.NewString(old...)) {
		return nil, false
	}

	return tags.UnsortedList(), true
}
