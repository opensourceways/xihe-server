package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type LikeMessageProducer interface {
	AddLike(*domain.ResourceObject) error
	RemoveLike(*domain.ResourceObject) error
}
