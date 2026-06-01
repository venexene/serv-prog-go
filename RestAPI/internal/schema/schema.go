package schema

import "time"

type LoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=3,max=100"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type PagedResp struct {
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type CreateAuthorReq struct {
	Name string `json:"name" binding:"required,min=1,max=400"`
}
type UpdateAuthorReq struct {
	Name string `json:"name" binding:"required,min=1,max=400"`
}
type AuthorResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateGameReq struct {
	Title       string `json:"title" binding:"required,min=1,max=400"`
	Genre       string `json:"genre"`
	Platform    string `json:"platform"`
	PublisherID int    `json:"publisher_id"`
	ReleaseDate string `json:"release_date"`
	NumPlayers  int    `json:"num_players"`
	AuthorIDs   []uint `json:"author_ids"`
}
type UpdateGameReq struct {
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Platform    string `json:"platform"`
	PublisherID int    `json:"publisher_id"`
	ReleaseDate string `json:"release_date"`
	NumPlayers  int    `json:"num_players"`
	AuthorIDs   []uint `json:"author_ids"`
}
type GameResp struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Genre       string       `json:"genre"`
	Platform    string       `json:"platform"`
	PublisherID int          `json:"publisher_id"`
	ReleaseDate string       `json:"release_date"`
	NumPlayers  int          `json:"num_players"`
	Authors     []AuthorResp `json:"authors,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateCustomerReq struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=200"`
	LastName  string `json:"last_name" binding:"required,min=1,max=200"`
	Email     string `json:"email" binding:"required,email,max=350"`
}
type UpdateCustomerReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type CustomerResp struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderResp struct {
	ID               uint              `json:"id"`
	OrderDate        time.Time         `json:"order_date"`
	CustomerID       uint              `json:"customer_id"`
	ShippingMethodID int               `json:"shipping_method_id"`
	DestAddressID    uint              `json:"dest_address_id"`
	OrderLines       []OrderLineResp   `json:"order_lines,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
type OrderLineResp struct {
	ID      uint      `json:"id"`
	OrderID uint      `json:"order_id"`
	GameID  uint      `json:"game_id"`
	Price   float64   `json:"price"`
	Game    *GameResp `json:"game,omitempty"`
}

type ErrResp struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}