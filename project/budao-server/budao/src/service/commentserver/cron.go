package commentserver

import (
	"common"
	"db"
	"fmt"
	"strconv"
	"time"

	"service/util"

	"github.com/sumaig/glog"
)

const (
	COMMENT_MODIFY_OFFSET = "comment_modify_offset"
)

func (s *Server) cronTask() {
	ticker := time.NewTicker(time.Second * time.Duration(s.cronTaskInternalSec))
	ticker2 := time.NewTicker(time.Minute * time.Duration(s.cronTaskInternalSec))
	for {
		if s.cronTaskDisable {
			break
		}
		select {
		case <-ticker.C:
			UpdateCommentWeight()
		case <-ticker2.C:
			UpdateCommentInfoByTraverseMysql()
		case <-time.After(time.Second):
		}
	}
}

func UpdateCommentWeight() {
	//1.read offset of comment_weight_modify_record in mysql from redis
	autoIncreStr, _ := db.GetString(COMMENT_MODIFY_OFFSET)
	autoIncre, _ := strconv.Atoi(autoIncreStr)

	//2.read modify record from mysql by offset
	querySql := fmt.Sprintf("select id, cid, vid, parentcomid, to_comment_id, old_weight, new_weight, reason from comment_weight_modify_record where id > %d", autoIncre)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		var cid, parentCid, toCid, reason string
		var vid uint64
		var oldWeight, newWeight int
		err = rows.Scan(&id, &cid, &vid, &parentCid, &toCid, &oldWeight, &newWeight, &reason)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
		//3.update hot comment rank of video by comment's new weight
		//3.1 update weight
		key := fmt.Sprintf("%s%d", util.COMMENT_DYNAMIC_PREFIX, vid%common.COMMENT_DYNAMIC_KEY_NUM)
		fields := make([]interface{}, 0)
		field := fmt.Sprintf("%s%s", util.WEIGHT_PREFIX, cid)
		fields = append(fields, field, newWeight)
		if _, err := db.HMSet(key, fields); err != nil {
			glog.Error("hmset failed. err:%v", err)
			continue
		}
		//3.2 update rank
		if err := UpdateCommentRank(cid, vid); err != nil {
			glog.Error("update comment rank failed. err:%v", err)
			continue
		}
	}

	//update offset
	_ = db.SetString(COMMENT_MODIFY_OFFSET, id)
}
