// # Functions
//
// # Convert C types to go types (internal; Use at entrypoint to Go libraries)
//
//	CStringToString(input *C.char) string{} //Convert a string to a c-compatible C-string (glorified alias for C.GoString)
//	CFloatArrayToSlice(cArray *C.float, length int) []float32{} // Converts a C array of floats to a slice of floats
//	CIntArrayToSlice(cArray *C.int, length int) []int{} // Takes a C integer array and coverts it to an integer slice
//	CStringArrayToSlice(cArray **C.char, numberOfStrings int) []string{} // Takes in an array of strings, and converts it to a slice of strings
//
// # Convert Go types to C types (external; Use to prep data to return to C)
//
//	StringToCString(data string) *C.char{} // Convert a string to a c-compatible C-string (glorified alias for C.CString)
//	StringSliceToCArray(data []string) *C.StringArrayResult{} // Return dynamically sized string array as a C-Compatible array
//	IntSliceToCArray(data []int) *C.IntArrayResult{} // Return dynamically sized int array as a C-Compatible array
//	FloatSliceToCArray(data []float32) *C.FloatArrayResult{} // Return dynamically float sized array as a C-Compatible array
//
// # Memory Freeing
//
//	FreeCString(data *C.char){} // Free's a C-string
//	FreeStringArray(inputArray **C.char, count C.int){} // Free's an array of strings
//	FreeIntArray(ptr *C.int){}  // Free's an array of integers
//	FreeFloatArray(ptr *C.float){} // Free's an array of floats
//
// # Debugging Functions
//
//	return_string(data *C.char) *C.char{} // Used to convert a C-compatible string to a C-compatible string, useful for debugging encoding issues
//	return_string_array(cArray **C.char, numberOfStrings int) *C.StringArrayResult{} // Used to convert a C-compatible string array to wrapper type
//	return_int_array(cArray *C.int, numberOfElements C.int) *C.IntArrayResult{} // Used to convert a C-compatible integer array to wrapper type
//	return_float_array(cArray *C.float, numberOfElements C.int) *C.FloatArrayResult{} // Used to convert a C-compatible float array to wrapper type
//	print_string(ptr *C.char){} // Prints the go representation of a C string, good for debugging encoding issues
//	print_string_array(cArray **C.char, numberOfString int){} // Prints the go representation of an array, good for debugging encoding issues
//	print_int_array(cArray *C.int, numberOfInts int){} // Prints the go representation of an array, good for debugging rounding/conversion issues
//	print_float_array(cArray *C.float, numberOfFloats int){} // Prints the go representation of an array, good for debugging rounding/conversion issues
//
// # Examples
//
// Create a function to show the internal go representation of an array of strings
//
//	 // Takes in a C string array, prints the go representation, then returns it
//	 //export print_string_array
//	 func print_string_array(cArray **C.char, numberOfStrings int) *C.StringArrayResult {
//		  internalRepresentation := CStringArrayToSlice(cArray, numberOfStrings)
//		  fmt.Printf("return_string_array() Go representation: %v\n", internalRepresentation)
//
//		  result := StringSliceToCArray(internalRepresentation)
//
//		  return result
//	 }
package main

/*
#include <stdlib.h>

typedef struct{
	int numberOfElements;
	char** data;
} StringArrayResult;

typedef struct {
    int numberOfElements;
    int* data;
} IntArrayResult;

typedef struct {
    int numberOfElements;
    float* data;
} FloatArrayResult;

*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ======== Convert Go types to C type ========

// Convert a string to a c-compatible C-string (glorified alias for C.CString)
func StringToCString(input string) *C.char {
	return C.CString(input)
}

// A function to take a slice and convert it to a StringArrayResult to be returned to C code
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
func StringSliceToCArray(data []string) *C.StringArrayResult {
	count := len(data)

	// Allocate memory for an array of C string pointers (char**)
	amountOfElements := C.size_t(count)
	sizeOfSingleElement := C.size_t(unsafe.Sizeof(uintptr(0)))
	amountOfMemory := amountOfElements * sizeOfSingleElement
	stringArray := (**C.char)(C.malloc(amountOfMemory))

	// Create Array of data
	for i, currentString := range data {
		// Calculate where to put the string
		locationOfArray := uintptr(unsafe.Pointer(stringArray)) // Starting point of first byte of slice
		offsetIntoArray := uintptr(i)                           // The offset for the current element
		sizeOfSingleElement := unsafe.Sizeof(uintptr(0))        // Size of a single string

		locationInMemory := (**C.char)(unsafe.Pointer(locationOfArray + offsetIntoArray*sizeOfSingleElement))
		*locationInMemory = C.CString(currentString) // Convert go string to C string and insert at location in array

	}

	// Allocate memory for the struct
	result := (*C.StringArrayResult)(C.malloc(C.size_t(unsafe.Sizeof(C.StringArrayResult{}))))
	result.numberOfElements = C.int(count)
	result.data = stringArray

	return result
}

// Return dynamically sized int array as a C-Compatible array
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
func IntSliceToCArray(data []int) *C.IntArrayResult {
	count := len(data)

	// Allocate memory in C for the int array
	amountOfMemory := C.size_t(count) * C.size_t(unsafe.Sizeof(C.int(0)))
	cArray := (*C.int)(C.malloc(amountOfMemory))

	// Fill in the values
	array := (*[1 << 30]C.int)(unsafe.Pointer(cArray))
	for i, val := range data {
		array[i] = C.int(val)
	}

	// Allocate the result struct
	result := (*C.IntArrayResult)(C.malloc(C.size_t(unsafe.Sizeof(C.IntArrayResult{}))))
	result.numberOfElements = C.int(count)
	result.data = cArray

	return result
}

// Return dynamically float sized array as a C-Compatible array
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
func FloatSliceToCArray(data []float32) *C.FloatArrayResult {
	count := len(data)

	// Allocate memory in C for the float array
	amountOfMemory := C.size_t(count) * C.size_t(unsafe.Sizeof(C.float(0)))
	cArray := (*C.float)(C.malloc(amountOfMemory))

	// Fill in the values
	array := (*[1 << 30]C.float)(unsafe.Pointer(cArray))
	for i, val := range data {
		array[i] = C.float(val)
	}

	// Allocate the result struct
	result := (*C.FloatArrayResult)(C.malloc(C.size_t(unsafe.Sizeof(C.FloatArrayResult{}))))
	result.numberOfElements = C.int(count)
	result.data = cArray

	return result
}

// ======== Convert C types to Go ========

// Convert a string to a c-compatible C-string (glorified alias for C.GoString)
func CStringToString(input *C.char) string {
	return C.GoString(input)
}

// Takes a C integer array and coverts it to an integer slice
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
//
// # Usage
//
//	var cIntArray *C.int // Assuming it's set in some line after this
//	goInts := CIntArrayToSlice(cIntArray, length).([]int)
func CIntArrayToSlice(cArray *C.int, length int) []int {
	// Setup buffer for array contents
	const bufferSize = 1 << 30 // Allocate a huge buffer
	slice := (*[bufferSize]C.int)(unsafe.Pointer(cArray))[:length:length]

	// Convert to []int
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = int(slice[i])
	}
	return result
}

// Converts a C array of floats to a slice of floats
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
//
// # Usage
//
//	var cFloatArray  *C.float // Assuming it's set in some line after this
//	goFloats := CFloatArrayToSlice(cFloatArray, length).([]float32)
func CFloatArrayToSlice(cArray *C.float, length int) []float32 {
	// Setup buffer for array contents
	const bufferSize = 1 << 30 // Allocate a huge buffer
	slice := (*[bufferSize]C.float)(unsafe.Pointer(cArray))[:length:length]

	// Convert to []float32
	result := make([]float32, length)
	for i := 0; i < length; i++ {
		result[i] = float32(slice[i])
	}
	return result
}

// Takes in an array of strings, and converts it to a slice of strings
// C array -> slice of strings
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
func CStringArrayToSlice(cArray **C.char, numberOfStrings int) []string {
	// Setup buffer for array contents
	const bufferSize = 1 << 30 // Allocate a huge buffer
	stringPointers := (*[bufferSize]*C.char)(unsafe.Pointer(cArray))[:numberOfStrings:numberOfStrings]

	result := make([]string, 0, numberOfStrings)
	for i := range numberOfStrings {
		result = append(result, C.GoString(stringPointers[i]))
	}
	return result
}

// Converts a C array of floats or ints to a slice of floats or ints
//
// # Notes
//
//   - This function DOES NOT clean memory of input array, that's up to others to clear
//
// # Usage
//
//	var cFloatArray  *C.float // Assuming it's set in some line after this
//	goFloats := NumberArrayToSlice(cFloatArray, length).([]float32)
func NumberArrayToSlice(ptr interface{}, length int) interface{} {
	if ptr == nil {
		return nil
	}

	switch p := ptr.(type) {
	case *C.int:
		return CIntArrayToSlice(p, length)

	case *C.float:
		return CFloatArrayToSlice(p, length)
	default:
		panic(fmt.Sprintf("NumberArrayToSlice(): unsupported type: %T", ptr))
	}
}

// ========== Debugging Functions ==========

// Used to convert a C-compatible string back to itself, good for debugging encoding issues
//
// # Notes
//
//   - Does not free any memory
//
//export return_string
func return_string(cString *C.char) *C.char {

	internalRepresentation := C.GoString(cString)
	result := StringToCString(internalRepresentation)
	return result
}

// Used to convert a C-compatible string array to wrapper type
//
// # Notes
//
//   - Does not free any memory
//
//export return_string_array
func return_string_array(cArray **C.char, numberOfStrings int) *C.StringArrayResult {

	internalRepresentation := CStringArrayToSlice(cArray, numberOfStrings)

	result := StringSliceToCArray(internalRepresentation)

	return result
}

// Used to convert a C-compatible integer array to wrapper type
//
// # Notes
//
//   - Does not free any memory
//
//export return_int_array
func return_int_array(cArray *C.int, numberOfElements C.int) *C.IntArrayResult {
	internalRepresentation := CIntArrayToSlice(cArray, int(numberOfElements))
	result := IntSliceToCArray(internalRepresentation)
	return result
}

// Used to convert a C-compatible float array to wrapper type
//
// # Notes
//
//   - Does not free any memory
//
//export return_float_array
func return_float_array(cArray *C.float, numberOfElements C.int) *C.FloatArrayResult {
	internalRepresentation := CFloatArrayToSlice(cArray, int(numberOfElements))
	result := FloatSliceToCArray(internalRepresentation)
	return result
}

// Prints the go representation of a C string, good for debugging encoding issues
//
//export print_string
func print_string(ptr *C.char) {
	if ptr != nil {
		fmt.Printf("print_string() Go representation: %s\n", C.GoString(ptr))
	} else {
		fmt.Println("print_string() received nil pointer")
	}
}

// Prints the go representation of an array, good for debugging encoding issues
//
//export print_string_array
func print_string_array(cArray **C.char, numberOfString int) {
	res := CStringArrayToSlice(cArray, numberOfString)
	fmt.Printf("print_string_array() Go representation: %v\n", res)
}

// Prints the go representation of an array, good for debugging rounding/conversion issues
//
//export print_int_array
func print_int_array(cArray *C.int, numberOfInts int) {
	res := CIntArrayToSlice(cArray, numberOfInts)

	fmt.Printf("print_int_array() Go representation: %v\n", res)
}

// Prints the go representation of an array, good for debugging rounding/conversion issues
//
//export print_float_array
func print_float_array(cArray *C.float, numberOfFloats int) {
	res := CFloatArrayToSlice(cArray, numberOfFloats)

	fmt.Printf("print_float_array() Go representation: %v\n", res)
}

// ========== Functions to free memory ==========

//export FreeCString
func FreeCString(ptr *C.char) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

//export FreeStringArray
func FreeStringArray(inputArray **C.char, count C.int) {
	for i := 0; i < int(count); i++ {
		// Calculate where to find the string
		locationOfArray := uintptr(unsafe.Pointer(inputArray)) // Starting point of first byte of slice
		offsetIntoArray := uintptr(i)                          // The offset for the current element
		memorySizeOfStruct := unsafe.Sizeof(uintptr(0))        // Size of a single struct

		ptr := *(**C.char)(unsafe.Pointer(locationOfArray + offsetIntoArray*memorySizeOfStruct))
		C.free(unsafe.Pointer(ptr))
	}
	C.free(unsafe.Pointer(inputArray))
}

//export FreeIntArray
func FreeIntArray(ptr *C.int) {
	C.free(unsafe.Pointer(ptr))
}

//export FreeFloatArray
func FreeFloatArray(ptr *C.float) {
	C.free(unsafe.Pointer(ptr))
}

//export free_string_array_result
func free_string_array_result(StringArrayResultReference *C.StringArrayResult) {
	FreeStringArray(StringArrayResultReference.data, StringArrayResultReference.numberOfElements)
	C.free(unsafe.Pointer(StringArrayResultReference))
}

//export free_int_array_result
func free_int_array_result(ptr *C.IntArrayResult) {
	C.free(unsafe.Pointer(ptr.data))
	C.free(unsafe.Pointer(ptr))
}

//export free_float_array_result
func free_float_array_result(ptr *C.FloatArrayResult) {
	C.free(unsafe.Pointer(ptr.data))
	C.free(unsafe.Pointer(ptr))
}

func main() {}
