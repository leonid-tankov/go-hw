package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Is  int = iota // ошибки, не связанные с валидацией значений в структуре
	As             // ошибки, связанные с валидацией значений в структуре
	Err            // ошибки при валидации самих тегов
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
		Phones []string        `validate:"len:11|regexp:^\\d+$"`
		Codes  []int           `validate:"in:400,211,301"`
		Some   []int           `validate:"min:5|max:10"`
		Region int             `validate:"in:12,15,28"`
		meta   json.RawMessage //nolint:unused
	}

	InvalidTag struct {
		Name string `validate:"len:36:21|min:18"`
	}

	InvalidTagLen struct {
		Name string `validate:"len:in"`
	}

	InvalidTagRegexp struct {
		Name string `validate:"regexp:["`
	}

	InvalidTagMin struct {
		Age int `validate:"min:min"`
	}

	InvalidTagMax struct {
		Age int `validate:"max:max"`
	}

	InvalidTagInInt struct {
		Age int `validate:"in:in1,in2"`
	}

	App struct {
		Version string `validate:"len:5"`
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
		expectedErr []error
		name        string
		testCase    int
	}{
		{
			User{
				ID:     "7673c766-99a3-45bf-a7fe-d9df8cc14259",
				Age:    20,
				Email:  "df@ds.ewi",
				Phones: []string{"21313329192", "21510932919"},
				Codes:  []int{400, 211},
				Role:   "admin",
				Some:   []int{6, 7},
				Region: 12,
			},
			[]error{nil},
			"valid struct",
			Is,
		},
		{
			"string",
			[]error{ErrNotStruct},
			"invalid input",
			Is,
		},
		{
			InvalidTag{},
			[]error{ErrInvalidRule},
			"invalid rule",
			Is,
		},
		{
			App{
				Version: "1.0",
			},
			[]error{ErrInvalidLen},
			"invalid len",
			As,
		},
		{
			User{
				Age:    51,
				Email:  "test.ru",
				Role:   "user",
				Region: 30,
			},
			[]error{ErrNotMatchRegex, ErrNotInSlice, ErrInvalidMin, ErrInvalidMax},
			"list of errors for int and string",
			As,
		},
		{
			User{
				Codes: []int{200, 211, 301},
			},
			[]error{ErrNotInSlice},
			"value not in slice",
			As,
		},
		{
			User{
				Phones: []string{"213a", "2q"},
				Some:   []int{11, 7, 8},
			},
			[]error{ErrNotMatchRegex, ErrInvalidLen, ErrInvalidMin, ErrInvalidMax},
			"list of errors for slices",
			As,
		},
		{
			InvalidTagLen{},
			nil,
			"invalid tag len",
			Err,
		},
		{
			InvalidTagRegexp{},
			nil,
			"invalid tag regexp",
			Err,
		},
		{
			InvalidTagMin{},
			nil,
			"invalid tag min",
			Err,
		},
		{
			InvalidTagMax{},
			nil,
			"invalid tag max",
			Err,
		},
		{
			InvalidTagInInt{},
			nil,
			"invalid tag in for int",
			Err,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			for i := range tt.expectedErr {
				switch tt.testCase {
				case Is:
					require.ErrorIs(t, err, tt.expectedErr[i])
				case As:
					require.ErrorAs(t, err, &tt.expectedErr[i])
				case Err:
					require.Error(t, err)
				}
			}
		})
	}
}
