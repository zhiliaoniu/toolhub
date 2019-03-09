package timelineserver

import (
	"common"
	"context"
	"db"
	"fmt"
	"service/topicserver"
	"testing"
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
					"listenAddr": ":8081"
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
		},
		"timelineConf": {
			"cronTaskDisable"  : false,
			"cronTaskInteranlSec" : 60,
			"maxRetVideoNum" : 10,
			"retPerTopicViewNum" : 10
		}
	}`
)

func Test_GetSubscribedTimeLine(t *testing.T) {
	topicserver.InitUserSubscribedTopic()

	//1.Direction_ALL
	s := GetServer()
	req := &pb.GetSubscribedTimeLineRequest{
		Header: &pb.Header{
			UserId: "34327992036",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "83262F5C-1D82-48D3-8E0B-1D1FC1638591",
			},
		},
		Direction: pb.Direction_ALL, //Direction_TOP Direction_BOTTOM
	}
	resp, err := s.GetSubscribedTimeLine(context.Background(), req)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Logf("resp:%v", resp)
}

func Test_GetTimeLine(t *testing.T) {

	//1.Direction_ALL
	s := GetServer()
	req := &pb.GetTimeLineRequest{
		Header: &pb.Header{
			UserId: "34327992036",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "83262F5C-1D82-48D3-8E0B-1D1FC1638591",
			},
		},
	}
	resp, err := s.GetTimeLine(context.Background(), req)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Logf("resp:%v", resp)
}

//测试从TestMain进入，依次执行测试用例，最后从TestMain退出
func TestMain(m *testing.M) {
	fmt.Println("begin")

	common.InitConfig(cfg)
	db.Start()
	db.InitRedisConnPool()
	LoadRecommendVIDIntoCache("video_0")

	m.Run()

	fmt.Println("end")
}
