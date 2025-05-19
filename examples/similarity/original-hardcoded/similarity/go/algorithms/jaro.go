package algorithms

// This file implements the Jaro Similarity of two strings
//
// # References
//  - https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance
//  - https://www.baseclass.io/newsletter/jaro-winkler
//  - https://pypi.org/project/jarowinkler/
//  - https://www.geeksforgeeks.org/jaro-and-jaro-winkler-similarity/

import "math"

// Calculates the Jaro similarity between two strings
//
// The Jaro similarity metric measures the similarity between two strings.
// It is especially effective for short strings such as names, and accounts for character matches and transpositions.
// The score ranges from 0.0 (no similarity) to 1.0 (exact match)
//
// # Parameters
//  inputString (string): The first string for comparison
//  targetString (string): The second string for comparison
//
// # Returns
//  float32: A value between 0 and 1 representing the Jaro similarity score
func JaroSimilarity(inputString, targetString string) float32 {
	// If the strings are equal
	if inputString == targetString {
		return 1.0
	}

	inputStringLength := len(inputString)
	targetStringLength := len(targetString)

	// How far to consider a letter a match (half the longest string - 1)
	max_match_distance := math.Floor(float64(max(inputStringLength, targetStringLength))/2.0) - 1

	matches := 0 // How many matches

	// setup empty matrices
	inputStringMatrix := make([]int, inputStringLength)
	targetStringMatrix := make([]int, targetStringLength)

	// Get number of matches
	for i := range inputStringLength {
		// Set ranges to iterate
		start := max(0, i-int(max_match_distance))
		end := min(targetStringLength, i+int(max_match_distance)+1)

		// Look for matches within range of current letter
		for j := start; j < end; j++ {
			if inputString[i] == targetString[j] && targetStringMatrix[j] == 0 {
				// Match found, update matrices and match counter
				inputStringMatrix[i] = 1
				targetStringMatrix[j] = 1
				matches += 1
				break
			}
		}
	}

	// No matches, so transpositions aren't possible
	if matches < 1 {
		return 0.0
	}
	transpositions := calculateTranspositions(inputString, targetString, inputStringMatrix, targetStringMatrix)

	// 1/3 * ((m/s1)+(m/s2)+((m-t)/m)) SEE: https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance#Jaro_similarity
	return ((float32(matches) / float32(inputStringLength)) +
		(float32(matches) / float32(targetStringLength)) +
		((float32(matches) - float32(transpositions)) / float32(matches))) /
		3.0
}

// Calculates the number of transpositions between two strings
//
// A transposition is when two characters are in the wrong order. This function assumes
// that the matches between the two strings are already known and marked using arrays
// (inputStringMatrix and targetStringMatrix), where non-zero values indicate matches
//
// # Parameters
//  inputString (string): The first string for comparison
//  targetString (string): The second string for comparison
//  inputStringMatrix ([]int): Array indicating matching characters in the input string (non-zero for matches)
//  targetStringMatrix ([]int): Array indicating matching characters in the target string (non-zero for matches)
//
// # Returns
//  int: The number of transpositions (half the number of character mismatches in matched positions)
func calculateTranspositions(inputString, targetString string, inputStringMatrix, targetStringMatrix []int) int {
	inputStringLength := len(inputString)

	transpositions := 0 // Count of mismatched matched characters
	marker := 0         // Marker to track position in target string

	for i := range inputStringLength {
		// Only check matched characters
		if inputStringMatrix[i] > 0 {
			// Move the target string pointer to the next matched character
			for targetStringMatrix[marker] == 0 {
				marker += 1
			}
			// Count the mismatch if characters are different at the matched position
			if inputString[i] != targetString[marker] {
				transpositions += 1
			}
			// Move to the next character in the target string
			marker += 1
		}
	}
	// Each transposition involves two characters, so divide the count by 2
	return transpositions / 2
}
