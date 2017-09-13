package serialized

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func loadJSON(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, b); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func assertEqualJSON(t *testing.T, b1, b2 []byte) {
	b1 = bytes.TrimSpace(b1)
	b2 = bytes.TrimSpace(b2)

	if !bytes.Equal(b1, b2) {
		var buf1 bytes.Buffer
		json.Indent(&buf1, b1, "", "\t")

		var buf2 bytes.Buffer
		json.Indent(&buf2, b2, "", "\t")

		t.Errorf("unexpected request body =\n%s\n\nwant =\n%s", buf1.String(), buf2.String())
	}
}
