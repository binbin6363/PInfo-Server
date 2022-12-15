package api

// LoginReq 登录请求
type LoginReq struct {
	UserName string `json:"mobile"`
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
	Uid      int64  `json:"uid"`
	UserName string `json:"username"`
	Sequence int64  `json:"sequence"` // 差量拉取的序列号
}

// TalkListRsp 会话列表回包
type TalkListRsp struct {
	TalkList []*TalkInfo `json:"items"`
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

// ModifyUsersSettingReq 修改用户设置
type ModifyUsersSettingReq struct {
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Motto    string `json:"motto"`
	NickName string `json:"nickname"`
}

// ModifyUsersSettingRsp 修改用户设置
type ModifyUsersSettingRsp struct {
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
	Token string
	Uid   int64
}

// ContactListRsp 联系人列表回包
type ContactListRsp struct {
	ContactList []*ContactInfo `json:"rows"`
}

type ContactInfo struct {
	Id           int64  `json:"id"`
	Nickname     string `json:"nickname"`
	Gender       int    `json:"gender"`
	Motto        string `json:"motto"`
	Avatar       string `json:"avatar"`
	FriendRemark string `json:"friend_remark"`
	IsOnline     int    `json:"is_online"`
}

type SendTextMsgReq struct {
	ClientMsgId int64  `json:"client_msg_id"` // 消息去重
	ReceiverId  int64  `json:"receiver_id"`
	TalkType    int    `json:"talk_type"`
	Text        string `json:"text"`
	Uid         int64
}

type SendTextMsgRsp struct {
	Content SendTextMsgContent `json:"content"`
}

type SendTextMsgEvtNotice struct {
	Event   string             `json:"event"`
	Content SendTextMsgContent `json:"content"`
}

type SendTextMsgContent struct {
	Data       SendTextMsgData `json:"data"`
	ReceiverId int64           `json:"receiver_id"`
	SenderId   int64           `json:"sender_id"`
	TalkType   int             `json:"talk_type"`
}

type SendTextMsgData struct {
	Id         int64  `json:"id"`
	Sequence   int64  `json:"sequence"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int64  `json:"user_id"`
	ReceiverId int64  `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

type MsgRecordsReq struct {
	Uid      int64
	MinMsgId int64 // record_id
	PeerId   int64 // receiver_id
	TalkType int   // talk_type
	Limit    int   // limit
}

type MsgRecordsRsp struct {
	Limit       int          `json:"limit"`
	MaxRecordId int64        `json:"record_id"`
	Rows        []MessageRow `json:"rows"`
}

type MessageRow struct {
	Id         int64  `json:"id"`
	Sequence   int64  `json:"sequence"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int64  `json:"user_id"`
	ReceiverId int64  `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// RegisterReq 注册用户
type RegisterReq struct {
	NickName string `json:"nickname"`
	UserName string `json:"mobile"`
	Password string `json:"password"`
	SmsCode  string `json:"sms_code"`
	Platform string `json:"platform"`
}

type ContactSearchReq struct {
	UserName string `json:"mobile"`
}

type ContactSearchRsp struct {
	Uid int64 `json:"id"`
}

type ContactDetailReq struct {
	Uid       int64 `json:"id"`
	ContactId int64 `json:"user_id"`
}

type ContactDetailRsp struct {
	Gender         int    `json:"gender"`
	FriendStatus   int    `json:"friend_status"`
	FriendApply    int    `json:"friend_apply"`
	NickNameRemark string `json:"nickname_remark"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	UserName       string `json:"mobile"`
	Motto          string `json:"motto"`
	NickName       string `json:"nickname"`
	Uid            int64  `json:"id"`
}

type AddContactReq struct {
	Uid        int64  `json:"uid"`
	UserName   string `json:"username"`
	ContactID  int64  `json:"friend_id"`
	RemarkName string `json:"remark"`
}

type ApplyAddContactReq struct {
	Uid      int64  `json:"uid"` // 我的ID
	UserName string `json:"username"`
	ApplyId  int64  `json:"apply_id"` // 对方ID
	Remark   string `json:"remark"`
}

type AddContactRsp struct {
}

type EditContactInfoReq struct {
	Uid        int64  `json:"id"`        // 我的ID
	UserName   string `json:"username"`  // 我的用户名
	ContactId  int64  `json:"friend_id"` // 好友ID
	RemarkName string `json:"remarks"`   // 好友备注名
}

// CommRsp 针对请求的ack
type CommRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
