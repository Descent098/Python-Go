import os
from platform import platform
from .helpers import get_library

# Check if dynamic library is compiled
if platform().lower().startswith("windows"):
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.dll")
else:
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.so")

source_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.go")

get_library(lib_path,source_path, compile=True )

# import python library
from .lib import Site