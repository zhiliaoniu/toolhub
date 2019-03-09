package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	pb "twirprpc"
)

func GetLocalUnix(t time.Time) uint64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"), loc)
	return uint64(theTime.Unix())
}

func TransStrArrToInterface(strArr []string) []interface{} {
	l := len(strArr)
	retArr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		retArr[i] = strArr[i]
	}
	return retArr
}

/**
 * 获取请求响应的初始状态
 */
func GetInitStatus() (status *pb.Status) {
	return &pb.Status{
		Code:       pb.Status_SERVER_ERR,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
}

/**
 * 获取md5字符串
 */
func GetMd5Str(str []byte) string {
	s := string(str) + "jxz_budao"
	hmd := md5.New()
	//hmd.Write([]byte(str))
	hmd.Write([]byte(s))
	md5Str := hmd.Sum(nil)

	return hex.EncodeToString(md5Str)
}

/**
 * sha1加密串
 */
func GetSha1Str(secret string, str string) []byte {

	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))

	return mac.Sum(nil)
}

/**
 * sha256加密串
 */
func GetSha256Str(secret string, str string) []byte {

	key := []byte(secret)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(str))

	return mac.Sum(nil)
}

/**
* 获取格林时间
 * Mon, 02 Jan 2006 07:04:05
  * @return string
*/
func Gmtime() string {
	local, _ := time.LoadLocation("PRC")
	timeFormat := "Mon, 02 Jan 2006 07:04:05 GMT"

	return time.Now().In(local).Format(timeFormat)
}

func IsVersionBigger(lhs, rhs string) (result bool) {
	larr := strings.Split(lhs, ".")
	rarr := strings.Split(rhs, ".")
	if len(larr) < 3 || len(rarr) < 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		lnum, _ := strconv.Atoi(larr[i])
		rnum, _ := strconv.Atoi(rarr[i])
		if lnum > rnum {
			return true
		} else if lnum < rnum {
			return false
		}
	}
	return true
}
