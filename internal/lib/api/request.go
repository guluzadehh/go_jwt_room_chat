package api

import (
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
)

func DecodeBody(log *slog.Logger, w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := render.DecodeJSON(r.Body, v)
	if err != nil {
		log.Error("can't decode from json", sl.Err(err))
		render.JSON(w, http.StatusBadRequest, Err("failed to parse request body"))
		return err
	}

	log.Info("request body decoded", slog.Any("body", v))

	return nil
}
