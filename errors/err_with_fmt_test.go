package errorutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	errors "github.com/khulnasoft-lab/utils/errors"
)

func TestErrWithFmt(t *testing.T) {
	errBase := errors.NewWithFmt("error: %s")
	errWithMsg1 := errBase.Msgf("test1")
	errWithMsg2 := errBase.Msgf("test2")

	require.Equal(t, "error: test1", errWithMsg1.Error())
	require.Equal(t, "error: test2", errWithMsg2.Error())
}
