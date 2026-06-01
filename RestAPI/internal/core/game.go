package core

import (
	"math"

	"gravity-game-store/internal/entity"
	"gravity-game-store/internal/schema"
	"gravity-game-store/internal/store"
)

type GameSvc struct{ s *store.GameStore }

func NewGameSvc(s *store.GameStore) *GameSvc { return &GameSvc{s: s} }

func (svc *GameSvc) List(page, limit int) (*schema.PagedResp, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	games, total, err := svc.s.List(page, limit)
	if err != nil {
		return nil, err
	}
	res := make([]schema.GameResp, len(games))
	for i, g := range games {
		res[i] = gameOut(&g)
	}
	return &schema.PagedResp{Data: res, Page: page, Limit: limit, Total: total, TotalPages: int64(math.Ceil(float64(total) / float64(limit)))}, nil
}

func (svc *GameSvc) Get(id uint) (*schema.GameResp, error) {
	g, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	r := gameOut(g)
	return &r, nil
}

func (svc *GameSvc) Authors(id uint) ([]schema.AuthorResp, error) {
	g, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	res := make([]schema.AuthorResp, len(g.Authors))
	for i, a := range g.Authors {
		res[i] = schema.AuthorResp{ID: a.ID, Name: a.Name, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
	}
	return res, nil
}

func (svc *GameSvc) Add(req *schema.CreateGameReq) (*schema.GameResp, error) {
	g := &entity.Game{
		Title: req.Title, Genre: req.Genre, Platform: req.Platform,
		PublisherID: req.PublisherID, ReleaseDate: req.ReleaseDate, NumPlayers: req.NumPlayers,
	}
	if err := svc.s.Create(g); err != nil {
		return nil, err
	}
	if len(req.AuthorIDs) > 0 {
		if err := svc.s.SetAuthors(g, req.AuthorIDs); err != nil {
			return nil, err
		}
	}
	created, err := svc.s.ByID(g.ID)
	if err != nil {
		return nil, err
	}
	r := gameOut(created)
	return &r, nil
}

func (svc *GameSvc) Upd(id uint, req *schema.UpdateGameReq) (*schema.GameResp, error) {
	g, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if req.Title != "" {
		g.Title = req.Title
	}
	if req.Genre != "" {
		g.Genre = req.Genre
	}
	if req.Platform != "" {
		g.Platform = req.Platform
	}
	if req.PublisherID > 0 {
		g.PublisherID = req.PublisherID
	}
	if req.ReleaseDate != "" {
		g.ReleaseDate = req.ReleaseDate
	}
	if req.NumPlayers > 0 {
		g.NumPlayers = req.NumPlayers
	}
	if err := svc.s.Update(g); err != nil {
		return nil, err
	}
	if req.AuthorIDs != nil {
		if err := svc.s.SetAuthors(g, req.AuthorIDs); err != nil {
			return nil, err
		}
	}
	upd, err := svc.s.ByID(g.ID)
	if err != nil {
		return nil, err
	}
	r := gameOut(upd)
	return &r, nil
}

func (svc *GameSvc) Del(id uint) error {
	if err := svc.s.Delete(id); err != nil {
		return ErrNotFound
	}
	return nil
}

func gameOut(g *entity.Game) schema.GameResp {
	r := schema.GameResp{
		ID: g.ID, Title: g.Title, Genre: g.Genre, Platform: g.Platform,
		PublisherID: g.PublisherID, ReleaseDate: g.ReleaseDate, NumPlayers: g.NumPlayers,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
	r.Authors = make([]schema.AuthorResp, len(g.Authors))
	for i, a := range g.Authors {
		r.Authors[i] = schema.AuthorResp{ID: a.ID, Name: a.Name, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
	}
	return r
}