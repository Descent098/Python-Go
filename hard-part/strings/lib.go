package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// A function that generates a string to greet someone
//
// # Parameters
//
// name(string): The name of the person being greeted
//
// # Returns
//
// string: a greeting
func GreetS(name string) string {
	return fmt.Sprintf("Hello %s\nHow's your day?\n", name)
}

// A function that generates a string to greet someone
//
// # Parameters
//
// name(*C.char): The name of the person being greeted
//
// # Returns
//
// *C.char: a greeting to display to the user
//
//export greet_string
func greet_string(name *C.char) *C.char {
	goName := C.GoString(name) // Convert input to go string
	result := GreetS(goName)   // Get a result as a go string
	return C.CString(result)   // Return a C-compatible string
}

// Used to free a C string after use
//
// # Parameters
//
// str (*C.char): The string to free
//
//export free_string
func free_string(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {
	// This has to stay here, but leave it empty
}
