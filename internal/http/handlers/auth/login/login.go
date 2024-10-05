package login

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/auth"
	"github.com/guluzadehh/go_chat/internal/lib/jwt"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
)

type LoginStorage interface {
	UserByUsername(username string) (*models.User, error)
}

func New(log *slog.Logger, config *config.Config, loginStorage LoginStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		var req Request
		err := api.DecodeBody(log, w, r, &req)
		if err != nil {
			return
		}

		user, err := loginStorage.UserByUsername(req.Username)
		if errors.Is(err, storage.UserNotFound) {
			log.Info(err.Error(), slog.String("username", req.Username))
			render.JSON(w, http.StatusUnauthorized, api.Err("invalid credentials."))
			return
		}
		if err != nil {
			log.Error("failed to get user by username from storage", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		if !auth.CheckPasswordHash(user.Password, req.Password) {
			log.Info("invalid credentials", slog.String("username", req.Username), slog.String("password", req.Password))
			render.JSON(w, http.StatusUnauthorized, api.Err("invalid credentials"))
			return
		}

		access, err := jwt.AccessToken(user.Username, config)
		if err != nil {
			log.Error("can't create jwt access token", slog.String("username", user.Username), sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("access token have been created", slog.String("username", user.Username))

		refresh, err := jwt.RefreshToken(user.Username, config)
		if err != nil {
			log.Error("can't create jwt refresh token", slog.String("username", user.Username), sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("refresh token have been created", slog.String("username", user.Username))

		encoded, err := auth.Encrypt(refresh, []byte(config.JWT.Refresh.EncryptSecretKey))
		if err != nil {
			log.Error("failed to encrypt token", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("encrypted the refresh token")

		http.SetCookie(w, &http.Cookie{
			Name:     config.JWT.Refresh.CookieName,
			Value:    encoded,
			SameSite: http.SameSiteNoneMode,
			Path:     "/api/refresh",
			HttpOnly: true,
			MaxAge:   int(config.JWT.Refresh.Expire.Seconds()),
		})
		log.Info("refresh cookie has been set", slog.String("refresh_cookie_name", config.JWT.Refresh.CookieName))

		render.JSON(w, http.StatusOK, Response{
			Response: api.Ok(),
			Data: Data{
				Token: access,
			},
		})
	})
}
