package core

import (
	"math"

	"gravity-game-store/internal/entity"
	"gravity-game-store/internal/schema"
	"gravity-game-store/internal/store"
)

type AuthorSvc struct{ s *store.AuthorStore }

func NewAuthorSvc(s *store.AuthorStore) *AuthorSvc { return &AuthorSvc{s: s} }

func (svc *AuthorSvc) List(page, limit int) (*schema.PagedResp, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	authors, total, err := svc.s.List(page, limit)
	if err != nil {
		return nil, err
	}
	res := make([]schema.AuthorResp, len(authors))
	for i, a := range authors {
		res[i] = *authorOut(&a)
	}
	return &schema.PagedResp{Data: res, Page: page, Limit: limit, Total: total, TotalPages: int64(math.Ceil(float64(total) / float64(limit)))}, nil
}

func (svc *AuthorSvc) Get(id uint) (*schema.AuthorResp, error) {
	a, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return authorOut(a), nil
}

func (svc *AuthorSvc) Games(id uint) ([]schema.GameResp, error) {
	a, err := svc.s.WithGames(id)
	if err != nil {
		return nil, ErrNotFound
	}
	res := make([]schema.GameResp, len(a.Games))
	for i, g := range a.Games {
		res[i] = gameOut(&g)
	}
	return res, nil
}

func (svc *AuthorSvc) Add(req *schema.CreateAuthorReq) (*schema.AuthorResp, error) {
	a := &entity.Author{Name: req.Name}
	if err := svc.s.Create(a); err != nil {
		return nil, err
	}
	return authorOut(a), nil
}

func (svc *AuthorSvc) Upd(id uint, req *schema.UpdateAuthorReq) (*schema.AuthorResp, error) {
	a, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	a.Name = req.Name
	if err := svc.s.Update(a); err != nil {
		return nil, err
	}
	return authorOut(a), nil
}

func (svc *AuthorSvc) Del(id uint) error {
	if err := svc.s.Delete(id); err != nil {
		return ErrNotFound
	}
	return nil
}

func authorOut(a *entity.Author) *schema.AuthorResp {
	return &schema.AuthorResp{ID: a.ID, Name: a.Name, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}