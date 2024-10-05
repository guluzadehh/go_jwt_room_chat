package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/http/handlers/auth/login"
	"github.com/guluzadehh/go_chat/internal/http/handlers/auth/logout"
	"github.com/guluzadehh/go_chat/internal/http/handlers/auth/refresh"
	"github.com/guluzadehh/go_chat/internal/http/handlers/auth/signup"
	"github.com/guluzadehh/go_chat/internal/http/handlers/chat"
	roomcreate "github.com/guluzadehh/go_chat/internal/http/handlers/room/create"
	roomdelete "github.com/guluzadehh/go_chat/internal/http/handlers/room/delete"
	roomlist "github.com/guluzadehh/go_chat/internal/http/handlers/room/list"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/authmdw"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/loggingmdw"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/storage/redis"
	"github.com/guluzadehh/go_chat/internal/storage/sqlite"
	"github.com/joho/godotenv"
)

const env_local = "local"
const env_dev = "dev"
const env_prod = "prod"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	// config
	config := config.MustLoad()

	// logger
	log := setupLogger(config.Env)
	log.Info("starting go-chat app", slog.String("env", config.Env))

	// storage
	sqliteStorage, err := sqlite.New(config.StoragePath)
	if err != nil {
		log.Error("failed to init sqlite", sl.Err(err))
		os.Exit(1)
	}

	redisStorage, err := redis.New(config)
	if err != nil {
		log.Error("failed to init redis", sl.Err(err))
		os.Exit(1)
	}

	// router
	router := mux.NewRouter()

	router.Use(requestmdw.AddRequestId)
	router.Use(loggingmdw.LogRequests(log))

	// Public routes
	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/login", login.New(log, config, sqliteStorage)).Methods("POST")
	api.Handle("/signup", signup.New(log, sqliteStorage)).Methods("POST")
	api.Handle("/refresh", refresh.New(log, config)).Methods("POST")

	// Protected routes
	apiAuth := api.NewRoute().Subrouter()
	apiAuth.Use(authmdw.Authorize(log, config, sqliteStorage))

	apiAuth.Handle("/logout", logout.New(log, config)).Methods("POST")
	apiAuth.Handle("/rooms", roomcreate.New(log, redisStorage)).Methods("POST")
	apiAuth.Handle("/rooms", roomlist.New(log, redisStorage, sqliteStorage)).Methods("GET")
	apiAuth.Handle("/rooms/{room_uuid}", roomdelete.New(log, redisStorage)).Methods("DELETE")

	apiAuth.Handle("/rooms/{room_uuid}/chat", chat.New(log, config, redisStorage)).Methods("GET")

	// run
	log.Info("starting server listener", slog.String("addr", config.HTTPServer.Address))
	if err := http.ListenAndServe(config.HTTPServer.Address, router); err != nil {
		log.Error("error while initializing the server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case env_local, env_dev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env_prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
