# Python To Go Helper

This is a library that contains python and go utilities for passing data between the two languages. Please note that while I've done my best due diligence I cannot guarentee there are no memory leaks in the current code. I've tested in a bunch of scenarios, but I would recommend being vigilant in any code that uses the library.

## Python

The python lib has the following API functions:

**Converting to ctypes**

- `prepare_string(data: str | bytes) -> c_char_p`: Takes in a string and returns a C-compatible string
- `prepare_string_array(data:list[str|bytes]) -> tuple[Array[c_char_p], int]`: Takes in a string list, and converts it to a C-compatible array
- `prepare_int_array(data:list[int]) -> tuple[Array[c_int], int]`: Takes in a int list, and converts it to a C-compatible array
- `prepare_float_array(data:list[float]) -> tuple[Array[c_float], int]`: Takes in a float list, and converts it to a C-compatible array

**Converting from ctypes**

- `string_to_str(pointer: c_char_p) -> str`: Takes in a pointer to a C string and returns a Python string
- `string_array_result_to_list(pointer:_CStringArrayResult) -> list[str]`: 
- `int_array_result_to_list(pointer: _CIntArrayResult) -> list[int]`: 
- `float_array_result_to_list(pointer: _CFloatArrayResult) -> list[float]`: 

**Debugging Functions**

- `return_string(text: str | bytes) -> str`: Debugging function that shows you the Go representation of a C string and returns the python string version
- `return_string_array(c_array:CStringArray, number_of_elements:int) ->list[str]`: Debugging function that shows you the Go representation of a C array and returns the python list version (does not free)
- `return_int_array(c_array: CIntArray, number_of_elements: int) -> list[int]`: Debugging function that shows you the Go representation of a C int array and returns a Python list
- `return_float_array(c_array: CFloatArray, number_of_elements: int) -> list[float]`: Debugging function that shows you the Go representation of a C float array and returns a Python list
- `print_string(text: str | bytes)`: Prints a string's go representation, useful to look for encoding issues
- `print_string_array(data:list[str|bytes])`: Prints a string array's go representation, useful to look for encoding issues
- `print_int_array(data:list[int])`: Prints a int array's go representation, useful to look for rounding/conversion issues
- `print_float_array(data:list[float])`: Prints a float array's go representation, useful to look for rounding/conversion issues

**Freeing Functions**

- `free_c_string(ptr: c_char_p)`: Frees a single C string returned from Go (allocated via C.CString).
- `free_string_array(ptr: CStringArray, count: int)`: Frees an array of C strings returned from Go.
- `free_int_array(ptr: CIntArray)`: Frees a C int array returned from Go.
- `free_float_array(ptr: CFloatArray)`: Frees a C float array returned from Go.
- `free_string_array_result(ptr: _CStringArrayResult)`: Frees a StringArrayResult (including the array of strings and struct itself).
- `free_int_array_result(ptr: _CIntArrayResult)`: Frees an IntArrayResult (including the array and the struct itself).
- `free_float_array_result(ptr: _CFloatArrayResult)`: Frees a FloatArrayResult (including the array and the struct itself).


### Tests

To run the tests first install pytest:

```bash
pip install pytest pytest-cov
```

To run tests install pytest and run:

```bash
pytest --ignore=__init__.py --cov-report term-missing --cov=. test_lib.py
```

This will run the test suite and let you know any coverage misses. There's ~%80 coverage currently due to some conditions not being possible (or I don't know how to make them happen)

## Go

Below are details for hooking up the go side of your code with the helper

### Setup

In your own go code import the package with:

```go
import (
	helpers "github.com/Descent098/cgo-python-helpers/helpers"
)
```

Then run:

```bash
go mod tidy
```

Here is an example:

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	helpers "github.com/Descent098/cgo-python-helpers"
)

func main() {
	// Sample data
	numbers := []int{1, 2, 3, 4, 5}

	// Convert Go slice to C-compatible struct
	cIntArray := helpers.IntSliceToCArray(numbers)
	fmt.Printf("Converted to C: %v elements\n", cIntArray.numberOfElements)

	// Convert back to Go slice
	goSlice := helpers.CIntArrayToSlice(cIntArray.data, int(cIntArray.numberOfElements))
	fmt.Printf("Back to Go: %v\n", goSlice)

	// Clean up memory
	helpers.FreeIntArray(cIntArray.data)
	C.free(unsafe.Pointer(cIntArray)) // or helpers.free_int_array_result(cIntArray) if exported
}
```

### API

The go lib has the following API functions:

**Convert C types to go types (internal; Use at entrypoint to Go libraries)**

- `CStringToString(input *C.char) string{}`: Convert a string to a c-compatible C-string (glorified alias for C.GoString)
- `CFloatArrayToSlice(cArray *C.float, length int) []float32{}`: Converts a C array of floats to a slice of floats
- `CIntArrayToSlice(cArray *C.int, length int) []int{}`: Takes a C integer array and coverts it to an integer slice
- `CStringArrayToSlice(cArray **C.char, numberOfStrings int) []string{}`: Takes in an array of strings, and converts it to a slice of strings


**Convert Go types to C types (external; Use to prep data to return to C)**

- `StringToCString(data string) *C.char{}`: Convert a string to a c-compatible C-string (glorified alias for C.CString)
- `StringSliceToCArray(data []string) *C.StringArrayResult{}`: Return dynamically sized string array as a C-Compatible array
- `IntSliceToCArray(data []int) *C.IntArrayResult{}`: Return dynamically sized int array as a C-Compatible array
- `FloatSliceToCArray(data []float32) *C.FloatArrayResult{}`: Return dynamically float sized array as a C-Compatible array

**Memory Freeing**

- `FreeCString(data *C.char){}`: Free's a C-string
- `FreeStringArray(inputArray **C.char, count C.int){}`: Free's an array of strings
- `FreeIntArray(ptr *C.int){}`: Free's an array of integers
- `FreeFloatArray(ptr *C.float){}`: Free's an array of floats

**Debugging Functions**

- `return_string(data *C.char) *C.char{}`: Used to convert a C-compatible string to a C-compatible string, useful for debugging encoding issues
- `return_string_array(cArray **C.char, numberOfStrings int) *C.StringArrayResult{}`: Used to convert a C-compatible string array to wrapper type
- `return_int_array(cArray *C.int, numberOfElements C.int) *C.IntArrayResult{}`: Used to convert a C-compatible integer array to wrapper type
- `return_float_array(cArray *C.float, numberOfElements C.int) *C.FloatArrayResult{}`: Used to convert a C-compatible float array to wrapper type
- `print_string(ptr *C.char){}`: Prints the go representation of a C string, good for debugging encoding issues
- `print_string_array(cArray **C.char, numberOfString int){}`: Prints the go representation of an array, good for debugging encoding issues
- `print_int_array(cArray *C.int, numberOfInts int){}`: Prints the go representation of an array, good for debugging rounding/conversion issues
- `print_float_array(cArray *C.float, numberOfFloats int){}`: Prints the go representation of an array, good for debugging rounding/conversion issues

### Tests

To run the tests use: 

```bash
go test
```
