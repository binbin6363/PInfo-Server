package model

// UserInfo 用户信息表
type UserInfo struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	UserName   string `gorm:"column:username"`
	PassHash   string `gorm:"column:passhash"`
	NickName   string `gorm:"column:nickname"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
	Avatar     string `gorm:"column:avatar"`
	Gender     int    `gorm:"column:gender"`
	UserTag    string `gorm:"column:user_tag"`
	Motto      string `gorm:"column:motto"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (UserInfo) TableName() string {
	return "user_infos"
}

type ContactType int

const (
	ContactStranger    ContactType = 1 // 陌生人
	ContactWaitMeApply ContactType = 2 // 好友申请等我审批
	ContactSendApply   ContactType = 3 // 我已发送好友申请
	ContactFriend      ContactType = 4 // 已通过申请，好友关系

	MsgTypeText     = 1  // 文本
	MsgTypeImg      = 2  // 图片
	MsgTypeRecord   = 3  // 会话记录消息
	MsgTypeCode     = 4  // 代码块消息
	MsgTypeVote     = 5  // 投票消息
	MsgTypeAudio    = 6  // 音频消息
	MsgTypeVideo    = 7  // 视频消息
	MsgTypeLogin    = 8  // 登录消息
	MsgTypeFile     = 9  // 文件消息
	MsgTypeLocation = 10 // 位置消息

	SingleTalkType = 1 // 单聊
	GroupTalkType  = 2 // 群聊

)

// Contacts 联系人好友
type Contacts struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	ContactID  int64  `gorm:"column:contact_id"`
	RemarkName string `gorm:"column:remark_name"`
	Status     int    `gorm:"column:status"` // 0非好友，1好友申请中，2已是好友
	Sequence   int64  `gorm:"column:sequence"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Contacts) TableName() string {
	return "contacts"
}

// UserContact 好友关系，联表查询结果，临时数据结构
type UserContact struct {
	Uid          int64  `gorm:"column:uid"`        // 我的ID
	ContactID    int64  `gorm:"column:contact_id"` // 好友ID
	Nickname     string `gorm:"column:nickname"`   // 好友昵称
	Gender       int    `gorm:"column:gender"`
	Motto        string `gorm:"column:motto"`
	Avatar       string `gorm:"column:avatar"`
	FriendRemark string `gorm:"column:remark_name"` // 好友备注
	Status       int    `gorm:"column:status"`
	CreateTime   int64  `gorm:"column:create_time"`
}

// Conversations 会话列表
type Conversations struct {
	ID                 int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid                int64  `gorm:"column:uid"`
	ContactID          int64  `gorm:"column:contact_id"`
	ConversationType   int    `gorm:"column:conversation_type"`
	ConversationName   string `gorm:"column:conversation_name"`
	ConversationStatus int    `gorm:"column:conversation_status"`
	Unread             int    `gorm:"column:unread"`
	MsgDigest          string `gorm:"column:msg_digest"`
	Sequence           int64  `gorm:"column:sequence"`
	CreateTime         int64  `gorm:"column:create_time"`
	UpdateTime         int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Conversations) TableName() string {
	return "conversations"
}

// ConversationDetails 会话列表详情，通过多表联合查询得到，临时数据结构
type ConversationDetails struct {
	ID               int64  `gorm:"column:id"`
	Uid              int64  `gorm:"column:uid"`
	ContactID        int64  `gorm:"column:contact_id"`
	ConversationType int    `gorm:"column:conversation_type"`
	Unread           int    `gorm:"column:unread"`
	MsgDigest        string `gorm:"column:msg_digest"`
	Sequence         int64  `gorm:"column:sequence"`
	UpdateTime       int64  `gorm:"column:update_time"`
	UserName         string `gorm:"column:username"`
	RemarkName       string `gorm:"column:remark_name"`
	ConversationName string `gorm:"column:conversation_name"`
	Motto            string `gorm:"column:motto"`
	Avatar           string `gorm:"column:avatar"`
}

// SingleMessages 单人消息
type SingleMessages struct {
	ID          int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid         int64  `gorm:"column:uid"`
	MsgID       int64  `gorm:"column:msgid"`
	ClientMsgID int64  `gorm:"column:client_msgid"`
	SenderID    int64  `gorm:"column:senderid"`
	ReceiverID  int64  `gorm:"column:receiverid"`
	MsgType     int    `gorm:"column:msg_type"`
	Content     string `gorm:"column:content"`
	MsgStatus   int    `gorm:"column:msg_status"`
	CreateTime  int64  `gorm:"column:create_time"`
	UpdateTime  int64  `gorm:"column:update_time"`
	MediaInfo   string `gorm:"column:media"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (SingleMessages) TableName() string {
	return "single_messages"
}

// GroupMessages 群消息
type GroupMessages struct {
	ID          int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	GroupID     int64  `gorm:"column:groupid"`
	MsgID       int64  `gorm:"column:msgid"`
	ClientMsgID int64  `gorm:"column:client_msgid"`
	SenderID    int64  `gorm:"column:senderid"`
	MsgType     int    `gorm:"column:msg_type"`
	Content     string `gorm:"column:content"`
	MsgStatus   int    `gorm:"column:msg_status"`
	CreateTime  int64  `gorm:"column:create_time"`
	UpdateTime  int64  `gorm:"column:update_time"`
	MediaInfo   string `gorm:"column:media"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (GroupMessages) TableName() string {
	return "group_messages"
}

// Groups 群信息列表
type Groups struct {
	ID            int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	GroupID       int64  `gorm:"column:group_id"`
	GroupName     string `gorm:"column:group_name"`
	GroupStatus   int    `gorm:"column:group_status"` // 1,正常;2,封禁;3,全员禁言
	GroupAvatar   string `gorm:"column:group_avatar"`
	GroupTag      string `gorm:"column:group_tag"`
	GroupAnnounce string `gorm:"column:group_announce"`
	Sequence      int64  `gorm:"column:sequence"`
	CreateTime    int64  `gorm:"column:create_time"`
	UpdateTime    int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Groups) TableName() string {
	return "groups"
}

type GroupInfoList struct {
	GroupID       int64  `gorm:"column:group_id"`
	GroupName     string `gorm:"column:group_name"`
	GroupAvatar   string `gorm:"column:group_avatar"`
	GroupAnnounce string `gorm:"column:group_announce"`
	Leader        int64  `gorm:"column:leader"`
}

// GroupMembers 群信息列表
type GroupMembers struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	GroupID    int64  `gorm:"column:group_id"`
	Uid        int64  `gorm:"column:uid"`
	UserRole   int    `gorm:"column:user_role"`   // 成员角色，1普通成员，2群管理员
	Disturb    int    `gorm:"column:disturb"`     // 该成员是否设置群消息免打扰，0否，1是
	RemarkName string `gorm:"column:remark_name"` // 用户自己备注在群里的名字
	Sequence   int64  `gorm:"column:sequence"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (GroupMembers) TableName() string {
	return "group_members"
}

// GroupMemberInfoList 群成员具体信息，临时数据结构
type GroupMemberInfoList struct {
	Uid          int64  `gorm:"column:uid"` // 我的ID
	Gender       int    `gorm:"column:gender"`
	Motto        string `gorm:"column:motto"`
	Avatar       string `gorm:"column:avatar"`
	NickName     string `gorm:"column:nickname"`    // 个人备注
	FriendRemark string `gorm:"column:remark_name"` // 群内个人备注
	Status       int    `gorm:"column:status"`
	UserRole     int    `gorm:"column:user_role"`
	CreateTime   int64  `gorm:"column:create_time"`
	Sequence     int64  `gorm:"column:sequence"`
}

// GroupDetailInfo 我的群具体信息，临时数据结构
type GroupDetailInfo struct {
	GroupId     int64  `gorm:"column:group_id"`
	GroupName   string `gorm:"column:group_name"`
	GroupAvatar string `gorm:"column:group_avatar"`
	Notice      string `gorm:"column:group_announce"`
	UserRole    int    `gorm:"column:user_role"`
	VisitCard   string `gorm:"column:remark_name"`
	IsDisturb   int    `gorm:"column:disturb"`
	CreatedAt   string `gorm:"column:create_time"`
	ManagerName string
	IsManager   bool
}

// Classes 文章分类表
type Classes struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	Flag       int    `gorm:"column:flag"` // 文章分类标识，0默认
	Name       string `gorm:"column:name"` // 分类名字
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Classes) TableName() string {
	return "classes"
}

// Tags 文章tag表
type Tags struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	Flag       int    `gorm:"column:flag"` // 文章tag标识，0默认
	Name       string `gorm:"column:name"` // tag名字
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Tags) TableName() string {
	return "tags"
}

// Articles 文章表
type Articles struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	ClassId    int64  `gorm:"column:class_id"`
	Title      string `gorm:"column:title"`      // 文章标题
	Content    string `gorm:"column:content"`    // html格式内容
	MdContent  string `gorm:"column:md_content"` // md格式内容
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Articles) TableName() string {
	return "articles"
}
