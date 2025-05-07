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
		memorySizeOfStruct := unsafe.Sizeof(C.User{})     // Size of a single struct

		// Calculate where to put the current struct
		startPoint := unsafe.Pointer(locationOfArray + offsetIntoArray*memorySizeOfStruct)

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
