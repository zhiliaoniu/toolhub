package transfer

import (
	"common"
	"database/sql"
	"db"
	"fmt"
	"testing"

	"github.com/sumaig/glog"
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
					"master": "221.228.79.244:8066",
					"passwd": "xywaD3kfz",
					"user": "jxz_db@jxz_bd",
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
		}
	}`
)

//func Test_UploadPic(t *testing.T) {
//	//filepath := "http://image1.pearvideo.com/cont/20180404/cont-1315106-11126803.png"
//	//filepath := "http://mat1.gtimg.com/www/qq2018/imgs/qq_logo_2018x2.png"
//	filepath := "http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaELHLVzc5LRtdYyh117SAUKQPB2m0iaHLduicdW3PrxDYWa8eFtFD0LlRr3Yj3qCwqWUetibUOibrgJFrw/132"
//	url, err := UploadPic(filepath)
//	if err != nil {
//		t.Errorf("err:%v", err)
//	}
//	t.Logf("url:%s, err:%v", url, err)
//}

func Test_UpdateOnlineVideoPic(t *testing.T) {
	//1.read url
	//querySql := fmt.Sprintf("select vid, coverurl from video_0 limit 10")
	likeStr, err := db.MysqlEscapeString("%jxzimg%")
	querySql := fmt.Sprintf("select vid, coverurl from video_0 where coverurl not like '%s'", likeStr)
	var rows *sql.Rows
	rows, err = db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()
	var vid, coverurl string
	num := 1
	for rows.Next() {
		err = rows.Scan(&vid, &coverurl)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
		glog.Debug("vid:%s, coverurl:%s", vid, coverurl)
		//2.upload
		url, err := UploadPic(coverurl)
		if err != nil {
			glog.Error("upload pic :%s failed. vid:%s, err:%v", coverurl, vid, err)
			continue
		}
		glog.Debug("num:%d\tnewUrl:%s\n\n", num, url)
		num++
		//3.update

		querySql := fmt.Sprintf("update video_0 set coverurl='%s' where vid=%s", url, vid)
		_, err = db.Exec(common.BUDAODB, querySql)
		if err != nil {
			glog.Error("update mysql failed. querySql:%s, vid:%s, err:%v", querySql, vid, err)
			continue
		}
	}
}

//测试从TestMain进入，依次执行测试用例，最后从TestMain退出
func TestMain(m *testing.M) {
	fmt.Println("begin")

	common.InitConfig(cfg)
	//db.Start()
	db.InitMysqlClient(common.GetConfig().DB.MySQL)
	//db.InitRedisConnPool()

	m.Run()

	fmt.Println("end")
}
