import os
import sys
import random
from platform import platform
from ctypes import ArgumentError, cdll, c_char_p, c_int, POINTER, c_float
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

from lib import *
from lib import _CStringArrayResult, _CIntArrayResult, _CFloatArrayResult

import pytest

# import dll library
if platform().lower().startswith("windows"):
    lib = cdll.LoadLibrary(os.path.join(os.path.dirname(os.path.realpath(__file__)), "lib.dll"))
else:
    lib = cdll.LoadLibrary(os.path.join(os.path.dirname(os.path.realpath(__file__)), "lib.so")) 

# Setup CGo functions
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

## Array-based functions

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

def cstring_checks(correct_content:str, data_to_test:c_char_p):
    """Checks that a c string is setup correctly"""
    assert data_to_test is not None # NULL check
    assert type(data_to_test) == c_char_p
    if correct_content: # Strings are each other (i.e. '' is '')
        assert correct_content is not data_to_test # Should not be the SAME object
        assert correct_content != data_to_test # data_to_test should be a pointer, and correct_content is a str
    assert correct_content == data_to_test.value.decode(errors="replace") 
    
def string_checks(correct_content:str, data_to_test:str):
    """Checks that a c string is setup correctly"""
    assert data_to_test is not None # NULL check
    assert type(data_to_test) == str
    assert correct_content == data_to_test

def test_python_to_c_functions():
    # Test Valid input for prepare_string
    ## Testing basic strings
    for test_input in ("","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"):
        result = prepare_string(test_input)
        cstring_checks(test_input, result)
    
    ## Testing bytes types
    for test_input in (b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8", errors="replace"), b"\x41", b"\n"):
        result = prepare_string(test_input)
        cstring_checks(test_input.decode(errors="replace"), result)
    
    # Test invalid input for prepare_string
    ## Only strings with null-termination characters in them should error
    for test_input in ("\0\0", "as\0\nasdf","\n\n\n\0", b"\0\0", b"as\0\nasdf"):
        result = prepare_string(test_input)
        with pytest.raises(AssertionError) as e:
            cstring_checks(test_input, result)
    
    # Test Valid input for prepare_string_array
    ## Testing normal strings
    for test_input in (
        ["","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"],
        ["Here", "are", "some", "other", "strings"],
        ["This", "is", "getting", "TeDiOuS"],
        ["Reeeee"],
        []):
        result_array, result_number_of_items = prepare_string_array(test_input)
        assert len(test_input) == result_number_of_items
        converted_result_array = [c_char_p(result_array[i]) for i in range(result_number_of_items)]
        for current_test_input_string, result_current_string in zip(test_input, converted_result_array):
            if current_test_input_string:
                cstring_checks(current_test_input_string, result_current_string)
            else:
                # '' is '' in python 
                assert current_test_input_string is result_current_string.value.decode(errors="replace")

    ## Testing Binary Strings
    for test_input in (
        [b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8"), b"\x41", b"\n"],
        [b"Here", b"are", b"some", b"other", b"strings"],
        [b"This", b"is", b"getting", b"TeDiOuS"],
        [b"Reeeee"],
        []):
        result_array, result_number_of_items = prepare_string_array(test_input)
        assert len(test_input) == result_number_of_items
        converted_result_array = [c_char_p(result_array[i]) for i in range(result_number_of_items)]
        for current_test_input_string, result_current_string in zip(test_input, converted_result_array):
            if current_test_input_string:
                cstring_checks(current_test_input_string.decode(errors="replace"), result_current_string)
            else:
                # '' is '' in python 
                assert current_test_input_string is result_current_string.value

    # Test invalid input for prepare_string_array
    for test_input in (["\0\0", "as\0\nasdf","\n\n\n\0", b"\0\0", b"as\0\nasdf"],):
        result_array, result_number_of_items = prepare_string_array(test_input)
        assert len(test_input) == result_number_of_items
        converted_result_array = [c_char_p(result_array[i]) for i in range(result_number_of_items)]
        for current_test_input_string, result_current_string in zip(test_input, converted_result_array):
            if current_test_input_string:
                with pytest.raises(AssertionError) as e:
                    cstring_checks(current_test_input_string, result_current_string)
            else:
                # '' is '' in python 
                assert current_test_input_string is result_current_string.value
    
    # Test input for prepare_int_array
    n = 1000
    test_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_int_array(test_input)
    assert len(test_input) == number_of_items
    for i in range(n):
        assert test_input[i] == c_array[i]
    
    # Test input for prepare_float_array
    n = 1000
    test_input = [random.uniform(-1000.0, 1000.0) for _ in range(n)]
    c_array, number_of_items = prepare_float_array(test_input)
    assert len(test_input) == number_of_items
    for i in range(n):
        if abs(test_input[i] - c_array[i]) > 1e-4: # Test to 4 decimal places of accuracy
            raise AssertionError(f"Values did not match: {original_float} != {returned_float}")   
    
def test_internal_lib_functions():
    """Testing functions used by the exposed API to make sure they function correctly"""
    # Testing return_string_array()/_CStringArrayResult
    ## Testing normal strings
    for test_input in (
        ["","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"],
        [
            random.choice(["Lorem", "ipsum", "dolor", "sit", "amet"]) 
            for _ in range(100)
        ],
        ["Here", "are", "some", "other", "strings"],
        ["This", "is", "getting", "TeDiOuS"],
        ["Reeeee"],
        []):
        result_c_array, result_number_of_items = prepare_string_array(test_input)
        result:_CStringArrayResult = lib.return_string_array(result_c_array, result_number_of_items)

        temp_data = result.contents
        assert len(test_input) == result_number_of_items == temp_data.numberOfElements
        for i in range(len(test_input)):
            string_checks(test_input[i], temp_data.data[i].decode(errors="replace"))
        free_string_array_result(result)

    ## Testing Binary Strings
    for test_input in (
        [b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8"), b"\x41", b"\n"],
        [
            random.choice([b"Lorem", b"ipsum", b"dolor", b"sit", b"amet"]) 
            for _ in range(100)
        ],
        [b"Here", b"are", b"some", b"other", b"strings"],
        [b"This", b"is", b"getting", b"TeDiOuS"],
        [b"Reeeee"],
        []):
        result_c_array, result_number_of_items = prepare_string_array(test_input)
        result:_CStringArrayResult = lib.return_string_array(result_c_array, result_number_of_items)

        temp_data = result.contents
        assert len(test_input) == result_number_of_items == temp_data.numberOfElements
        for i in range(len(test_input)):
            string_checks(test_input[i].decode(errors="replace"), temp_data.data[i].decode(errors="replace"))
        free_string_array_result(result)
    
    # Testing return_int_array()/_CIntArrayResult
    ## Valid input
    n = 1000
    original_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_int_array(original_input)
    test_input = lib.return_int_array(c_array, number_of_items)
    test_input_data = test_input.contents
    assert test_input_data.numberOfElements == number_of_items == len(original_input)
    for i in range(n):
        assert original_input[i] == test_input_data.data[i]
    free_int_array_result(test_input)
    
    ## Invalid input
    test_input = ["A"]
    c_array, number_of_items = prepare_string_array(test_input)
    with pytest.raises(ArgumentError):
        lib.return_int_array(c_array, number_of_items)
        return_int_array(c_array, number_of_items)
    
    
    # Testing return_float_array()/_CFloatArrayResult
    n = 1000
    original_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_float_array(original_input)
    test_input = lib.return_float_array(c_array, number_of_items)
    test_input_data = test_input.contents
    assert test_input_data.numberOfElements == number_of_items == len(original_input)
    for i in range(n):
        if abs(original_input[i] - test_input_data.data[i]) > 1e-4: # Test to 4 decimal places of accuracy
            raise AssertionError(f"Values did not match: {original_float} != {returned_float}")
    free_float_array_result(test_input)

def test_go_to_python_functions():
    # Test Valid input for string_to_str
    ## Testing basic strings
    for test_input in ("","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"):
        temp = prepare_string(test_input)
        result = string_to_str(temp)
        string_checks(test_input, result)
    
    ## Testing bytes types
    for test_input in (b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8", errors="replace"), b"\x41", b"\n"):
        temp = prepare_string(test_input)
        result = string_to_str(temp)
        string_checks(test_input.decode(errors="replace"), result)
    
    # Test Valid input for string_array_result_to_list
    ## Testing normal strings
    for test_input in (
        ["","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"],
        [
            random.choice(["Lorem", "ipsum", "dolor", "sit", "amet"]) 
            for _ in range(100)
        ],
        ["Here", "are", "some", "other", "strings"],
        ["This", "is", "getting", "TeDiOuS"],
        ["Reeeee"],
        []):
        result_c_array, result_number_of_items = prepare_string_array(test_input)
        intermediate:_CStringArrayResult = lib.return_string_array(result_c_array, result_number_of_items)
        result = string_array_result_to_list(intermediate)

        assert len(test_input) == result_number_of_items == len(result)
        for input_string, result_string in zip(test_input, result):
            string_checks(input_string, result_string)

    ## Testing Binary Strings
    for test_input in (
        [b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8"), b"\x41", b"\n"],
        [
            random.choice([b"Lorem", b"ipsum", b"dolor", b"sit", b"amet"]) 
            for _ in range(100)
        ],
        [b"Here", b"are", b"some", b"other", b"strings"],
        [b"This", b"is", b"getting", b"TeDiOuS"],
        [b"Reeeee"],
        []):
        result_c_array, result_number_of_items = prepare_string_array(test_input)
        intermediate:_CStringArrayResult = lib.return_string_array(result_c_array, result_number_of_items)
        result = string_array_result_to_list(intermediate)

        assert len(test_input) == result_number_of_items == len(result)
        for input_string, result_string in zip(test_input, result):
            string_checks(input_string.decode(errors="replace"), result_string)
    
    
    # Testing int_array_result_to_list()
    n = 1000
    original_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_int_array(original_input)
    temp = lib.return_int_array(c_array, number_of_items)
    
    test_input = int_array_result_to_list(temp)
    assert len(test_input) == number_of_items == len(original_input)
    for original_number, test_number in zip(original_input, test_input):
        assert original_number == test_number
    
    # Testing float_array_result_to_list()
    n = 1000
    original_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_float_array(original_input)
    temp = lib.return_float_array(c_array, number_of_items)
    test_input = float_array_result_to_list(temp)
    assert len(test_input) == number_of_items == len(original_input)
    for original_number, test_number in zip(original_input, test_input):
        assert original_number == test_number
    
    for original_number, test_number in zip(original_input, test_input):
        if abs(original_number - test_number) > 1e-4: # Test to 4 decimal places of accuracy
            raise AssertionError(f"Values did not match: {original_float} != {returned_float}")


def test_debugging_functions(capsys:pytest.CaptureFixture[str]):
    # Test Valid input for return_string
    ## Testing basic strings
    for test_input in ("","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"):
        result = return_string(test_input)
        string_checks(test_input, result)
    ## Testing bytes types
    for test_input in (b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8", errors="replace"), b"\x41", b"\n"):
        result = return_string(test_input)
        string_checks(test_input.decode(errors="replace"), result)
    # Test invalid input for prepare_string
    ## Only strings with null-termination characters in them should error
    for test_input in ("\0\0", "as\0\nasdf","\n\n\n\0", b"\0\0", b"as\0\nasdf"):
        result = prepare_string(test_input)
        with pytest.raises(AssertionError) as e:
            string_checks(test_input, result)


    # Test Valid input for return_string_array
    ## Testing normal strings
    for test_input in (
        ["","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"],
        ["Here", "are", "some", "other", "strings"],
        ["This", "is", "getting", "TeDiOuS"],
        ["Reeeee"],
        []):
        result_array, result_number_of_items = prepare_string_array(test_input)
        converted_result_array = return_string_array(result_array, result_number_of_items)
        
        assert len(converted_result_array) == result_number_of_items == len(test_input)
        
        for current_test_input_string, result_current_string in zip(test_input, converted_result_array):
            if current_test_input_string:
                string_checks(current_test_input_string, result_current_string)
            else:
                # '' is '' in python 
                assert current_test_input_string is result_current_string

    ## Testing Binary Strings
    for test_input in (
        [b"",b"Hello World!", b"!@$#^%!#@@%*!", b"AWDsadfSA", bytes("\u2764", encoding="utf-8"), b"\x41", b"\n"],
        [b"Here", b"are", b"some", b"other", b"strings"],
        [b"This", b"is", b"getting", b"TeDiOuS"],
        [b"Reeeee"],
        []):
        result_array, result_number_of_items = prepare_string_array(test_input)
        converted_result_array = return_string_array(result_array, result_number_of_items)
        
        assert len(converted_result_array) == result_number_of_items == len(test_input)
        
        for current_test_input_string, result_current_string in zip(test_input, converted_result_array):
            if current_test_input_string:
                string_checks(current_test_input_string.decode(errors="replace"), result_current_string)
            else:
                # '' is '' in python 
                assert current_test_input_string.decode(errors="replace") is result_current_string

    # Test input for return_int_array
    n = 1000
    test_input = [random.randint(-1000, 1000) for _ in range(n)]
    c_array, number_of_items = prepare_int_array(test_input)
    converted_result_array = return_int_array(c_array, number_of_items)
    assert len(test_input) == number_of_items
    for test_num, result_num in zip(test_input, converted_result_array):
        assert test_num == result_num
        
        
    # Test input for return_float_array
    n = 1000
    test_input = [random.uniform(-1000.00, 1000.00) for _ in range(n)]
    c_array, number_of_items = prepare_float_array(test_input)
    converted_result_array = return_float_array(c_array, number_of_items)
    assert len(test_input) == number_of_items
    for test_num, result_num in zip(test_input, converted_result_array):
        if abs(test_num - result_num) > 1e-4: # Test to 4 decimal places of accuracy
            raise AssertionError(f"Values did not match: {original_float} != {returned_float}") 
        
    # Test print_string
    for test_input_array in (
        ["","Hello World!", "!@$#^%!#@@%*!", "AWDsadfSA", "\u2764", "\x41", "\n"],
        ["Here", "are", "some", "other", "strings"],
        ["This", "is", "getting", "TeDiOuS"],
        ["Reeeee"],
        []):
        for test_input in test_input_array:
            print_string(test_input)
            print_string_array(test_input)
            # Cannot capture pytest, so this is just making sure the function doesn't crash
            # https://docs.pytest.org/en/7.1.x/how-to/capture-stdout-stderr.html does not work
    
    n = 1000
    test_input = [random.randint(-1000, 1000) for _ in range(n)]
    print_int_array(test_input)
    
    n = 1000
    test_input = [random.uniform(-1000.00, 1000.00) for _ in range(n)]
    print_float_array(test_input)
