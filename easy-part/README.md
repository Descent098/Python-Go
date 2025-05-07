# Easy Part

This section includes an example of how to use copy-by-value types (`int`'s, `float`'s). For this demo we're getting the [factorial](https://www.freecodecamp.org/news/what-is-a-factorial/) of a number using go. 

## Folder Structure

The structure of the folder is 
```
📂easy-part
├─ 📄go.mod
├─ 📄lib.go
├─ 📄lib.dll or 📄lib.so
├─ 📄lib.h
└──📄testing.py
```

- `📄lib.go`: The Go code that has the `Factorial()` implementation
- `📄lib.dll` or `📄lib.so`: The generated file that is the compiled form of the go library
- `📄go.mod`: The file that allows you to compile go
- `📄lib.h`: A generated file that tells C how to use your `.dll` or `.so` file
- `📄testing.py`: The python code that consumes the go library

## Running

Once you have your environment setup, you can use the below commands to build:

<details><summary>Linux/Mac</summary>

```bash
go build -buildmode=c-shared -o lib.so lib.go 
```

</details>

<details><summary>Windows</summary>

```bash
go build -buildmode=c-shared -o lib.dll lib.go
```

</details>
<br>

Then run the python code using:


<details><summary>Linux/Mac</summary>

```bash
python3 testing.py
```

</details>

<details><summary>Windows</summary>

```bash
python testing.py
```

</details>



