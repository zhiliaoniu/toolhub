package main

import (
	"base"
	"common"
	"context"
	"db"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"service"
	pb "twirprpc"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version == true {
		fmt.Println(common.VERSION)
		os.Exit(0)
	}

	common.ParseConfig(*cfg)
	base.InitLogger(common.GetConfig().LoggerConf)

	runtime.GOMAXPROCS(runtime.NumCPU())
	db.InitMysqlClient(common.GetConfig().DB.MySQL)

	base.InitLogCollector(common.GetConfig().LogCollectorConf)

	client := pb.NewTransferProtobufClient(common.GetConfig().Extension["transferUrl"], &http.Client{})
	//抖音
	rows, err := db.Query(service.SPIDERDB, `
select id,i_id,video_cover,play_url,video_title,video_duration,video_height,video_width from (select v_source_id,id,i_id,video_cover,play_url,video_title,video_duration,video_height,video_width from video_data where source=1 and status=0 and create_time>='2018-04-01 00:00:00')a left join   (select v_source_id from comment_data where source=1 group by v_source_id having count(id)>=10 ) c  on a.v_source_id =c.v_source_id where c.v_source_id is not null;
`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var (
			id       string
			vid      string
			Cover    string
			Url      string
			Title    string
			duration int32
			height   int32
			width    int32
		)
		rows.Scan(&id, &vid, &Cover, &Url, &Title, &duration, &height, &width)
		rpcReq := &pb.PostVideosRequest{}
		tVideos := make([]*pb.PostVideo, 0)
		postVideo := &pb.PostVideo{}
		postVideo.SourceVid = vid
		postVideo.Pic = Cover
		postVideo.VideoUrl = Url
		postVideo.Title = Title
		postVideo.Width = func() uint32 {
			if width <= 0 {
				return 540
			}
			return uint32(width)
		}()
		postVideo.Height = func() uint32 {
			if height <= 0 {
				return 960
			}
			return uint32(height)
		}()
		postVideo.Duration = func() uint32 {
			if duration <= 0 {
				return 30
			}
			return uint32(duration)
		}()
		postVideo.EVideoSource = pb.EVideoSource(1)
		postVideo.EVideoParseRule = pb.VideoParseRule_PARSE_RULE_DOUYIN
		tVideos = append(tVideos, postVideo)
		rpcReq.TVideos = tVideos

		rpcResp, rpcErr := client.PostVideos(context.Background(), rpcReq)
		if rpcErr != nil {
			fmt.Println(err)
			continue
		}
		if len(rpcResp.GetTResults()) > 0 {
			result := rpcResp.GetTResults()[0]
			//fmt.Println(result)
			data := struct {
				Vid string `json:"vid"`
			}{result.Vid}
			if data.Vid == "" {
				continue
			}

			updateSql := fmt.Sprintf(`update video_data set status=1,post_vid=%v where id=%v`, data.Vid, id)
			db.Exec(service.SPIDERDB, updateSql)
		} else {
			fmt.Println("return status error")
			continue
		}
	}

}
