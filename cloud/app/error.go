package app

import "errors"

const (
	errorResourceBusy        = "cloud_resource_busy"
	errorNotAllowed          = "cloud_not_allowed"
	errorWhitelistNotAllowed = "not_allowed"
)

var (
	ErrCloudReleased   = errors.New("cloud was released")
	ErrCloudNotAllowed = errors.New("not allowed")
	ErrPodNotFound     = errors.New("not found")
)
