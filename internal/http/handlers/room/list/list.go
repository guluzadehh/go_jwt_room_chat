package roomlist

import (
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/types"
)

type RoomStorage interface {
	Rooms() ([]*models.Room, error)
}

type UserStorage interface {
	UsersWithIds(ids []int64) (map[int64]*models.User, error)
}

func New(log *slog.Logger, roomStorage RoomStorage, userStorage UserStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.room.list.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		rooms, err := roomStorage.Rooms()
		if err != nil {
			log.Error("failed to get the list of rooms", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		owner_ids := make([]int64, 0)
		for _, room := range rooms {
			owner_ids = append(owner_ids, room.OwnerId)
		}

		owners, err := userStorage.UsersWithIds(owner_ids)
		if err != nil {
			log.Error("failed to get the owners of rooms", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		roomsResponse := make([]*types.RoomView, 0)
		for _, room := range rooms {
			roomsResponse = append(roomsResponse, types.NewRoom(room, owners[room.OwnerId]))
		}

		render.JSON(w, http.StatusOK, Response{
			Response: api.Ok(),
			Data: Data{
				RoomsResponse: roomsResponse,
				Size:          len(rooms),
			},
		})
	})
}
