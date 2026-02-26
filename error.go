package humane

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
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
}

type humaneError struct {
	message string
	advice  []string
	cause   error
}

// Option is a functional option that configures a humane error
// created with [Newf] or [Wrapf]. Options provide a flexible,
// extensible way to attach metadata such as advice to an error.
type Option interface {
	apply(*humaneError)
}

type adviceOption struct {
	advice []string
}

func (o *adviceOption) apply(e *humaneError) {
	e.advice = append(e.advice, o.advice...)
}

// WithAdvice returns an [Option] that attaches one or more pieces of advice
// to a humane error. Multiple pieces of advice can be provided in a single
// call, and multiple WithAdvice options will append to the error's advice list.
func WithAdvice(advice ...string) Option {
	return &adviceOption{advice: advice}
}

// parseArgs separates format arguments from Options in a mixed variadic slice.
func parseArgs(args []any) (fmtArgs []any, opts []Option) {
	for _, a := range args {
		if o, ok := a.(Option); ok {
			opts = append(opts, o)
		} else {
			fmtArgs = append(fmtArgs, a)
		}
	}
	return
}

// applyOptions applies a slice of Options to a humaneError.
func applyOptions(e *humaneError, opts []Option) {
	for _, o := range opts {
		o.apply(e)
	}
}

// New constructs a humane error without a defined cause.
// It is commonly used when creating an initial error, in contrast
// to [Wrap] which is used to wrap an existing error.
func New(message string, advice ...string) Error {
	return &humaneError{message, advice, nil}
}

// Newf constructs a humane error with a formatted message and functional options.
// Format arguments and [Option] values can be freely intermixed in the variadic
// args; any argument implementing [Option] is extracted and applied to the error,
// while the remaining arguments are passed to [fmt.Sprintf] along with the
// format string.
func Newf(format string, args ...any) Error {
	fmtArgs, opts := parseArgs(args)

	e := &humaneError{
		message: fmt.Sprintf(format, fmtArgs...),
	}
	applyOptions(e, opts)

	return e
}

// Wrap constructs a humane error which wraps an existing error.
// It is commonly used when another error has resulted in something that
// should be bubbled up to the user to handle.
// If nil is provided as the cause, this function will return nil as well,
// allowing it to safely wrap optional errors and curry their presence correctly.
func Wrap(cause error, message string, advice ...string) Error {
	if cause == nil {
		return nil
	}

	return &humaneError{message, advice, cause}
}

// Wrapf constructs a humane error which wraps an existing error, with a
// formatted message and functional options. Format arguments and [Option]
// values can be freely intermixed in the variadic args; any argument
// implementing [Option] is extracted and applied to the error, while the
// remaining arguments are passed to [fmt.Sprintf] along with the format string.
// If nil is provided as the cause, this function will return nil as well,
// allowing it to safely wrap optional errors and curry their presence correctly.
func Wrapf(cause error, format string, args ...any) Error {
	if cause == nil {
		return nil
	}

	fmtArgs, opts := parseArgs(args)

	e := &humaneError{
		message: fmt.Sprintf(format, fmtArgs...),
		cause:   cause,
	}
	applyOptions(e, opts)

	return e
}

// Display returns a structured error string which provides the user with
// advice on how best to respond, along with information about the underlying
// cause of the error.
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

// The Unwrap method returns the error which resulted in this error being
// raised. It may be a linked list of errors, in which case calling Unwrap
// on each resulting error will return the parent and so on. This method
// behaves identically to the `Cause()` method but provides interop with Go's
// `error.Unwrap()` functionality.
func (e *humaneError) Unwrap() error {
	return e.cause
}
