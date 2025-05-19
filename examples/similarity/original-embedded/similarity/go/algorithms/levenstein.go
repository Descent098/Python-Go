package algorithms

// Calculates the Levenstein similarity of two strings
//
// # Parameters
//
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  float32: The similarity (between 0-1, closer to 1 is more similar)
func LevensteinSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, LevensteinDistance)
	return similarity
}

// Calculates the Levenstein distance of two strings
//
// # Notes
//  - This solution utilizes the dynamic programming approach, not the recursive one
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenstein distance (add, edit, delete distance)
func LevensteinDistance(inputString, targetString string) int {
	return DynamicLevenshtein(inputString, targetString)
}

// Calculates the Levenstein distance of two strings recursively
//
// # Notes
//  - Heavily inspired by the recursive haskel implementation on wikipedia https://en.wikipedia.org/wiki/Levenshtein_distance#Recursive
//  - Very slow, roughly O(3^n)
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenstein distance (add, edit, delete distance)
func RecursiveLevenstein(inputString, targetString string) int {
	if len(inputString) == 0 {
		return len(targetString)
	}
	if len(targetString) == 0 {
		return len(inputString)
	}

	firstInputChar := inputString[0]
	firstTargetChar := targetString[0]
	restInputString := inputString[1:]
	restTargetString := targetString[1:]

	if firstInputChar == firstTargetChar {
		return LevensteinDistance(restInputString, restTargetString)
	} else {
		return 1 + min(
			LevensteinDistance(inputString, restTargetString),     // Add
			LevensteinDistance(restInputString, targetString),     // Delete
			LevensteinDistance(restInputString, restTargetString), // Edit/replace
		)
	}
}

// A dynamic-programming based implementation of Levenstein distance
//
// # Notes
//  - Relies on Wagner–Fischer algorithm https://en.wikipedia.org/wiki/Wagner%E2%80%93Fischer_algorithm#Calculating_distance
//  - More details: https://gist.github.com/Descent098/401c2ca6bdf3fa655738e7a1ddf1aeee
//  - Faster than the recursive solution, runs in roughly O(m*n) where m and n are the size of strings
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenstein distance (add, edit, delete distance)
func DynamicLevenshtein(inputString, targetString string) int {
	// Convert to runes to avoid weird encoding issues
	inputStringRunes := []rune(inputString)
	targetStringRunes := []rune(targetString)

	inputStringLength := len(inputStringRunes)
	targetStringLength := len(targetStringRunes)

	// Create a 2D matrix
	matrix := make([][]int, inputStringLength+1)
	for i := range matrix {
		matrix[i] = make([]int, targetStringLength+1)
	}

	// Initialize base cases
	for i := 0; i <= inputStringLength; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= targetStringLength; j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= inputStringLength; i++ {
		for j := 1; j <= targetStringLength; j++ {
			if inputStringRunes[i-1] == targetStringRunes[j-1] {
				// Characters match, no cost added
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = 1 + min(
					matrix[i][j-1],   // Add
					matrix[i-1][j],   // Delete
					matrix[i-1][j-1], // Edit/replace
				)
			}
		}
	}

	return matrix[inputStringLength][targetStringLength]
}
