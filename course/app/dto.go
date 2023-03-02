package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/course/domain"
)

// player apply
type PlayerApplyCmd domain.Player

func (cmd *PlayerApplyCmd) Validate() error {
	b := cmd.Student.Account != nil &&
		cmd.Student.Name != nil &&
		cmd.Student.Email != nil &&
		cmd.Student.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *PlayerApplyCmd) toPlayer() (p domain.Player) {
	return *(*domain.Player)(cmd)
}
