# Scraping

This demo will be using go to do some web scraping, then dumping the results to python, while keeping the python API simple to use. Essentially functions that generate this go struct:

```go
type Site struct {
	url         string // the raw URL
	domain      string // The domain the URL is hosted at
	server      string // The value of the server header
	protocol    string // The protocl of the site (http or https)
	contentType string // The content type of the body (i.e. "text/html")
	body        string // The body of the url
	port        int    // The port the url is on
}
```

Then packing it into this python class:

```python
@dataclass
class Site:
    url:str         # the raw URL
    domain:str      # The domain the URL is hosted at
    server:str      # The value of the server header
    protocol:str    # The protocl of the site (http or https)
    contentType:str # The content type of the body (i.e. "text/html")
    body:str        # The body of the url
    port:int        # The port the url is on
```

API wise it's usable via:

```python
from scraping import Site

Site.from_str("https://kieranwood.ca")

Site.from_urls(["https://google.ca", "https://cloudflare.ca"])
```

## Running

You should be able to run by just running `testing.py`, if you have your go and c compiler setup it will compile the lib and run it for you, or if it fails it will give you the command(s) to run.

## Folder Structure

Here is the folder structure for this folder, each version will have details about it's implementation in the README:

```
📂scraping/
├─ 📂original/
└──📂with-helper/
```


- `📂original`: The original implementation without the helper library
- `📂with-helper`: The version with the helper library
