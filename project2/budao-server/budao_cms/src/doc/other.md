####添加图片

- url:/other/uploadPic
- 传输方式: form-data  post
| 字段名称  | 字段注释   | 是否必须 | 备注  |
| :--------   | -----:  | -----: | :----:  |
| picName     | 本地图片名称/网络图片地址 |  是 | eg:like.jpg  图片类型支持 .jpg .jpeg .png  |
| fileType     | picName的类型 |  否 | 类型是网络图片地址的时候，fileType = text 类型是本地图片时候 fileType为空 |

返回值
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "picUrl": "http://jxzimg.bs2ul.yy.com/ec9db9b4eee528cd2872be689950ccf0"
    }
}
```

####app统计每日用户相关数-导出execl

- url:/other/user/xlsx
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
    "data":nil
}
```
####app统计每日视频相关数-导出execl

- url:/other/video/xlsx
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
    "data": nil
}
```