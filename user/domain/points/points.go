package points

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type Points interface {
	Points(types.Account) (int, error)
}
