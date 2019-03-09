package main

import (
	"base"
	"common"
	"db"
	"flag"
	"fmt"
	"github.com/sumaig/glog"
	"math/rand"
	"os"
	"runtime"
	"service"
	"strings"
	"time"
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

	//查询即刻
	querysql := `SELECT v_source_id,post_vid FROM video_data where source=4 and status=1 and  source_type like "%33%" `
	rows, err := db.Query(service.SPIDERDB, querysql)
	if err != nil {
		fmt.Println(err)
		return
	}

	rs := rand.NewSource(time.Now().Unix())
	r := rand.New(rs)
	fmt.Println("ok")
	defer rows.Close()
	for rows.Next() {
		comments := make([]service.Comment, 0, 200)
		var (
			vsid string
			vid  string
		)
		err := rows.Scan(&vsid, &vid)
		if err != nil {
			continue
		}

		qsql := fmt.Sprintf(`select id,v_source_id,content,source,favor_num,user_id,create_time from comment_data where v_source_id='%v' and status=0`, vsid)
		tempr, err := db.Query(service.SPIDERDB, qsql)
		if err != nil {
			glog.Error(err)
			continue
		}
		for tempr.Next() {
			temp := service.Comment{}

			err := tempr.Scan(&temp.PId, &temp.VSourceID, &temp.Content, &temp.Source, &temp.FavorNum, &temp.Uid, &temp.CTime)
			if err != nil {
				glog.Error(err)
				continue
			}

			temp.CTime = common.GetTimeStr(temp.CTime)

			comments = append(comments, temp)
		}
		tempr.Close()
		for _, temp := range comments {
			cid, autoIncreId, tableName, err := common.GetItemId(service.SPIDERDB, service.COMMENT_)
			if temp.Uid == "0" {
				temp.Uid = fmt.Sprintf(`%v`, r.Int63())
			}
			u := service.GetRandUserByHash(temp.Source + "_" + temp.Uid)
			temp.Content = strings.Replace(temp.Content, `"`, `\"`, -1)
			u.Name = strings.Replace(u.Name, `"`, `\"`, -1)
			updateSql := fmt.Sprintf(`update %v set state=2,content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v",favor_num=%v,create_time='%v' where id=%v`, tableName, temp.Content, cid, vid, u.Uid, u.Name, u.Photo, temp.FavorNum, temp.CTime, autoIncreId)
			_, err = db.Exec(service.BUDAODB, updateSql)
			if err != nil {
				glog.Error(err)
				continue
			}
			updateSpiderCommentSql := fmt.Sprintf(`update comment_data set status=1,post_cid=%v where id=%v`, cid, temp.PId)
			db.Exec(service.SPIDERDB, updateSpiderCommentSql)
		}

	}

}
