package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCovertJSON(t *testing.T) {
	j, err := ConvertJSON(map[string]interface{}{
		"test": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, ok := j.Int64("test")
	if !ok {
		t.Fail()
	}
	require.Equal(t, int64(1), n)
}
