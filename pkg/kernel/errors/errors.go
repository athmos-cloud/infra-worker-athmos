package errors

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"reflect"
	"strconv"
)

type Error struct {
	Code    int
	Message string
}

// New : msg string, err error, code int
func New(message string, code string) Error {
	strCode, err := strconv.Atoi(code)
	if err != nil {
		return Error{
			Message: fmt.Sprintf("Error code %s is not a number", code),
			Code:    500,
		}
	}
	return Error{
		Message: message,
		Code:    strCode,
	}
}

func (e *Error) WithMessage(msg interface{}) Error {
	if reflect.TypeOf(msg).Kind() == reflect.String {
		return Error{
			Message: fmt.Sprintf("%s: %s", e.Message, msg.(string)),
			Code:    e.Code,
		}
	}
	if reflect.TypeOf(msg).Kind() == reflect.TypeOf(e).Kind() {
		return Error{
			Message: fmt.Sprintf("%s: %s", e.Message, msg.(Error).Message),
			Code:    e.Code,
		}
	}
	if reflect.TypeOf(msg).Kind() == reflect.TypeOf(errors.New("")).Kind() {
		return Error{
			Message: fmt.Sprintf("%s %s", e.Message, msg.(error).Error()),
			Code:    e.Code,
		}
	}
	return *e
}

func (e *Error) ToString() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Message)
}

func (e *Error) Error(errorInput option.Option) error {
	if errorInput.SetType(reflect.String.String()); !errorInput.Validate() {
		return errors.New(errorInput.Value.(string))
	}
	if errorInput.SetType(reflect.TypeOf(e).String()); errorInput.Validate() {
		return errors.New(fmt.Sprintf("[%d] %s", e.Code, e.Message))
	}
	return errors.New(e.ToString())
}

func (e *Error) Get() Error {
	if reflect.DeepEqual(e, &OK) {
		return Error{Code: 200, Message: "No error"}
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

func RaiseError(ctx goContext.Context, err any) {
	if reflect.TypeOf(err) == reflect.TypeOf(Error{}) {
		errorRaised := err.(Error)
		ctx = goContext.WithValue(ctx, context.ResponseKey, errorRaised)
	} else {
		ctx = goContext.WithValue(ctx, context.ResponseKey, InternalError.WithMessage(err.(string)))
	}
	ctx.Done()
}

var (
	OK                   Error
	Created              = New("Created", "201")
	Accepted             = New("Accepted", "202")
	NoContent            = New("No content", "204")
	BadRequest           = New("Bad request", "400")
	InvalidArgument      = New("Invalid argument", "400")
	Conflict             = New("Conflict", "409")
	NotFound             = New("Not found", "404")
	InternalError        = New("Internal error", "500")
	InvalidOption        = New("Invalid option", "500")
	KubernetesError      = New("Kubernetes error", "500")
	ConversionError      = New("Conversion error", "500")
	ExternalServiceError = New("External application error", "500")
)
