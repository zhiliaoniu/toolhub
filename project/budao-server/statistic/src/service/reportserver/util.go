package reportserver

import (
	"common"
	"db"
	"fmt"
	"time"

	"github.com/sumaig/glog"
)

const (
	DATE_FORMAT  = "2006-01-02"
	DATE_FORMAT2 = "2006-01-02 15:04:05"
)

func GetYesterday() (yesterday string) {
	now := time.Now()
	yesterday = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format(DATE_FORMAT)
	return
}

func GetYesterdayBegin() (yesterdayBegin string) {
	now := time.Now()
	yesterdayBegin = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format(DATE_FORMAT2)
	return
}

func GetYesterdayEnd() (yesterdayEnd string) {
	now := time.Now()
	yesterdayEnd = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, -1, time.Local).Format(DATE_FORMAT2)
	return
}

func GetYesterdayRecordNum(tableName string) (num int, err error) {
	return GetYesterdayRecordNumWithName(tableName, "create_time")
}

func GetYesterdayRecordNumWithName(tableName string, timeFieldName string) (num int, err error) {
	return GetRecordNumByTime(tableName, timeFieldName, GetYesterdayBegin(), GetYesterdayEnd())
}

func GetRecordNumByTime(tableName, timeFieldName, beginTime, endTime string) (num int, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count int
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select count(1) from %s where %s > '%s' and %s < '%s'", tableName, timeFieldName, beginTime, timeFieldName, endTime)
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		num += count
	}
	return
}

func GetYesterdayDintinctRecordNum(tableName, distinctName string) (num int, err error) {
	return GetDintinctRecordNumByTime(tableName, distinctName, "create_time", GetYesterdayBegin(), GetYesterdayEnd())
}

func GetBeforeYesterdayDintinctRecordNum(tableName, distinctName string) (num int, err error) {
	return GetDintinctRecordNumByTime(tableName, distinctName, "create_time", "", GetYesterdayBegin())
}

func GetDintinctRecordNumByTime(tableName, distinctName, timeFieldName, beginTime, endTime string) (num int, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count int
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select count(distinct %s) from %s where %s > '%s' and %s < '%s'", distinctName, tableName, timeFieldName, beginTime, timeFieldName, endTime)
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		num += count
	}
	return
}

func GetAllRecordNum(tableName string) (num int, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count int
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select count(1) from %s", tableName)
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		num += count
	}
	return
}

func GetDictinctAllRecordNum(tableName, distinctName string) (num int, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count int
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select count(distinct %s) from %s", distinctName, tableName)
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		num += count
	}
	return
}

func GetYesterdayVideoRecordNumWithFlag(tableName string, flag int) (num int, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count int
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select count(1) from %s where flag=%d and create_time > '%s' and create_time < '%s'", tableName, flag, GetYesterdayBegin(), GetYesterdayEnd())
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		num += count
	}
	return
}

func GetYesterdayVideoViewTimeWithFlag(tableName string, flag int) (viewTime float32, err error) {
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[tableName])
	var count float32
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", tableName, i)
		querySql := fmt.Sprintf("select sum(video_play_time) from %s where flag=%d and create_time > '%s' and create_time < '%s'", tableName, flag, GetYesterdayBegin(), GetYesterdayEnd())
		tempRow, _ := db.QueryRow(common.BUDAODB, querySql)
		err = tempRow.Scan(&count)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
			return
		}
		viewTime += count
	}
	return
}
