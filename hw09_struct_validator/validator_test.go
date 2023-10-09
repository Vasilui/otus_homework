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

	Response struct {
		Code int    `validate:"in:200,404,500"`
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
			expectedErr: ValidationErrors{{Field: "App.Version", Err: joinErrors([]error{ErrInvalidMinLength, ErrNotContains})}},
		},
		{
			in:          App{Version: "hello, world"},
			expectedErr: ValidationErrors{{Field: "App.Version", Err: joinErrors([]error{ErrInvalidMaxLength, ErrNotContains})}},
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
