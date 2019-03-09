package common

import "sync"

// VERSION define program version.
const VERSION = "1.0"

var WG sync.WaitGroup

const (
	DUPLICATE_ENTRY = "Duplicate entry"
	REDIS_RET_NIL   = "redigo: nil returned"
)
