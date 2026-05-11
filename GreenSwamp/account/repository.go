package account

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&IdentityUser{}, &AuthSession{})
}

func (r *Repository) FindByLogin(ctx context.Context, login string) (*IdentityUser, error) {
	var user IdentityUser
	err := r.db.WithContext(ctx).
		Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)", login, login).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*IdentityUser, error) {
	var user IdentityUser
	err := r.db.WithContext(ctx).
		Where("LOWER(email) = LOWER(?)", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*IdentityUser, error) {
	var user IdentityUser
	err := r.db.WithContext(ctx).
		Where("LOWER(username) = LOWER(?)", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *IdentityUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&IdentityUser{}).
		Where("user_id = ?", userID).
		Update("last_login_at", now).Error
}

func (r *Repository) CreateSession(ctx context.Context, userID uint, rawToken string, expiresAt time.Time) error {
	session := AuthSession{
		UserID:    userID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(&session).Error
}

func (r *Repository) FindSession(ctx context.Context, rawToken string) (*AuthSession, error) {
	var s AuthSession
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", hashToken(rawToken), time.Now()).
		First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) RevokeSession(ctx context.Context, rawToken string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&AuthSession{}).
		Where("token_hash = ? AND revoked_at IS NULL", hashToken(rawToken)).
		Update("revoked_at", now).Error
}

func (r *Repository) CurrentUserFromRequest(ctx context.Context, req *http.Request, cookieName string) (*IdentityUser, bool) {
	cookie, err := req.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return nil, false
	}

	session, err := r.FindSession(ctx, cookie.Value)
	if err != nil {
		return nil, false
	}

	var user IdentityUser
	if err := r.db.WithContext(ctx).First(&user, "user_id = ?", session.UserID).Error; err != nil {
		return nil, false
	}

	return &user, true
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}