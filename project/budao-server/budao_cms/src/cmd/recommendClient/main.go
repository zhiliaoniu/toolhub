package main

import (
	"base"
	"common"
	"context"
	"db"
	"flag"
	"fmt"
	"github.com/sumaig/glog"
	"net/http"
	"os"
	"service"
	"strconv"
	"strings"
	"time"
	pb "twirprpc"
)

func init() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version == true {
		fmt.Println(common.VERSION)
		os.Exit(0)
	}

	common.ParseConfig(*cfg)
	base.InitLogger(common.GetConfig().LoggerConf)
	indexPath = common.GetConfig().Extension["indexPath"]
	i, err := strconv.ParseInt(common.GetConfig().Extension["interval"], 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	interval = time.Duration(i) * time.Second
	db.InitMysqlClient(common.GetConfig().DB.MySQL)

	base.InitLogCollector(common.GetConfig().LogCollectorConf)
}

var interval time.Duration
var indexPath string

func main() {
	strconv.ParseUint(common.GetConfig().Extension["interval"], 10, 64)
	client := pb.NewRecommendServiceProtobufClient(common.GetConfig().Extension["recommendUrl"], &http.Client{})
	for {
		queryVideo := `select vid,"title",type,duration,view_num, favor_num, comment_num,"tag",UNIX_TIMESTAMP(create_time),topic_id from video_0 where state=2 and op_state=0 and duration < 1800 and parse_type in (1,2,3,4,5,6,8,9,10)`
		rows, err := db.Query(service.BUDAODB, queryVideo)
		if err != nil {
			glog.Error(err)
			continue
		}
		allData := make([]*VideoMsq, 0, 5000)
		allVid := make([]string, 0, 5000)
		videoMap := make(map[string]*VideoMsq)
		for rows.Next() {
			temp := new(VideoMsq)
			err := rows.Scan(&temp.Vid, &temp.Title, &temp.Type, &temp.Duration, &temp.ViewNum, &temp.FaveNum, &temp.CommentNum, &temp.Tag, &temp.CTime, &temp.TopicID)
			if err != nil {
				glog.Error(err)
				continue
			}
			videoMap[temp.Vid] = temp
			allData = append(allData, temp)
			allVid = append(allVid, temp.Vid)
		}
		rows.Close()

		queryVideo = fmt.Sprintf(`select post_vid,ifnull(play_count,0),praise_count from video_data where post_vid in (%v)`, strings.Join(allVid, `,`))
		rows, err = db.Query(service.SPIDERDB, queryVideo)
		if err != nil {
			glog.Error(err)
			continue
		}
		for rows.Next() {
			var vid string
			var pv int
			var praise_count int
			err := rows.Scan(&vid, &pv, &praise_count)
			if err != nil {
				glog.Error(err)
				continue
			}
			temp := videoMap[vid]
			temp.SourcePV = pv
			if temp.SourcePV <= 0 {
				temp.SourcePV = praise_count
			}
		}
		rows.Close()
		filePath := indexPath + "/recommend_index_" + fmt.Sprintf("%v", time.Now().Format("20060102150405"))
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0664)
		if err != nil {
			glog.Error(err)
			continue
		}
		for _, v := range allData {
			_, err := f.WriteString(v.String() + "\n")
			if err != nil {
				glog.Error(err)
				continue
			}
		}
		f.Close()
		req := pb.ReloadIndexRequest{filePath}
		resp, err := client.ReloadIndex(context.Background(), &req)
		if err != nil {
			glog.Error(err)
			continue
		}
		glog.Info(resp.Status)
		<-time.After(interval)
	}
}

type VideoMsq struct {
	Vid        string
	Title      string
	Type       string
	Duration   string
	ViewNum    string
	FaveNum    int
	CommentNum string
	Tag        string
	CTime      string
	TopicID    string
	SourcePV   int
}

func (v VideoMsq) String() string {
	return fmt.Sprintf(`{"%v":["%v",%v,%v,%v,%v,%v,"%v",%v,"%v",%v]}`,
		v.Vid, v.Title, v.Type, v.Duration, v.ViewNum, v.FaveNum, v.CommentNum, v.Tag, v.CTime, v.TopicID, v.SourcePV)
}
