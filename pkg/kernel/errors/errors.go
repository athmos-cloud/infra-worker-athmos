package errors

import (
	"errors"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"reflect"
	"strconv"
)

type Error struct {
	Err  string
	Code int
}

//msg string, err error, code int

func New(args ...interface{}) Error {
	invalidArgument := Error{Err: fmt.Sprintf("code must be a number, got : %s", args[1].(string))}
	if len(args) == 2 && option.NewList(reflect.String, args).Validate() {
		codeInt, err := strconv.Atoi(args[1].(string))
		if err != nil {
			return invalidArgument
		}
		return Error{
			Err:  args[0].(string),
			Code: codeInt,
		}
	} else if len(args) == 3 && option.NewList(reflect.String, args).Validate() {
		return Error{
			Err:  fmt.Sprintf("%s: %s", args[0].(string), args[1].(error).Error()),
			Code: args[2].(int),
		}
	} else {
		return invalidArgument
	}
}

func (e *Error) WithMessage(msg interface{}) Error {
	if reflect.TypeOf(msg).Kind() == reflect.String {
		return Error{
			Err:  fmt.Sprintf("%s: %s", e.Err, msg.(string)),
			Code: e.Code,
		}
	}
	if reflect.TypeOf(msg).Kind() == reflect.TypeOf(e).Kind() {
		return Error{
			Err:  fmt.Sprintf("%s: %s", e.Err, msg.(Error).Err),
			Code: e.Code,
		}
	}
	if reflect.TypeOf(msg).Kind() == reflect.TypeOf(errors.New("")).Kind() {
		return Error{
			Err:  fmt.Sprintf("%s: %s", e.Err, msg.(error).Error()),
			Code: e.Code,
		}
	}
	return *e
}

func (e *Error) ToString() string {
	return fmt.Sprintf("%s \n Code: %d", e.Err, e.Code)
}

func (e *Error) Error(errorInput option.Option) error {
	if errorInput.SetType(reflect.String); !errorInput.Validate() {
		return errors.New(errorInput.Value.(string))
	}
	if errorInput.SetType(reflect.TypeOf(e).Kind()); errorInput.Validate() {
		return errors.New(fmt.Sprintf("[%d] %s", e.Code, e.Err))
	}
	return errors.New(e.ToString())
}

func (e *Error) Get() Error {
	if reflect.DeepEqual(e, &OK) {
		return Error{Code: 200, Err: "No error"}
	}
	return *e
}

func (e *Error) Equals(err Error) bool {
	return reflect.DeepEqual(*e, err)
}

func (e *Error) IsOk() bool {
	return reflect.DeepEqual(e, &OK)
}

var (
	OK                   Error
	InvalidArgument      = New("Invalid argument", "400")
	NotFound             = New("Not found", "404")
	ConfigError          = New("Config error", "500")
	ParseError           = New("Parse error", "500")
	IOError              = New("IO error", "500")
	ConversionError      = New("Conversion error", "500")
	WrongConfig          = New("Wrong config", "500")
	ExternalServiceError = New("External service error", "500")
)
