//go:build go1.18 && !go1.21
// +build go1.18,!go1.21

package mathext

import (
	"math"

	constraintsext "github.com/khulnasoft-lab/utils/constraints"
)

// Max returns the larger value.
//
// NOTE: this function does not check for difference in floats of 0/zero vs -0/negative zero using Signbit.
func Max[N constraintsext.Number](x, y N) N {
	// special case for floats
	// IEEE 754 says that only NaNs satisfy f != f.
	if x != x || y != y {
		return N(math.NaN())
	}

	if x > y {
		return x
	}
	return y
}
