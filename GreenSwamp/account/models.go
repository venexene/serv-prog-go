package account

import "time"

type IdentityUser struct {
	UserID       uint      `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Email        string    `gorm:"column:email;uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"column:username;uniqueIndex;not null" json:"username"`
	DisplayName  string    `gorm:"column:display_name;not null" json:"display_name"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (IdentityUser) TableName() string { return "identity_users" }

type AuthSession struct {
	SessionID uint      `gorm:"column:session_id;primaryKey;autoIncrement" json:"session_id"`
	UserID    uint      `gorm:"column:user_id;index;not null" json:"user_id"`
	TokenHash string    `gorm:"column:token_hash;uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at" json:"revoked_at,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (AuthSession) TableName() string { return "auth_sessions" }