package main

/*
#include <stdlib.h>

typedef struct{
	char* name;
	int age;
	char* email;
} User;
*/
import "C"
import (
	"unsafe"

	"github.com/brianvoe/gofakeit"
)

func Fib(n int) int {
	if n < 2 {
		return 1
	} else {
		return Fib(n-2) + Fib(n-1)
	}
}

func FibSequence(n int) []int {
	results := make([]int, 0, n)
	for i := 0; i < int(n); i++ {
		results = append(results, Fib(i))
	}
	return results
}

//export fib_sequence
func fib_sequence(n C.int) *C.int {
	results := FibSequence(int(n))

	// Allocate memory for an array of C ints (int*)
	sizeOfArray := C.size_t(n)
	sizeOfEachElement := C.size_t(unsafe.Sizeof(C.int(0)))
	amountOfMemory := sizeOfArray * sizeOfEachElement
	cArray := (*C.int)(C.malloc(amountOfMemory))

	// Create Array of data
	for i, currentNumber := range results {
		// Calculate where to put the string
		locationOfArray := uintptr(unsafe.Pointer(cArray)) // Starting point of first byte of slice
		offsetIntoArray := uintptr(i)                      // The offset for the current element
		sizeOfEachElement := unsafe.Sizeof(C.int(0))       // Size of a single value

		locationInMemory := (*C.int)(unsafe.Pointer(locationOfArray + offsetIntoArray*sizeOfEachElement))
		*locationInMemory = C.int(currentNumber) // Convert go int to C int and insert at location in array

	}
	return cArray
}

//export free_int_array
func free_int_array(array *C.int) {
	C.free(unsafe.Pointer(array))
}

func MultiplyString(inputString string, count int) []string {
	result := make([]string, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, inputString)
	}

	return result
}

//export multiply_string
func multiply_string(inputString *C.char, count C.int) **C.char {
	res := MultiplyString(C.GoString(inputString), int(count))

	// Allocate memory for an array of C string pointers (char**)
	amountOfElements := C.size_t(count)
	sizeOfSingleElement := C.size_t(unsafe.Sizeof(uintptr(0)))
	amountOfMemory := amountOfElements * sizeOfSingleElement
	stringArray := (**C.char)(C.malloc(amountOfMemory))

	// Create Array of data
	for i, currentString := range res {
		// Calculate where to put the string
		locationOfArray := uintptr(unsafe.Pointer(stringArray)) // Starting point of first byte of slice
		offsetIntoArray := uintptr(i)                           // The offset for the current element
		sizeOfSingleElement := unsafe.Sizeof(uintptr(0))        // Size of a single string

		locationInMemory := (**C.char)(unsafe.Pointer(locationOfArray + offsetIntoArray*sizeOfSingleElement))
		*locationInMemory = C.CString(currentString) // Convert go string to C string and insert at location in array

	}
	return stringArray
}

//export free_string_array
func free_string_array(inputArray **C.char, count C.int) {
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

type User struct {
	name  string
	age   int
	email string
}

func createRandomUser() *User {
	return &User{gofakeit.Name(), gofakeit.Number(13, 90), gofakeit.Email()}
}

func CreateRandomUsers(count int) []*User {
	result := make([]*User, count)
	for i := 0; i < count; i++ {
		result[i] = createRandomUser()
	}
	return result
}

//export create_user
func create_user(name *C.char, age C.int, email *C.char) *C.User {
	// allocate the struct
	memoryFootprint := unsafe.Sizeof(C.User{})
	CMemoryFootprint := C.size_t(memoryFootprint)
	user := (*C.User)(C.malloc(CMemoryFootprint))

	// Create the struct and it's poitners
	*user = C.User{
		name:  C.CString(C.GoString(name)),
		age:   age,
		email: C.CString(C.GoString(email)),
	}
	return user
}

//export create_random_user
func create_random_user() *C.User {
	// Create a random go version of the user
	res := createRandomUser()

	// Expand go versions of variables
	goName := res.name
	goEmail := res.email
	age := res.age

	// Create C-compatible versions of variables
	cName := C.CString(goName)
	cEmail := C.CString(goEmail)

	// Allocate necessary memory
	user := (*C.User)(C.malloc(C.size_t(unsafe.Sizeof(C.User{}))))

	// Assign values to freshly created struct
	user.name = cName
	user.age = C.int(age)
	user.email = cEmail

	return user
}

//export create_random_users
func create_random_users(count C.int) *C.User {
	// Create a random go version of the user
	res := CreateRandomUsers(int(count))

	users := (*C.User)(C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.User{}))))

	// Create Array of data
	for i, user := range res {
		locationOfArray := uintptr(unsafe.Pointer(users)) // Starting point of first byte of slice
		offsetIntoArray := uintptr(i)                     // The offset for the current element
		sizeOfSingleStruct := unsafe.Sizeof(C.User{})     // Size of a single struct

		// Calculate where to put the current struct
		startPoint := unsafe.Pointer(locationOfArray + offsetIntoArray*sizeOfSingleStruct)

		// Get pointer location for current struct
		currentUser := (*C.User)(startPoint)

		// Assign values to new C.User struct
		currentUser.name = C.CString(user.name)
		currentUser.age = C.int(user.age)
		currentUser.email = C.CString(user.email)

	}
	return users
}

//export free_user
func free_user(userReference *C.User) {
	C.free(unsafe.Pointer(userReference.name))
	C.free(unsafe.Pointer(userReference.email))
	C.free(unsafe.Pointer(userReference))
}

//export free_users
func free_users(users *C.User, count C.int) {
	for i := range int(count) {
		currentUserPointer := (*C.User)(unsafe.Pointer(uintptr(unsafe.Pointer(users)) + uintptr(i)*unsafe.Sizeof(C.User{})))

		// Clear strings
		C.free(unsafe.Pointer(currentUserPointer.name))
		C.free(unsafe.Pointer(currentUserPointer.email))
	}
	C.free(unsafe.Pointer(users))
}

func main() {
	// Do nothing
}
