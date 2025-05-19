import os
import subprocess
from platform import platform
from ctypes import cdll, c_char_p, Structure, POINTER, c_float

# import library
if platform().lower().startswith("windows"):
    library_location = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "similarity.dll")
else:
    library_location = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "similarity.so")

source_location = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.go")

if not os.path.exists(library_location):
    if platform().lower().startswith("windows"):
        additional_flags = "set GOTRACEBACK=system &&"
    else:
        additional_flags = "env GOTRACEBACK=system"
    
    command = f"{additional_flags} go build -ldflags \"-s -w\" -buildmode=c-shared -o \"{library_location}\""
    print("\nRequired shared library is not available, building...")
    try:
        subprocess.run(command, shell=True, check=True, cwd=os.path.dirname(source_location))
    except Exception as e:
        if isinstance(e, FileNotFoundError):
            print("Unable to find Go install, please install it and try again\n")
        else:
            print(f"Ran into error while trying to build shared library, make sure go, and a compatible compiler are installed, then try building manually using:\n\t{command}\nExiting with error:\n\t{e}")
        
        raise ValueError(f"Linked Library is not available or compileable: {library_location}")
lib = cdll.LoadLibrary(library_location)


# Define the C-compatible User struct in Python
class CSuggestion(Structure):
    _fields_ = [
        ("word", c_char_p),
        ("likelihood", c_float),
    ]

# Setup functions

lib.free_suggestion.argtypes = [POINTER(CSuggestion)]

lib.check_dictionary_similarity.argtypes = [c_char_p]
lib.check_dictionary_similarity.restype = POINTER(CSuggestion)

lib.check_dictionary_similarity_levenstein.argtypes = [c_char_p]
lib.check_dictionary_similarity_levenstein.restype = POINTER(CSuggestion)

def spellcheck(word:str|bytes) -> tuple[str, float]:
    """Takes in a word and returns a suggestion and likelihood

    Parameters
    ----------
    word : str | bytes
        The word to check

    Notes
    -----
    - Likelihood will be 0.0 if word is a valid word

    Returns
    -------
    tuple[str, float]
        [suggestion, likelihood] form where likelihood is a % of how likely (will be 0.0 if word is a valid word)

    Examples
    --------
    ## Invalid word
    ```
    from similarity import spellcheck

    incorrect_word = "almni"

    suggested_word, likelihood = spellcheck(incorrect_word)

    # Prints: almni is likely amni with a likelihood of %88.88888955116272
    print(f"{incorrect_word} is likely {suggested_word} with a likelihood of %{likelihood}")
    ```
    
    ## Valid word
    ```
    from similarity import spellcheck

    real_word = "move"

    suggested_word, likelihood = spellcheck(real_word)

    # Prints: move is likely move with a likelihood of %0.0
    print(f"{real_word} is likely {suggested_word} with a likelihood of %{likelihood}")
    ```

    """
    if type(word) == str:
        word = word.strip().lower().encode()
    pointer = lib.check_dictionary_similarity(word)
    try:
        data = pointer.contents
        temp = data.word.decode(errors="replace")

        # Copy out values before clearing
        result = ""
        for character in temp:
            result += str(character)

        # Copy out 2 digits of accuracy for the float
        _, right = str(data.likelihood).split(".")
        likelihood = float(f"{right[0:2]}.{right[2:]}")
        return result, likelihood
    except Exception as e:
        raise Exception(f"Something went wrong during processing: {e}")
    finally:
        lib.free_suggestion(pointer)
        
    
    

if __name__ == "__main__":
    import time
    t1 = time.time()
    word = "almni"
    suggested_word, likelihood = spellcheck("mve")
    print(f"The suggested word for {word} is {suggested_word} with a likelihood of {likelihood}")
    t2 = time.time()
    print(f"Took {(t2-t1)*1000:.3f} miliseconds")
    
