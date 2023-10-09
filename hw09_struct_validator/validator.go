package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidInputData        = errors.New("invalid input data")
	ErrReflectValueIsNotStruct = errors.New("reflect value is not struct")
	ErrInvalidLength           = errors.New("invalid length")
	ErrInvalidMaxLength        = errors.New("invalid max length")
	ErrInvalidMinLength        = errors.New("invalid min length")
	ErrNotContains             = errors.New("not contains")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for _, err := range v {
		res += fmt.Sprintf("%s is incorrect: %s\n", err.Field, err.Err.Error())
	}

	return res
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
		if rv.Field(i).Kind() == reflect.Struct {
			err := validateStruct(rv.Field(i), rt.Field(i).Type, strings.Join([]string{baseName, rt.Field(i).Name}, "."))
			if err != nil {
				t := ValidationErrors{}
				if errors.As(err, &t) {
					res = append(res, t...)
				} else {
					res = append(res, ValidationError{
						Field: fmt.Sprintf("%s.%s", baseName, rt.Field(i).Name),
						Err:   err,
					})
				}
			}
		}

		if rv.Field(i).Kind() == reflect.Slice {
			_ = checkSliceField(rv.Field(i), rt.Field(i))
		}

		err := checkStructField(rv.Field(i), rt.Field(i))
		if err != nil {
			res = append(res, ValidationError{
				Field: fmt.Sprintf("%s.%s", baseName, rt.Field(i).Name),
				Err:   err,
			})
		}
	}

	return res
}

func checkStructField(val reflect.Value, sf reflect.StructField) error {
	validateTags := sf.Tag.Get("validate")
	if validateTags == "" {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		return joinErrors(validateString(val.String(), strings.Split(validateTags, "|")))
	}

	return nil
}

func checkSliceField(sv reflect.Value, sf reflect.StructField) error {
	fmt.Println("slice: ", string(sv.Bytes()), sf.Tag.Get("validate"))
	//for i := 0; i < sf.NumField(); i++ {
	//	fmt.Println(sf.Field(i).Type())
	//}

	return nil
}

func validateString(val string, validators []string) []error {
	var res []error
	for _, v := range validators {
		item := strings.Split(v, ":")

		if item[0] == "len" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				res = append(res, err)
			} else if len(val) != length {
				res = append(res, ErrInvalidLength)
			}
		}

		if item[0] == "max" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				res = append(res, err)
			} else if len(val) > length {
				res = append(res, ErrInvalidMaxLength)
			}
		}

		if item[0] == "min" {
			length, err := strconv.Atoi(item[1])
			if err != nil {
				res = append(res, err)
			} else if len(val) < length {
				res = append(res, ErrInvalidMinLength)
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
	}

	return res
}

func joinErrors(e []error) error {
	if len(e) == 0 {
		return nil
	}

	var data []string
	for _, err := range e {
		data = append(data, err.Error())
	}
	resultError := strings.Join(data, "; ")

	return errors.New(resultError)
}
