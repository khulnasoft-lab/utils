//go:build go1.21

package mathext

import (
	constraintsext "github.com/khulnasoft-lab/utils/constraints"
)

// Max returns the larger value.
//
// NOTE: this function does not check for difference in floats of 0/zero vs -0/negative zero using Signbit.
//
// Deprecated: use the new std library `max` instead.
func Max[N constraintsext.Number](x, y N) N {
	return max(x, y)
}
