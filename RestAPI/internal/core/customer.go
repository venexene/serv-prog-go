package core

import (
	"math"

	"gravity-game-store/internal/entity"
	"gravity-game-store/internal/schema"
	"gravity-game-store/internal/store"
)

type CustomerSvc struct{ s *store.CustomerStore }

func NewCustomerSvc(s *store.CustomerStore) *CustomerSvc { return &CustomerSvc{s: s} }

func (svc *CustomerSvc) List(page, limit int) (*schema.PagedResp, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	cust, total, err := svc.s.List(page, limit)
	if err != nil {
		return nil, err
	}
	res := make([]schema.CustomerResp, len(cust))
	for i, c := range cust {
		res[i] = *custOut(&c)
	}
	return &schema.PagedResp{Data: res, Page: page, Limit: limit, Total: total, TotalPages: int64(math.Ceil(float64(total) / float64(limit)))}, nil
}

func (svc *CustomerSvc) Get(id uint) (*schema.CustomerResp, error) {
	c, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return custOut(c), nil
}

func (svc *CustomerSvc) Orders(id uint) ([]schema.OrderResp, error) {
	c, err := svc.s.WithOrders(id)
	if err != nil {
		return nil, ErrNotFound
	}
	res := make([]schema.OrderResp, len(c.Orders))
	for i, o := range c.Orders {
		res[i] = orderOut(&o)
	}
	return res, nil
}

func (svc *CustomerSvc) Add(req *schema.CreateCustomerReq) (*schema.CustomerResp, error) {
	c := &entity.Customer{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}
	if err := svc.s.Create(c); err != nil {
		return nil, err
	}
	return custOut(c), nil
}

func (svc *CustomerSvc) Upd(id uint, req *schema.UpdateCustomerReq) (*schema.CustomerResp, error) {
	c, err := svc.s.ByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if req.FirstName != "" {
		c.FirstName = req.FirstName
	}
	if req.LastName != "" {
		c.LastName = req.LastName
	}
	if req.Email != "" {
		c.Email = req.Email
	}
	if err := svc.s.Update(c); err != nil {
		return nil, err
	}
	return custOut(c), nil
}

func (svc *CustomerSvc) Del(id uint) error {
	if err := svc.s.Delete(id); err != nil {
		return ErrNotFound
	}
	return nil
}

func custOut(c *entity.Customer) *schema.CustomerResp {
	return &schema.CustomerResp{ID: c.ID, FirstName: c.FirstName, LastName: c.LastName, Email: c.Email, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

func orderOut(o *entity.CustOrder) schema.OrderResp {
	r := schema.OrderResp{
		ID: o.ID, OrderDate: o.OrderDate, CustomerID: o.CustomerID,
		ShippingMethodID: o.ShippingMethodID, DestAddressID: o.DestAddressID,
		CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
	}
	r.OrderLines = make([]schema.OrderLineResp, len(o.OrderLines))
	for i, ol := range o.OrderLines {
		olr := schema.OrderLineResp{ID: ol.ID, OrderID: ol.OrderID, GameID: ol.GameID, Price: ol.Price}
		if ol.Game.ID != 0 {
			g := gameOut(&ol.Game)
			olr.Game = &g
		}
		r.OrderLines[i] = olr
	}
	return r
}