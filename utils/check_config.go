package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// baseError is an error type that all other error types embed.
type baseError struct {
	defaultErrString string
	errInfo          string
}

func (e baseError) choseErrString() string {
	if e.errInfo != "" {
		return e.errInfo
	}

	return e.defaultErrString
}

// errMissingInput is the error when input is required in a particular
// situation but not provided by the user
type errMissingInput struct {
	baseError

	errArgument string
}

func (e errMissingInput) Error() string {
	e.defaultErrString = fmt.Sprintf("Missing input for argument [%s]", e.errArgument)

	return e.choseErrString()
}

// CheckConfig
func CheckConfig(opts interface{}, parent string) error {
	optsValue := reflect.ValueOf(opts)
	if optsValue.Kind() == reflect.Ptr {
		optsValue = optsValue.Elem()
	}

	optsType := reflect.TypeOf(opts)
	if optsType.Kind() == reflect.Ptr {
		optsType = optsType.Elem()
	}

	if optsValue.Kind() != reflect.Struct {
		return fmt.Errorf("Options type is not a struct.")
	}

	fieldChain := func(s string) string {
		if parent == "" {
			return s
		}

		return parent + "." + s
	}

	isStruct := func(v *reflect.Value) bool {
		return v.Kind() == reflect.Struct || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct)
	}

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)

		if f.Tag.Get("json") == "-" || f.Name != strings.Title(f.Name) {
			continue
		}

		if v.Kind() == reflect.Slice || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice) {
			sliceValue := v
			if sliceValue.Kind() == reflect.Ptr {
				sliceValue = sliceValue.Elem()
			}

			for i := 0; i < sliceValue.Len(); i++ {
				element := sliceValue.Index(i)

				if isStruct(&element) {
					if err := CheckConfig(element.Interface(), fieldChain(f.Name)); err != nil {
						return err
					}
				}
			}
		}

		if isStruct(&v) {
			if err := CheckConfig(v.Interface(), fieldChain(f.Name)); err != nil {
				return err
			}
		}

		if t := f.Tag.Get("required"); t == "true" && isZero(v) {
			return errMissingInput{
				errArgument: fieldChain(f.Name),
			}
		}
	}

	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return false
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		var t time.Time

		if v.Type() == reflect.TypeOf(t) {
			if v.Interface().(time.Time).IsZero() {
				return true
			}
			return false
		}
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}

	// Compare other types directly:
	z := reflect.Zero(v.Type())

	return v.Interface() == z.Interface()
}
