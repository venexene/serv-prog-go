package core

import (
	"time"

	"gravity-game-store/internal/conf"
	"gravity-game-store/internal/entity"
	"gravity-game-store/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthSvc struct {
	users *store.UserStore
	cfg   *conf.Cfg
}

func NewAuthSvc(u *store.UserStore, c *conf.Cfg) *AuthSvc { return &AuthSvc{users: u, cfg: c} }

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *AuthSvc) Register(username, password string) (*entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &entity.User{Username: username, Password: string(hash), Role: "user"}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthSvc) Login(username, password string) (*TokenPair, error) {
	u, err := s.users.ByUsername(username)
	if err != nil {
		return nil, ErrBadCreds
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrBadCreds
	}
	access, err := s.genToken(u, s.cfg.JWTAccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := s.genToken(u, s.cfg.JWTRefreshTTL)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken: access, RefreshToken: refresh,
		TokenType: "Bearer", ExpiresIn: int(s.cfg.JWTAccessTTL.Seconds()),
	}, nil
}

func (s *AuthSvc) Validate(raw string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrBadToken
		}
		return []byte(s.cfg.JWTKey), nil
	})
	if err != nil {
		return nil, ErrBadToken
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, ErrBadToken
	}
	return claims, nil
}

func (s *AuthSvc) genToken(u *entity.User, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: u.ID, Username: u.Username, Role: u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTKey))
}