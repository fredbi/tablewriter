package tablewrappers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordWrapper(t *testing.T) {
	t.Parallel()

	matrix := buildLengthsMatrix([]string{"characters", "too", "long"}, 1)
	require.EqualValues(t,
		[][]int{{10, 14, 19}, {0, 3, 8}, {0, 0, 4}},
		matrix,
	)
}
