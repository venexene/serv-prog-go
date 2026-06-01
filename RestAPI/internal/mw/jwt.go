package mw

import (
	"net/http"
	"strings"

	"gravity-game-store/internal/core"
	"gravity-game-store/internal/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ClaimsCtx = "claims"

func Auth(auth *core.AuthSvc, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if hdr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, schema.ErrResp{Error: "missing auth header"})
			return
		}
		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, schema.ErrResp{Error: "bad auth header format"})
			return
		}
		claims, err := auth.Validate(parts[1])
		if err != nil {
			log.WithError(err).Warn("bad jwt")
			c.AbortWithStatusJSON(http.StatusUnauthorized, schema.ErrResp{Error: "bad token"})
			return
		}
		c.Set(ClaimsCtx, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *core.Claims {
	v, ok := c.Get(ClaimsCtx)
	if !ok {
		return nil
	}
	claims, _ := v.(*core.Claims)
	return claims
}