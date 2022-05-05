package humane

import (
	"fmt"
	"io"
	"os"
)

// The Format method is implemented to support Go 2's error printing
// proposal, which can be viewed [here](https://go.googlesource.com/proposal/+/master/design/go2draft-error-printing.md).
// It allows an error to be printed with additional context, if requested.
func (e *humaneError) Format(p interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Detail() bool
}) (next error) {
	p.Print(e.message)
	if p.Detail() {
		if len(e.advice) > 0 {
			p.Print("Advice:")
			for _, advice := range e.advice {
				p.Printf(" - %s", advice)
			}
		}
	}

	return e.cause
}

// The Print method will print an error to your `Stdout` stream.
// If the error implements a `Display()` method, it will be used
// to format the error, falling back on the standard `Error()` method.
func Print(err error) {
	if err == nil {
		return
	}

	if h, ok := err.(interface {
		Display() string
	}); ok {
		_, _ = fmt.Println(h.Display())
	} else {
		_, _ = fmt.Println(err.Error())
	}
}

// The Eprint method will print an error to your `Stderr` stream.
// If the error implements a `Display()` method, it will be used
// to format the error, falling back on the standard `Error()` method.
func Eprint(err error) {
	if err == nil {
		return
	}

	if h, ok := err.(interface {
		Display() string
	}); ok {
		_, _ = fmt.Fprintln(os.Stderr, h.Display())
	} else {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

// The Fprint method will print an error to the provided stream.
// If the error implements a `Display()` method, it will be used
// to format the error, falling back on the standard `Error()` method.
func Fprint(w io.Writer, err error) error {
	if err == nil {
		return nil
	}

	if h, ok := err.(interface {
		Display() string
	}); ok {
		_, err := fmt.Fprintln(w, h.Display())
		return err
	} else {
		_, err := fmt.Fprintln(w, err.Error())
		return err
	}
}
