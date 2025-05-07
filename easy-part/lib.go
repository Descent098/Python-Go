package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
)

// A function to greet someone
//
//export Greeting
func Greeting() {
	fmt.Println("Hello from Go!")
}

// A function to calculate the factorial of a number n
//
// # Parameters
//
// n (int): The integer to calculate the factorial of
//
// # Returns
//
// int: The factorial of n
func Factorial(n int) int {
	result := n
	lastVal := n - 1
	for range int(n) {
		if lastVal > 0 {
			result *= lastVal
			lastVal -= 1
		}
	}
	return result
}

// The cgo binding to call the Factorial Function through
//
// # Parameters
//
// n (C.int): The integer to calculate the factorial of
//
// # Returns
//
// C.int: The factorial of n
//
//export factorial
func factorial(n C.int) C.int {
	goN := int(n)            // Convert to go integer
	result := Factorial(goN) // Get go integer result
	r := C.int(result)       // Convert to C integer
	return r
}

func main() {
	// This has to stay here, but leave it empty
}
