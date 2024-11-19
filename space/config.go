package space

import "github.com/opensourceways/xihe-server/space/infrastructure"

type Config struct {
	Topics infrastructure.Topics `json:"topics"`
}
