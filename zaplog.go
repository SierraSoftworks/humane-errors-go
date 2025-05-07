package humane

import (
	"errors"

	"go.uber.org/zap"
)

// Zap converts a humane error into []zap.Field for logging purposes.
// It is intended to be used in conjuction with Zap's `.Error(...)` method to
// capture information about an error which has occurred within your application.
//
//	    err := humane.Wrap(
//		      os.Remove("nonexistent.txt"),
//		      "We couldn't remove the nonexistent.txt file from the current directory.",
//		      "Ensure that the file exists in the current directory.",
//		      "Ensure you have write permissions to the file.",
//	    )
//
//	    zap.L().Error("file deletion failed unexpectedly", humane.Zap(err)...)
func Zap(e error) []zap.Field {
	context := struct {
		Advice []string
		Causes []error
	}{
		Advice: []string{},
		Causes: []error{},
	}

	err := e
	for err != nil {
		context.Causes = append(context.Causes, err)
		if c, ok := err.(Error); ok {
			context.Advice = append(context.Advice, c.Advice()...)
		}

		err = errors.Unwrap(err)
	}

	return []zap.Field{
		zap.Error(errors.New(e.Error())),
		zap.Strings("advice", context.Advice),
		zap.Errors("causes", context.Causes[1:]),
	}
}
