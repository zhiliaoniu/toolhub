####微信回调地址

-method:post
- url:/api.WechatService/WechatRedirect
- input:
```param
    code

```

- output:
```json
{
    "code":"200",
    "msg":"OK",
    "data":{
        "token":"4otEkYAcGiWvWqT1BhaGv+nyF1IpmtgYwLI6XP6JqNAGZP1LpM7cUDAPu9e+LmLaZFWjhq/ohyyZvktgfIxVdg==",
        "userId":"HouShanJie"
        }
}
```

####用户token验证
header设置 X-Automatic-Token

{
    "code": "unauthenticated",
    "msg": "授权失败"
}











