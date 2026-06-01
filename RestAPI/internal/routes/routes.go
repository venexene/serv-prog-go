package routes

import (
	"net/http"

	"gravity-game-store/internal/api"
	"gravity-game-store/internal/core"
	"gravity-game-store/internal/mw"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gravity-game-store/docs"
)

func Build(
	authCtrl *api.AuthCtrl,
	authorCtrl *api.AuthorCtrl,
	gameCtrl *api.GameCtrl,
	customerCtrl *api.CustomerCtrl,
	authSvc *core.AuthSvc,
	log *logrus.Logger,
) *gin.Engine {
	r := gin.New()

	r.Use(mw.WithLogging(log))
	r.Use(mw.WithCORS())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
			auth.POST("/register", authCtrl.Register)
		}

		authors := v1.Group("/authors")
		{
			authors.GET("", authorCtrl.GetAll)
			authors.GET("/:id", authorCtrl.GetByID)
			authors.GET("/:id/games", authorCtrl.GetGames)
		}

		games := v1.Group("/games")
		{
			games.GET("", gameCtrl.GetAll)
			games.GET("/:id", gameCtrl.GetByID)
			games.GET("/:id/authors", gameCtrl.GetAuthors)
		}

		customers := v1.Group("/customers")
		{
			customers.GET("", customerCtrl.GetAll)
			customers.GET("/:id", customerCtrl.GetByID)
			customers.GET("/:id/orders", customerCtrl.GetOrders)
		}

		protected := v1.Group("")
		protected.Use(mw.Auth(authSvc, log))
		{
			protected.POST("/authors", authorCtrl.Create)
			protected.PUT("/authors/:id", authorCtrl.Update)
			protected.DELETE("/authors/:id", authorCtrl.Delete)

			protected.POST("/games", gameCtrl.Create)
			protected.PUT("/games/:id", gameCtrl.Update)
			protected.DELETE("/games/:id", gameCtrl.Delete)

			protected.POST("/customers", customerCtrl.Create)
			protected.PUT("/customers/:id", customerCtrl.Update)
			protected.DELETE("/customers/:id", customerCtrl.Delete)
		}
	}

	return r
}