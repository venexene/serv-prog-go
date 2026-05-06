package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	internal "github.com/venexene/serv-prog-go/greenswamp/internal"
	"github.com/venexene/serv-prog-go/greenswamp/posts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := internal.LoadConfig()
	logger := log.New(os.Stdout, "[greenswamp] ", log.LstdFlags)

	db, err := gorm.Open(sqlite.Open("data/greenswamp.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	if err := posts.AutoMigrate(db); err != nil {
		logger.Fatal(err)
	}

	app := internal.CreateApp(cfg, logger)

	mux := http.NewServeMux()
	app.Routes(mux)

	posts.RegisterRoutes(mux, db, logger, posts.Config{
		BasePath:     "/posts",
		TemplatesDir: "templates/posts",
	})

	handler := app.WithMiddleware(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go shutdownHandler(server, logger, 10*time.Second)

	logger.Println("Server started on", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func shutdownHandler(srv *http.Server, logger *log.Logger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Println("Shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Println("Shutdown error:", err)
	}
}