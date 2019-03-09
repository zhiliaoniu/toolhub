package topicserver

import (
	"bufio"
	"common"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"db"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"service/util"
	"strconv"
	"strings"
	"time"
	pb "twirprpc"

	"github.com/sumaig/glog"
)

var GIOSAuditTopicIds []string

const (
	uploadKey    = "ak_yey"
	uploadSecret = "97cbac39d755ee2c3c63ac79944a2a477fd5faff"
	uploadHost   = "jxzimg.bs2ul.yy.com"
	bucketName   = "jxzimg"
)

// GetBannerItems get bannerItem for banner
func GetBannerItems() (bannerItems []*pb.BannerItem, err error) {
	sqlString := "select id, pic_url, link from banner where status != 1 limit 3"
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query banner failed. err:%v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      uint64
			picURL  string
			linkURL string
		)
		rows.Scan(&id, &picURL, &linkURL)
		body, err := LoadPicture(picURL)
		if err != nil {
			glog.Error("get picture from cdn failed. err:%v", err)
			continue
		}
		bannerItem := &pb.BannerItem{
			BannerId: strconv.FormatUint(id, 10),
			PicUrl:   picURL,
			PicMd5:   body,
			LinkUrl:  linkURL,
		}
		bannerItems = append(bannerItems, bannerItem)
	}

	return
}

// LoadPicture load picture from internet
func LoadPicture(pictureURL string) (body string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", pictureURL, nil)
	if err != nil {
		glog.Error("newRequest error %v", err)
		return
	}

	//加密过程
	fileName := strings.Split(pictureURL, "yy.com/")
	if len(fileName) != 2 {
		return
	}
	fname := fileName[1]
	expires := strconv.FormatInt(time.Now().Unix()+7200, 10)
	hmacStr := "GET\njxzimg\n" + bucketName + "\n" + fname + "\n" + expires + "\n"
	hm := GetSha1Str(uploadSecret, hmacStr)
	base := base64.URLEncoding.EncodeToString(hm)

	request.Header.Set("Host", uploadHost)
	request.Header.Set("Date", Gmtime())
	request.Header.Set("Authorization", uploadKey+":"+base+":"+expires)

	resp, err := client.Do(request)
	if err != nil {
		glog.Error("http request failed. err:%v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		glog.Error("respstatusCode not 200")
		return
	}
	glog.Debug("respstatusCode:%v\n", resp.StatusCode)
	respBody, err := ioutil.ReadAll(resp.Body)

	h := md5.New()
	h.Write([]byte(respBody))
	tempByte := h.Sum(nil)
	body = hex.EncodeToString(tempByte)

	return body, err
}

// GetSha1Str sha1加密串
func GetSha1Str(secret string, str string) []byte {
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))

	return mac.Sum(nil)
}

// Gmtime 获取格林时间 Mon, 02 Jan 2006 07:04:05
func Gmtime() string {
	local, _ := time.LoadLocation("PRC")
	timeFormat := "Mon, 02 Jan 2006 07:04:05 GMT"

	return time.Now().In(local).Format(timeFormat)
}

func InitIOSOnlineAuditTopicIds() {
	GIOSAuditTopicIds = make([]string, 0)
	var fileName string
	fileName = common.GetConfig().IOSAuditConf.TopicIdsFileName
	if fileName == "" {
		fileName = "./conf/ios_audit_topicids"
	}
	f, err := os.Open(fileName)
	if err != nil {
		glog.Error("open ios audit topicids file failed. filename:%s, err:%v", fileName, err)
		panic(err)
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				glog.Debug("init ios online audit topicid success. GIOSAuditTopicIds.len:%d", len(GIOSAuditTopicIds))
				//glog.Debug("GIOSAuditTopicIds:%v", GIOSAuditTopicIds)
				break
			}
			panic(err)
		}
		line = strings.TrimSpace(line)
		GIOSAuditTopicIds = append(GIOSAuditTopicIds, line)
	}
}

func InitTopicDynamicInfo() {
	//1.read mysql
	sqlString := fmt.Sprintf("select topic_id, user_num from topic")
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query topic table failed. sqlString:%s, err:%v", sqlString, err)
		return
	}
	defer rows.Close()

	fields := make([]interface{}, 0)
	var topicId, userNum string
	for rows.Next() {
		err = rows.Scan(&topicId, &userNum)
		if err != nil {
			glog.Error("scan topic dynamic info failed. err:%v", err)
			return
		}
		fields = append(fields, topicId, userNum)
	}

	if _, err := db.HMSet(common.TOPIC_DYNAMIC, fields); err != nil {
		glog.Error("hmset failed. err:%v", err)
		return
	}
	glog.Debug("update topic dynamic info success")
}

func GetUserSubscribedTopicId(userId string) (topicIds []string, err error) {
	key := fmt.Sprintf("%s%s", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, userId)
	topicIds, err = db.ZRevRange(key, 0, -1)
	return
}

func GetIdsByLastId(key, lastId string, returnNum int) (ids []string, err error) {
	//1.get lastCommentId rank in zset
	offset := 0
	if lastId != "" {
		offset, err = db.ZRevRank(key, lastId)
		if err != nil {
			return
		}
		offset += 1
	}
	ids, err = db.ZRevRange(key, int64(offset), int64(offset+returnNum))
	return
}

//InitUserSubscribedTopic update user subscribe topic list
func InitUserSubscribedTopic() {
	userTopicMap := make(map[uint64][]interface{}) //<uid <topicid, update_time>>
	//1.get user subscribe topic from mysql
	tableNUM := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_follow_topic_"]
	var i uint64
	for i = 0; i < tableNUM; i++ {
		tableName := fmt.Sprintf("user_follow_topic_%d", i)
		sqlString := fmt.Sprintf("select uid, topic_id, update_time from %s", tableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. sqlString:%s, err:%v", tableName, sqlString, err)
			continue
		}
		defer rows.Close()

		var (
			uid        uint64
			topicId    uint64
			updateTime time.Time
		)
		for rows.Next() {
			if err := rows.Scan(&uid, &topicId, &updateTime); err != nil {
				rows.Close()
				break
			}
			userTopicMap[uid] = append(userTopicMap[uid], updateTime.Unix(), topicId)
		}
		//close mysql conn active
		rows.Close()
	}

	//2.set user subscribe topic to redis by zset
	for uid, topicIds := range userTopicMap {
		key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, uid)
		_, err := db.ZAddMulti(key, topicIds)
		if err != nil {
			glog.Error("zadd multi failed. err:%v", err)
			continue
		}
	}
	glog.Debug("UpdateUserTopic success")
}

func GetTopicItemsWithExplicit(userId uint64, topicIds []string, explicitNum int) (topicIdMap map[string]*pb.TopicItem, err error) {
	topicIdMap = make(map[string]*pb.TopicItem, 0)
	//1.get vids
	vids := make([]string, 0)
	for _, topicId := range topicIds {
		key := fmt.Sprintf("%s%s", common.TOPICWITHVIDEOID, topicId)
		topicVids, err := GetIdsByLastId(key, "", explicitNum-1)
		if err != nil {
			if strings.Contains(err.Error(), "redigo: nil returned") {
				glog.Debug("topicId:%s has no vid. err:%v", topicId, err)
				continue
			}
			glog.Error("get vids of topicid:%s failed. err:%v", topicId, err)
			continue
		}
		if len(topicVids) == 0 {
			glog.Error("topic:%s has no video", topicId)
			continue
		}
		vids = append(vids, topicVids...)
	}
	if len(vids) == 0 {
		glog.Error("topicIds:%v has no video", topicIds)
		//err = errors.New("topic has no video")
		//TODO add alarm
		return
	}

	//2.getlistItem
	listItems, err := util.GetListItemsNotReal(vids, userId)
	if err != nil {
		glog.Error("get listItems failed. err:%v", err)
		return
	}
	if len(listItems) == 0 {
		glog.Error("user:%d get empty listItem.vids:%s from topicIds:%v", userId, vids, topicIds)
		err = errors.New("topic has no video")
		return
	}

	//3.compose topicItems
	for _, listItem := range listItems {
		if listItem.TopicItem == nil {
			glog.Debug("listItem:%s has no topicItem", listItem.VideoItem.VideoId)
			continue
		}
		topicId := listItem.TopicItem.TopicId
		topicItem, ok := topicIdMap[topicId]
		if !ok {
			topicItem = listItem.TopicItem
			topicIdMap[topicId] = topicItem
		}
		topicItem.VideoItems = append(topicItem.VideoItems, listItem.VideoItem)
	}

	return
}
