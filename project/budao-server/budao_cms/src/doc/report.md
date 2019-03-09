####今日发布视频及问题说

- url:/api.ReportService/ReportOpToday
- output:
```json
{
    "allPush":30,//今日发布视频数
    "hasQuestion": 1 //有问题的数量
}
```

####push的vv数量

- url:/api.ReportService/ReportPushVV
- input:
```json
{
   "sTime":"2018-05-10 00:00:00",//开始时间
   "eTime":"2018-05-16 23:59:59",//结束时间
   "photoType":"ios"//平台 ios或安卓 值对应的是 ios/adr
}
```
- output
```json
[
    {
        "date": "2018-05-14",//日期
        "sumVV": "20"//vv数量
    }
]
```
####一天push的vv列表

- url:/api.ReportService/ReportPushVVList
- input:
```json
{
   "date":"2018-05-16",//日期
   "photoType":"adr"//平台 ios或安卓 值对应的是 ios/adr
}
```
- output
```json
[
     {
            "vid": "123",
            "sumVV": "20",
            "title": "123",//视频标题
            "duration": "0"//时长
     }
]
```
####内容数据  发布数据

- url:/api.ReportService/ReportContentPostData
- input:
```json
 {
    "filter":{
         "sTime":"2018-05-10 00:00:00",//开始时间
         "eTime":"2018-05-16 23:59:59",//结束时间
         "source":"1", //来源  西瓜、即刻。。。
         "topicID":"12725472897"  //话题id
    }
    "sort":"cTime-asc",
    "num":1,
    "size":10
 }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "postVideoCount": 5612,           //发布视频数
        "postVideoQueCount": 0            //发布视频存在问题的视频数
    }
}
```

####内容数据  话题数据

- url:/api.ReportService/ReportContentTopicData
- input:
```json
  {
        	"filter":{
        	    无
            },
            "sort":"cTime-asc",
        	"num":1,
        	"size":10
    }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "data": [
            {
                "queCount": "0",                            //问题数
                "topicName": "星座运势了解一下",               //话题名称
                "userCount": "0",                           //用户参与数
                "videoCount": "1"                           //视频数
            }
        ],
        "total": 3,                                        //话题数
        "tOCount": 0,                                       //话题开启数
        "tCCount": 0                                        //话题禁用数
    }
}
```
####app统计每日用户相关数

- url:/api.ReportService/ReportUserStatistics
- input:
```json
 {
 	"filter":{
 		"sTime":"",       //起始时间
 		"eTime":"",       //终止时间
 		"totalNum":"",    //当前用户总数
 		"activeNum":"",   //前一天活跃用户数
 		"newNum":""       //前一天新增用户数
 	},
 	"sort":"sTime-desc",
 	"num":1,
     "size":10
 }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "data": [
            {
                "id": "1",
                "totalNum": "10",
                "activeNum": "1",
                "newNum": "3",
                "cTime": "2018-06-22 11:59:53"
            }
        ],
        "total": 1
    }
}
```
####app统计每日视频相关数

- url:/api.ReportService/ReportVideoStatistics
- input:
```json
 {
 	"filter":{
 		"sTime":"",             //起始时间
 		"eTime":"",             //终止时间
 		"vExposeNum":"",        //视频曝光数
 		"vClictNum":"",         //视频点击
 		"vViewNum":"",          //视频观看
 		"vFavorNum":"",         //视频点赞
 		"commFavorNum":"",      //评论点赞
 		"commNum":"",           //评论发表
 		"tFollowNum":""         //话题订阅数
 	},
 	"sort":"sTime-desc",
 	"num":1,
     "size":10
 }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "data": [
            {
                "id": "1",
                "vExposeNum": "2",
                "vClictNum": "3",
                "vViewNum": "4",
                "vFavorNum": "5",
                "commFavorNum": "6",
                "commNum": "7",
                "tFollowNum": "8",
                "cTime": "2018-06-22 12:05:37"
            }
        ],
        "total": 1
    }
}
```

####app统计每日用户相关数-导出execl

- url:/api.ReportService/ReportUserStatisticsExecl
- input:
```json
 {
 	"filter":{
 		"sTime":"",       //起始时间
 		"eTime":"",       //终止时间
 		"totalNum":"",    //当前用户总数
 		"activeNum":"",   //前一天活跃用户数
 		"newNum":""       //前一天新增用户数
 	}
 }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data":"http://116.31.122.113/data/houshanjie/report_user_2015-06-15 08:52:32.xlsx"
}
```
####app统计每日视频相关数-导出execl

- url:/api.ReportService/ReportVideoStatisticsExecl
- input:
```json
 {
 	"filter":{
 		"sTime":"",             //起始时间
 		"eTime":"",             //终止时间
 		"vExposeNum":"",        //视频曝光数
 		"vClictNum":"",         //视频点击
 		"vViewNum":"",          //视频观看
 		"vFavorNum":"",         //视频点赞
 		"commFavorNum":"",      //评论点赞
 		"commNum":"",           //评论发表
 		"tFollowNum":""         //话题订阅数
 	}
 }
```
- output
```json
{
    "code": "200",
    "msg": "OK",
    "data": "http://116.31.122.113/data/houshanjie/report_video_2015-06-15 08:52:32.xlsx"
}
```