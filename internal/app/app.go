package app

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"songs/internal/api"
	songHandler "songs/internal/api/songs"
	"songs/internal/usecase"
	"songs/internal/usecase/repo/postgres"
	"songs/pkg/config"
	"songs/pkg/router"
	"syscall"
	"time"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
)

func Run() {
	//TODO: change type of date

	cfg := config.ConfigLoad()

	logger := setupLogger(cfg.LogLevel)

	logger.Info("starting server")
	logger.Debug("debug messages are enabled")

	db, err := sql.Open("postgres", cfg.StoragePath)
	if err != nil {
		logger.Error("failed to connect db: ", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("failed to check db connection")
	} else {
		logger.Info("successful connect to db")
	}

	songs := postgres.NewSongRepo(db)

	songUseCase := usecase.NewSongUseCase(*songs)

	sHandler := songHandler.NewSongsHandler(*songUseCase, logger)

	r := router.New(logger)

	r.Add(
		router.NewGroup("/song",
			router.POST("/add", sHandler.Add),
			router.POST("/delete/", sHandler.Delete),
			router.POST("/update/", sHandler.Update),
			router.GET("/lyric/", sHandler.GetLyric),
		).SetErrHandler(api.ErrHandler),
	)

	r.Add(router.POST("/songs", sHandler.GetSongs).SetErrHandler(api.ErrHandler))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		srv.ListenAndServe()
	}()

	logger.Info("server started")
	logger.Debug("listen on", "addr", srv.Addr)

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server: ", err)

		return
	}

	logger.Info("server stopped")
}

func setupLogger(level string) *slog.Logger {
	var log *slog.Logger

	switch level {
	case LevelInfo:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case LevelDebug:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
