/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

// CommonListRequest is a struct that holds common parameters for list requests.
type CommonListRequest struct {
	SortBy       string `form:"sort_by"`
	Count        bool   `form:"count"`
	PageNum      int    `form:"page_num"`
	CountPerPage int    `form:"count_per_page"`
}
