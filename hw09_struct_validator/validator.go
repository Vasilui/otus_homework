package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidInputData        = errors.New("invalid input data")
	ErrReflectValueIsNotStruct = errors.New("reflect value is not struct")
	ErrInvalidLength           = errors.New("invalid length")
	ErrInvalidMax              = errors.New("invalid max")
	ErrInvalidMin              = errors.New("invalid min")
	ErrNotContains             = errors.New("not contains")
	ErrInvalidValidator        = errors.New("invalid validator")
	ErrNoMatched               = errors.New("failed matched")
	ErrInvalidBuildError       = errors.New("failed build error string")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var res strings.Builder

	for _, err := range v {
		if _, e := fmt.Fprintf(&res, "%s is incorrect: %s\n", err.Field, err.Err.Error()); e != nil {
			return ErrInvalidBuildError.Error()
		}
	}

	return res.String()
}

func Validate(v interface{}) error {
	if v == nil {
		return ErrInvalidInputData
	}

	val := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	return validateStruct(val, rt, rt.Name())
}

func validateStruct(rv reflect.Value, rt reflect.Type, baseName string) error {
	if rv.Kind() != reflect.Struct {
		return ErrReflectValueIsNotStruct
	}

	res := ValidationErrors{}

	for i := 0; i < rv.NumField(); i++ {
		err := checkStructField(rv.Field(i), rt.Field(i), baseName)
		if err != nil {
			t := ValidationErrors{}
			switch {
			case errors.Is(err, ErrInvalidValidator):
				return err
			case errors.As(err, &t):
				res = append(res, t...)
			default:
				res = append(res, ValidationError{
					Field: fmt.Sprintf("%s.%s", baseName, rt.Field(i).Name),
					Err:   err,
				})
			}
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func checkStructField(val reflect.Value, sf reflect.StructField, baseName string) error {
	validateTags, ok := sf.Tag.Lookup("validate")
	if !sf.IsExported() || !ok {
		return nil
	}

	switch val.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(val.String(), strings.Split(validateTags, "|"))
	case reflect.Int:
		return validateInt(val.Int(), strings.Split(validateTags, "|"))
	case reflect.Slice:
		return checkSliceField(val, strings.Split(validateTags, "|"), fmt.Sprintf("%s.%s", baseName, sf.Name))
	case reflect.Struct:
		if validateTags != "nested" {
			break
		}
		res := ValidationErrors{}
		err := validateStruct(val, sf.Type, strings.Join([]string{baseName, sf.Name}, "."))
		if err != nil {
			t := ValidationErrors{}
			switch {
			case errors.Is(err, ErrInvalidValidator):
				return err
			case errors.As(err, &t):
				res = append(res, t...)
			default:
				res = append(res, ValidationError{
					Field: fmt.Sprintf("%s.%s", baseName, sf.Name),
					Err:   err,
				})
			}
		}
		if len(res) == 0 {
			return nil
		}
		return res
	default:
		return nil
	}

	return nil
}

func checkSliceField(val reflect.Value, validators []string, baseName string) error {
	res := ValidationErrors{}

	for i := 0; i < val.Len(); i++ {
		var err error
		switch val.Index(i).Kind() { //nolint:exhaustive
		case reflect.String:
			err = validateString(val.Index(i).String(), validators)
		case reflect.Int:
			err = validateInt(val.Index(i).Int(), validators)
		default:
			return nil
		}

		if err != nil {
			if errors.Is(err, ErrInvalidValidator) {
				return err
			}
			res = append(res, ValidationError{Field: fmt.Sprintf("%s.[%d]", baseName, i), Err: err})
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func validateString(val string, validators []string) error { //nolint:gocognit
	var res []error
	for _, v := range validators {
		item := strings.Split(v, ":")
		if len(item) < 2 {
			return ErrInvalidValidator
		}

		if item[0] == "len" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				return ErrInvalidValidator
			} else if len(val) != length {
				res = append(res, ErrInvalidLength)
			}
		}

		if item[0] == "max" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				return ErrInvalidValidator
			} else if len(val) > length {
				res = append(res, ErrInvalidMax)
			}
		}

		if item[0] == "min" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				return ErrInvalidValidator
			} else if len(val) < length {
				res = append(res, ErrInvalidMin)
			}
		}

		if item[0] == "in" {
			data := strings.Split(item[1], ",")
			contains := false
			for _, word := range data {
				if word == val {
					contains = true
					break
				}
			}
			if !contains {
				res = append(res, ErrNotContains)
			}
		}

		if item[0] == "regexp" {
			re, err := regexp.Compile(item[1])
			if err != nil {
				return ErrInvalidValidator
			}
			if !re.MatchString(val) {
				res = append(res, ErrNoMatched)
			}
		}
	}

	return joinErrors(res)
}

func validateInt(val int64, validators []string) error { //nolint:gocognit
	var res []error
	for _, v := range validators {
		item := strings.Split(v, ":")
		if len(item) < 2 {
			return ErrInvalidValidator
		}

		if item[0] == "len" {
			i := strconv.FormatInt(val, 10)
			length, err := strconv.Atoi(item[1])
			if err != nil {
				return ErrInvalidValidator
			} else if len(i) != length {
				res = append(res, ErrInvalidLength)
			}
		}

		if item[0] == "max" {
			i, err := strconv.ParseInt(item[1], 10, 64)
			if err != nil {
				return ErrInvalidValidator
			} else if val > i {
				res = append(res, ErrInvalidMax)
			}
		}

		if item[0] == "min" {
			i, err := strconv.ParseInt(item[1], 10, 64)
			if err != nil {
				return ErrInvalidValidator
			} else if val < i {
				res = append(res, ErrInvalidMin)
			}
		}

		if item[0] == "in" {
			data := strings.Split(item[1], ",")
			contains := false
			for _, s := range data {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return ErrInvalidValidator
				}
				if i == val {
					contains = true
					break
				}
			}
			if !contains {
				res = append(res, ErrNotContains)
			}
		}
	}

	return joinErrors(res)
}

func joinErrors(e []error) error {
	if len(e) == 0 {
		return nil
	}

	data := make([]string, 0)
	for _, err := range e {
		data = append(data, err.Error())
	}
	resultError := strings.Join(data, "; ")

	return errors.New(resultError)
}
