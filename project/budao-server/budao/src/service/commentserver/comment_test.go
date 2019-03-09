package commentserver

import (
	"common"
	"context"
	"db"
	"fmt"
	"testing"
	"time"
	pb "twirprpc"
)

//go test -run=文件名字 -bench=bench名字 -cpuprofile=生产的cprofile文件名称 文件夹
//go test -test.bench=".*" -count=5，可以看到如下结果： （使用-count可以指定执行多少次）

const (
	cfg = `{
		"appName": "budao_server_test",
		"httpAddr": ":8080",
		"db": {
			"mysql": {
				"budao": {
					"master": "61.147.187.183:6301",
					"slave": "183.36.121.66:8066",
					"passwd": "V4hwwxZqu",
					"user": "jxz_bd_test",
					"tableDesc": {
						"video_": 1,
						"comment_": 1,
						"user_": 1,
						"user_favor_comment_": 1,
						"user_favor_video_": 1,
						"user_shared_video_": 1,
						"user_question_": 1,
						"user_follow_topic_": 1,
						"topic_video_": 1
					}
				}
			},
			"redis": {
				"addr": "61.160.36.168:6395",
				"maxIdle": 50,
				"db": 0,
				"connTimeout": 5000,
				"readTimeout": 5000,
				"writeTimeout": 5000
			}
		},
		"auditconf": {
			"video": {
				"auditAddr": "221.228.107.63:10015",
				"appid": "301068",
				"chid": "301068",
				"secretId": "301068",
				"secretKey": "3Y302vwDgxqcPvNKPfq6MwIM8FtJy2Zb",
				"punishConf": {
					"listenAddr": ":8082"
				},
				"isDisable": true
			},
			"comment": {
				"auditAddr": "http://twapi.yy.com/txt/api",
				"appid": "999970035",
				"secretId": "7a7ca6c3",
				"secretKey": "ce9e0527cddc386ada3b07a7ca6c3131",
				"punishConf": {
					"listenAddr": ":8081"
				}
			},
			"title": {
				"auditAddr": "http://twapi.yy.com/txt/api",
				"appid": "999970033",
				"secretId": "73b4dcb2",
				"secretKey": "3f6c4c9d0573b4dcb2591ec4b095f673",
				"punishConf": {
					"listenAddr": ":8081"
				}
			}
		},
		"logconf": {
			"path": "/data/budao-server/logs/budao-server.log",
			"debug": true,
			"maxInterMSec": 500,
			"maxRemainNum": 2000,
			"reportFormat": {
				"appVersion": "V3.0",
				"scale": [50,100,150,200,250,300,350,400,450,500,550,600,650,700,750,800,850,900,950,1000]
			}
		},
		"commentConf": {
			"cronTaskDisable"  : true,
			"cronTaskInteranlSec" : 60,
			"returnHotCommentNum" : 2,
			"returnCommentNum" : 10,
			"returnReplyNum"   : 10,
			"explicitReplyNum" : 2
		}
	}`
)

var num int = 1
var commentId string

func Test_CommentVideo(t *testing.T) {
	var resp *pb.CommentVideoResponse
	var err error
	for i := 0; i < num; i++ {
		data := time.Now().Format("2006-01-02 15:04:05")
		content := "comment video from test unit. time: " + data
		s := GetServer()
		req := &pb.CommentVideoRequest{
			Header: &pb.Header{
				UserId: "47046836863550",
				Token:  "2902e346a0c2cd4364e286ba3a3b0782",
				DeviceInfo: &pb.DeviceInfo{
					DeviceId: "deviceid_123_ysz",
				},
			},
			VideoId: "60234336695331",
			Content: content,
		}

		resp, err = s.CommentVideo(context.Background(), req)
		if err != nil {
			t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
			break
		}
		if resp.Status.Code != pb.Status_OK {
			t.Errorf("code failed. resp:%v", resp) //记录一些你期望记录的信息
			break
		}
	}
	t.Logf("resp:%v", resp) //记录一些你期望记录的信息
	commentId = resp.CommentItem.CommentId
}

func Test_CommentVideo_IllegalContent(t *testing.T) {
	var resp *pb.CommentVideoResponse
	var err error
	//data := time.Now().Format("2006-01-02 15:04:05")
	//content := "comment video from test unit. time: " + data
	content := "习近平，共产党"
	s := GetServer()
	req := &pb.CommentVideoRequest{
		Header: &pb.Header{
			UserId: "47046836863550",
			Token:  "2902e346a0c2cd4364e286ba3a3b0782",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "deviceid_123_ysz",
			},
		},
		VideoId: "60234336695331",
		Content: content,
	}
	resp, err = s.CommentVideo(context.Background(), req)
	if err != nil {
		t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
		return
	}
	if resp.Status.Code != pb.Status_OK {
		t.Logf("code failed. resp:%v", resp) //记录一些你期望记录的信息
		t.Logf("illegal content is find. resp:%v", resp)
		return
	}
	t.Errorf("illegal content not find. resp:%v", resp) //记录一些你期望记录的信息
}

var replyId string

func Test_ReplyComment(t *testing.T) {
	var resp *pb.ReplyCommentResponse
	var err error
	for i := 0; i < num; i++ {
		data := time.Now().Format("2006-01-02 15:04:05")
		content := "comment video from test unit. time: " + data
		s := GetServer()
		req := &pb.ReplyCommentRequest{
			Header: &pb.Header{
				UserId: "47046836863550",
				Token:  "2902e346a0c2cd4364e286ba3a3b0782",
				DeviceInfo: &pb.DeviceInfo{
					DeviceId: "deviceid_123_ysz",
				},
			},
			Content: content,
			Type:    &pb.ReplyCommentRequest_CommentId{commentId},
		}
		resp, err = s.ReplyComment(context.Background(), req)
		if err != nil {
			t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
			break
		}
		if resp.Status.Code != pb.Status_OK {
			t.Errorf("code failed. resp:%v", resp) //记录一些你期望记录的信息
			break
		}
	}
	t.Logf("resp:%v", resp) //记录一些你期望记录的信息
	replyId = resp.ReplyItem.ReplyId
}

func Test_ReplyComment2(t *testing.T) {
	var resp *pb.ReplyCommentResponse
	var err error
	for i := 0; i < 10; i++ {
		data := time.Now().Format("2006-01-02 15:04:05")
		content := "comment video from test unit. time: " + data
		s := GetServer()
		req := &pb.ReplyCommentRequest{
			Header: &pb.Header{
				UserId: "47046836863550",
				Token:  "2902e346a0c2cd4364e286ba3a3b0782",
				DeviceInfo: &pb.DeviceInfo{
					DeviceId: "deviceid_123_ysz",
				},
			},
			Content: content,
			Type:    &pb.ReplyCommentRequest_ReplyId{replyId},
		}
		resp, err = s.ReplyComment(context.Background(), req)
		if err != nil {
			t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
		}
		if resp.Status.Code != pb.Status_OK {
			t.Errorf("code failed. resp:%v", resp) //记录一些你期望记录的信息
		}
	}
	t.Logf("resp:%v", resp) //记录一些你期望记录的信息
}

var hotCommentId string

func Test_GetVideoCommentList(t *testing.T) {
	var resp *pb.GetVideoCommentListResponse
	var err error
	var lastCommentId string
	for {
		s := GetServer()
		req := &pb.GetVideoCommentListRequest{
			Header: &pb.Header{
				UserId: "47046836863550",
				Token:  "2902e346a0c2cd4364e286ba3a3b0782",
				DeviceInfo: &pb.DeviceInfo{
					DeviceId: "deviceid_123_ysz",
				},
			},
			VideoId:       "60234336695331",
			LastCommentId: lastCommentId,
			//Content: content,
			//Type: &pb.ReplyCommentRequest_ReplyId{replyId},
		}
		resp, err = s.GetVideoCommentList(context.Background(), req)
		if err != nil {
			t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
		}
		if resp.Status.Code != pb.Status_OK {
			t.Errorf("code failed. resp:%v", resp) //记录一些你期望记录的信息
		}
		if hotCommentId == "" {
			hotCommentId = resp.HotCommentItems[0].CommentId
		}
		l := len(resp.CommentItems)
		t.Logf("##############################hot len:%d len:%d", len(resp.HotCommentItems), l)
		if resp.HasMore == false {
			break
		}
		lastCommentId = resp.CommentItems[l-1].CommentId
	}
	t.Logf("resp:%v", resp) //记录一些你期望记录的信息
}

func Test_GetCommentReplyList(t *testing.T) {
	var resp *pb.GetCommentReplyListResponse
	var err error
	var lastReplyId string
	for {
		s := GetServer()
		req := &pb.GetCommentReplyListRequest{
			Header: &pb.Header{
				UserId: "47046836863550",
				Token:  "2902e346a0c2cd4364e286ba3a3b0782",
				DeviceInfo: &pb.DeviceInfo{
					DeviceId: "deviceid_123_ysz",
				},
			},
			CommentId:   hotCommentId,
			LastReplyId: lastReplyId,
		}
		resp, err = s.GetCommentReplyList(context.Background(), req)
		if err != nil {
			t.Errorf("err:%v", err) // 如果不是如预期的那么就报错
		}
		if resp.Status.Code != pb.Status_OK {
			t.Errorf("code failed. resp:%v", resp) //记录一些你期望记录的信息
		}
		l := len(resp.ReplyItems)
		t.Logf("############################## len:%d ", l)
		if resp.HasMore == false {
			break
		}
		lastReplyId = resp.ReplyItems[l-1].ReplyId
	}
	t.Logf("resp:%v", resp) //记录一些你期望记录的信息
}

//func Benchmark_Test2(b *testing.B) {
//	for i := 0; i < b.N; i++ { //use b.N for looping
//		glog.Debug("index:", i)
//	}
//}

//测试从TestMain进入，依次执行测试用例，最后从TestMain退出
func TestMain(m *testing.M) {
	fmt.Println("begin")

	common.InitConfig(cfg)
	common.WG.Add(1)
	db.InitMysqlClient(common.GetConfig().DB.MySQL)
	db.InitRedisConnPool(common.GetConfig().DB.Redis)

	m.Run()

	fmt.Println("end")
}
