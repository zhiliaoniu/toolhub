####添加banner

- url:/api.BannerService/BannerAdd
- input:
```json
  {
  	"picUrl":"wwww",                    //图片路径
  	"clickUrl":"sssssss",               //点击路径
  	"position":"push",                  //banner位置　　eg:topic...
  	"description":"shishasia"           //描述
  	"fromTime":"2018-06-23 19:23:33",   //banner的开始时间
    "toTime":"2018-06-24 19:23:33"      //banner的结束时间
  }
```

- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": null
}
```
####banner列表

- url:/api.BannerService/BannerList
- input:
```json
{
  "filter":{
	"status":"",　　　　　　　　　　//'状态　0,上线　1,下架',
	"position":"",              //banner位置
	"description":"",            //描述
	"fromTime":"",             //banner开始时间
	"toTime":""                //banner结束时间
	},
	"sort":"fromTime-desc",
	"num":1,
    "size":10
}
 ```
- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "data": [
            {
                "id": "1",
                "picUrl": "wwww",                            //图片路径
                "clickUrl": "sssssss",                  　　　//点击路径　
                "position": "push",                          //banner位置
                "description": "shishasia",                  //描述
                "status": "0",                               //状态　'状态　0,上线　1,下架'
                "createdAt": "2018-06-26 15:15:52",          //创建时间
                "updatedAt": "2018-06-26 15:15:52"           //更新时间
                "fromTime": "2018-06-24 15:15:52"            //banner的开始时间
                "toTime": "2018-06-24 15:15:52"              //banner的结束时间
            }
        ],
        "total": 1
    }
}
```

####获取banner信息

- url:/api.BannerService/BannerInfo
- input:
```json
   {
   	"id":"1"                 //id
   }
```
- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": {
        "data": {
            "id": "1",
            "picUrl": "wwww",
            "clickUrl": "sssssss",
            "position": "push",
            "description": "shishasia",
            "status": "0",
            "createdAt": "",
            "updatedAt": ""
            "fromTime": "2018-06-24 15:15:52"            //banner的开始时间
            "toTime": "2018-06-24 15:15:52"              //banner的结束时间
        }
    }
}
```

####编辑banner信息

- url:/api.BannerService/BannerModify
- input:
```json
   {
        "id":"1",                           //id
        "status":"0",                       //状态　'状态　0,上线　1,下架'
    	"picUrl":"wwww",                    //图片路径
    	"clickUrl":"sssssss",               //点击路径
    	"position":"push",                  //banner位置　　eg:topic...
    	"description":"shishasia"           //描述
    	"fromTime":"2018-06-23 19:23:33",   //banner的开始时间
        "toTime":"2018-06-24 19:23:33"      //banner的结束时间
    }
```
- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": null
}
```
####删除banner信息

- url:/api.BannerService/BannerDel
- input:
```json
   {
      "id":"1"                 //id
   }
```
- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": null
}
```














