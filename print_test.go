package humane_test

import (
	"fmt"
	"testing"

	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	err := humane.Wrap(fmt.Errorf("internal error"), "Something went wrong", "Try turning it off...permanently.")
	assert.Equal(t, `Something went wrong`, fmt.Sprintf("%s", err))
	assert.Equal(t, `Something went wrong`, fmt.Sprintf("%v", err))

	// TODO: We should add a test for the %+v formatter once that is supported
}
