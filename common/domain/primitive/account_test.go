/*
Copyright (c) Huawei Technologies Co., Ltd. 2023-2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"regexp"
	"testing"
)

// TestNewTokenName Test the rules when creating token name
func TestNewTokenName(t *testing.T) {
	tokenConfig.MaxNameLength = 50
	tokenConfig.MinNameLength = 1
	tokenConfig.regexp, _ = regexp.Compile("^[a-zA-Z0-9_-]*[a-zA-Z_-]+[a-zA-Z0-9_-]*$")

	if tokenName, _ := NewTokenName("0b1001"); tokenName != nil {
		t.Fatalf("it should return nil when name is 0b1001")
	}

	if tokenName, _ := NewTokenName("0o123"); tokenName != nil {
		t.Fatalf("it should return nil when name is 0o123")
	}

	if tokenName, _ := NewTokenName("0x4d9"); tokenName != nil {
		t.Fatalf("it should return nil when name is 0x4d9")
	}

	if tokenName, _ := NewTokenName("0b_test"); tokenName == nil {
		t.Fatalf("it should return a new token name when name is 0b_test")
	}
}
