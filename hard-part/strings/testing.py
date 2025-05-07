from ctypes import cdll, c_char_p
from platform import platform

# import library
if platform().lower().startswith("windows"):
    lib = cdll.LoadLibrary("./lib.dll")
else:
    lib = cdll.LoadLibrary("./lib.so") 

# Setup functions
lib.greet_string.argtypes = [c_char_p]
lib.greet_string.restype = c_char_p

lib.free_string.argtypes = [c_char_p]
lib.free_string.restype = None

# Manually cleaning up memory
name2 = "Kieran".encode()
result2:str = lib.greet_string(name2)
print(f"The result is: {result2.decode(errors='replace')}")
lib.free_string(result2)

# Class-based approach with __del__()
class EvilString(str):
    def __init__(self:'EvilString', value:bytes) -> 'EvilString':
        self.value = value
        
    def __str__(self) -> str:
        return self.value.decode(errors='replace')
        
    def __del__(self):
        lib.free_string(self.value)

name = "Kieran".encode()

result = EvilString(lib.greet_string(name))
print(f"The result is: {result}")


# Below code errors out 
# name2 = "Kieran2".encode()
# result2:str = str(EvilString(lib.greet_string(name2)))

# print(f"The result is: {result2}")

