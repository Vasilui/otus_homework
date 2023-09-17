package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrInternal      = errors.New("internal error")
)

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
			err := isCorrectSymbol(s.el, &backslash, &elByOutString, &sBuilder, &count)
			if err != nil {
				return "", err
			}
		}

		if s.IsNumber() {
			err := isNumberSymbol(s.el, &backslash, &elByOutString, &sBuilder, &count)
			if err != nil {
				return "", err
			}
		}

		if s.IsBackslash() {
			err := isBackslashSymbol(&backslash, &elByOutString, &sBuilder)
			if err != nil {
				return "", err
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

func isCorrectSymbol(s uint8, backslash *bool, elByOutString *string, sBuilder *strings.Builder, count *int) error {
	if *backslash {
		return ErrInvalidString
	}

	if *elByOutString != "" {
		_, err := fmt.Fprintf(sBuilder, "%s", *elByOutString)
		if err != nil {
			return ErrInternal
		}
	}

	*count = 0
	*elByOutString = string(s)

	return nil
}

func isNumberSymbol(s uint8, backslash *bool, elByOutString *string, sBuilder *strings.Builder, count *int) error {
	if *backslash {
		*elByOutString = string(s)
		*backslash = false
		return nil
	}

	*count = int(s - 48)

	if *elByOutString == "" || *count == -1 {
		return ErrInvalidString
	}

	_, err := fmt.Fprintf(sBuilder, "%s", strings.Repeat(*elByOutString, *count))
	if err != nil {
		return ErrInternal
	}

	*count = -1
	*elByOutString = ""

	return nil
}

func isBackslashSymbol(backslash *bool, elByOutString *string, sBuilder *strings.Builder) error {
	if *elByOutString != "" {
		if _, err := fmt.Fprintf(sBuilder, "%s", *elByOutString); err != nil {
			return ErrInternal
		}
		*elByOutString = ""
	}
	if *backslash {
		*elByOutString = "\\"
		*backslash = false
		return nil
	}

	*backslash = true

	return nil
}

func (s Symbol) IsNumber() bool {
	return s.el > 47 && s.el < 58
}

func (s Symbol) IsCorrect() bool {
	return s.el < 33 || s.el > 64 && s.el < 91 || s.el > 96 && s.el < 123
}

func (s Symbol) IsBackslash() bool {
	return s.el == 92
}
