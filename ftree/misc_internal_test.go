package ftree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNumDigits test an interal numDigits function
func TestNumDigits(t *testing.T) {
	require.Equal(t, 1, numDigits(0))
	require.Equal(t, 1, numDigits(9))
	require.Equal(t, 2, numDigits(23))
	require.Equal(t, 3, numDigits(100))
}
