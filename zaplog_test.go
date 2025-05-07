package humane_test

import (
	"fmt"
	"testing"

	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapLog(t *testing.T) {
	t.Run("with Humane Error", func(t *testing.T) {
		expected := []zap.Field{
			zap.Error(fmt.Errorf("Something went wrong")),
			zap.Strings("advice", []string{"Try turning it off...permanently."}),
			zap.Errors("causes", []error{fmt.Errorf("internal error")}),
		}

		err := humane.Wrap(fmt.Errorf("internal error"), "Something went wrong", "Try turning it off...permanently.")
		zapFields := humane.Zap(err)
		assert.Equal(t, 3, len(zapFields), "the zap log should have two fields")
		assert.Equal(t, expected, zapFields, "the zap log should have the expected fields")
	})

	t.Run("with Normal Error", func(t *testing.T) {
		expected := []zap.Field{
			zap.Error(fmt.Errorf("Something went wrong: internal error")),
			zap.Strings("advice", []string{}),
			zap.Errors("causes", []error{fmt.Errorf("internal error")}),
		}

		err := fmt.Errorf("Something went wrong: %w", fmt.Errorf("internal error"))
		zapFields := humane.Zap(err)
		assert.Equal(t, 3, len(zapFields), "the zap log should have two fields")
		assert.Equal(t, expected, zapFields, "the zap log should have the expected fields")
	})
}
