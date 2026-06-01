package api

import (
	"net/http"

	"gravity-game-store/internal/core"
	"gravity-game-store/internal/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomerCtrl struct {
	svc *core.CustomerSvc
	log *logrus.Logger
}

func NewCustomerCtrl(s *core.CustomerSvc, l *logrus.Logger) *CustomerCtrl {
	return &CustomerCtrl{svc: s, log: l}
}

func (h *CustomerCtrl) GetAll(c *gin.Context) {
	p, l := pageLimit(c)
	r, err := h.svc.List(p, l)
	if err != nil {
		h.log.WithError(err).Error("list customers fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *CustomerCtrl) GetByID(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	cust, err := h.svc.Get(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get customer fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *CustomerCtrl) GetOrders(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	orders, err := h.svc.Orders(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get orders fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *CustomerCtrl) Create(c *gin.Context) {
	var req schema.CreateCustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad customer body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	cust, err := h.svc.Add(&req)
	if err != nil {
		h.log.WithError(err).Error("create customer fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", cust.ID).Info("customer created")
	c.JSON(http.StatusCreated, cust)
}

func (h *CustomerCtrl) Update(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	var req schema.UpdateCustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad customer body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	cust, err := h.svc.Upd(id, &req)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("update customer fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", cust.ID).Info("customer updated")
	c.JSON(http.StatusOK, cust)
}

func (h *CustomerCtrl) Delete(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	if err := h.svc.Del(id); err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("delete customer fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", id).Info("customer deleted")
	c.Status(http.StatusNoContent)
}