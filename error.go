package humane

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

var tmpl = template.Must(template.New("").Parse(strings.TrimSpace(`
{{ .Message }}
{{- if .Advice }}

To fix this, you can try:
{{- range .Advice }}
 - {{ . }}
{{- end -}}
{{- end -}}
{{- if .Causes }}

This was caused by:
{{- range .Causes }}
 - {{ .Error }}
{{- end -}}
{{- end -}}
`)))

// The Error interface describes what one can do with a humane error.
//
// Humane errors provide the ability to report information on the suspected
// cause of an error, as well as any advice that may be offered to the user
// to support them in responding to the problem.
//
// This advice and cause information can be expanded upon by wrapping errors,
// providing the user with multiple pieces of advice which may be increasingly
// less specific as the stack is popped.
type Error interface {
	// The Display method returns a structured error string which provides the
	// user with advice on how best to respond.
	Display() string

	// The Error method returns the error message associated with this error.
	// This method allows humane errors to implement Go's `error` type and
	// return a sensible message when used in this capacity.
	//
	// **NOTE**: In general you should use the `Display()` method when presenting
	// an error to a user.
	Error() string

	// The Advice method provides a list of advice which is intended to help
	// the user recover from the error.
	Advice() []string

	// The Cause method returns the error which caused this one to be raised.
	// This may be a linked list of errors pointing to the underlying failure.
	Cause() error

	// IntoZapLog converts the error into []zap.Field for logging purposes.
	IntoZapLog() []zap.Field
}

type humaneError struct {
	message string
	advice  []string
	cause   error
}

// The New method constructs a humane error without a defined cause.
// It is commonly used when creating an initial error, in contrast
// to the `Wrap` method which is used to wrap an existing error.
func New(message string, advice ...string) Error {
	return &humaneError{message, advice, nil}
}

// The Wrap method constructs a human error which wraps an existing error.
// It is commonly used when another error has resulted in something that should
// be bubbled up to the user to handle.
// If `nil` is provided as the cause, this method will return `nil` as well,
// allowing it to safely wrap optional errors and currying their presence correctly.
func Wrap(cause error, message string, advice ...string) Error {
	if cause == nil {
		return nil
	}

	return &humaneError{message, advice, cause}
}

// The Display method implements Go's `error` interface and returns a templated
// error message describing what a user may do to respond to the issue.
func (e *humaneError) Display() string {
	b := bytes.NewBufferString("")

	context := struct {
		Message string
		Advice  []string
		Causes  []error
	}{
		Message: e.message,
		Advice:  e.advice,
		Causes:  []error{},
	}

	cause := e.cause
	for cause != nil {
		context.Causes = append(context.Causes, cause)

		if cause, ok := e.cause.(interface {
			Advice() []string
		}); ok {
			context.Advice = append(cause.Advice(), context.Advice...)
		}

		cause = errors.Unwrap(cause)
	}

	err := tmpl.Execute(b, context)
	if err != nil {
		return fmt.Sprintf("%s (error building nice message: %s)", e.message, err.Error())
	}

	return b.String()
}

// The Is method tests whether this error is equivalent to another.
// We use message equality to determine whether errors are equivalent,
// and will attempt to match against the cause of the humane error as well
// if it is present.
func (e *humaneError) Is(err error) bool {
	if err == nil {
		return false
	}

	if is, ok := e.cause.(interface {
		Is(error) bool
	}); ok {
		return e.message == err.Error() || is.Is(err)
	}

	return e.message == err.Error() || e.cause == err
}

// The Error method returns the message describing the error that has occurred.
func (e *humaneError) Error() string {
	return e.message
}

// The Advice method returns a list of pieces of advice which should be
// provided to the user to help them respond to the failure.
func (e *humaneError) Advice() []string {
	return e.advice
}

// The Cause method returns the error which resulted in this error being
// raised. It may be a linked list of errors, in which case their advice
// will be aggregated when rendering.
func (e *humaneError) Cause() error {
	return e.cause
}

// IntoZapLog converts the error into []zap.Field for logging purposes.
func (e *humaneError) IntoZapLog() []zap.Field {
	context := struct {
		Advice []string
		Causes []error
	}{
		Advice: e.advice,
		Causes: []error{},
	}

	cause := e.cause
	for cause != nil {
		context.Causes = append(context.Causes, cause)

		if cause, ok := e.cause.(interface {
			Advice() []string
		}); ok {
			context.Advice = append(cause.Advice(), context.Advice...)
		}

		cause = errors.Unwrap(cause)
	}

	zapFields := make([]zap.Field, 0)
	zapFields = append(zapFields, zap.Errors("causes", context.Causes))
	zapFields = append(zapFields, zap.Strings("advice", context.Advice))
	zapFields = append(zapFields, zap.Error(errors.New(e.message)))
	return zapFields
}
