package cache

import (
	"github.com/flower-corp/rosedb"
	"github.com/tinode/chat/server/logs"
	"path/filepath"
)

var DB *rosedb.RoseDB

func init() {
	path := filepath.Join("./tmp", "rosedb")
	opts := rosedb.DefaultOptions(path)
	var err error
	DB, err = rosedb.Open(opts)
	if err != nil {
		logs.Err.Fatalf("open rosedb err: %v", err)
	}
}
