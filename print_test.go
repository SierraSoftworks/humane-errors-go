package humane_test

import (
	"fmt"
	"testing"

	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFormatter(t *testing.T) {
	err := humane.Wrap(fmt.Errorf("internal error"), "Something went wrong", "Try turning it off...permanently.")
	assert.Equal(t, `Something went wrong`, fmt.Sprintf("%s", err))
	assert.Equal(t, `Something went wrong`, fmt.Sprintf("%v", err))

	// TODO: We should add a test for the %+v formatter once that is supported
}

func TestZapLog(t *testing.T) {
	expected := make([]zap.Field, 0)
	expected = append(expected, zap.Errors("causes", []error{fmt.Errorf("internal error")}))
	expected = append(expected, zap.Strings("advice", []string{"Try turning it off...permanently."}))
	expected = append(expected, zap.Error(fmt.Errorf("Something went wrong")))

	err := humane.Wrap(fmt.Errorf("internal error"), "Something went wrong", "Try turning it off...permanently.")
	zapFields := err.IntoZapLog()
	assert.Equal(t, 3, len(zapFields), "the zap log should have two fields")
	assert.Equal(t, expected, zapFields, "the zap log should have the expected fields")
}
