package posts

import "time"

type User struct {
	UserID      uint      `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Username    string    `gorm:"column:username;uniqueIndex;not null" json:"username"`
	DisplayName string    `gorm:"column:display_name;not null" json:"display_name"`
	AvatarURL   *string   `gorm:"column:avatar_url" json:"avatar_url"`
	Bio         *string   `gorm:"column:bio" json:"bio"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"is_active"`

	Posts []Post `gorm:"foreignKey:UserID;references:UserID"`
}

func (User) TableName() string { return "users" }

type Post struct {
	PostID       uint      `gorm:"column:post_id;primaryKey;autoIncrement" json:"post_id"`
	UserID       uint      `gorm:"column:user_id;not null" json:"user_id"`
	Content      string    `gorm:"column:content;not null" json:"content"`
	PostType     string    `gorm:"column:post_type;not null" json:"post_type"`
	MediaURL     *string   `gorm:"column:media_url" json:"media_url"`
	MediaType    *string   `gorm:"column:media_type" json:"media_type"`
	AltText      *string   `gorm:"column:alt_text" json:"alt_text"`
	ThumbnailURL *string   `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	ParentPostID *uint     `gorm:"column:parent_post_id" json:"parent_post_id"`

	User             User          `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Event            *Event        `gorm:"foreignKey:PostID;references:PostID" json:"event,omitempty"`
	Interactions     []Interaction `gorm:"foreignKey:PostID;references:PostID" json:"interactions,omitempty"`
	Tags             []Tag         `gorm:"many2many:post_tags;joinForeignKey:PostID;JoinReferences:TagID" json:"tags,omitempty"`
	ViewCount        int64         `gorm:"column:view_count;default:0"`
	InteractionCount int64         `gorm:"-" json:"interaction_count"`
	LikeCount        int64         `gorm:"-" json:"like_count"`
	ReribbCount      int64         `gorm:"-" json:"reribb_count"`
}

func (Post) TableName() string { return "posts" }

type Interaction struct {
	InteractionID   uint      `gorm:"column:interaction_id;primaryKey;autoIncrement" json:"interaction_id"`
	UserID          uint      `gorm:"column:user_id;not null" json:"user_id"`
	PostID          uint      `gorm:"column:post_id;not null" json:"post_id"`
	InteractionType string    `gorm:"column:interaction_type;not null" json:"interaction_type"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Content         *string   `gorm:"column:content" json:"content"`
}

func (Interaction) TableName() string { return "interactions" }

type Event struct {
	EventID     uint      `gorm:"column:event_id;primaryKey;autoIncrement" json:"event_id"`
	PostID      uint      `gorm:"column:post_id;uniqueIndex" json:"post_id"`
	EventTime   time.Time `gorm:"column:event_time;not null" json:"event_time"`
	Location    string    `gorm:"column:location;not null" json:"location"`
	HostOrg     *string   `gorm:"column:host_org" json:"host_org"`
	RSVPCount   int       `gorm:"column:rsvp_count;default:0" json:"rsvp_count"`
	MaxCapacity *int      `gorm:"column:max_capacity" json:"max_capacity"`
}

func (Event) TableName() string { return "events" }

type Tag struct {
	TagID      uint      `gorm:"column:tag_id;primaryKey;autoIncrement" json:"tag_id"`
	TagName    string    `gorm:"column:tag_name;uniqueIndex;not null" json:"tag_name"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UsageCount int       `gorm:"column:usage_count;default:0" json:"usage_count"`
}

func (Tag) TableName() string { return "tags" }

type PostTag struct {
	PostID uint `gorm:"column:post_id;primaryKey"`
	TagID  uint `gorm:"column:tag_id;primaryKey"`
}

func (PostTag) TableName() string { return "post_tags" }

type TrendingPond struct {
	TagID       uint   `gorm:"column:tag_id" json:"tag_id"`
	TagName     string `gorm:"column:tag_name" json:"tag_name"`
	RecentPosts int64  `gorm:"column:recent_posts" json:"recent_posts"`
	URL         string `gorm:"-" json:"url"`
}