# Hard Part

...

## Folder Structure

The structure of an individual folder is 
```
📂slices/ or 📂strings/ or 📂structs/
├─ 📄go.mod
├─ 📄go.sum
├─ 📄lib.go
├─ 📄lib.dll or 📄lib.so
├─ 📄lib.h
└──📄testing.py
```

- `📄lib.go`: The Go code that has the go implementation
- `📄lib.dll` or `📄lib.so`: The generated file that is the compiled form of the go library
- `📄go.mod`: The file that allows you to compile go
- `📄go.sum`: For folders with third party dependencies this tells golang how to download them
- `📄lib.h`: A generated file that tells C how to use your `.dll` or `.so` file
- `📄testing.py`: The python code that consumes the go library

## Running

Once you have your environment setup, you can use the below commands to build:

<details><summary>Linux/Mac</summary>

```bash
go mod tidy
go build -buildmode=c-shared -o lib.so lib.go 
```

</details>

<details><summary>Windows</summary>

```bash
go mod tidy
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



