package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"min:4|max:6|in:foo,hello!"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Phone struct {
		Number string `validate:"regexp:\\d{11}"`
	}

	Slices struct {
		TestString []string `validate:"len:5"`
		TestInt    []int    `validate:"max:10"`
	}

	Nested struct {
		NameApp string
		App     App `validate:"nested"`
	}

	NoExported struct {
		name string `validate:"len5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500|min:100|max:600"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "hello"},
			expectedErr: ValidationErrors{{Field: "App.Version", Err: ErrNotContains}},
		},
		{
			in:          App{Version: "hel"},
			expectedErr: ValidationErrors{{Field: "App.Version", Err: joinErrors([]error{ErrInvalidMin, ErrNotContains})}},
		},
		{
			in:          App{Version: "hello, world"},
			expectedErr: ValidationErrors{{Field: "App.Version", Err: joinErrors([]error{ErrInvalidMax, ErrNotContains})}},
		},
		{
			in:          Phone{Number: "9876543210"},
			expectedErr: ValidationErrors{{Field: "Phone.Number", Err: ErrNoMatched}},
		},
		{
			in:          Phone{Number: "98765432100"},
			expectedErr: nil,
		},
		{
			in:          Nested{NameApp: "application", App: App{Version: "hello"}},
			expectedErr: ValidationErrors{{Field: "Nested.App.Version", Err: ErrNotContains}},
		},
		{
			in:          NoExported{name: "Vasilii"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 201},
			expectedErr: ValidationErrors{{Field: "Response.Code", Err: ErrNotContains}},
		},
		{
			in:          Response{Code: 90},
			expectedErr: ValidationErrors{{Field: "Response.Code", Err: joinErrors([]error{ErrNotContains, ErrInvalidMin})}},
		},
		{
			in:          Response{Code: 999},
			expectedErr: ValidationErrors{{Field: "Response.Code", Err: joinErrors([]error{ErrNotContains, ErrInvalidMax})}},
		},
		{
			in:          Slices{TestString: []string{"Vasilii"}},
			expectedErr: ValidationErrors{{Field: "Slices.TestString.[0]", Err: ErrInvalidLength}},
		},
		{
			in:          Slices{TestInt: []int{11}},
			expectedErr: ValidationErrors{{Field: "Slices.TestInt.[0]", Err: ErrInvalidMax}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)

			_ = tt
		})
	}
}
