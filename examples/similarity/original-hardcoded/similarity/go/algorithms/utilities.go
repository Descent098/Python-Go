package algorithms

type Suggestion struct {
	Likelihood float32
	Word       string
}

type DistanceAlgorithm func(inputString, targetString string) int
type SimilarityAlgorithm func(inputString, targetString string) float32

// Function that calculates the similarity of two strings using a distance algortithm
//
// # Parameters
//
//	inputString (string): The first string to use for the comparison
//	targetString (string): The second string to use for the comparison
//	algorithm (DistanceAlgorithm): The algorithm to use to calculate the distance
//
// # Returns
//
//	float32: The similarity (between 0-1, closer to 1 is more similar)
func CalculateSimilarity(inputString, targetString string, algorithm DistanceAlgorithm) float32 {
	inputLength := len(inputString)
	targetLength := len(targetString)

	// Get the distance
	distance := algorithm(inputString, targetString)

	// Normalize your distance across the lengths of the inputs
	normalized_distance := float32(distance) / (float32(inputLength) + float32(targetLength))

	// Get the final similarity and return it
	similarity := 1 - normalized_distance
	return similarity
}

// Function that suggests the highest similarity word to the input string
//
// # Parameters
//
//	inputString (string): The first string to use for the comparison
//	validStrings ([]string): The valid words to check against
//	algorithm (SimilarityAlgorithm): The algorithm to run and generate the similarity for
//
// # Returns
//
//	float32: The similarity (between 0-1, closer to 1 is more similar)
func SuggestWord(inputString string, validStrings []string, algorithm SimilarityAlgorithm) Suggestion {
	var (
		highestRatio float32
		result       string
	)

	for _, currentString := range validStrings {
		likelihood := algorithm(inputString, currentString)
		if likelihood > highestRatio {
			highestRatio = likelihood
			result = currentString
		}
	}

	return Suggestion{highestRatio, result}
}

// Function that suggests the highest similarity word to the input string
//
// # Parameters
//
//	inputString (string): The first string to use for the comparison
//	validStrings ([]string): The valid words to check against
//	threshold (float32): The minimum threshold of likelihood the result needs to be at
//	algorithm (SimilarityAlgorithm): The algorithm to run and generate the similarity for
//
// # Returns
//
//	string: The most likely word, will be a blank string if no likely word was found over the threshold
func SuggestWordWithThreshold(inputString string, validStrings []string, threshold float32, algorithm SimilarityAlgorithm) string {
	suggested := SuggestWord(inputString, validStrings, algorithm)

	if !(suggested.Likelihood > threshold) {
		return ""
	} else {
		return suggested.Word
	}
}
