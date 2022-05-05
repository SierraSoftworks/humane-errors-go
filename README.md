# Humane Errors
**Errors which make your users' lives easier**

Most of the time, errors are the kind of thing your development team tries to
avoid thinking about, and they're something your users dread receiving. The
trouble is, your developers are the people best equipped to understand why an
error might occur and provide advice for how to deal with it.

This module provides an error type which is designed to help your developers
share that knowledge directly with your users, all while keeping your code
simple and readable.

## Features
 - Integrates with Go's native `error` types and the `errors` package, including `Unwrap()`.
 - Provides context specific advice to your users to help them respond to an error.
 - Simple APIs designed to make your life easy.
 - Proactive support for [Go 2.0's error spec](https://go.googlesource.com/proposal/+/master/design/go2draft.md). 

## Example

```go
package main

import (
	humane "github.com/sierrasoftworks/humane-errors-go"
	"io/ioutil"
	"os"
)

func main() {
	err := humane.Wrap(
		ioutil.WriteFile("demo.txt", []byte("This is an example"), os.ModePerm),
		"We couldn't write the demo.txt file to the current directory.",
		"Ensure you have write permissions to the current directory.",
		"Make sure you have free space on your disk.",
	)
	
	humane.Eprint(err)
}
```