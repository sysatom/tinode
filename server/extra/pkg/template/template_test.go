package template

import (
	"github.com/stretchr/testify/require"
	"github.com/tinode/chat/server/extra/types"
	"testing"
)

func TestParse(t *testing.T) {
	txt, err := Parse(types.Context{}, "Welcome $1 $2", "user", "user2")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "Welcome user user2", txt)
}
