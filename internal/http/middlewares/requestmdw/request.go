package requestmdw

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIdKey string

var requestId requestIdKey = "requestId"

func GetReqId(r *http.Request) string {
	return r.Context().Value(requestId).(string)
}

func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestId, uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
