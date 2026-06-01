package api

import (
	"net/http"
	"strconv"

	"gravity-game-store/internal/core"
	"gravity-game-store/internal/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthorCtrl struct {
	svc *core.AuthorSvc
	log *logrus.Logger
}

func NewAuthorCtrl(s *core.AuthorSvc, l *logrus.Logger) *AuthorCtrl { return &AuthorCtrl{svc: s, log: l} }

func (h *AuthorCtrl) GetAll(c *gin.Context) {
	p, l := pageLimit(c)
	r, err := h.svc.List(p, l)
	if err != nil {
		h.log.WithError(err).Error("list authors fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *AuthorCtrl) GetByID(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	a, err := h.svc.Get(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get author fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AuthorCtrl) GetGames(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	g, err := h.svc.Games(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get author games fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *AuthorCtrl) Create(c *gin.Context) {
	var req schema.CreateAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad author body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	a, err := h.svc.Add(&req)
	if err != nil {
		h.log.WithError(err).Error("create author fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", a.ID).Info("author created")
	c.JSON(http.StatusCreated, a)
}

func (h *AuthorCtrl) Update(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	var req schema.UpdateAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad author body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	a, err := h.svc.Upd(id, &req)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("update author fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", a.ID).Info("author updated")
	c.JSON(http.StatusOK, a)
}

func (h *AuthorCtrl) Delete(c *gin.Context) {
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
		h.log.WithError(err).Error("delete author fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", id).Info("author deleted")
	c.Status(http.StatusNoContent)
}

func paramUint(c *gin.Context, key string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(v), err
}

func pageLimit(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}