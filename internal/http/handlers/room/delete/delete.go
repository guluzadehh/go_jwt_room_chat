package roomdelete

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/authmdw"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
)

type RoomStorage interface {
	RoomByUuid(uuid string) (*models.Room, error)
	DeleteRoom(uuid string) error
}

func New(log *slog.Logger, roomStorage RoomStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.room.delete.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		roomUuid := mux.Vars(r)["room_uuid"]

		room, err := roomStorage.RoomByUuid(roomUuid)
		if err == storage.RoomNotFound {
			log.Info("couldn't find the room to delete")
			render.JSON(w, http.StatusNotFound, api.Err("room doesn't exist"))
			return
		}
		if err != nil {
			log.Error("failed to get a room", slog.String("room_uuid", roomUuid), sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		user := authmdw.User(r)
		if room.OwnerId != user.Id {
			log.Info("unauthorized access to delete the room", sl.User(user), slog.Any("room", room))
			render.JSON(w, http.StatusForbidden, api.Err("you are not allowed"))
			return
		}

		if err := roomStorage.DeleteRoom(roomUuid); err != nil {
			log.Error("failed to delete the room", slog.Any("room", room), sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		render.JSON(w, http.StatusNoContent, api.Ok())
	})
}
