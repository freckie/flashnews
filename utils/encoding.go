package utils

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

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

func ReadCP949Bytes(data []byte) (string, error) {
	br := bytes.NewReader(data)
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

func ReadISO88591(data string) (string, error) {
	iso8859_1_buf := []byte(data)

	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}

	dec := charmap.ISO8859_1.NewEncoder()
	out, err := dec.Bytes([]byte(string(buf)))
	if err != nil {
		return "", err
	}

	result, err := ReadCP949Bytes(out)
	if err != nil {
		return "", err
	}

	return string(result), nil
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
	result := strings.TrimSpace(value)
	for _, str := range stringsToRemove {
		result = strings.Replace(result, str, "", -1)
	}
	return result
}
