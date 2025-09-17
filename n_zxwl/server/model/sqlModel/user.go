package sqlModel

import (
	"time"
)

// UserProfile 用户模型
type UserProfile struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string     `gorm:"size:50;not null;unique" json:"username"`
	Email          string     `gorm:"size:255;not null;unique" json:"email"`
	PasswordHash   string     `gorm:"size:255;not null" json:"-"`
	DisplayName    string     `gorm:"size:100;not null" json:"displayName"`
	AvatarURL      string     `gorm:"size:500" json:"avatarUrl"`
	Gender         int8       `gorm:"default:0" json:"gender"`
	BirthYear      int        `gorm:"type:year" json:"birthYear"`
	Location       string     `gorm:"size:255" json:"location"`
	Bio            string     `gorm:"type:text" json:"bio"`
	IsOnline       bool       `gorm:"default:false" json:"isOnline"`
	LastOnlineTime *time.Time `json:"lastOnlineTime"`
	DeviceInfo     string     `gorm:"type:json" json:"deviceInfo"`
	AuthToken      string     `gorm:"size:255" json:"-"`
	TokenExpiry    *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`

	// 关联关系
	Meetings []ChatMeetingParticipant `gorm:"foreignKey:UserID" json:"-"`
	//Friends          []UserFriend             `gorm:"foreignKey:UserID" json:"-"`
	SentMessages     []ChatMessage `gorm:"foreignKey:SenderID" json:"-"`
	ReceivedMessages []ChatMessage `gorm:"foreignKey:ReceiverID" json:"-"`
}

// ChatMeetingRoom 会议室模型
type ChatMeetingRoom struct {
	ID                 string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	HostID             int        `gorm:"not null" json:"hostId"`
	Title              string     `gorm:"size:255;not null" json:"title"`
	Description        string     `gorm:"type:text" json:"description"`
	Password           string     `gorm:"size:100" json:"-"`
	MaxParticipants    int        `gorm:"default:10" json:"maxParticipants"`
	IsRecording        bool       `gorm:"default:false" json:"isRecording"`
	Status             string     `gorm:"type:enum('scheduled','ongoing','ended');default:'scheduled'" json:"status"`
	ScheduledStartTime *time.Time `json:"scheduledStartTime"`
	ActualStartTime    *time.Time `json:"actualStartTime"`
	EndTime            *time.Time `json:"endTime"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`

	// 关联关系
	Host         UserProfile              `gorm:"foreignKey:HostID" json:"host"`
	Participants []ChatMeetingParticipant `gorm:"foreignKey:MeetingID" json:"participants"`
	Messages     []ChatMessage            `gorm:"foreignKey:MeetingID" json:"-"`
}

// ChatMeetingParticipant 会议参与者模型
type ChatMeetingParticipant struct {
	ID           int        `gorm:"primaryKey;autoIncrement" json:"id"`
	MeetingID    string     `gorm:"type:varchar(36);not null" json:"meetingId"`
	UserID       int        `gorm:"not null" json:"userId"`
	JoinTime     time.Time  `json:"joinTime"`
	LeaveTime    *time.Time `json:"leaveTime"`
	Role         string     `gorm:"type:enum('host','co-host','participant');default:'participant'" json:"role"`
	IsMuted      bool       `gorm:"default:false" json:"isMuted"`
	IsVideoOff   bool       `gorm:"default:false" json:"isVideoOff"`
	IsHandRaised bool       `gorm:"default:false" json:"isHandRaised"`
	ConnectionID string     `gorm:"size:255" json:"-"`
	DeviceType   string     `gorm:"size:50" json:"deviceType"`

	// 关联关系
	User    UserProfile     `gorm:"foreignKey:UserID" json:"user"`
	Meeting ChatMeetingRoom `gorm:"foreignKey:MeetingID" json:"meeting"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	MeetingID   *string   `gorm:"type:varchar(36)" json:"meetingId"`
	SenderID    int       `gorm:"not null" json:"senderId"`
	ReceiverID  *int      `gorm:"default:null" json:"receiverId"`
	MessageType string    `gorm:"type:enum('text','image','file','system');default:'text'" json:"messageType"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	IsRead      bool      `gorm:"default:false" json:"isRead"`
	CreatedAt   time.Time `json:"createdAt"`

	// 关联关系
	Meeting  ChatMeetingRoom `gorm:"foreignKey:MeetingID" json:"meeting"`
	Sender   UserProfile     `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver *UserProfile    `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
func (ChatMeetingRoom) TableName() string {
	return "chat_meeting_room"
}
func (ChatMeetingParticipant) TableName() string {
	return "chat_meeting_participant"
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
