package render

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		log.Printf("json encode error: %s\n", err)
		http.Error(w, "failed to return json", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
