package util

import (
	"common"
	"crypto/md5"
	"db"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"github.com/sumaig/glog"
)

const (
	SALT = "J3%#^4X4Z^%(@6#,)*?&"
)

//check userid, commentid, videoid, topicid
func CheckIDValid(id uint64, token string) (valid bool, err error) {
	valid = false
	//0.check uid first
	if id == 0 {
		err = errors.New("uid is 0")
		return
	}
	//1.check 0-9, must < 1000
	tableNum := id >> 54
	curTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_"]
	if tableNum >= curTableNum {
		err = errors.New("user table num not correct")
		return
	}

	//2.check 48-63
	crcNum := (id << 48) >> 48
	beforeStr := strconv.FormatUint((id >> 16), 10)
	if crcNum != uint64(crc32.ChecksumIEEE([]byte(beforeStr))>>16) {
		err = errors.New("user crc not correct")
		return
	}

	//3.check token
	uidStr := strconv.FormatUint(id, 10)
	if GenerateToken(uidStr) != token {
		err = errors.New("user token not correct")
		return
	}
	valid = true

	return
}

func GenerateToken(userId string) (token string) {
	hash := md5.New()
	hash.Write([]byte(userId))
	hash.Write([]byte(SALT))
	token = hex.EncodeToString(hash.Sum(nil))
	return
}

func CheckTokenValid(userId, token string) bool {
	if GenerateToken(userId) == token {
		return true
	}
	return false
}

func CheckTokenExpire(userId, token string) (expire bool, err error) {
	expire = true
	//1.check token is expired
	key := fmt.Sprintf("user_act_%s", userId)
	expireTime, err := db.HGetInt(key, common.TOKEN_PREFIX+token)
	if err != nil {
		if strings.Contains(err.Error(), common.REDIS_RET_NIL) {
			glog.Debug("user:%s have no token now", userId)
			err = nil
		} else {
			glog.Error("check user:%s token failed", userId)
			return
		}
	}
	now := time.Now().Unix()
	if expireTime != 0 && int64(expireTime) > now+common.USER_TOKEN_EXPIRE_SEC {
		expire = false
		return
	}
	return
}

func ExtendTokenExpire(userId, token string) (err error) {
	key := fmt.Sprintf("user_act_%s", userId)
	now := time.Now().Unix()
	if _, err = db.HSet(key, common.TOKEN_PREFIX+token, now+common.USER_TOKEN_EXPIRE_SEC); err != nil {
		glog.Error("extend user:%s token:%s time failed. err:%v", userId, token, err)
		return
	}
	return
}
