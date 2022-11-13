package locker

import "sync"

var Mux = sync.Mutex{}
