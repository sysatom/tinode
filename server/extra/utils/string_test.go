package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHasHan(t *testing.T) {
	require.True(t, HasHan("hello, 世界"))
	require.False(t, HasHan("hello"))
}

func TestMasker(t *testing.T) {
	require.Equal(t, "qwerty**********7890", Masker("qwertyuiop1234567890", 2))
}
