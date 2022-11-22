package cache

import (
	"fmt"
	"github.com/flower-corp/rosedb"
	"path/filepath"
)

var DB *rosedb.RoseDB

func InitCache() {
	path := filepath.Join("./tmp", "rosedb")
	opts := rosedb.DefaultOptions(path)
	var err error
	DB, err = rosedb.Open(opts)
	if err != nil {
		panic(fmt.Sprintf("open rosedb err: %v", err))
	}
}
