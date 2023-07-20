package hw02unpackstring

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	if str == "" {
		return "", nil
	}
	if isDigit(rune(str[0])) {
		return "", ErrInvalidString
	}
	var prev rune
	for _, s := range str {
		if prev == 0 {
			prev = s
			continue
		}
		if isDigit(s) && isDigit(prev) {
			return "", ErrInvalidString
		}
		if isDigit(s) {
			repeat, err := strconv.Atoi(string(s))
			if err != nil {
				log.Println(err)
				return "", ErrInvalidString
			}
			result.WriteString(strings.Repeat(string(prev), repeat))
			prev = s
		} else {
			if isDigit(prev) {
				prev = s
				continue
			}
			result.WriteRune(prev)
			prev = s
		}
	}
	if !isDigit(prev) {
		result.WriteRune(prev)
	}
	return result.String(), nil
}

func isDigit(symbol rune) bool {
	if symbol >= 48 && symbol <= 57 {
		return true
	}
	return false
}
