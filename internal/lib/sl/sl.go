package sl

import (
	"log/slog"

	"github.com/guluzadehh/go_chat/internal/models"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func User(u *models.User) slog.Attr {
	return slog.Attr{
		Key: "user",
		Value: slog.AnyValue(struct {
			Id       int64
			Username string
		}{
			Id:       u.Id,
			Username: u.Username,
		}),
	}
}

func ForHandler(log *slog.Logger, op, requestId string) *slog.Logger {
	return log.With(
		slog.String("op", op),
		slog.String("request_id", requestId),
	)
}
