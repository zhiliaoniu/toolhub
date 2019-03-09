package common

// VERSION define program version.
const VERSION = "1.0"

//video state
const (
	VIDEOSTATE_WAIT_AUDIT     = 0
	VIDEOSTATE_AUDITING       = 1
	VIDEOSTATE_PASS_AUDIT     = 2
	VIDEOSTATE_NOT_PASS_AUDIT = 3
	VIDEOSTATE_DELETED        = 4
)
