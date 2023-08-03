package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const ValidationTag = "validate"

type Rule int

const (
	Len Rule = iota
	Regexp
	In
	Min
	Max
)

var (
	ErrNotStruct     = errors.New("input value not a struct")
	ErrInvalidRule   = errors.New("invalid rule")
	ErrInvalidLen    = errors.New("invalid len of string")
	ErrNotMatchRegex = errors.New("not match regex")
	ErrNotInSlice    = errors.New("value not in slice ")
	ErrInvalidMin    = errors.New("value less than the min")
	ErrInvalidMax    = errors.New("value more than the max")
)

type Validators struct {
	Kind       reflect.Kind
	validators []Validator
}

type Validator struct {
	Rule  Rule
	Value string
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	errorList := fmt.Errorf("errors: ")
	for _, err := range v {
		errorList = fmt.Errorf("%w for field '%s' error is %s", errorList, err.Field, err.Err)
	}
	return errorList.Error()
}

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var validationErrors ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := reflect.ValueOf(v).FieldByName(field.Name)
		if tag, ok := field.Tag.Lookup(ValidationTag); ok {
			typeValue := fieldValue.Type()
			if fieldValue.Kind() == reflect.Slice {
				typeValue = fieldValue.Type().Elem()
			}
			validators, err := validateTag(tag, typeValue)
			if err != nil {
				return err
			}
			if len(validators.validators) == 0 {
				continue
			}
			errList := validateField(fieldValue, validators)
			for _, err := range errList {
				validationErrors = append(validationErrors, ValidationError{field.Name, err})
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateField(value reflect.Value, validators Validators) []error {
	if validators.Kind == reflect.String {
		return validateStringField(value, validators.validators)
	}
	if validators.Kind == reflect.Int {
		return validateIntField(value, validators.validators)
	}
	return nil
}

func validateStringField(value reflect.Value, validators []Validator) []error {
	var errorList []error
	for _, validator := range validators {
		var values []string
		if value.Kind() == reflect.Slice {
			for i := 0; i < value.Len(); i++ {
				values = append(values, value.Index(i).String())
			}
		} else {
			values = append(values, value.String())
		}
		switch validator.Rule { //nolint:exhaustive
		case Len:
			length, _ := strconv.Atoi(validator.Value)
			if err := validateLen(values, length); err != nil {
				errorList = append(errorList, err)
			}
		case Regexp:
			re := regexp.MustCompile(validator.Value)
			if err := validateRegexp(values, re); err != nil {
				errorList = append(errorList, err)
			}
		case In:
			if err := validateIn(values, strings.Split(validator.Value, ",")); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return errorList
}

func validateIntField(value reflect.Value, validators []Validator) []error {
	var errorList []error
	var values []int64
	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			values = append(values, value.Index(i).Int())
		}
	} else {
		values = append(values, value.Int())
	}
	for _, validator := range validators {
		switch validator.Rule { //nolint:exhaustive
		case Min:
			min, _ := strconv.Atoi(validator.Value)
			if err := validateMin(values, int64(min)); err != nil {
				errorList = append(errorList, err)
			}
		case Max:
			max, _ := strconv.Atoi(validator.Value)
			if err := validateMax(values, int64(max)); err != nil {
				errorList = append(errorList, err)
			}
		case In:
			if err := validateIn(intSliceToString(values), strings.Split(validator.Value, ",")); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return errorList
}

func validateMin(values []int64, min int64) error {
	for _, value := range values {
		if value < min {
			return ErrInvalidMin
		}
	}
	return nil
}

func validateMax(values []int64, max int64) error {
	for _, value := range values {
		if value > max {
			return ErrInvalidMax
		}
	}
	return nil
}

func validateLen(values []string, length int) error {
	for _, value := range values {
		if len(value) != length {
			return ErrInvalidLen
		}
	}
	return nil
}

func validateRegexp(values []string, re *regexp.Regexp) error {
	for _, value := range values {
		if !re.Match([]byte(value)) {
			return ErrNotMatchRegex
		}
	}
	return nil
}

func validateIn[T comparable](values []T, expected []T) error {
	for _, value := range values {
		var ok bool
		for _, e := range expected {
			if e == value {
				ok = true
			}
		}
		if !ok {
			return ErrNotInSlice
		}
	}
	return nil
}

func validateTag(tag string, value reflect.Type) (Validators, error) {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateStringTag(tag)
	case reflect.Int:
		return validateIntTag(tag)
	default:
		return Validators{}, nil
	}
}

func validateStringTag(tag string) (Validators, error) {
	rules := strings.Split(tag, "|")
	validators := Validators{Kind: reflect.String}
	for _, rule := range rules {
		fields := strings.Split(rule, ":")
		if len(fields) != 2 {
			return Validators{}, ErrInvalidRule
		}
		switch fields[0] {
		case "len":
			_, err := strconv.Atoi(fields[1])
			if err != nil {
				return Validators{}, err
			}
			validators.validators = append(validators.validators, Validator{Len, fields[1]})
		case "regexp":
			_, err := regexp.Compile(fields[1])
			if err != nil {
				return Validators{}, err
			}
			validators.validators = append(validators.validators, Validator{Regexp, fields[1]})
		case "in":
			validators.validators = append(validators.validators, Validator{In, fields[1]})
		default:
			continue
		}
	}
	return validators, nil
}

func validateIntTag(tag string) (Validators, error) {
	rules := strings.Split(tag, "|")
	validators := Validators{Kind: reflect.Int}
	for _, rule := range rules {
		fields := strings.Split(rule, ":")
		if len(fields) != 2 {
			return Validators{}, ErrInvalidRule
		}
		switch fields[0] {
		case "min":
			_, err := strconv.Atoi(fields[1])
			if err != nil {
				return Validators{}, err
			}
			validators.validators = append(validators.validators, Validator{Min, fields[1]})
		case "max":
			_, err := strconv.Atoi(fields[1])
			if err != nil {
				return Validators{}, err
			}
			validators.validators = append(validators.validators, Validator{Max, fields[1]})
		case "in":
			numList := strings.Split(fields[1], ",")
			for _, num := range numList {
				_, err := strconv.Atoi(num)
				if err != nil {
					return Validators{}, err
				}
			}
			validators.validators = append(validators.validators, Validator{In, fields[1]})
		default:
			continue
		}
	}
	return validators, nil
}

func intSliceToString(values []int64) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strconv.Itoa(int(value)))
	}
	return result
}
