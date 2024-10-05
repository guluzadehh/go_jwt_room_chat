package logout

import (
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
)

func New(log *slog.Logger, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.logout.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		http.SetCookie(w, &http.Cookie{
			Name:     config.JWT.Refresh.CookieName,
			Value:    "",
			SameSite: http.SameSiteNoneMode,
			Path:     "/api/refresh",
			HttpOnly: true,
			MaxAge:   -1,
		})
		log.Info("refresh cookie has been deleted", slog.String("refresh_cookie_name", config.JWT.Refresh.CookieName))

		render.JSON(w, http.StatusOK, api.Ok())
	})
}
