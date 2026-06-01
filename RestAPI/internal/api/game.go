package api

import (
	"net/http"

	"gravity-game-store/internal/core"
	"gravity-game-store/internal/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GameCtrl struct {
	svc *core.GameSvc
	log *logrus.Logger
}

func NewGameCtrl(s *core.GameSvc, l *logrus.Logger) *GameCtrl { return &GameCtrl{svc: s, log: l} }

func (h *GameCtrl) GetAll(c *gin.Context) {
	p, l := pageLimit(c)
	r, err := h.svc.List(p, l)
	if err != nil {
		h.log.WithError(err).Error("list games fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *GameCtrl) GetByID(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	g, err := h.svc.Get(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get game fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *GameCtrl) GetAuthors(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	a, err := h.svc.Authors(id)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("get game authors fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *GameCtrl) Create(c *gin.Context) {
	var req schema.CreateGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad game body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	g, err := h.svc.Add(&req)
	if err != nil {
		h.log.WithError(err).Error("create game fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", g.ID).Info("game created")
	c.JSON(http.StatusCreated, g)
}

func (h *GameCtrl) Update(c *gin.Context) {
	id, err := paramUint(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad id"})
		return
	}
	var req schema.UpdateGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad game body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	g, err := h.svc.Upd(id, &req)
	if err != nil {
		if err == core.ErrNotFound {
			c.JSON(http.StatusNotFound, schema.ErrResp{Error: "not found"})
			return
		}
		h.log.WithError(err).Error("update game fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", g.ID).Info("game updated")
	c.JSON(http.StatusOK, g)
}

func (h *GameCtrl) Delete(c *gin.Context) {
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
		h.log.WithError(err).Error("delete game fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	h.log.WithField("id", id).Info("game deleted")
	c.Status(http.StatusNoContent)
}