package handler

import (
	"testing"
)

func TestFirstMissing(t *testing.T) {
	type data struct {
		input  []int64
		output int64
	}

	d := []data{
		{
			input:  []int64{0, 2, 3},
			output: 1,
		},
		{
			input:  []int64{1, 3, 7},
			output: 2,
		},
		{
			input:  []int64{1, 2, 3},
			output: 4,
		},
	}

	for idx, entry := range d {
		if x := firstMissing(entry.input); x != entry.output {
			t.Error("x is not valid", "x", x, "idx", idx)
		}
	}
}
