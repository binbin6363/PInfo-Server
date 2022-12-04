# PIM协议设计

*politewang，20221111*



## 协议请求鉴权

校验：
http请求头HEAD包含如下信息
{
"name": "Authorization",
"value": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6Imx1bWVuLWltIiwiZXhwIjoxNjY4MzU0NzkwLCJpYXQiOjE2NjgyNjgzOTAsImp0aSI6IjIwNTQifQ.y8yEESSsRrglsCESNercoXdljlHhChWsrgqjaKVc7zA"
}



## 协议请求列表

### POST     "/api/v1/auth/logout"



### POST     "/api/v1/auth/login"
请求：

```json
{
"mobile": "18798272054",
"password": "admin123", // todo: 要改。有点可怕，这里明文传密码
"platform": "web"
}
```

回包：

```json
{
"code": 200,
"message": "success",
"data": {
"type": "Bearer",
"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6Imx1bWVuLWltIiwiZXhwIjoxNjY4MzU0NzkwLCJpYXQiOjE2NjgyNjgzOTAsImp0aSI6IjIwNTQifQ.y8yEESSsRrglsCESNercoXdljlHhChWsrgqjaKVc7zA",
"expires_in": 86400
}
}
```

### GET      "/api/v1/talk/list"
请求：
回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"items": [{
			"id": 231,
			"talk_type": 1,
			"receiver_id": 1887,
			"is_top": 0,
			"is_disturb": 0,
			"is_online": 0,
			"is_robot": 0,
			"name": "9T3l7MalRA",
			"avatar": "http://im-serve0.gzydong.club/static/image/sys-head/u=2427656680,2114648815\u0026fm=26\u0026gp=0.jpg",
			"remark_name": "",
			"unread_num": 0,
			"msg_text": "????",
			"updated_at": "2022-11-11 22:33:40"
		}, {
			"id": 429,
			"talk_type": 2,
			"receiver_id": 147,
			"is_top": 0,
			"is_disturb": 0,
			"is_online": 0,
			"is_robot": 0,
			"name": "666",
			"avatar": "",
			"remark_name": "",
			"unread_num": 0,
			"msg_text": "@flora",
			"updated_at": "2022-09-29 14:06:23"
		}]
	}
}
```

### GET      "/api/v1/users/setting"
请求：
回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"setting": {
			"keyboard_event_notify": "",
			"notify_cue_tone": "",
			"theme_bag_img": "",
			"theme_color": "",
			"theme_mode": ""
		},
		"user_info": {
			"avatar": "https://im.gzydong.club/public/media/image/avatar/20221110/d587fecfe84d5b625c810b495b1393a1_200x200.png",
			"email": "837215079@qq.com",
			"gender": 2,
			"is_qiye": true,
			"mobile": "18798272054",
			"motto": "客服",
			"nickname": "vv_vv9",
			"uid": 2054
		}
	}
}
```

### GET      "/api/v1/users/detail"
请求：
回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"id": 2054,
		"mobile": "18798272054",
		"nickname": "vv_vv9",
		"avatar": "https://im.gzydong.club/public/media/image/avatar/20221110/d587fecfe84d5b625c810b495b1393a1_200x200.png",
		"gender": 2,
		"motto": "客服",
		"email": "837215079@qq.com",
		"birthday": "2022-11-09"
	}
}
```

### GET      "/api/v1/contact/apply/unread-num"
请求：

回包：

```json
{
    "code": 200,
    "message": "success",
    "data": {
      "unread_num": 0
    }
}
```

### GET      "/api/v1/contact/apply/records?page=1&page_size=10000"【demo未实现】

请求：

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"items": [],
		"paginate": {
			"page": 1,
			"size": 1000,
			"total": 0
		}
	}
}
```


### GET      "/api/v1/note/class/list"

请求：

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"items": [{
			"id": 3353,
			"class_name": "阿萨达大fds",
			"is_default": 0,
			"count": 2
		}, {
			"id": 3352,
			"class_name": "defau",
			"is_default": 0,
			"count": 0
		],
		"paginate": {
			"page": 1,
			"size": 100000,
			"total": 2
		}
	}
}
```


### GET      "/api/v1/note/tag/list"

请求：

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"tags": [{
			"id": 1,
			"tag_name": "docker2222",
			"count": 6
		}, {
			"id": 3,
			"tag_name": "Swoole",
			"count": 5
		}]
	}
}
```


### GET      "/api/v1/note/article/list?page=1&keyword=&find_type=1&cid="

请求：

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"items": [{
			"abstract": "大范甘迪管道符个",
			"class_id": 3353,
			"class_name": "阿萨达大fds",
			"created_at": "2022-11-12 15:49:25",
			"id": 434,
			"image": "",
			"is_asterisk": 0,
			"status": 1,
			"tags_id": "",
			"title": "123123",
			"updated_at": "2022-11-12 15:49:25"
		}, {
			"abstract": "12312313",
			"class_id": 3353,
			"class_name": "阿萨达大fds",
			"created_at": "2022-11-12 15:49:10",
			"id": 433,
			"image": "",
			"is_asterisk": 0,
			"status": 1,
			"tags_id": "",
			"title": "请编辑标题！！！",
			"updated_at": "2022-11-12 15:49:10"
		}],
		"paginate": {
			"page": 1,
			"size": 1000,
			"total": 2
		}
	}
}
```


### GET      "/api/v1/group/member/invites?group_id=0"
请求：
```json
{
  "name": "group_id",
  "value": "0"
}
```

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": [{
		"id": 1887,
		"nickname": "9T3l7MalRA",
		"gender": 0,
		"motto": "",
		"avatar": "http://im-serve0.gzydong.club/static/image/sys-head/u=2427656680,2114648815\u0026fm=26\u0026gp=0.jpg",
		"friend_remark": "",
		"is_online": 0
	}, {
		"id": 2022,
		"nickname": "bvXqS4dIgV",
		"gender": 0,
		"motto": "",
		"avatar": "http://im-serve0.gzydong.club/static/image/sys-head/886e1e38bcaf0fce870976dcdf9c090.jpg",
		"friend_remark": "",
		"is_online": 0
	}]
}
```


###  GET      "/api/v1/contact/search?mobile=wefwefweq"
请求：

回包：


### GET      "/api/v1/contact/list"
请求：


回包：

```json
{
	"code": 200,
	"message": "success",
	"data": [{
		"id": 1887,
		"nickname": "9T3l7MalRA",
		"gender": 0,
		"motto": "",
		"avatar": "http://im-serve0.gzydong.club/static/image/sys-head/u=2427656680,2114648815\u0026fm=26\u0026gp=0.jpg",
		"friend_remark": "",
		"is_online": 0
	}, {
		"id": 2022,
		"nickname": "bvXqS4dIgV",
		"gender": 0,
		"motto": "",
		"avatar": "http://im-serve0.gzydong.club/static/image/sys-head/886e1e38bcaf0fce870976dcdf9c090.jpg",
		"friend_remark": "",
		"is_online": 0
	}]
}
```


### GET      "/api/v1/group/list"
请求：

回包：

```json
{
	"code": 200,
	"message": "success",
	"data": {
		"rows": [{
			"id": 375,
			"group_name": "12121221",
			"avatar": "",
			"profile": "sadfasdf",
			"leader": 2,
			"is_disturb": 0
		}, {
			"id": 374,
			"group_name": "test",
			"avatar": "",
			"profile": "",
			"leader": 2,
			"is_disturb": 0
		}]
	}
}
```


### POST     "/api/v1/auth/register"

请求：

```json
// TODO: 明文密码需要改成密文
{
	"nickname": "wdqwdqw",
	"mobile": "12212233212",
	"password": "1234",
	"sms_code": "qw",
	"platform": "web"
}
```
回包：

```json
{
	
}
```



## 消息回调

### 消息发送回调

```json
{
	"event": "event_talk",
	"content": {
		"data": {
			"id": 75824,
			"talk_type": 1,
			"msg_type": 1,
			"user_id": 2054,
			"receiver_id": 2055,
			"nickname": "vv_vv9",
			"avatar": "https://im.gzydong.club/public/media/image/avatar/20221110/d587fecfe84d5b625c810b495b1393a1_200x200.png",
			"is_revoke": 0,
			"is_mark": 0,
			"is_read": 0,
			"content": "/nice",
			"created_at": "2022-11-13 21:30:20"
		},
		"receiver_id": 2055,
		"sender_id": 2054,
		"talk_type": 1
	}
}
```

### 消息阅读回调

```json
{
	"event": "event_talk_read",
	"content": {
		"ids": [75824],
		"receiver_id": 2054,
		"sender_id": 2055
	}
}
```

### 消息输入回调
```json
{
  "event": "event_talk_keyboard",
  "content": {
    "receiver_id": 2054,
    "sender_id": 2055
  }
}
```

```json
{
  "event": "event_talk",
  "content": {
    "data": {
      "id": 76446,
      "sequence": 6,
      "talk_type": 1,
      "msg_type": 1,
      "user_id": 2054,
      "receiver_id": 4101,
      "nickname": "wlb",
      "avatar": "https://im.gzydong.club/public/media/image/avatar/20221124/ea1bf7400e61fad835ad72c2c9e985b1_200x200.png",
      "is_revoke": 0,
      "is_mark": 0,
      "is_read": 0,
      "content": "12345",
      "created_at": "2022-11-28 13:28:40"
    },
    "receiver_id": 4101,
    "sender_id": 2054,
    "talk_type": 1
  }
}
```
