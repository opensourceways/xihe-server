package domain

import (
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/domain"
)

type SpaceIndex struct {
	Name  domain.ResourceName
	Owner domain.Account
	Id    primitive.Identity
}
