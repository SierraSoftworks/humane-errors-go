package humane_test

import (
	"errors"
	"fmt"
	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNilError(t *testing.T) {
	err := humane.Wrap(nil, "No error occurred here")
	assert.Nil(t, err, "the nil error should be curried")
}

func TestSimpleError(t *testing.T) {
	err := humane.New("Something exploded unexpectedly", "Avoid exploding things...")
	assert.Equal(t, "Something exploded unexpectedly", err.Error())
	assert.Contains(t, err.Advice(), "Avoid exploding things...")
	assert.Nil(t, err.Cause())

	assert.Equal(t, `Something exploded unexpectedly

To fix this, you can try:
 - Avoid exploding things...`, err.Display())
}

func TestErrorNoAdvice(t *testing.T) {
	err := humane.New("Something went wrong")
	assert.Equal(t, "Something went wrong", err.Display())
}

func TestErrorWithCause(t *testing.T) {
	err := humane.Wrap(
		humane.New("This code is broken", "Ask us to not write broken code"),
		"Something broke",
		"Go yell at us on Twitter",
	)

	assert.Equal(t, `Something broke

To fix this, you can try:
 - Ask us to not write broken code
 - Go yell at us on Twitter

This was caused by:
 - This code is broken`, err.Display())
}

func TestErrorIs(t *testing.T) {
	base := fmt.Errorf("base error")

	err := humane.Wrap(base, "this is a wrapping error")
	assert.True(t, errors.Is(err, base), "the humane error should be equivalent to the base error")
	assert.True(t, errors.Is(err, fmt.Errorf("this is a wrapping error")), "the humane error should be equivalent to an error with the same message")
}
