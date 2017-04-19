##sync项目背景
sync项目主要就是用来定时同步数据库的数据到redis这个缓存数据（该项目主要是同步用户登陆的数据到redis），这样就可以大大的提交数据库查询的效率。

####项目结构说明：

**如图所示：**

![](http://i.imgur.com/xnG2ukV.png)

**说明：**


- cmd 程序可行文件主目录
- sync.go 程序main函数文件
- app/options 程序配置文件处理
- server.go 实际业务处理程序文件
- docs 项目文档文件夹
- pkg/logger 日志处理封装
- types 具体同步处理的结构体以及相关结构体的处理
- util 一些常用工具的封装

####处理流程：

程序开启一个定时器（定时器的执行周期可配置，在config.json文件中配置），定时的重指定的mysql（数据类型可扩展）数据库同步指定（目前只是user，wxuser）数据表查询数据，然后根据指定的key的生成规则，将所有用户的数据存入redis（数据没有失效时间）。

####配置文件说明：

```
{
  "mysql": {
    "dsn": "root:root@(192.168.0.40:3306)/test?charset=utf8"
  },
  "redis": {
    "addr": "192.168.0.40:6379",
    "password": "",
    "requiredPW": true,
    "db": 0,
    "poolSize": 100
  },
  "interval": 60
}
```

 - mysql mysql的链接信息
 - redis redis的链接信息
 - addr redis的链接地址
 - password redis的password
 - requiredPW redis是否需要密码
 - db 用的是redis中的那个数据库（redis默认有15个数据库）
 - poolSize redis链接池的链接数的大小
 - interval  定时器执行的时间周期，单位为秒
 
 
####redis中数据结构说明：

系统用户的key的前缀分两种：
 1. `user_mail_`   对应key = `user_mail_`+mail
 2. `user_mob_`     对应key = `user_mob_`+mob
 
系统用户：
```
key  
user_mail_449264675@qq.com

value 

{
    "corp_id": 1,
    "tye": "0",
    "pwd": "123456",
    "nickname": "nick",
    "realname": "just",
    "mob": "15989511262",
    "mail": "449264675@qq.com",
    "sex": 0,
    "qq": "449264675",
    "headimgurl": "",
    "is_del": 0,
    "create_at": "2017-03-09 11:07:13",
    "revise_at": "2017-03-09 11:07:16",
    "id": "01"
}
```

```
key  
user_mob_15989511262

value 

{
    "corp_id": 1,
    "tye": "0",
    "pwd": "123456",
    "nickname": "nick",
    "realname": "just",
    "mob": "15989511262",
    "mail": "449264675@qq.com",
    "sex": 0,
    "qq": "449264675",
    "headimgurl": "",
    "is_del": 0,
    "create_at": "2017-03-09 11:07:13",
    "revise_at": "2017-03-09 11:07:16",
    "id": "01"
}
```

微信用户：

```
key  
wxuser_123456

value 

{
    "corp_id": 1,
    "tye": "0",
    "pwd": "123456",
    "nickname": "nick",
    "realname": "just",
    "mob": "15989511262",
    "mail": "449264675@qq.com",
    "sex": 0,
    "qq": "449264675",
    "headimgurl": "",
    "is_del": 0,
    "create_at": "2017-03-09 11:07:13",
    "revise_at": "2017-03-09 11:07:16",
    "id": "01",
    "wxapp_id": "jufucx",
    "city": "深圳",
    "country": "中国",
    "province": "广东",
    "subscribe_time": "2017-03-09 11:41:51",
    "remark": "",
    "user_id": "01",
    "is_attn": 0,
    "open_id": "123456"
}
```