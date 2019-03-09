####获取算法生成评论

- url:/api.NlpService/SimComment
- input:
```json
  {
  	"content":"补刀"
  }
```

- output:
```json
{
    "code": "200",
    "msg": "OK",
    "data": [
        {
            "content": "补刀补刀。。",
            "sim": 0.97598
        },
        {
            "content": "这后补刀。",
            "sim": 0.96959
        },
        {
            "content": "不带这么补刀的😂",
            "sim": 0.95461
        },
        {
            "content": "大难不死必补刀，一刀还比一刀靠",
            "sim": 0.94218
        },
        {
            "content": "神补刀。",
            "sim": 0.93528
        },
        {
            "content": "神补刀！😏",
            "sim": 0.93528
        },
        {
            "content": "看出来了这位就是来补刀的",
            "sim": 0.93035
        },
        {
            "content": "最后一句神补刀啊",
            "sim": 0.92312
        },
        {
            "content": "就和补刀一样舒服",
            "sim": 0.9146
        },
        {
            "content": "老师刀补得好啊",
            "sim": 0.90104
        },
        {
            "content": "补刀能手",
            "sim": 0.89961
        }
    ]
}
```

















