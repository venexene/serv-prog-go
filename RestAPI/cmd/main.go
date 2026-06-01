package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gravity-game-store/internal/api"
	"gravity-game-store/internal/conf"
	"gravity-game-store/internal/core"
	"gravity-game-store/internal/db"
	"gravity-game-store/internal/routes"
	"gravity-game-store/internal/store"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := conf.LoadCfg()
	if err != nil {
		panic("config: " + err.Error())
	}

	log := logrus.New()
	log.SetLevel(cfg.LogLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05",
	})

	log.WithField("port", cfg.Port).Info("Gravity Game Store")

	d, err := db.Connect(cfg.DBPath, cfg.LogLevel)
	if err != nil {
		log.Fatal("db connect: ", err)
	}
	log.WithField("path", cfg.DBPath).Info("db connected")

	if err := db.Migrate(d); err != nil {
		log.Fatal("migrate: ", err)
	}
	log.Info("migration done")

	if err := db.Seed(d, log); err != nil {
		log.Fatal("seed: ", err)
	}

	authorStore := store.NewAuthorStore(d)
	gameStore := store.NewGameStore(d)
	customerStore := store.NewCustomerStore(d)
	userStore := store.NewUserStore(d)

	authSvc := core.NewAuthSvc(userStore, cfg)
	authorSvc := core.NewAuthorSvc(authorStore)
	gameSvc := core.NewGameSvc(gameStore)
	customerSvc := core.NewCustomerSvc(customerStore)

	authCtrl := api.NewAuthCtrl(authSvc, log)
	authorCtrl := api.NewAuthorCtrl(authorSvc, log)
	gameCtrl := api.NewGameCtrl(gameSvc, log)
	customerCtrl := api.NewCustomerCtrl(customerSvc, log)

	r := routes.Build(authCtrl, authorCtrl, gameCtrl, customerCtrl, authSvc, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("listening on :%s", cfg.Port)
		log.Infof("swagger: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("serve: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutdown: ", err)
	}
	log.Info("bye")
}