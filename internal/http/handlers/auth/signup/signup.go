package signup

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/lib/auth"
	"github.com/guluzadehh/go_chat/internal/lib/render"
	"github.com/guluzadehh/go_chat/internal/lib/sl"
	"github.com/guluzadehh/go_chat/internal/lib/validators"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
	"github.com/guluzadehh/go_chat/internal/types"
)

type SignupStorage interface {
	CreateUser(username, password string) (*models.User, error)
}

func New(log *slog.Logger, signupStorage SignupStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.signup.New"

		log := sl.ForHandler(log, op, requestmdw.GetReqId(r))

		var body Request
		err := api.DecodeBody(log, w, r, &body)
		if err != nil {
			return
		}

		v := validator.New()
		v.RegisterValidation("passwordpattern", validators.PasswordPatternValidator)

		if err := v.Struct(body); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Info("invalid request", sl.Err(err))
			render.JSON(w, http.StatusBadRequest, api.ValidationError(validateErr))
			return
		}

		hashedPassword, err := auth.HashPassword(body.Password)
		if err != nil {
			log.Error("can't hash password", slog.String("password", body.Password), sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		user, err := signupStorage.CreateUser(body.Username, hashedPassword)
		if errors.Is(err, storage.UsernameExists) {
			log.Info(err.Error(), slog.String("username", body.Username))
			render.JSON(w, http.StatusConflict,
				api.ErrD("username exists", []api.ErrDetail{
					{
						Field:   "username",
						Message: "username is already taken",
					},
				}),
			)
			return
		}
		if err != nil {
			log.Info("failed to create user", sl.Err(err))
			render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		log.Info("user has been created", sl.User(user))

		render.JSON(w, http.StatusCreated, Response{
			Response: api.Ok(),
			Data: Data{
				User: types.NewUser(user),
			},
		})
	})
}
