package chat

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/authmdw"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/roomchat"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
)

type RoomStorage interface {
	RoomByUuid(uuid string) (*models.Room, error)
}

func New(log *slog.Logger, config *config.Config, roomStorage RoomStorage) http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	hub := roomchat.NewHub(config)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.chat.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		roomUuid := mux.Vars(r)["room_uuid"]

		room, err := roomStorage.RoomByUuid(roomUuid)
		if err != nil {
			if errors.Is(err, storage.RoomNotFound) {
				log.Info("room doesn't exist", slog.String("uuid", roomUuid))
				render.JSON(w, http.StatusNotFound, api.Err("room is not found"))
				return
			}

			log.Error("failed to get the room", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("failed to upgrade connection", sl.Err(err))
			return
		}

		user := authmdw.User(r)

		if room.IsPrivate() {
			var msg struct {
				Password string `json:"password"`
			}

			msgType, rcv, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Warn("error while reading password message", sl.Err(err))
				}
				conn.Close()
				return
			}

			failMsg := "failed to grant access"

			if msgType == websocket.CloseMessage {
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				conn.Close()
				return
			}

			if err := json.Unmarshal(rcv, &msg); err != nil {
				log.Info("failed to read password", slog.String("rcv", string(rcv)), sl.Err(err))
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, failMsg))
				conn.Close()
				return
			}

			if msg.Password != room.Password {
				log.Warn("invalid password for room", slog.Any("room", room), sl.User(user), slog.String("password", msg.Password))
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "wrong password"))
				conn.Close()
				return
			}

			log.Info("gained access to the room", sl.User(user), slog.Any("room", room))
		}

		chatRoom := hub.GetOrCreateRoom(room)
		member, err := chatRoom.NewMember(conn, user)
		if errors.Is(err, roomchat.RoomIsFull) {
			log.Info("full room join attempt", sl.User(user), slog.Any("room", room))

			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "room is full"))
			conn.Close()

			return
		}
		log.Info("member is created", sl.User(user), slog.Any("room", room))

		go member.ReadPump()
	})
}
