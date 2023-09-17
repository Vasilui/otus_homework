package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")
var ErrInternal = errors.New("internal error")

type Symbol struct {
	el uint8
}

func Unpack(data string) (string, error) {
	sBuilder := strings.Builder{}
	elByOutString := ""
	count := -1
	backslash := false

	for i := range data {
		s := Symbol{el: data[i]}
		if s.IsCorrect() {
			if backslash {
				return "", ErrInvalidString
			}
			if elByOutString != "" {
				_, err := fmt.Fprintf(&sBuilder, "%s", elByOutString)
				if err != nil {
					return "", ErrInternal
				}
			}

			count = 0
			elByOutString = string(s.el)
		}

		if s.IsNumber() {
			if backslash {
				elByOutString = string(s.el)
				backslash = false
				continue
			} else {
				count = int(s.el - 48)
			}

			if elByOutString == "" || count == -1 {
				return "", ErrInvalidString
			}

			_, err := fmt.Fprintf(&sBuilder, "%s", strings.Repeat(elByOutString, count))
			if err != nil {
				return "", ErrInternal
			}

			count = -1
			elByOutString = ""
		}

		if s.IsBackslash() {
			if elByOutString != "" {
				_, err := fmt.Fprintf(&sBuilder, "%s", elByOutString)
				if err != nil {
					return "", ErrInternal
				}
				elByOutString = ""
			}
			if backslash {
				elByOutString = "\\"
				backslash = false
			} else {
				backslash = true
			}
		}
	}

	if elByOutString != "" {
		_, err := fmt.Fprintf(&sBuilder, "%s", elByOutString)
		if err != nil {
			return "", ErrInternal
		}
	}

	return sBuilder.String(), nil
}

func (s Symbol) IsNumber() bool {
	if s.el > 47 && s.el < 58 {
		return true
	}

	return false
}

func (s Symbol) IsCorrect() bool {
	if s.el < 33 || s.el > 64 && s.el < 91 || s.el > 96 && s.el < 123 {
		return true
	}

	return false
}

func (s Symbol) IsBackslash() bool {

	return s.el == 92
}
