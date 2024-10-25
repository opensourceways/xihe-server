package domain

import "github.com/opensourceways/xihe-server/domain"

type SpaceIndex struct {
	Name  domain.ResourceName
	Owner domain.Account
	Id    domain.Identity
}
