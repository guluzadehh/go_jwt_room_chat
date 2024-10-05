package refresh

import (
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/auth"
	"github.com/guluzadehh/go_chat/internal/lib/jwt"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
)

func New(log *slog.Logger, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.refresh.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		cookie, err := r.Cookie(config.JWT.Refresh.CookieName)
		if err == http.ErrNoCookie {
			log.Info("refresh cookie doesn't exist", slog.String("refresh_cookie_name", config.JWT.Refresh.CookieName))
			render.JSON(w, http.StatusUnauthorized, refreshInvalidResponse())
			return
		}

		refreshStr, err := auth.Decrypt(cookie.Value, []byte(config.JWT.Refresh.EncryptSecretKey))
		if err != nil {
			log.Error("failed to decrypt token", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("refresh token decrypted")

		refresh, err := jwt.Verify(refreshStr, config)
		if err != nil {
			log.Info("refresh token is invalid", slog.String("invalid_refresh_token", refreshStr), sl.Err(err))
			render.JSON(w, http.StatusUnauthorized, refreshInvalidResponse())
			return
		}
		log.Info("refresh token is verified")

		username, err := refresh.Claims.GetSubject()
		if err != nil {
			log.Error("error while getting the subject from refresh token", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		access, err := jwt.AccessToken(username, config)
		if err != nil {
			log.Error("can't create jwt access token", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("access token have been created", slog.String("username", username))

		render.JSON(w, http.StatusOK, Response{
			Response: api.Ok(),
			Data: Data{
				Token: access,
			},
		})
	})
}

func refreshInvalidResponse() api.Response {
	return api.Err("refresh token is invalid")
}
