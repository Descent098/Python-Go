# Pure python implementation of indel similarity

def indel_similarity(inputString:str, targetString:str) -> float:
    """Calculates the indel similarity of two strings

    Parameters
    ----------
    inputString : str
        The string you want to compare
    targetString : str
        The string you want to compare to

    Returns
    -------
    float
        The similarity (a normalized weight of the distance)
    """
    inputLength, targetLength = len(inputString), len(targetString)

    matrix = [[0 for _ in range(targetLength + 1)] for _ in range(inputLength + 1)]

    for inputIndex in range(inputLength + 1):
        matrix[inputIndex][0] = inputIndex
    for targetIndex in range(targetLength + 1):
        matrix[0][targetIndex] = targetIndex
    
    for inputIndex in range(1, inputLength + 1):
        for targetIndex in range(1, targetLength + 1):
            if inputString[inputIndex - 1] == targetString[targetIndex - 1]:
                matrix[inputIndex][targetIndex] = matrix[inputIndex - 1][targetIndex - 1]
            else:
              matrix[inputIndex][targetIndex] = min(matrix[inputIndex - 1][targetIndex] + 1, matrix[inputIndex][targetIndex - 1] + 1)
    
    distance = matrix[inputLength][targetLength]
    normalized_distance = distance / (inputLength + targetLength)
    similarity = 1 - normalized_distance
    return similarity

def suggestWord(input_string:str, valid_words:list[str]) -> list[str, float]:
    """Takes in a word and suggests a word from the list of valid words

    Parameters
    ----------
    input_string : str
        The string to compare against valid words
    valid_words : list[str]
        The words that are considered valid

    Returns
    -------
    list[str, float]
        The suggested word, and the likelihood
    """
    highest = 0.0
    suggested_word = ""
    for word in valid_words:
        similarity = indel_similarity(input_string, word)
        if similarity > highest:
            highest = similarity
            suggested_word = word
    return suggested_word, highest

def load_words() -> list[str]:
    """Load a set of words from a text file

    Returns
    -------
    list[str]
        The list of words
    """
    res = []
    with open("words.txt", "r") as f:
        for line in f.read().split():
            res.append(line)
    return res

if __name__ == "__main__":
    import time
    t1 = time.time()
    word = "almni"
    words = load_words()
    print(f"Number of words loaded: {len(words)}")
    suggested_word, highest = suggestWord(word, words)
    print(f"The suggested word for {word} is {suggested_word} with a ratio of {highest}")
    t2 = time.time()
    print(f"Took {(t2-t1)*1000:.3f} miliseconds")
    
