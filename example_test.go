package humane_test

import (
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
