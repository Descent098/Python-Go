# Similarity

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
similarity/
├─ 📂original-embedded/
├─ 📂original-hardcodded/
└──📂with-helper/
```

- `📂original-embedded`: The original implementation without the helper library that uses embedding (will not compile on some systems)
- `📂original-hardcodded`: The original implementation without the helper library that uses hardcoded dictionaries
- `📂with-helper`: The version with the helper library


