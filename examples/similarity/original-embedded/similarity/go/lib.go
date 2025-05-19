// A library to help calculate the similarity between two strings
//
// # Functions
//
//	CheckSimilarity(): Check the similarity of a word to a corpus of validWords
//	LoadWords(): Returns a slice of a huge number of words (>350,000)
//
// # DLL/Shared Lib functions
//
//	check_dictionary_similarity(): Compare a string to a large corpus of real words
//	check_dictionary_similarity_levenstein(): Compare a string to a large corpus of words w/Levenstein distance
package main

/*
#include <stdlib.h>

typedef struct{
	char* word;
	float likelihood;
} Suggestion;
*/
import "C"
import (
	_ "embed"
	"fmt"
	"lib/algorithms"
	"strings"
	"time"
	"unsafe"
)

//go:embed words.txt
var wordsContent string

// Helper function to load a default corpus of over 350,000 words
//
// # Returns
//
//	[]string: A slice with the words in the corpus
func LoadWords() []string {
	result := strings.Split(wordsContent, "\r\n")
	return result
}

// Utility function that does the heavy lifting, essentially a general solution that will take in any SimilarityAlgorithm
//
// # Parameters
//
//	inputWord (*C.char): The word to find a similar word for
//	algorithm (algorithms.SimilarityAlgorithm): The algorithm to use to calculate the similarity of the words
//	validWords ([]string): A slice with the words in the corpus
//
// # Returns
//
//	Suggestion: A pointer to the suggestion struct with the word and it's likelihood
func CheckSimilarity(word string, algorithm algorithms.SimilarityAlgorithm, validWords []string) algorithms.Suggestion {
	result := algorithms.SuggestWord(word, validWords, algorithm)
	return result
}

// Utility function that does the heavy lifting for the DLL side of the API
//
// # Parameters
//
//	inputWord (*C.char): The word to find a similar word for
//	algorithm (algorithms.SimilarityAlgorithm): The algorithm to use to calculate the similarity of the words
//	validWords ([]string): A slice with the words in the corpus
//
// # Returns
//
//	*C.Suggestion: A pointer to the suggestion object with the word and it's likelihood
func cCheckSimilarity(inputWord *C.char, algorithm algorithms.SimilarityAlgorithm, validWords []string) *C.Suggestion {
	word := C.GoString(inputWord)
	r := algorithms.SuggestWord(word, validWords, algorithm)

	// allocate the struct
	memoryFootprint := unsafe.Sizeof(C.Suggestion{})
	CMemoryFootprint := C.size_t(memoryFootprint)
	result := (*C.Suggestion)(C.malloc(CMemoryFootprint))

	*result = C.Suggestion{
		word:       C.CString(r.Word),
		likelihood: C.float(r.Likelihood),
	}
	return result
}

// Checks the similarity of the word with one in the corpus based on the indel similarity
//
// # Notes
//   - This uses indel similarity, making it the fastest, but not necesssarily most thorough option
//
// # Parameters
//
//	inputWord (*C.char): The word to find a similar word for
//
// # Returns
//
//	*C.Suggestion: A pointer to the suggestion object with the word and it's likelihood
//
//export check_dictionary_similarity
func check_dictionary_similarity(inputWord *C.char) *C.Suggestion {
	words := LoadWords()
	return cCheckSimilarity(inputWord, algorithms.IndelSimilarity, words)
}

// Checks the similarity of the word with one in the corpus based on full levenstein distance
//
// # Notes
//
//	This is typically much slower than check_dictionary_similarity()
//
// # Parameters
//
//	inputWord (*C.char): The word to find a similar word for
//
// # Returns
//
//	*C.Suggestion: A pointer to the suggestion object with the word and it's likelihood
//
//export check_dictionary_similarity_levenstein
func check_dictionary_similarity_levenstein(inputWord *C.char) *C.Suggestion {
	words := LoadWords()
	return cCheckSimilarity(inputWord, algorithms.LevensteinSimilarity, words)
}

//export free_suggestion
func free_suggestion(suggestionReference *C.Suggestion) {
	if suggestionReference == nil {
		return
	}
	C.free(unsafe.Pointer(suggestionReference.word))
	C.free(unsafe.Pointer(suggestionReference))
}

func main() {

	word := "almni"
	t1 := time.Now()
	words := LoadWords()
	result := algorithms.SuggestWord(word, words, algorithms.IndelSimilarity)
	if !(1 > len(result.Word)) {
		fmt.Printf("\nThe most likely word in the group is %q with %.2f likelihood\n", result.Word, result.Likelihood)
	}

	t2 := time.Now()
	fmt.Printf("The functions took %d miliseconds", t2.UnixMilli()-t1.UnixMilli())

	t1 = time.Now()
	result = algorithms.SuggestWord(word, words, algorithms.JaroSimilarity)
	if !(1 > len(result.Word)) {
		fmt.Printf("\nThe most likely word in the group is %q with %.2f likelihood\n", result.Word, result.Likelihood)
	}

	t2 = time.Now()
	fmt.Printf("The functions took %d miliseconds", t2.UnixMilli()-t1.UnixMilli())
}
