from ctypes import cdll, c_int
from platform import platform

# import library
if platform().lower().startswith("windows"):
    lib = cdll.LoadLibrary("./lib.dll")
else:
    lib = cdll.LoadLibrary("./lib.so") 
    
# Simple idempotent function call
lib.Greeting()

# Variadic function (with arguments/returns)
lib.factorial.argtypes = [c_int]
lib.factorial.restype = c_int

n = 10

print(f"The factorial of {n} is {lib.factorial(n)} {type(lib.factorial(n))}")
