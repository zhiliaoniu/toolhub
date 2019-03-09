package service

import (
	"bytes"
	"common"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"hash/fnv"
	"io/ioutil"
	"math/rand"
	"net/http"
	"service/api"
	"strconv"
	"strings"
	"sync"
	"time"
	pb "twirprpc"
)

const (
	LOGIN_EXPIRE = 8
	SUCCESS_CODE = "200"
	FAIL_CODE    = "400"
)

/**
 * 记录运营记录
 *
 * 1.用户名　　2.路由　　3.动作　　4.状态码　　5.操作对象
 */
func RerecordOpLog(code string, ctx context.Context) error {
	route := ctx.Value("ROUTE") //路由
	if route != nil {
		route = route.(string)
	}
	obj := ctx.Value("PARAMS") //请求参数
	action, _ := twirp.MethodName(ctx)
	ip := ctx.Value("X-Remote-Addr")
	userId := ctx.Value("USERID")

	insertSql := fmt.Sprintf(`insert into record_op_log(userId, route, ip, action, code, obj) 
values ('%v', '%v', '%v', '%v', '%v', '%v')`, userId, route, ip, action, code, obj)

	_, err := db.Exec(BUDAODB, insertSql)

	return err
}

/**
 * 记录运营记录
 *
 * 1.用户名　　2.路由　　3.动作　　4.状态码　　5.操作对象
 */
func RecordOpLog(code, action string, obj interface{}, ctx context.Context) error {
	var b []byte
	if temp, ok := obj.(proto.Message); ok {
		var buf bytes.Buffer
		marshaler := &jsonpb.Marshaler{OrigName: false}
		_ = marshaler.Marshal(&buf, temp)
		b = buf.Bytes()
	} else {
		b, _ = json.Marshal(obj)
	}

	var packageName, serviceName, methodName string
	if packName, ok := twirp.PackageName(ctx); ok {
		packageName = packName
	}
	if sName, ok := twirp.ServiceName(ctx); ok {
		serviceName = sName
	}
	if mName, ok := twirp.MethodName(ctx); ok {
		methodName = mName
	}

	ip := ctx.Value("X-Remote-Addr")
	userId := ctx.Value("USERID")
	route := "/" + packageName + "." + serviceName + "/" + methodName

	insertSql := fmt.Sprintf(`insert into record_op_log(userId, route, ip, action, code, obj) 
values ('%v', '%v', '%v', '%v', '%v', '%v')`, userId, route, ip, action, code, string(b))
	_, err := db.Exec(BUDAODB, insertSql)

	return err
}

func GetSqlParam(param api.QueryListRequest, filterMap map[string]string) (filterSql, sortSql, pageSql string, err error) {
	//分页
	pageSql = turnSql(param.Size, param.Num)
	//排序
	if param.Sort != "" {
		ss := strings.Split(param.Sort, "-")
		if col, ok := filterMap[ss[0]]; ok && ss[1] != "" {
			sortSql += fmt.Sprintf(`,%v %v`, strings.Split(col, `,`)[0], ss[1])
		}

		if sortSql != "" {
			sortSql = `ORDER by ` + sortSql[1:]
		}
	}

	//过滤
	for k, v := range param.Filter {
		if col, ok := filterMap[k]; ok && v != "" {
			col := strings.Split(col, `,`)
			filterSql += fmt.Sprintf(` and %v `+col[1], col[0], v)
		}
	}
	return
}

func turnSql(Size, Num int32) string {
	if Size < 0 || Num < 0 {
		return ""
	}
	if Num == 0 {
		Num = 0
	} else {
		Num -= 1
	}
	if Size == 0 {
		Size = 10
	}
	return fmt.Sprintf(` LIMIT  %v,%v`, Num*Size, Size)
}

func GetUpdateSql(buff []byte, templateMap map[string]string) (param map[string]string, updateSql string, err error) {
	err = json.Unmarshal(buff, &param)
	if err != nil {
		return
	}
	for k, v := range param {
		if template, ok := templateMap[k]; ok && v != "" {
			ss := strings.Split(template, ",")
			updateSql += fmt.Sprintf(`,%v `+ss[1], ss[0], v)
		}
	}
	if updateSql != "" {
		updateSql = updateSql[1:]
	} else {
		err = fmt.Errorf("input param is err")
		return
	}
	return
}

var (
	rpcReq  *pb.PostVideosRequest
	rpcResp *pb.PostVideosResponse
	rpcErr  error
)

//发布视频
func PostVideo(v api.Video) *pb.PostVideosResponse {
	if v.VideoHeight == 0 || v.VideoWidth == 0 || v.VideoDuration == 0 {
		glog.Info(v)
		glog.Debug("video height or weight or duration 0")
		//更新video_data发布失败原因
		//updateVideoData(v.SourceVid, "video height or weight or duration 0", 3)
		return nil
	}
	rpcReq = &pb.PostVideosRequest{}
	tVideos := make([]*pb.PostVideo, 0)
	postVideo := &pb.PostVideo{}

	if len(v.VSourceId) <= 0 {
		glog.Info("v_source_id is empty")
		//更新video_data发布失败原因
		//updateVideoData(v.SourceVid, "v_source_id is empty = "+ v.VSourceID, 3)
		return nil
	}

	postVideo.VsourceVid = v.VSourceId
	postVideo.SourceVid = v.SourceVid
	postVideo.Pic = v.VideoCover
	//postVideo.VideoUrl = v.VideoUrl
	postVideo.Title = v.Title
	postVideo.Width = v.VideoWidth
	postVideo.Height = v.VideoHeight
	postVideo.Duration = uint32(v.VideoDuration)
	postVideo.EVideoSource = pb.EVideoSource(v.Source)
	postVideo.EVideoParseRule = pb.VideoParseRule(v.ParseType)
	//视频的地址:playUrl 视频的真实地址    videoUrl 视频的页面地址
	if len(v.PlayUrl) > 0 {
		postVideo.VideoUrl = v.PlayUrl
	} else {
		postVideo.VideoUrl = v.VideoUrl
	}

	tVideos = append(tVideos, postVideo)
	rpcReq.TVideos = tVideos

	client := pb.NewTransferProtobufClient(common.GetConfig().Extension["transferUrl"], &http.Client{})

	rpcResp, rpcErr = client.PostVideos(context.Background(), rpcReq)
	if rpcErr != nil {
		//更新video_data发布失败原因
		//updateVideoData(v.SourceVid, fmt.Sprintf("rpc return err: %v", rpcErr), 3)
		if twerr, ok := rpcErr.(twirp.Error); ok {
			if twerr.Meta("retryable") != "" {
				// Log the error and go again.
				glog.Debug("got error %q, retrying", twerr)
			}
		}
		return nil
	}

	glog.Info(rpcResp)
	return rpcResp
}

const SPIDERDB = "spider"
const BUDAODB = "budao"

//发布视频下的评论
func PostVideoComment(vSourceVid string, vid string) {
	//1.获取视频下的所有评论
	comSql := fmt.Sprintf(`select cid, ifnull(v_source_id,""), content, favor_num, user_id, user_name, user_photo, reply_num, create_time from comment_data where v_source_id = '%v' and status=0`, vSourceVid)
	crows, err := db.Query(SPIDERDB, comSql)
	glog.Error(err)
	defer crows.Close()
	for crows.Next() {
		com := Comment{}
		crows.Scan(&com.Cid, &com.VSourceID, &com.Content, &com.FavorNum, &com.Uid, &com.UName, &com.UPhoto, &com.ReplyNum, &com.CTime)
		cid, autoIncreId, tableName, err := common.GetItemId(BUDAODB, COMMENT_)
		updateSql := fmt.Sprintf(`update %v set content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v",create_time='%v',favor_num=%v,reply_num=%v where id=%v`,
			tableName, com.Content, cid, vid, com.Uid, com.UName, com.UPhoto, com.CTime, com.FavorNum, com.ReplyNum, autoIncreId)
		_, err = db.Exec(BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			return
		}
		//更新comment_data的状态
		db.Exec(SPIDERDB, fmt.Sprintf(`update comment_data set status=1,post_cid=%v where v_source_id='%v' and cid='%v'`, cid, com.VSourceID, com.Cid))
	}
}

//将视频存入话题下
func PostVideoTopic(topicId string, vid string, ruleId string) {
	tID, _ := strconv.ParseUint(topicId, 10, 64)
	tableName, err := db.GetTableName(TOPIC_VIDEO_, tID)
	glog.Error(err)

	//获取topic
	querySql := fmt.Sprintf(`select topic_id from video_0 where vid =%v`, vid)
	row, _ := db.QueryRow(BUDAODB, querySql)
	var topics string
	err = row.Scan(&topics)
	glog.Error(err)

	//存入topic_video_X表
	insertSql := fmt.Sprintf(`insert into %v (topic_id, vid, rule_id) values (%v,%v,%v)`, tableName, topicId, vid, ruleId)
	_, err = db.Exec(BUDAODB, insertSql)
	glog.Error(err)

	//更新话题下的视频数
	updateSql := fmt.Sprintf(`update topic set video_num=video_num+1 where topic_id=%v`, topicId)
	db.Exec(BUDAODB, updateSql)

	//更新topic
	topics = strings.Trim(topics+`,`+topicId, ",")
	updateSql = fmt.Sprintf(`update video_0 set topic_id='%v' where vid=%v`, topics, vid)
	db.Exec(BUDAODB, updateSql)
}

//视频评论缓存
var videoCommentCache *commentCache = func() *commentCache {

	cache := &commentCache{
		data: make(map[string]*commentCacheDate),
	}
	go func() {
		for {
			now := time.Now()
			for vid, v := range cache.data {
				if now.Sub(v.inTime) > 5*time.Minute {
					delete(cache.data, vid)
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	return cache
}()

type commentCache struct {
	sync.Mutex
	data map[string]*commentCacheDate
}

func (c commentCache) get(vid string) ([]*Comment, bool) {
	c.Lock()
	defer c.Unlock()
	data, ok := c.data[vid]
	if !ok {
		return nil, ok
	}
	now := time.Now()
	if now.Sub(data.inTime) > 5*time.Minute {
		delete(c.data, vid)
		return nil, false
	}
	return data.data, ok
}
func (c commentCache) del(vid string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, vid)
}

type commentCacheDate struct {
	inTime time.Time
	data   []*Comment
}

func (c commentCache) put(vid string, data []*Comment) {
	c.Lock()
	defer c.Unlock()
	cacheDate := commentCacheDate{
		inTime: time.Now(),
		data:   data,
	}
	c.data[vid] = &cacheDate
}

//获取资源图片内容
func getHttpPic(picUrl string) (string, []byte) {
	response, err := http.Get(picUrl)
	defer response.Body.Close()
	if err != nil {
		return fmt.Sprintf("respose http Get error %s", err), nil
	}
	buff, err := ioutil.ReadAll(response.Body)

	return "", buff
}

//内部用户缓存
var IuCache *internalUserCache = &internalUserCache{
	inTime: time.Now().Add(-6 * time.Minute),
	data:   make([]user, 0),
}

type internalUserCache struct {
	sync.RWMutex
	inTime time.Time
	data   []user
}
type user struct {
	Uid   string
	Photo string
	Name  string
}

func (cache *internalUserCache) Get() []user {
	cache.RLock()
	if time.Now().Sub(cache.inTime) > 5*time.Minute {
		cache.RUnlock()
		cache.Lock()
		defer cache.Unlock()
		if time.Now().Sub(cache.inTime) > 5*time.Minute {
			cache.refresh()
		}
	} else {
		defer cache.RUnlock()
	}
	return cache.data

}
func (cache *internalUserCache) refresh() {
	count := common.GetConfig().DB.MySQL[BUDAODB].TableDesc["user_"]
	data := make([]user, 0, 100)
	for i := uint64(0); i < count; i++ {
		querySql := fmt.Sprintf(`select uid,name,photo from user_%v where login_method=10`, i)
		rows, err := db.Query(BUDAODB, querySql)
		if err != nil {
			glog.Error(err)
			return
		}
		for rows.Next() {
			u := user{}
			rows.Scan(&u.Uid, &u.Name, &u.Photo)
			data = append(data, u)
		}
	}
	cache.inTime = time.Now()
	cache.data = data
}

func GetRandUser(leng int) []user {
	return getRandUser(leng)
}
func getRandUser(leng int) []user {
	s := rand.NewSource(time.Now().Unix())
	ra := rand.New(s)
	allUser := IuCache.Get()
	max := len(allUser) - 1
	result := make([]user, 0, leng)
	set := make(map[int]bool)
	for len(result) < leng {
		i := ra.Intn(max)
		if _, ok := set[i]; !ok {
			set[i] = true
			result = append(result, allUser[i])
		}
	}
	return result
}

func GetRandUserByHash(soucre string) user {
	num := hash(soucre)
	allUser := IuCache.Get()
	mod := len(allUser)
	index := num % uint32(mod)
	return allUser[index]
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
