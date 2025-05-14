"""A package to help with building Go-python libraries

Converting to ctypes
--------------------
- prepare_string(data: str | bytes) -> c_char_p: Takes in a string and returns a C-compatible string
- prepare_string_array(data:list[str|bytes]) -> tuple[Array[c_char_p], int]: Takes in a string list, and converts it to a C-compatible array
- prepare_int_array(data:list[int]) -> tuple[Array[c_int], int]: Takes in a int list, and converts it to a C-compatible array
- prepare_float_array(data:list[float]) -> tuple[Array[c_float], int]: Takes in a float list, and converts it to a C-compatible array

Converting from ctypes
----------------------
- string_to_str(pointer: c_char_p) -> str: Takes in a pointer to a C string and returns a Python string
- string_array_result_to_list(pointer:_CStringArrayResult) -> list[str]: 
- int_array_result_to_list(pointer: _CIntArrayResult) -> list[int]: 
- float_array_result_to_list(pointer: _CFloatArrayResult) -> list[float]: 

Debugging Functions
-------------------
- return_string(text: str | bytes) -> str: Debugging function that shows you the Go representation of a C string and returns the python string version
- return_string_array(c_array:CStringArray, number_of_elements:int) ->list[str]: Debugging function that shows you the Go representation of a C array and returns the python list version (does not free)
- return_int_array(c_array: CIntArray, number_of_elements: int) -> list[int]: Debugging function that shows you the Go representation of a C int array and returns a Python list
- return_float_array(c_array: CFloatArray, number_of_elements: int) -> list[float]: Debugging function that shows you the Go representation of a C float array and returns a Python list
- print_string(text: str | bytes): Prints a string's go representation, useful to look for encoding issues
- print_string_array(data:list[str|bytes]): Prints a string array's go representation, useful to look for encoding issues
- print_int_array(data:list[int]): Prints a int array's go representation, useful to look for rounding/conversion issues
- print_float_array(data:list[float]): Prints a float array's go representation, useful to look for rounding/conversion issues

Freeing Functions
-----------------
- free_c_string(ptr: c_char_p): Frees a single C string returned from Go (allocated via C.CString).
- free_string_array(ptr: CStringArray, count: int): Frees an array of C strings returned from Go.
- free_int_array(ptr: CIntArray): Frees a C int array returned from Go.
- free_float_array(ptr: CFloatArray): Frees a C float array returned from Go.
- free_string_array_result(ptr: _CStringArrayResult): Frees a StringArrayResult (including the array of strings and struct itself).
- free_int_array_result(ptr: _CIntArrayResult): Frees an IntArrayResult (including the array and the struct itself).
- free_float_array_result(ptr: _CFloatArrayResult): Frees a FloatArrayResult (including the array and the struct itself).
"""
import os
from platform import platform

# Check if dynamic library is compiled
if platform().lower().startswith("windows"):
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "lib.dll")
else:
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "lib.so")
    
if not os.path.exists(lib_path):
    # Try to build dynamic library since it's not present
    import subprocess
    lib_folder_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go")
    
    if platform().lower().startswith("windows"):
        filename = "lib.dll"
        additional_flags = "set GOTRACEBACK=system &&"
    else:
        filename = "lib.so"
        additional_flags = "env GOTRACEBACK=system"
    command = f"{additional_flags} go build -ldflags \"-s -w\" -buildmode=c-shared -o \"{os.path.join(lib_folder_path, filename)}\" \"{os.path.join(lib_folder_path, 'lib.go')}\""
    print("\nRequired shared library is not available, building...")
    try:
        subprocess.run(command, shell=True, check=True)
    except Exception as e:
        if isinstance(e, FileNotFoundError):
            print("Unable to find Go install, plkease install it and try again\n")
        else:
            print(f"Ran into error while trying to build shared library, make sure go, and a compatible compiler are installed, then try building manually using:\n\t{command}\nExiting with error:\n\t{e}")
        exit(-1)
    
# Exported functions
from .lib import (
    prepare_string,
    prepare_string_array,
    prepare_int_array,
    prepare_float_array,
    string_array_result_to_list,
    int_array_result_to_list,
    float_array_result_to_list,
    return_string,
    return_string_array,
    return_int_array,
    return_float_array,
    print_string,
    print_string_array,
    print_int_array,
    print_float_array,
    free_c_string,
    free_string_array,
    free_int_array,
    free_float_array,
    free_string_array_result,
    free_int_array_result,
    free_float_array_result,
)