/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "testing"

// TestMSDDescEqual tests the equality of two MSDDesc instances.
func TestMSDDescEqual(t *testing.T) {
	var desc1 MSDDesc
	var desc2 MSDDesc

	// case1
	if desc1 != desc2 {
		t.Fatalf("empty interfaces should be equal")
	}

	// case2
	desc2 = CreateMSDDesc("desc")
	if desc1 == desc2 {
		t.Fatalf("it should be unequal when desc1=%v, desc2=%v", desc1, desc2)
	}

	// case3
	desc1 = CreateMSDDesc("desc")
	if desc1 != desc2 {
		t.Fatalf("it should be equal when desc1=%v, desc2=%v", desc1, desc2)
	}

	// case4
	desc1 = CreateMSDDesc("abc")
	if desc1 == desc2 {
		t.Fatalf("it should be unequal when desc1=%v, desc2=%v", desc1, desc2)
	}
}
