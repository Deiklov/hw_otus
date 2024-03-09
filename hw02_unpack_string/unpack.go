package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	// for simplify if cases
	str := []rune(input + " ")

	// only numbers case
	if unicode.IsDigit(str[0]) {
		return "", ErrInvalidString
	}

	var result strings.Builder

	atoi := func(str string) int {
		val, _ := strconv.Atoi(str)
		return val
	}

	var repData string

	for i, v := range input {
		if unicode.IsLetter(v) {
			switch {
			case !unicode.IsDigit(str[i+1]):
				repData = strings.Repeat(string(v), 1)
			case unicode.IsDigit(str[i+1]) && !unicode.IsDigit(str[i+2]):
				repData = strings.Repeat(string(v), atoi(string(str[i+1])))
			default:
				return "", ErrInvalidString
			}
			result.WriteString(repData)
		}

		if unicode.IsDigit(v) {
			continue
		}
	}

	return result.String(), nil
}
