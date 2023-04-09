package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaxI64(t *testing.T) {
	assert.Equal(t, int64(2), MaxI64(1, 2))
	assert.Equal(t, int64(2), MaxI64(2, 1))
}

func Test_MaxUI32(t *testing.T) {
	assert.Equal(t, uint32(2), MaxUI32(1, 2))
	assert.Equal(t, uint32(2), MaxUI32(2, 1))
}

func Test_MaxF64(t *testing.T) {
	assert.Equal(t, float64(2), MaxF64(1, 2))
	assert.Equal(t, float64(2), MaxF64(2, 1))
}

func Test_RandIntInRange(t *testing.T) {
	assert.GreaterOrEqual(t, uint32(3), RandIntInRange(1, 3))
	assert.LessOrEqual(t, uint32(1), RandIntInRange(2, 3))
}

func Test_Init(t *testing.T) {
	assert.NotPanics(t, Init)
}

func TestRoundFloat(t *testing.T) {
	testCases := []struct {
		input     float64
		precision uint
		expected  float64
	}{
		{12.3456789, 2, 12.35},
		{12.3456789, 3, 12.346},
		{12.3456789, 4, 12.3457},
		{12.3456789, 5, 12.34568},
	}

	for _, tc := range testCases {
		output := roundFloat(tc.input, tc.precision)
		if output != tc.expected {
			t.Errorf(
				"roundFloat(%v, %v) = %v; want %v",
				tc.input,
				tc.precision,
				output,
				tc.expected,
			)
		}
	}
}
