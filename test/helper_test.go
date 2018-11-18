package test

import (
	"bytes"
	"encoding/json"
)

func equalsJSON(got, want []byte) bool {
	var gotbuf bytes.Buffer
	var wantbuf bytes.Buffer

	json.Compact(&gotbuf, got)
	json.Compact(&wantbuf, want)

	return bytes.Equal(gotbuf.Bytes(), wantbuf.Bytes())
}
