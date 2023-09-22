package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type LikeMessageProducer interface {
	AddLike(*domain.ResourceLikedEvent) error
	RemoveLike(*domain.ResourceLikedEvent) error
}
