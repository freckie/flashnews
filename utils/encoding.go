package utils

import (
	"bytes"
	"strings"

	"github.com/suapapa/go_hangul/encoding/cp949"
)

func ReadCP949(data string) (string, error) {
	br := bytes.NewReader([]byte(data))
	r, err := cp949.NewReader(br)
	if err != nil {
		return "", err
	}

	b := make([]byte, 10*1024)

	c, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b[:c])), nil
}
