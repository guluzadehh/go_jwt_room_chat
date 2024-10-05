package render

import (
	"encoding/json"
	"io"
)

func DecodeJSON(r io.Reader, dst interface{}) error {
	defer io.Copy(io.Discard, r)
	return json.NewDecoder(r).Decode(dst)
}
