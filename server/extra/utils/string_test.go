package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHasHan(t *testing.T) {
	require.True(t, HasHan("hello, 世界"))
	require.False(t, HasHan("hello"))
}
