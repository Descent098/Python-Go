package algorithms

// This file implements the Indel Distance of two strings
//
// # References
//  - https://arxiv.org/abs/2410.09877
//  - https://cs.stackexchange.com/questions/136529/calculating-the-string-similarity-of-an-optimal-alignment
//  - https://cran.r-project.org/web/packages/RapidFuzz/RapidFuzz.pdf
//  - https://arxiv.org/pdf/1909.12877
//  - https://en.wikipedia.org/wiki/Indel

// Calculates the Indel similarity of two strings
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  float32: The indel distance (between 0-1, closer to 1 is more similar)
func IndelSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, IndelDistance)
	return similarity
}

// Calculates the Indel distance of two strings
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The indel distance (edit, delete distance)
func IndelDistance(inputString, targetString string) int {
	inputLength := len(inputString)
	targetLength := len(targetString)

	// Construct initial matrix
	matrix := make([][]int, inputLength+1)
	for i := range matrix {
		matrix[i] = make([]int, targetLength+1)
	}

	for inputIndex := 0; inputIndex <= inputLength; inputIndex++ {
		matrix[inputIndex][0] = inputIndex
	}
	for targetIndex := 0; targetIndex <= targetLength; targetIndex++ {
		matrix[0][targetIndex] = targetIndex
	}

	// Calculate indel matrix
	for i := 1; i < inputLength+1; i++ {
		for j := 1; j < targetLength+1; j++ {
			if inputString[i-1] == targetString[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min((matrix[i-1][j] + 1), (matrix[i][j-1] + 1))
			}
		}
	}

	// Determine distance for the provided input
	distance := matrix[inputLength][targetLength]
	return distance
}
