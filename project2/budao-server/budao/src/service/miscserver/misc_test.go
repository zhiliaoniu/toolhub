package miscserver

import (
	"common"
	"context"
	"db"
	"fmt"
	"testing"
	pb "twirprpc"
)

//go test -run=文件名字 -bench=bench名字 -cpuprofile=生产的cprofile文件名称 文件夹
//go test -test.bench=".*" -count=5，可以看到如下结果： （使用-count可以指定执行多少次）

//test
//"master": "61.147.187.183:6301",
//"passwd": "V4hwwxZqu",
//"user": "jxz_bd_test",

// online
//"master": "221.228.79.244:8066",
//"passwd": "xywaD3kfz",
//"user": "jxz_db@jxz_bd",

const (
	cfg = `{
		"appName": "budao_server_test",
		"httpAddr": ":8080",
		"db": {
			"mysql": {
				"budao": {
					"master": "61.147.187.183:6301",
					"passwd": "V4hwwxZqu",
					"user": "jxz_bd_test",
					"slave": "183.36.121.66:8066",
					"tableDesc": {
						"video_": 1,
						"comment_": 1,
						"user_": 1,
						"user_favor_comment_": 1,
						"user_favor_video_": 1,
						"user_shared_video_": 1,
						"user_question_": 1,
						"user_follow_topic_": 1,
						"topic_video_": 1,
						"user_token_": 1
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
		}
	}`
)

func Test_PostDeviceToken(t *testing.T) {
	s := &Server{}
	req := &pb.PostDeviceTokenRequest{
		Header: &pb.Header{
			UserId: "34327992036",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "83262F5C-1D82-48D3-8E0B-1D1FC1638591",
			},
		},
		DeviceToken: "db4e573046fb44c3bd0a5ece6e6f44c6d4f05dccd5e01c7234d2015ce4998d94",
	}
	resp, err := s.PostDeviceToken(context.Background(), req)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Logf("resp:%v, err:%v", resp, err)
}

func Test_GetVideoInfo(t *testing.T) {
	s := &Server{}
	req := &pb.GetVideoInfoRequest{
		Header: &pb.Header{
			UserId: "34327992036",
			DeviceInfo: &pb.DeviceInfo{
				DeviceId: "83262F5C-1D82-48D3-8E0B-1D1FC1638591",
			},
		},
		VideoId: "10976845541522",
	}
	resp, err := s.GetVideoInfo(context.Background(), req)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Logf("resp:%v, err:%v", resp, err)
}

//测试从TestMain进入，依次执行测试用例，最后从TestMain退出
func TestMain(m *testing.M) {
	fmt.Println("begin")

	common.InitConfig(cfg)
	db.Start()
	db.InitRedisConnPool()

	m.Run()

	fmt.Println("end")
}
