package app

import "github.com/opensourceways/xihe-server/infrastructure/authingimpl"

const (
	errorNoUserId      = "user_no_userid"
	errorNoAccessToken = "user_no_accesstoken"
)

func isCodeUserDuplicateBind(code string) bool {
	return code == authingimpl.ErrorUserDuplicateBind
}
