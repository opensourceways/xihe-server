/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "fmt"

type AvatarInfo struct {
	User        Account
	FileName    string
	Path        string
	Bucket      string
	CdnEndpoint string
	TmpPath     string
}

func (c *AvatarInfo) GetAvatarURL() string {
	return fmt.Sprintf("%s%s/%s/%s", c.CdnEndpoint, c.Path, c.User.Account(), c.FileName)
}

func (c *AvatarInfo) GetObsPath() string {
	return fmt.Sprintf("%s/%s/%s", c.Path, c.User.Account(), c.FileName)
}
