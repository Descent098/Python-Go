package main

// Tests for the helper lib, keep in mind debugging
// Functions are not tested because their constituent
// Functions are used by those tested

import (
	"math/rand/v2"
	"testing"
	"unsafe"
)

func TestStringConversions(t *testing.T) {
	// StringToCString <--> CStringToString
	for _, test_input := range []string{"", "Hello World", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"} {
		r := StringToCString(test_input)
		defer FreeCString(r)
		temp := CStringToString(r)

		if !(temp == test_input) {
			t.Errorf(`TestStringConversions:StringToCString("%s"): %s!=%s\n`, test_input, test_input, temp)
		}
	}

	// StringSliceToCArray <--> CStringArrayToSlice
	for _, test_input := range [][]string{
		{"", "Hello World", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"},
		{"Here", "are", "some", "other", "strings"},
		{"This", "is", "getting", "TeDiOuS"},
		{"Reeeee"},
		{},
	} {
		r := StringSliceToCArray(test_input)
		defer free_string_array_result(unsafe.Pointer(r))

		temp := CStringArrayToSlice(unsafe.Pointer(r.data), int(r.numberOfElements))

		for i := range len(test_input) {
			if !(temp[i] == test_input[i]) {
				t.Errorf(`TestStringConversions:StringSliceToCArray("%s"): %s!=%s\nTestStringConversions:StringSliceToCArray("%s"): %s!=%s\n`, test_input[i], test_input[i], temp[i], test_input, test_input, temp)
			}
		}

	}
}

func TestNumberConversions(t *testing.T) {
	// IntSliceToCArray <--> CIntArrayToSlice
	for range 10 {
		test_input := make([]int, 100)
		for i := range 100 {
			n := rand.IntN(10_000)
			modifier := rand.IntN(1)
			if modifier == 0 {
				modifier = -1
			}

			test_input[i] = n * modifier
		}
		r := IntSliceToCArray(test_input)
		defer free_int_array_result(unsafe.Pointer(r))
		temp := CIntArrayToSlice(unsafe.Pointer(r.data), int(r.numberOfElements))

		for i := range len(test_input) {
			if !(temp[i] == test_input[i]) {
				t.Errorf(`TestNumberConversions:IntSliceToCArray("%d"): %d!=%d\TestNumberConversions:IntSliceToCArray("%v"): %v!=%v\n`, test_input[i], test_input[i], temp[i], test_input, test_input, temp)
			}
		}
	}

	// FloatSliceToCArray <--> CFloatArrayToSlice

	for range 10 {
		test_input := make([]float32, 100)
		for i := range 100 {
			n := rand.Float32() * float32(rand.IntN(10_000))
			modifier := rand.IntN(1)
			if modifier == 0 {
				modifier = -1
			}

			test_input[i] = n * float32(modifier)
		}
		r := FloatSliceToCArray(test_input)
		defer free_float_array_result(unsafe.Pointer(r))
		temp := CFloatArrayToSlice(unsafe.Pointer(r.data), int(r.numberOfElements))

		for i := range len(test_input) {
			if !(temp[i] == test_input[i]) {
				t.Errorf(`TestNumberConversions:IntSliceToCArray("%f"): %f!=%f\TestNumberConversions:IntSliceToCArray("%v"): %v!=%v\n`, test_input[i], test_input[i], temp[i], test_input, test_input, temp)
			}
		}
	}
}
