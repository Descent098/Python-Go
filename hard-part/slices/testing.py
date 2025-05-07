import traceback
from dataclasses import dataclass
from platform import platform
from ctypes import cdll, c_char_p, c_int, POINTER, Structure

# import library
if platform().lower().startswith("windows"):
    lib = cdll.LoadLibrary("./lib.dll")
else:
    lib = cdll.LoadLibrary("./lib.so") 

# Setup functions
lib.fib_sequence.argtypes = [c_int]
lib.fib_sequence.restype = POINTER(c_int)

lib.free_int_array.argtypes = [POINTER(c_int)]

lib.multiply_string.argtypes = [c_char_p, c_int]
lib.multiply_string.restype = POINTER(c_char_p)

lib.free_string_array.argtypes = [POINTER(c_char_p), c_int]

n = 10
ptr = lib.fib_sequence(n)

try:
    results = []
    for i in range(n):
        results.append(ptr[i])
    print(results)
finally:
    lib.free_int_array(ptr)



count = 5
r = lib.multiply_string("Hello".encode(), count)
try:
    result = []
    for i in range(count):
        result.append(r[i].decode(errors='replace'))
    print(result)
    print("GEERE")
finally:
    lib.free_string_array(r, count)


print("GEERE")
# User demo
class CUser(Structure):
    _fields_ = [
        ("name", c_char_p),
        ("age", c_int),
        ("email", c_char_p),
    ]

lib.free_user.argtypes = [POINTER(CUser)]

lib.create_random_user.restype = POINTER(CUser)

lib.create_user.argtypes = [c_char_p, c_int, c_char_p]
lib.create_user.restype = POINTER(CUser)

lib.create_random_users.argtypes = [c_int]
lib.create_random_users.restype = POINTER(CUser)

lib.free_users.argtypes = [POINTER(CUser), c_int]



@dataclass
class User:
    name:str
    age:int
    email:str
    
    @classmethod
    def create_user_from_C(cls:'User', name:str, age:int, email:str) -> 'User':
        pointer = lib.create_user(name.encode(encoding="utf-8"), age, email.encode(encoding="utf-8"))
        data = pointer.contents
        try:
            assert data.name.decode() == name
            assert data.age == age
            assert data.email.decode() == email
            return User(data.name.decode(errors="replace"), data.age, data.email.decode(errors="replace"))
        except (AssertionError, UnicodeDecodeError) as e:
            raise ValueError(f"Could not instantiate User\n\t{repr(traceback.format_exception(e))}")
        finally:
            # Something went wrong, free the memory
            lib.free_user(pointer)
    
    @classmethod
    def create_random_user(cls:'User') -> 'User':
        pointer = lib.create_random_user()
        data = pointer.contents
        try:
            return User(data.name.decode(errors="replace"), data.age, data.email.decode(errors="replace"))
        except (AssertionError, UnicodeDecodeError) as e:
            raise ValueError(f"Could not instantiate User\n\t{repr(traceback.format_exception(e))}")
        finally:
            # Something went wrong, free the memory
            lib.free_user(pointer)
            
    @classmethod
    def create_random_users(cls:'User', count:int) -> list['User']:
        results = []
        pointer = lib.create_random_users(count)

        if not pointer:
            raise ValueError("Failed to parse URLs")
        try:
            for i in range(count):
                site_ptr = pointer[i]
                data = site_ptr
                
                try:
                    results.append(cls(
                        name=data.name.decode(errors="replace"),
                        age=data.age,
                        email=data.email.decode(errors="replace"),
                    ))
                except AttributeError:
                    continue # No data
        finally:
            cls.free_sites(pointer[0],count )
        return results
        
    @staticmethod
    def free_sites(array_pointer: CUser, count:int):
        if not array_pointer:
            return
        lib.free_users(array_pointer, count)
        
        
users = User.create_random_users(10)

print(users)
