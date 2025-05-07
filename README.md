# Python + Go
A repository for the examples from my series about integrating python with Go

Articles: 

1. [Python + Go: Introduction](https://kieranwood.ca/tech/blog/python-plus-go-intro/)
2. [Python + Go: The basics](https://kieranwood.ca/tech/blog/python-plus-go-basics/)
3. [Python + Go: Examples](https://kieranwood.ca/tech/blog/python-plus-go-examples/)

## Setup

For details on getting everything setup read the articles :)

In particular make sure you setup:

- [Go](https://go.dev/doc/install)
- [Python](https://www.python.org/downloads/)
- gcc/[zig cc](https://ziglang.org/download/)
- [CGo](https://kieranwood.ca/tech/blog/python-plus-go-basics/#prepping-cgo)

## Folder Structure

Each folder is setup to corespond to a part of an article. Within each folder there will be details about it's structure, but the folders correspond to the parts of the articles linked below:

- [`📂easy-part/`](https://kieranwood.ca/tech/blog/python-plus-go-basics/#easy-parts)
- [`📂hard-part/`](https://kieranwood.ca/tech/blog/python-plus-go-basics/#the-hard-part)
- [`📂examples/`](https://kieranwood.ca/tech/blog/python-plus-go-examples/)

## Conversion tables

The tables I put together are qutie handy for doing conversions, so I broke them out here to make them easy to find:

| go type | C type |cgo type| Python type | 
|---------|--------|--------|-------------|
| `string` | `char` | `C.char`  | `str` |
| `string` | `signed char` | `C.schar` | `str` |
| `string` | `unsigned char` | `C.uchar` | `str` |
| `int16` or `int` | `short` |`C.short`| `int` |
| `uint16` or `int` | `unsigned short` | `C.ushort`| `int` |
| `int` | `int` | `C.int` | `int` |
| `uint` | `unsigned int` | `C.uint` | `int` |
| `int32` or `int` | `long` | `C.long` | `int` |
| `uint32` or `int` | `unsigned long` | `C.ulong` | `int` |
| `int64` or `int` | `long long` | `C.longlong` | `int` |
| `uint64` or `int` | `unsigned long long` | `C.ulonglong` | `int` |
| `float32` | `float` | `C.float` | `float` |
| `float64` | `double` | `C.double` | `float` | 
| `struct` | `struct` |  `C.struct_<name_of_C_Struct>` | `class` |
| `struct` | `union` |  `C.union_<name_of_C_Union>` | `class` |
| `struct` | `enum` |  `C.enum_<name_of_C_Enum>` | `class` |
| `unsafe.pointer` | `void*` | `unsafe.pointer` | N/A |

\**Please note that unions and enums are probably better left out of your code as much as you can, they're very finicky*

![conversion table](https://kieranwood.ca/tech/blog/python-plus-go/full-pipeline-conversions.excalidraw.png)

