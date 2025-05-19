from similarity import spellcheck

incorrect_word = "almni"

suggested_word, likelihood = spellcheck(incorrect_word)

print(f"{incorrect_word} is likely {suggested_word} with a likelihood of %{likelihood}")

real_word = "move"

suggested_word, likelihood = spellcheck(real_word)

print(f"{real_word} is likely {suggested_word} with a likelihood of %{likelihood}")