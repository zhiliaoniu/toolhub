package configserver

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"

	"common"
	"db"
	pb "twirprpc"
)

const (
	TAB_KEY     = "global_tab"
	CHANNEL_KEY = "global_channel"
)

type Server struct {
	cronTaskDisable     bool //定时任务是否可用
	cronTaskInternalSec int  //定时任务执行间隔
}

func GetServer() *Server {
	server := &Server{}
	server.initServer()
	go server.cronTask()
	return server
}

func (s *Server) initServer() {
	s.cronTaskDisable = false
	s.cronTaskInternalSec = 10
}

func (s *Server) Close() {
	glog.Debug("config server close")
	s.cronTaskDisable = true
}

func (s *Server) cronTask() {
	ticker := time.NewTicker(time.Second * time.Duration(s.cronTaskInternalSec))
	for {
		if s.cronTaskDisable {
			glog.Debug("exit config server crontab")
			break
		}
		select {
		case <-ticker.C:
			s.updateTabCache()
			s.updateChannelCache()
		case <-time.After(time.Second):
		}
	}
}

func (s *Server) updateTabCache() {
	querySql := fmt.Sprintf("select id, name from tab where disable = 0")
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()

	tabItems := make([]*pb.TabItem, 0)
	for rows.Next() {
		tabItem := &pb.TabItem{}
		err = rows.Scan(&tabItem.TabId, &tabItem.Name)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
		tabItems = append(tabItems, tabItem)
	}

	resp := &pb.GetTabListResponse{}
	resp.Status = &pb.Status{
		Code:       pb.Status_OK,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
	resp.TabItems = tabItems
	respStr, err := proto.Marshal(resp)
	if err != nil {
		glog.Error("proto marshal failed. err:%v", err)
		return
	}
	if err := db.SetString(TAB_KEY, string(respStr)); err != nil {
		glog.Error("set string failed. err:%v", err)
	}
	//glog.Debug("update tab success.resp:%v", resp)
}

func (s *Server) updateChannelCache() {
	querySql := fmt.Sprintf("select id, name from channel where disable = 0")
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()

	channelItems := make([]*pb.ChannelItem, 0)
	for rows.Next() {
		channelItem := &pb.ChannelItem{}
		err = rows.Scan(&channelItem.ChannelId, &channelItem.Name)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
		if common.GetConfig().TestChannelShow == "false" {
			if channelItem.Name == "测试" {
				continue
			}
		}
		channelItems = append(channelItems, channelItem)
	}

	resp := &pb.GetChannelListResponse{}
	resp.Status = &pb.Status{
		Code:       pb.Status_OK,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
	resp.ChannelItems = channelItems
	respStr, err := proto.Marshal(resp)
	if err != nil {
		glog.Error("proto marshal failed. err:%v", err)
		return
	}
	if err := db.SetString(CHANNEL_KEY, string(respStr)); err != nil {
		glog.Error("set string failed. err:%v", err)
	}

	//glog.Debug("update channel success.resp:%v", resp)
}
