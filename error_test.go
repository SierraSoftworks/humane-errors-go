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

func TestNewf(t *testing.T) {
	err := humane.Newf("failed to resolve user %s", "alice@example.com",
		humane.WithAdvice("ensure the user exists"),
	)

	assert.Equal(t, "failed to resolve user alice@example.com", err.Error())
	assert.Equal(t, []string{"ensure the user exists"}, err.Advice())
	assert.Nil(t, err.Cause())
}

func TestNewfNoFormatArgs(t *testing.T) {
	err := humane.Newf("static message",
		humane.WithAdvice("some advice"),
	)

	assert.Equal(t, "static message", err.Error())
	assert.Equal(t, []string{"some advice"}, err.Advice())
}

func TestNewfMultipleFormatArgs(t *testing.T) {
	err := humane.Newf("user %s has %d items", "alice", 42,
		humane.WithAdvice("check inventory"),
	)

	assert.Equal(t, "user alice has 42 items", err.Error())
}

func TestNewfNoOptions(t *testing.T) {
	err := humane.Newf("failed to resolve user %s", "alice@example.com")

	assert.Equal(t, "failed to resolve user alice@example.com", err.Error())
	assert.Empty(t, err.Advice())
	assert.Nil(t, err.Cause())
}

func TestNewfMultipleWithAdvice(t *testing.T) {
	err := humane.Newf("something failed",
		humane.WithAdvice("first advice"),
		humane.WithAdvice("second advice"),
	)

	assert.Equal(t, []string{"first advice", "second advice"}, err.Advice())
}

func TestNewfVariadicWithAdvice(t *testing.T) {
	err := humane.Newf("something failed",
		humane.WithAdvice("advice 1", "advice 2"),
	)

	assert.Equal(t, []string{"advice 1", "advice 2"}, err.Advice())
}

func TestNewfMixedWithAdvice(t *testing.T) {
	err := humane.Newf("something failed",
		humane.WithAdvice("advice 1", "advice 2"),
		humane.WithAdvice("advice 3"),
	)

	assert.Equal(t, []string{"advice 1", "advice 2", "advice 3"}, err.Advice())
}

func TestWrapfNilCause(t *testing.T) {
	err := humane.Wrapf(nil, "failed to resolve %s", "alice")
	assert.Nil(t, err, "wrapf with nil cause should return nil")
}

func TestWrapf(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := humane.Wrapf(cause, "failed to connect to %s:%d", "localhost", 5432,
		humane.WithAdvice("check that the database is running"),
	)

	assert.Equal(t, "failed to connect to localhost:5432", err.Error())
	assert.Equal(t, []string{"check that the database is running"}, err.Advice())
	assert.Equal(t, cause, err.Cause())
}

func TestWrapfDisplay(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := humane.Wrapf(cause, "failed to connect to %s", "database",
		humane.WithAdvice("check that the database is running"),
		humane.WithAdvice("verify your connection string"),
	)

	assert.Equal(t, `failed to connect to database

To fix this, you can try:
 - check that the database is running
 - verify your connection string

This was caused by:
 - connection refused`, err.Display())
}

func TestWrapfErrorIs(t *testing.T) {
	base := fmt.Errorf("base error")
	err := humane.Wrapf(base, "wrapped %s", "error",
		humane.WithAdvice("some advice"),
	)

	assert.True(t, errors.Is(err, base), "wrapf error should be equivalent to the base error")
	assert.True(t, errors.Is(err, fmt.Errorf("wrapped error")), "wrapf error should be equivalent to an error with the same message")
}

func TestErrorIs(t *testing.T) {
	base := fmt.Errorf("base error")

	err := humane.Wrap(base, "this is a wrapping error")
	assert.True(t, errors.Is(err, base), "the humane error should be equivalent to the base error")
	assert.True(t, errors.Is(err, fmt.Errorf("this is a wrapping error")), "the humane error should be equivalent to an error with the same message")
}
