package utils

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/suapapa/go_hangul/encoding/cp949"
)

var /* const */ stringsToRemove = []string{"  ", "\t", "\n"}

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

func StringSplit(data string, length int) string {
	b := []byte(data)
	idx := 0
	for i := 0; i < length; i++ {
		_, size := utf8.DecodeRune(b[idx:])
		idx += size
	}
	return data[:idx]
}

func TrimAll(value string) string {
	result := value
	for _, str := range stringsToRemove {
		result = strings.Replace(result, str, "", -1)
	}
	return result
}
