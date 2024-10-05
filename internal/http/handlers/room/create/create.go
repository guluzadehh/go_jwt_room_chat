package roomcreate

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/authmdw"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/types"
)

type RoomStorage interface {
	CreateRoom(name, password string, owner_id int64) (*models.Room, error)
}

func New(log *slog.Logger, roomStorage RoomStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.room.create.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		var body Request
		err := api.DecodeBody(log, w, r, &body)
		if err != nil {
			return
		}

		v := validator.New()
		if err := v.Struct(body); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Info("invalid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, api.ValidationError(validateErr))
			return
		}

		user := authmdw.User(r)

		room, err := roomStorage.CreateRoom(body.Name, body.Password, user.Id)
		if err != nil {
			log.Error("failed to create a room", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}
		log.Info("room has been created", slog.Any("room", room))

		render.JSON(w, http.StatusCreated, Response{
			Response: api.Ok(),
			Data: Data{
				Room: types.NewRoom(room, user),
			},
		})
	})
}
