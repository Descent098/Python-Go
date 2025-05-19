"""A package to help with building Go-python libraries"""
import os
import subprocess
from platform import platform
from ctypes import CDLL, Array, cdll, c_char_p, c_int, POINTER, c_float, Structure, string_at 

# ========== Helper Functions  ============
def get_library(dll_path:str,source_path:str="", compile:bool=False) -> CDLL:
    """Get's the DLL specified, will compile if not found and flag is specified

    Parameters
    ----------
    dll_path : str
        The path to the DLL file, if compile is specified this will be the output path

    source_path:str, optional
        The path to the source go file, only needed if compile is true, by default ""

    compile : bool, optional
        Specify if you should try to compile DLL if not in path, by default False

    Raises
    ------
    ValueError:
        If linked library is not available and/or compilable (if compile is specified)

    Returns
    -------
    CDLL
        The linked library
    """
    if not os.path.exists(dll_path):
        if not compile:
            raise ValueError(f"Linked Library is not available: {dll_path}")
        if platform().lower().startswith("windows"):
            additional_flags = "set GOTRACEBACK=system &&"
        else:
            additional_flags = "env GOTRACEBACK=system"
        command = f"{additional_flags} go build -ldflags \"-s -w\" -buildmode=c-shared -o \"{dll_path}\" \"{source_path}\""
        if compile:
            print("\nRequired shared library is not available, building...")
            try:
                subprocess.run(command, shell=True, check=True)
            except Exception as e:
                if isinstance(e, FileNotFoundError):
                    print("Unable to find Go install, please install it and try again\n")
                else:
                    print(f"Ran into error while trying to build shared library, make sure go, and a compatible compiler are installed, then try building manually using:\n\t{command}\nExiting with error:\n\t{e}")
                raise ValueError(f"Linked Library is not available or compileable: {dll_path}")
    return cdll.LoadLibrary(dll_path)

# ========== C Structs ==========
class _CStringArrayResult(Structure):
    _fields_ = [
        ("numberOfElements", c_int),
        ("data", POINTER(c_char_p)),
    ]
    
class _CIntArrayResult(Structure):
    _fields_ = [
        ("numberOfElements", c_int),
        ("data", POINTER(c_int)),
    ]

class _CFloatArrayResult(Structure):
    _fields_ = [
        ("numberOfElements", c_int),
        ("data", POINTER(c_float)),
    ]

# ========== Setup CGo functions ==========

# import library
dll_source_file = os.path.join(os.path.dirname(os.path.realpath(__file__)), "lib.go")
if platform().lower().startswith("windows"):
    dll_file = os.path.join(os.path.dirname(os.path.realpath(__file__)),"lib.dll")
    lib = get_library(dll_file, dll_source_file, True)
else:
    dll_file = os.path.join(os.path.dirname(os.path.realpath(__file__)),"lib.so")
    lib = get_library(dll_file, dll_source_file, True)

lib.print_string_array.argtypes =  [POINTER(c_char_p), c_int]
lib.FreeStringArray.argtypes = [POINTER(c_char_p), c_int]

lib.print_int_array.argtypes =  [POINTER(c_int), c_int]
lib.FreeIntArray.argtypes = [POINTER(c_int)]

lib.print_float_array.argtypes =  [POINTER(c_float), c_int]
lib.FreeFloatArray.argtypes =  [POINTER(c_float)]

lib.return_string.argtypes = [c_char_p]
lib.return_string.restype = c_char_p

lib.FreeCString.argtypes = [c_char_p]

lib.print_string.argtypes = [c_char_p]

## ========== Array-based functions ==========

lib.FreeStringArray.argtypes = [POINTER(c_char_p), c_int]
lib.free_string_array_result.argtypes = [POINTER(_CStringArrayResult)]
lib.return_string_array.argtypes = [POINTER(c_char_p), c_int] 
lib.return_string_array.restype = POINTER(_CStringArrayResult)

lib.return_int_array.argtypes = [POINTER(c_int), c_int]
lib.return_int_array.restype = POINTER(_CIntArrayResult)
lib.free_int_array_result.argtypes = [POINTER(_CIntArrayResult)]

lib.return_float_array.argtypes = [POINTER(c_float), c_int]
lib.return_float_array.restype = POINTER(_CFloatArrayResult)
lib.free_float_array_result.argtypes = [POINTER(_CFloatArrayResult)]

# ========== Nice Typehints/Type Aliases ==========
CIntArray = Array[c_int]
CFloatArray = Array[c_float]
CStringArray = Array[c_char_p]

# ========== Python types to C ============
def prepare_string(data: str | bytes) -> c_char_p:
    """Takes in a string and returns a C-compatible string
    
    Notes
    -----
    - Does not prune null terminators (\\0 characters)

    Parameters
    ----------
    data : str | bytes
        The string to prepare

    Returns
    -------
    c_char_p
        The resulting pointer to the string
    """
    if not data:
        return c_char_p(bytes())
    if type(data) == str:
        return c_char_p(data.encode())
    return c_char_p(bytes(data))

def prepare_string_array(data:list[str|bytes]) -> tuple[CStringArray, int]:
    """Takes in a string list, and converts it to a C-compatible array

    Parameters
    ----------
    data : list[str | bytes]
        The list to convert

    Returns
    -------
    Array[c_char_p], int
        The resulting array, and the number of items
        
    Notes
    -----
    - Because the data is allocated in python, python will free the memory afterwords
    - Does not prune null terminators (\\0 characters)

    Examples
    --------
    ```
    lib = cdll.LoadLibrary("path/to/library.dll") # Load Library

    # Function that takes in string array, and number of items, then prints them in C
    lib.print_string_array.argtypes =  [POINTER(c_char_p), c_int]

    # Prep data using function
    data = ["Hello", "World", "!"]
    c_array, number_of_items = prepare_string_array(data)

    # Use data in C
    lib.print_string_array(c_array, number_of_items)
    ```
    """
    # Encode items to bytes
    data = [
            c_char_p(item.encode())
        if type(item) == str
        else
            c_char_p(bytes(item))
        for item in data
    ] 
    number_of_items = len(data)
    array_type = c_char_p * number_of_items # Create a C array of char* (aka **char)
    c_array = array_type(*data)
    return c_array, number_of_items

def prepare_int_array(data:list[int]) -> tuple[CIntArray, int]:
    """Takes in an int list, and converts it to a C-compatible array

    Parameters
    ----------
    data : list[int]
        The list of integers to convert to an array

    Returns
    -------
    Array[c_int], int
        The resulting array, and the number of items
    
    Notes
    -----
    - Because the data is allocated in python, python will free the memory afterwords
        
    Examples
    --------
    ```
    lib = cdll.LoadLibrary("path/to/library.dll") # Load Library

    # Function that takes in int array, and number of items, then prints them in C
    lib.print_int_array.argtypes =  [POINTER(c_int), c_int]

    # Prep data using function
    data = [1,2,3,4]
    c_array, number_of_items = prepare_int_array(data)

    # Use data in C
    lib.print_int_array(c_array, number_of_items)
    ```
    """
    data = [c_int(item) for item in data] # Force an error if wrong type
    number_of_items = len(data)
    array_type = c_int * number_of_items # Create a C array of int*
    c_array = array_type(*data)
    return c_array, number_of_items

def prepare_float_array(data:list[float]) -> tuple[CFloatArray, int]:
    """Takes in an float list, and converts it to a C-compatible array

    Parameters
    ----------
    data : list[float]
        The list of integers to convert to an array

    Returns
    -------
    Array[c_float], int
        The resulting array, and the number of items
    
    Notes
    -----
    - Because the data is allocated in python, python will free the memory afterwords
    - The data is only accurate up to ~4 decimals (i.e. if value is -790.5207366698761 you might get -790.520751953125)
        
    Examples
    --------
    ```
    lib = cdll.LoadLibrary("path/to/library.dll") # Load Library

    # Function that takes in float array, and number of items, then prints them in C
    lib.print_float_array.argtypes =  [POINTER(c_float), c_int]

    # Prep data using function
    data = [1.0,2.604,3.14159,4.964]
    c_array, number_of_items = prepare_float_array(data)

    # Use data in C
    lib.print_float_array(c_array, number_of_items)
    ```
    """
    data = [c_float(item) for item in data]  # Force an error if wrong type
    number_of_items = len(data)
    array_type = c_float * number_of_items # Create a C array of float*
    c_array = array_type(*data)
    return c_array, number_of_items

# ========== Convert C types to python ============
def string_to_str(pointer: c_char_p) -> str:
    """Takes in a pointer to a C string and returns a Python string

    Parameters
    ----------
    pointer : c_char_p
        A C-style string pointer returned from Go

    Notes
    -----
    - Assumes the pointer is a valid null-terminated UTF-8 encoded string
    - Does NOT free the pointer automatically, you must call `lib.FreeCString(pointer)` if needed

    Returns
    -------
    str
        The decoded Python string representation of the C string
        
    Examples
    --------
    ```
    c_str = prepare_string(b"Hello from Python!")
    result: str = string_to_str(c_str)
    lib.FreeCString(c_str)
    ```
    """
    if pointer:
        return pointer.value.decode("utf-8", errors="replace")
    return ""

def string_array_result_to_list(pointer:_CStringArrayResult) -> list[str]:
    """Takes in a pointer to a string result and returns a list of strings

    Parameters
    ----------
    pointer : _CStringArrayResult
        A pointer to a CString Result

    Notes
    -----
    - free's the original pointer

    Returns
    -------
    list[str]
        The list of strings the pointer pointed to
        
    Examples
    --------
    ```
    data = [
        random.choice(["Lorem", "ipsum", "dolor", "sit", "amet"]) 
        for _ in range(100)
    ]
    
    c_array, number_of_elements = prepare_string_array(data)
    
    pointer = return_string_array(c_array, number_of_elements)
    
    result:list[str] = string_array_result_to_list(pointer)
    ```
    """
    try:
        result_data = pointer.contents
        results = []
        for i in range(result_data.numberOfElements):
            results.append(result_data.data[i].decode(errors='replace'))
        return results
    finally:
        lib.free_string_array_result(pointer)

def int_array_result_to_list(pointer: _CIntArrayResult) -> list[int]:
    """Converts C int result struct to a Python list, and frees memory."""
    try:
        result_data = pointer.contents
        return [result_data.data[i] for i in range(result_data.numberOfElements)]
    finally:
        lib.free_int_array_result(pointer)

def float_array_result_to_list(pointer: _CFloatArrayResult) -> list[float]:
    """Converts C float result struct to a Python list, and frees memory."""
    try:
        result_data = pointer.contents
        return [result_data.data[i] for i in range(result_data.numberOfElements)]
    finally:
        lib.free_float_array_result(pointer)

# ========== Debugging Functions ==========

def return_string(text: str | bytes) -> str:
    """Debugging function that shows you the Go representation of a C string and returns the python string version

    Parameters
    ----------
    text : str | bytes
        The text to get the representation of

    Returns
    -------
    str
        The returned string
    """
    c_input = prepare_string(text)
    result = lib.return_string(c_input)

    if not result:
        return ""

    copied_bytes = string_at(result)
    decoded = copied_bytes.decode(errors="replace")

    return decoded

def return_string_array(c_array:CStringArray, number_of_elements:int) ->list[str]:
    """Debugging function that shows you the Go representation of a C array and returns the python list version

    Parameters
    ----------
    c_array : Array[c_char_p]
        The array to print and convert
    number_of_elements : int
        The number of elements in the array

    Notes
    -----
    - DOES NOT FREE INPUT ARRAY
    - This function returns the PYTHON list version, do not reassign input variable or it'll never free (i.e. c_array = return_string_array(c_array, number_of_elements))

    Returns
    -------
    list[str]
        The python string representation of the array
        
    Examples
    --------
    ```
    data = [
        random.choice(["Lorem", "ipsum", "dolor", "sit", "amet"]) 
        for _ in range(100)
    ]
    
    c_array, number_of_elements = prepare_string_array(data)
    
    result:list[str] = return_string_array(c_array, number_of_elements)
    
    lib.free_string_array_result(c_array, number_of_elements)
    ```
    """
    pointer = lib.return_string_array(c_array, number_of_elements)

    result_data = pointer.contents
    results = []
    for i in range(result_data.numberOfElements):
        results.append(result_data.data[i].decode(errors='replace'))
    return results

def return_int_array(c_array: CIntArray, number_of_elements: int) -> list[int]:
    """Debugging function that shows you the Go representation of a C int array and returns a Python list

    Notes
    -----
    - DOES NOT FREE INPUT ARRAY
    - Frees input array ONLY on exception
    - Returns the PYTHON list version, do not reassign input variable

    Returns
    -------
    list[int]
    """
    pointer = lib.return_int_array(c_array, number_of_elements)
    try:
        result_data = pointer.contents
        return [result_data.data[i] for i in range(result_data.numberOfElements)]
    except Exception as e:
        print(f"return_int_array(): Ran into error, freeing memory. Error: {e}")
        lib.free_int_array_result(c_array)  # In case you define a similar freeing function for input
        raise e
    finally:
        lib.free_int_array_result(pointer)

def return_float_array(c_array: CFloatArray, number_of_elements: int) -> list[float]:
    """Debugging function that shows you the Go representation of a C float array and returns a Python list

    Notes
    -----
    - DOES NOT FREE INPUT ARRAY ON SUCCESS
    - Frees input array ONLY on exception
    - Returns the PYTHON list version, do not reassign input variable
    - The data is only accurate up to ~4 decimals (i.e. if value is -790.5207366698761 you might get -790.520751953125)

    Returns
    -------
    list[float]
    """
    pointer = lib.return_float_array(c_array, number_of_elements)
    try:
        result_data = pointer.contents
        return [result_data.data[i] for i in range(result_data.numberOfElements)]
    except Exception as e:
        print(f"return_float_array(): Ran into error, freeing memory. Error: {e}")
        lib.free_float_array_result(c_array)  # In case you define a similar freeing function for input
        raise e
    finally:
        lib.free_float_array_result(pointer)

def print_string(text: str | bytes):
    """Prints a string's go representation, useful to look for encoding issues

    Parameters
    ----------
    text : str | bytes
        The data you want to see the go representation of
    """
    c_input = prepare_string(text)
    lib.print_string(c_input)

def print_string_array(data:list[str|bytes]):
    """Prints a string array's go representation, useful to look for encoding issues

    Notes
    -----
    - Does not free because everything is allocated in python, so GC will take care of it

    Parameters
    ----------
    data : list[str | bytes]
        The data you want to see the go representation of
    """
    c_array, number_of_items = prepare_string_array(data)

    lib.print_string_array(c_array, number_of_items)

def print_int_array(data:list[int]):
    """Prints a int array's go representation, useful to look for rounding/conversion issues

    Notes
    -----
    - Does not free because everything is allocated in python, so GC will take care of it

    Parameters
    ----------
    data : list[int]
        The data you want to see the go representation of
    """
    c_array, number_of_items = prepare_int_array(data)

    lib.print_int_array(c_array, number_of_items)

def print_float_array(data:list[float]):
    """Prints a float array's go representation, useful to look for rounding/conversion issues

    Notes
    -----
    - Does not free because everything is allocated in python, so GC will take care of it

    Parameters
    ----------
    data : list[float]
        The data you want to see the go representation of
    """
    c_array, number_of_items = prepare_float_array(data)
    lib.print_float_array(c_array, number_of_items)

# ========== Free Functions ==========
def free_c_string(ptr: c_char_p):
    """Frees a single C string returned from Go (allocated via C.CString)."""
    lib.FreeCString(ptr)

def free_string_array(ptr: CStringArray, count: int):
    """Frees an array of C strings returned from Go."""
    lib.FreeStringArray(ptr, count)

def free_int_array(ptr: CIntArray):
    """Frees a C int array returned from Go."""
    lib.FreeIntArray(ptr)
    
def free_float_array(ptr: CFloatArray):
    """Frees a C float array returned from Go."""
    lib.FreeFloatArray(ptr)

def free_string_array_result(ptr: _CStringArrayResult):
    """Frees a StringArrayResult (including the array of strings and struct itself)."""
    lib.free_string_array_result(ptr)

def free_int_array_result(ptr: _CIntArrayResult):
    """Frees an IntArrayResult (including the array and the struct itself)."""
    lib.free_int_array_result(ptr)

def free_float_array_result(ptr: _CFloatArrayResult):
    """Frees a FloatArrayResult (including the array and the struct itself)."""
    lib.free_float_array_result(ptr)
