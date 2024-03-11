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

	var repData string

	for i, v := range input {
		if unicode.IsLetter(v) {
			switch {
			case !unicode.IsDigit(str[i+1]):
				repData = string(v)
			case unicode.IsDigit(str[i+1]) && !unicode.IsDigit(str[i+2]):
				val, err := strconv.Atoi(string(str[i+1]))
				if err != nil {
					return "", err
				}
				repData = strings.Repeat(string(v), val)
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
