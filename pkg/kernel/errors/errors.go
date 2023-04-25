package errors

import (
	"errors"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"reflect"
	"strconv"
)

type Error struct {
	Err  string
	Code int
}

// New : msg string, err error, code int
func New(message string, code string) Error {
	strCode, err := strconv.Atoi(code)
	if err != nil {
		return Error{
			Err:  fmt.Sprintf("Error code %s is not a number", code),
			Code: 500,
		}
	}
	return Error{
		Err:  message,
		Code: strCode,
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
	return fmt.Sprintf("[%d]: %s", e.Code, e.Err)
}

func (e *Error) Error(errorInput option.Option) error {
	if errorInput.SetType(reflect.String.String()); !errorInput.Validate() {
		return errors.New(errorInput.Value.(string))
	}
	if errorInput.SetType(reflect.TypeOf(e).String()); errorInput.Validate() {
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
	return reflect.DeepEqual(e, &OK) ||
		reflect.DeepEqual(e, &Created) ||
		reflect.DeepEqual(e, &Accepted) ||
		reflect.DeepEqual(e, &NoContent)
}

func (e *Error) IsNotFound() bool {
	return reflect.DeepEqual(e, &NotFound)
}

var (
	OK                   Error
	Created              = New("Created", "201")
	Accepted             = New("Accepted", "202")
	NoContent            = New("No content", "204")
	InvalidArgument      = New("Invalid argument", "400")
	ValidationError      = New("Validation error", "400")
	Conflict             = New("Already exists", "409")
	NotFound             = New("Not found", "404")
	ConfigError          = New("Config error", "500")
	ParseError           = New("Parse error", "500")
	InternalError        = New("Internal error", "500")
	IOError              = New("IO error", "500")
	ConversionError      = New("Conversion error", "500")
	WrongConfig          = New("Wrong config", "500")
	ExternalServiceError = New("External application error", "500")
)
