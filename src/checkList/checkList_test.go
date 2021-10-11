package checkList

import (
	"testing"
)

func TestCheckList(t *testing.T) {
	for _, one := range []struct {
		List []int
		Ok   bool
		Next int
	}{
		{[]int{3, 5, 7, 9, 11}, true, 13},
		{[]int{2, 4, 8, 16, 32}, true, 64},
		{[]int{2, 15, 41, 80}, true, 132},
		{[]int{1, 2, 6, 15, 31}, true, 56},
		{[]int{1, 1, 3, 15, 105, 945}, true, 10395},
		{[]int{2, 14, 64, 202, 502, 1062, 2004}, true, 3474},
		{[]int{1, 1, 2, 6}, true, 15},
		{[]int{1, 16, 81, 256}, true, 625},
		{[]int{32, 16, 8, 4, 2}, false, 0},
		{[]int{1, 3, 5, 13, 85}, true, 1237}, //求差-->求商-->开方等差
		{[]int{1, 4, 9}, true, 16},
		{[]int{2, 2, 2, 2, 3}, false, 0},
		{[]int{3, 2, 2, 2, 2}, false, 0},
		{[]int{0, 0, 0, 0, 0, 0}, true, 0},
		{[]int{0, -1, 0, -1, 0, -1, 0, -1}, true, 0},
	} {
		next, ok := checkList(one.List)
		if ok == one.Ok && next == one.Next {
			t.Log("correct", "input", one.List, "get", next, ok, "want", one.Next, one.Ok)
		} else {
			t.Error("wrong", "input", one.List, "get", next, ok, "want", one.Next, one.Ok)
		}
	}
}
