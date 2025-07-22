package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func IsStructEmpty(s interface{}) bool {
	return reflect.DeepEqual(s, reflect.Zero(reflect.TypeOf(s)).Interface())
}

func Valid(us interface{}) error{
	validate := validator.New()
	return validate.Struct(us)
}