package api

type LoginRsp struct {
	Token    string   `json:"access_token"`
	Expire   int      `json:"expires_in"`
	UserInfo UserInfo `json:"userInfo"`
}

type UserInfo struct {
	Uid       int64  `json:"uid"`
	NickName  string `json:"nickname"`
	Signature string `json:"signature"`
	Avatar    string `json:"avatar"`
}
