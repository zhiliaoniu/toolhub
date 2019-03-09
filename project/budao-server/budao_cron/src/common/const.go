package common

import "sync"

// VERSION define program version.
const VERSION = "3.0.3V"

var WG sync.WaitGroup

//video state
const (
	VIDEOSTATE_WAIT_AUDIT     = 0
	VIDEOSTATE_AUDITING       = 1
	VIDEOSTATE_PASS_AUDIT     = 2
	VIDEOSTATE_NOT_PASS_AUDIT = 3
	VIDEOSTATE_DELETED        = 4
)

const (
	VIDEO_TABLE_PREFIX   = "video_"
	COMMNET_TABLE_PREFIX = "comment_"
	USER_TABLE_PREFIX    = "user_"
)

const (
	// FULLVIDEOHASHKEYREDIS hash
	FULLVIDEOHASHKEYREDIS = "full_video_hash_key_redis"

	// USEREXPOSURELISTPREFIX sort set
	USEREXPOSURELISTPREFIX = "exposure_list_deviceid_"
)

const (
	LIKE_COUNT       = "like_count"
	COMMENT_COUNT    = "comment_count"
	VIEW_COUNT       = "view_count"
	SHARE_COUNT      = "share_count"
	COMMENT_DISABLED = "comment_disabled"
	LIKE_DISABLED    = "like_disabled"
	SHARE_DISABLED   = "share_disabled"
)

const (
	// TOPICSORTSET sort set
	TOPICSORTSET = "topic_sort_set"

	// TOPICWITHUSERID set
	TOPICWITHUSERID = "topic_userid_"

	// TOPICWITHVIDEOID sort set
	TOPICWITHVIDEOID = "topic_video_"

	// TOPICHASH hash
	TOPICHASH = "topic_hash"

	//zset
	USER_SUBSCRIBE_TOPIC_WITH_TIME = "ustwt_"

	//hash
	USER_VIEW_TOPIC = "user_view_topic_"
)

//dynamic
const (
	VIDEO_DYNAMIC     = "video_dynamic"
	VIDEO_HOT_COMMENT = "video_hot_comment" //hash field:vid value pb.ArtibleItem
	QUESTION_DYNAMIC  = "question_dynamic"
	TOPIC_DYNAMIC     = "topic_dynamic"
	//变化的listitem
	VID_DYNAMIC     = "vid_dynamic"
	VID_DYNAMIC_TMP = "vid_dynamic_tmp"

	//user act
	USER_ACT_PREFIX        = "user_act_"
	VID_FAVOR_SUFFIX       = "_vfavor"
	CID_FAVOR_SUFFIX       = "_cfavor"
	TOPIC_SUBSCRIBE_SUFFIX = "_subscribe"
	TOKEN_PREFIX           = "token_"
)

//TODO data use conf file
const (
	// SEPERATESEVENDAYSTIMESTAMP expire time
	SEPERATESEVENDAYSTIMESTAMP = 5 * 24 * 60 * 60

	// COMMENT_DYNAMIC_KEY_NUM dyncmic num
	COMMENT_DYNAMIC_KEY_NUM = 1000

	//用户token过期时间，一年
	USER_TOKEN_EXPIRE_SEC = 365 * 24 * 60 * 60

	READ_MYSQL_MAX_ROWS = 200
)

const (
	DUPLICATE_ENTRY = "Duplicate entry"
	REDIS_RET_NIL   = "redigo: nil returned"
)

const (
	BUDAODB = "budao"
)
