package api

import (
	"net/http"

	"gravity-game-store/internal/core"
	"gravity-game-store/internal/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthCtrl struct {
	svc *core.AuthSvc
	log *logrus.Logger
}

func NewAuthCtrl(s *core.AuthSvc, l *logrus.Logger) *AuthCtrl { return &AuthCtrl{svc: s, log: l} }

func (h *AuthCtrl) Login(c *gin.Context) {
	var req schema.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad login body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	tok, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		h.log.WithError(err).Warn("login fail")
		c.JSON(http.StatusUnauthorized, schema.ErrResp{Error: "bad creds"})
		return
	}
	c.JSON(http.StatusOK, tok)
}

func (h *AuthCtrl) Register(c *gin.Context) {
	var req schema.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("bad reg body")
		c.JSON(http.StatusBadRequest, schema.ErrResp{Error: "bad body", Message: err.Error()})
		return
	}
	_, err := h.svc.Register(req.Username, req.Password)
	if err != nil {
		h.log.WithError(err).Warn("reg fail")
		c.JSON(http.StatusConflict, schema.ErrResp{Error: "username taken"})
		return
	}
	tok, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		h.log.WithError(err).Error("login after reg fail")
		c.JSON(http.StatusInternalServerError, schema.ErrResp{Error: "internal"})
		return
	}
	c.JSON(http.StatusCreated, tok)
}