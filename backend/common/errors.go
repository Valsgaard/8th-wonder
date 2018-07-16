package common

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// Error ...
type Error interface {
	// Dependencies
	String() string
	Error() string

	// Core
	Message() string
	Code() string
	SetInternal(internal interface{}) Error

	// HTTP Specific
	StatusCode() int
	SetStatusCode(status int) Error

	// Logging utilities
	WithField(key string, val interface{}) Error
	Log(log *logrus.Entry)
}

type baseError struct {
	code       string
	message    string
	internal   string
	httpStatus int
	fields     map[string]interface{}
}

var _ Error = new(baseError)

func (e *baseError) String() string {
	return fmt.Sprintf("%s - %s", e.code, e.message)
}

func (e *baseError) Error() string {
	return e.String()
}

func (e *baseError) Code() string {
	return e.code
}

func (e *baseError) Message() string {
	return e.message
}

func (e *baseError) StatusCode() int {
	return e.httpStatus
}

func (e *baseError) SetStatusCode(status int) Error {
	e.httpStatus = status
	e.WithField("httpStatus", status)
	return e
}

func (e *baseError) WithField(key string, val interface{}) Error {
	e.fields[key] = val
	return e
}

func (e *baseError) Log(log *logrus.Entry) {
	log.WithFields(e.fields).Error(e.message)
}

func (e *baseError) SetInternal(internal interface{}) Error {
	var str string
	switch i := internal.(type) {
	case fmt.Stringer:
		str = i.String()
	case error:
		str = i.Error()
	case string:
		str = i
	default:
		panic("Non-string type")
	}

	e.internal = str
	e.fields["internal"] = str
	return e
}

// Helper functions for error creation
func (e *baseError) setCode(code string) {
	e.code = code
	e.fields["code"] = code
}

// ErrorTemplate is used for creating prepared errors
type ErrorTemplate struct {
	code    string
	message string

	httpStatus int
}

// PrepareError creates a template which can be used to create specific
// instances of errors.
func PrepareError(code, message string) *ErrorTemplate {
	return &ErrorTemplate{
		code:    code,
		message: message,
	}
}

// SetStatusCode sets the http status code for the template
func (e *ErrorTemplate) SetStatusCode(status int) *ErrorTemplate {
	e.httpStatus = status
	return e
}

// NewError creates a new error using a previous entity, or from scratch
func NewError(err interface{}, msg string) Error {
	newErr := new(baseError)
	newErr.fields = make(map[string]interface{})

	if err == nil {
		newErr.message = msg
		return newErr
	}

	switch e := err.(type) {
	case *ErrorTemplate:
		// Use the template
		newErr.message = e.message
		newErr.setCode(e.code)
		newErr.SetStatusCode(e.httpStatus)

		// Overwrite template message if any is given
		if msg != "" {
			newErr.message = msg
		}

	case Error:
		// Continue with the given error, but change the message
		newErr = e.(*baseError) // TODO: check for validity
		newErr.message = msg

	case string:
		// New error with an unprepared code
		newErr.setCode(e)
		newErr.message = msg

	case error:
		// Continue with an error
		newErr.SetInternal(e.Error())
		newErr.message = msg
		if msg == "" {
			newErr.message = e.Error()
		}

	default:
		panic("Uknown error input type")
	}

	return newErr
}
