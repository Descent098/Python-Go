import os
from platform import platform
from dataclasses import dataclass
from ctypes import Structure, c_char_p, c_int, POINTER

import os
from platform import platform
from .helpers import get_library, prepare_string_array

# Check if dynamic library is compiled
if platform().lower().startswith("windows"):
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.dll")
else:
    lib_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.so")

source_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "go", "lib.go")

lib = get_library(lib_path,source_path, compile=True )

class _CSite(Structure):
    """The C compatible Site structure, DO NOT USE DIRECTLY, use Site instead"""
    _fields_ = [
        ("url", c_char_p),
        ("domain", c_char_p),
        ("server", c_char_p),
        ("protocol", c_char_p),
        ("contentType", c_char_p),
        ("body", c_char_p),
        ("port", c_int),
    ]

# Set function return/arg types
lib.scrape_single_url.argtypes = [c_char_p]
lib.scrape_single_url.restype = POINTER(_CSite)

lib.parse_urls.argtypes = [POINTER(c_char_p), c_int]
lib.parse_urls.restype  = POINTER(_CSite)

lib.free_sites.argtypes = [POINTER(_CSite), c_int]
lib.free_sites.restype  = None

@dataclass
class Site:
    """A class representing a single site
    
    # Class Methods
    
    - from_str(url:str) -> Site: Parse site data into a Site instance from a Url
    - from_urls(urls:list[str]) -> list[Site]: Parse list of urls into Site instances 
    """
    url:str         # the raw URL
    domain:str      # The domain the URL is hosted at
    server:str      # The value of the server header
    protocol:str    # The protocl of the site (http or https)
    contentType:str # The content type of the body (i.e. "text/html")
    body:str        # The body of the url
    port:int        # The port the url is on
    
    @classmethod
    def from_str(cls:'Site', url:str) -> 'Site':
        """Create a Site instance from a URL

        Parameters
        ----------
        url : str
            The url to parse

        Returns
        -------
        Site
            The resulting Site instance

        Raises
        ------
        ValueError
            If an unrecoverable error occurs while parsing
        """
        pointer = lib.scrape_single_url(url.encode("utf-8"))
        if not pointer:
            raise ValueError(f"Failed to scrape: {url}")
        try:
            data = pointer.contents
            result = cls(
                url=data.url.decode(errors="replace"),
                domain=data.domain.decode(errors="replace"),
                server=data.server.decode(errors="replace"),
                protocol=data.protocol.decode(errors="replace"),
                contentType=data.contentType.decode(errors="replace"),
                body=data.body.decode(errors="replace"),
                port=data.port,
            )
        except Exception as e:
            from traceback import format_tb
            tb = "".join(format_tb(e))
            print(f"Error while fetching: \n\t{tb}")
        finally:
            if not pointer:
                raise ValueError(f"Failed to scrape: {url}")
            lib.free_site(pointer)
        return result

    @classmethod
    def from_urls(cls:'Site', urls:list[str], fail_on_error:bool=False) -> list['Site']:
        """Takes in a list of URL's and parses them to Site objects

        Returns
        -------
        list[Site]
            A list of the resulting Site objects

        Raises
        ------
        ValueError
            If an unrecoverable error occurs while parsing
        """
        # Preprocess variables to hand off to parse_urls
        url_array, count = prepare_string_array(urls)
        
        # Parse URL's and get a pointer to the CSite array resulting from parsing

        pointer = lib.parse_urls(url_array, count)
        
        if not pointer:
            raise ValueError("Failed to parse URLs")

        # Turn responses into proper Site instances
        results:list[Site] = []
        try:
            for i in range(count):
                site_ptr = pointer[i]
                data = site_ptr
                try:
                    results.append(cls(
                        url=data.url.decode(errors="replace"),
                        domain=data.domain.decode(errors="replace"),
                        server=data.server.decode(errors="replace"),
                        protocol=data.protocol.decode(errors="replace"),
                        contentType=data.contentType.decode(errors="replace"),
                        body=data.body.decode(errors="replace"),
                        port=data.port
                    ))
                except AttributeError as e:
                    if fail_on_error:
                        raise ValueError(f"Provided URL {data.url.decode(errors='replace')} errored")
                    continue # No data
        finally:
            cls.free_sites(pointer,count )
        return results

    @staticmethod
    def free_sites(array_pointer: _CSite, count:int):
        """Free's a C array of sites

        Parameters
        ----------
        array_pointer : _CSite
            A pointer to a CSite array
        count : int
            The number of items in the array
        """
        if not array_pointer:
            return
        lib.free_sites(array_pointer, count)