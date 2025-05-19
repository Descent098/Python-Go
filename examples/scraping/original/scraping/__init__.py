import os
from platform import platform

# Check if dynamic library is compiled
if platform().lower().startswith("windows"):
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.dll")
else:
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.so")
    
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

# import python library
from .lib import Site