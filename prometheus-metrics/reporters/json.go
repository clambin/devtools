package reporters

import (
	"encoding/json"
	"io"
)

func NewJSONEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}
