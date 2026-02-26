package humane_test

import (
	"fmt"
	"github.com/sierrasoftworks/humane-errors-go"
	"os"
)

// This example demonstrates how one might go about wrapping an error
// in a humane error, providing advice, and then displaying that to the
// user.
func ExampleError() {
	err := humane.Wrap(
		os.Remove("nonexistent.txt"),
		"We couldn't remove the nonexistent.txt file from the current directory.",
		"Ensure that the file exists in the current directory.",
		"Ensure you have write permissions to the file.",
	)

	humane.Print(err)

	// Output:
	// We couldn't remove the nonexistent.txt file from the current directory.
	//
	// To fix this, you can try:
	//  - Ensure that the file exists in the current directory.
	//  - Ensure you have write permissions to the file.
	//
	// This was caused by:
	//  - remove nonexistent.txt: no such file or directory
	//  - no such file or directory
}

// This example demonstrates creating a new humane error with a formatted
// message and advice using functional options.
func ExampleNewf() {
	email := "alice@example.com"

	err := humane.Newf("failed to resolve user %s", email,
		humane.WithAdvice("ensure the user exists in the directory"),
		humane.WithAdvice("check that the email address is spelled correctly"),
	)

	humane.Print(err)

	// Output:
	// failed to resolve user alice@example.com
	//
	// To fix this, you can try:
	//  - ensure the user exists in the directory
	//  - check that the email address is spelled correctly
}

// This example demonstrates wrapping an existing error with a formatted
// message and advice using functional options.
func ExampleWrapf() {
	cause := fmt.Errorf("connection refused")

	err := humane.Wrapf(cause, "failed to connect to %s on port %d", "db.example.com", 5432,
		humane.WithAdvice("check that the database server is running"),
		humane.WithAdvice("verify your connection string is correct"),
	)

	humane.Print(err)

	// Output:
	// failed to connect to db.example.com on port 5432
	//
	// To fix this, you can try:
	//  - check that the database server is running
	//  - verify your connection string is correct
	//
	// This was caused by:
	//  - connection refused
}

// This example demonstrates providing multiple pieces of advice
// in a single WithAdvice call.
func ExampleWithAdvice() {
	err := humane.Newf("configuration file not found",
		humane.WithAdvice(
			"create a config.yaml file in the current directory",
			"run the init command to generate a default configuration",
		),
	)

	humane.Print(err)

	// Output:
	// configuration file not found
	//
	// To fix this, you can try:
	//  - create a config.yaml file in the current directory
	//  - run the init command to generate a default configuration
}
