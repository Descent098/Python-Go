# Hardcoded Version

This version takes the content of a ~370,000 word dictionary, and injects it as a hardcoded string slice into the code. This approach was taken to get around limitations/security on windows.

## Background

This demo will allow you to take in a string, compare it to a corpus/dictionary of valid words, and tell you which word it is most likely the user was trying to type. API wise it's usable via python with:

```python
from similarity import spellcheck

incorrect_word = "almni"

suggested_word, likelihood = spellcheck(incorrect_word)

# Prints: almni is likely amni with a likelihood of %88.88888955116272
print(f"{incorrect_word} is likely {suggested_word} with a likelihood of %{likelihood}")

real_word = "move"

suggested_word, likelihood = spellcheck(real_word)

# Prints: move is likely move with a likelihood of %0.0
print(f"{real_word} is likely {suggested_word} with a likelihood of %{likelihood}")
```

## Running

You should be able to run by just running `testing.py`, if you have your go and c compiler setup it will compile the lib and run it for you, or if it fails it will give you the command(s) to run.

## Folder Structure

Here is the folder structure for this folder, each version will have details about it's implementation in the README:

```
├─ 📂similarity/
|   ├─ 📂go/
|   |   ├─ 📂algorithms/
|   |   |   ├─ 📄indel.go
|   |   |   ├─ 📄jaro.go
|   |   |   ├─ 📄levenstein.go
|   |   |   └──📄utilities.go
|   |   ├─ 📄go.mod
|   |   └──📄lib.go
|   ├─ 📄__init__.py
|   ├─ 📄user_library.py
|   └──📄words.txt
├─ 📄indel.py
└──📄testing.py
```

- `📄indel.py`: A pure python implementation of the same functionality
- `📄testing.py`: A version in python that runs the go code
- `📂similarity/`: The folder that contains the library code that makes the API function
- `📂similarity/📄words.txt`: The dictionary of valid words
- `📂similarity/📄__init__.py`:  The file that sets up the python package
- `📂similarity/📄user_library.py`: The python package/libary that people can use (via `import similarity`)
- `📂similarity/📂go/`: The go side of the library
- `📂similarity/📂go/📄go.mod`: The file that lists dependencies and allows compilation
- `📂similarity/📂go/📄lib.go`: The main entrypoint file that stitches together everything on the go side
- `📂similarity/📂go/📂algorithms/`: The various algorithms that were implemented for the package
- `📂similarity/📂go/📂algorithms/📄indel.go`: The implementation of [indel distance](https://en.wikipedia.org/wiki/Indel) and similarity (the default)
- `📂similarity/📂go/📂algorithms/📄jaro.go`: The implementation of [jaro similarity](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance)
- `📂similarity/📂go/📂algorithms/📄levenstein.go`: The implementation of [levenstein distance](https://en.wikipedia.org/wiki/Levenshtein_distance) and similarity
- `📂similarity/📂go/📂algorithms/📄utilities.go`: Utilities to help generalize and build the library
