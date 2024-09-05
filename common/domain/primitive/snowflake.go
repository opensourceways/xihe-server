/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

// GetId generates a unique identifier as an int64 value.
func GetId() int64 {
	return node.Generate().Int64()
}
