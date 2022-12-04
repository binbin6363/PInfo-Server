package api

// LoginReq 登录请求
type LoginReq struct {
	Mobile   string `json:"mobile"`
	PassWord string `json:"password"`
	Platform string `json:"platform"`
}

// LoginRsp 登录回包
type LoginRsp struct {
	Type     string        `json:"type"` // Bearer
	Token    string        `json:"access_token"`
	Expire   int           `json:"expires_in"`
	UserInfo UserBasicInfo `json:"userInfo"`
}

// UserBasicInfo .
type UserBasicInfo struct {
	Uid       int64  `json:"uid"`
	NickName  string `json:"nickname"`
	Signature string `json:"signature"`
	Avatar    string `json:"avatar"`
}

// TalkListReq 会话列表请求
type TalkListReq struct {
}

// TalkListRsp 会话列表回包
type TalkListRsp struct {
	TalkList []TalkInfo `json:"items"`
}

// TalkInfo .
type TalkInfo struct {
	ID         int64  `json:"id"`
	Type       int    `json:"talk_type"`
	ReceiverId int64  `json:"receiver_id"`
	IsTop      int    `json:"is_top"`
	IsDisturb  int    `json:"is_disturb"`
	IsOnline   int    `json:"is_online"`
	IsRobot    int    `json:"is_robot"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	RemarkName string `json:"remark_name"`
	UnreadNum  int    `json:"unread_num"`
	MsgText    string `json:"msg_text"`
	UpdatedAt  string `json:"updated_at"`
}

// UserSettingReq 用户设置请求
type UserSettingReq struct {
}

// SettingInfo 用户设置
type SettingInfo struct {
	KeyboardEventNotify string `json:"keyboard_event_notify"`
	NotifyCueTone       string `json:"notify_cue_tone"`
	ThemeBagImg         string `json:"theme_bag_img"`
	ThemeColor          string `json:"theme_color"`
	ThemeMode           string `json:"theme_mode"`
}

// UserDetailInfo 用户详细信息
type UserDetailInfo struct {
	IsQiYe   bool   `json:"is_qiye"`
	Gender   int    `json:"gender"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
	Motto    string `json:"motto"`
	NickName string `json:"nickname"`
	Uid      int64  `json:"uid"`
}

// UserSettingRsp 用户设置回包
type UserSettingRsp struct {
	SettingInfo SettingInfo    `json:"setting"`
	UserInfo    UserDetailInfo `json:"user_info"`
}

// UserDetailReq 用户详情请求
type UserDetailReq struct {
}

// UserDetailRsp 用户详情回包
type UserDetailRsp struct {
	UserDetailInfo
}

// UnReadNumReq 消息未读数请求
type UnReadNumReq struct {
}

// UnReadNumRsp 消息未读数回包
type UnReadNumRsp struct {
	UnreadNum int `json:"unread_num"`
}

// ContactListReq 联系人列表请求
type ContactListReq struct {
}

// ContactListRsp 联系人列表回包
type ContactListRsp struct {
	ContactList []ContactInfo
}

type ContactInfo struct {
	Id           int    `json:"id"`
	Nickname     string `json:"nickname"`
	Gender       int    `json:"gender"`
	Motto        string `json:"motto"`
	Avatar       string `json:"avatar"`
	FriendRemark string `json:"friend_remark"`
	IsOnline     int    `json:"is_online"`
}

type SendTextMsgReq struct {
	ReceiverId int    `json:"receiver_id"`
	TalkType   int    `json:"talk_type"`
	Text       string `json:"text"`
}

type SendTextMsgRsp struct {
	Id         int    `json:"id"`
	TalkType   int    `json:"talk_type"`
	ReceiverId int    `json:"receiver_id"`
	SenderId   int    `json:"sender_id"`
	Name       string `json:"name"`
	RemarkName string `json:"remark_name"`
	Avatar     string `json:"avatar"`
	IsDisturb  int    `json:"is_disturb"`
	IsTop      int    `json:"is_top"`
	IsOnline   int    `json:"is_online"`
	IsRobot    int    `json:"is_robot"`
	UnreadNum  int    `json:"unread_num"`
	Content    string `json:"content"`
	DraftText  string `json:"draft_text"`
	MsgText    string `json:"msg_text"`
	IndexName  string `json:"index_name"`
	CreatedAt  string `json:"created_at"`
}

type SendTextMsgEvtRsp struct {
	Event   string             `json:"event"`
	Content SendTextMsgContent `json:"content"`
}

type SendTextMsgContent struct {
	Data       SendTextMsgData `json:"data"`
	ReceiverId int             `json:"receiver_id"`
	SenderId   int             `json:"sender_id"`
	TalkType   int             `json:"talk_type"`
}

type SendTextMsgData struct {
	Id         int    `json:"id"`
	Sequence   int    `json:"sequence"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int    `json:"user_id"`
	ReceiverId int    `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
