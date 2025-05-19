package main

/*
#include <stdlib.h>

typedef struct{
	char* url;
	char* domain;
	char* server;
	char* protocol;
	char* contentType;
	char* body;
	int port;
} Site;
*/
import "C"
import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
	"unsafe"

	helpers "github.com/Descent098/cgo-python-helpers"
)

type Site struct {
	url         string // the raw URL
	domain      string // The domain the URL is hosted at
	server      string // The value of the server header
	protocol    string // The protocl of the site (http or https)
	contentType string // The content type of the body (i.e. "text/html")
	body        string // The body of the url
	port        int    // The port the url is on
}

// Retrieves a value from HTTP headers or returns a default if not found
//
// # Parameters
//
//	headers (http.Header): The HTTP headers from the response
//	key (string): The header key to look for
//	defaultValue (string): The default value to return if the key is not found
//
// # Returns
//
//	string: The value of the header or the default value
func getFromHeaders(headers http.Header, key string, defaultValue string) string {
	values, ok := headers[key]
	if !ok || len(values) == 0 {
		return defaultValue
	}
	return values[0]
}

// Scrape metadata from from a single URL
//
// # Parameters
//
//	rawUrl (string): The raw URL string to fetch
//
// # Returns
//
//	*Site: A pointer to a Site struct containing metadata
//	error: An error if the request or parsing fails
func scrapeSite(rawUrl string) (*Site, error) {
	var result Site
	result.url = rawUrl
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return &result, err
	}

	protocol := parsedURL.Scheme
	domain := parsedURL.Hostname()
	port := 80
	if protocol == "https" {
		port = 443
	}
	if parsedURL.Port() != "" {
		p, err := net.LookupPort("tcp", parsedURL.Port())
		if err == nil {
			port = p
		}
	}

	result.protocol = protocol
	result.domain = domain
	result.port = port

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     5 * time.Second,
			TLSHandshakeTimeout: 2 * time.Second,
			DisableKeepAlives:   true,
			ForceAttemptHTTP2:   false,
		},
	}

	resp, err := client.Get(rawUrl)
	if err != nil {
		return &result, err
	}
	defer resp.Body.Close()

	contentType := getFromHeaders(resp.Header, "Content-Type", "text/plain")
	server := getFromHeaders(resp.Header, "Server", "")
	result.contentType = contentType
	result.server = server

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &result, err
	}

	result.body = string(bodyBytes)

	return &result, nil
}

// Takes in a list of URL's and parses their content to Site's
//
// # Parameters
//
//	urls ([]string): A slice of raw URLs to scrape
//
// # Notes
//
//   - Do not include repeat URL's or this function will panic
//
// # Returns
//
//	[]*Site: A slice of Site pointers containing parsed metadata
func ParseURLs(urls []string) []*Site {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	result := make([]*Site, len(urls))
	const MAXCONCURRENTGOROUTINES = 50

	var (
		wg             sync.WaitGroup
		semaphore      = make(chan int, MAXCONCURRENTGOROUTINES)
		writeArrayLock sync.Mutex
	)

	for index, url := range urls {
		semaphore <- 1
		wg.Add(1)
		go func(url string, index int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			site, err := scrapeSite(url)
			if err != nil {
				fmt.Printf("Error while processing %s: %v\n", url, err)
				site = &Site{
					url,
					"", "", "", "", "", 80,
				}
			}

			writeArrayLock.Lock()
			result[index] = site
			writeArrayLock.Unlock()
		}(url, index)

	}
	wg.Wait()

	if len(result) != len(urls) {
		// This should never happen
		fmt.Printf("Incorrect number of sites %d/%d", len(result), len(urls))
		panic("URL results may cause memory misalignment, exiting")
	}
	return result
}

// C-callable wrapper that parses multiple URLs and returns C structs
//
// # Parameters
//
//	cUrls (**C.char): An array of C strings (URLs)
//	cCount (C.int): The number of URLs
//
// # Returns
//
//	*C.Site: A pointer to the first element of an array of C.Site structs
//
//export parse_urls
func parse_urls(cUrls **C.char, cCount C.int) *C.Site {
	goURLs := helpers.CStringArrayToSlice(unsafe.Pointer(cUrls), int(cCount))

	sitesData := ParseURLs(goURLs)

	sites := PrepareSitesForExport(sitesData)

	return sites
}

// Takes in a slice of Site instances, and returns a C-exportable array of C.Site's
//
// # Parameters
//
//	sitesData ([]*Site): An slice of pointers to Site instances
//
// # Returns
//
//	*C.Site: A pointer to the first element of an array of C.Site structs
func PrepareSitesForExport(sitesData []*Site) *C.Site {
	count := len(sitesData)
	// Allocate one big C array for all Site structs
	size := C.size_t(count) * C.size_t(unsafe.Sizeof(C.Site{}))
	sites := (*C.Site)(C.malloc(size))

	// Copy each non-nil Site into the C array
	for i, site := range sitesData {
		if site == nil {
			// leave this slot zeroed
			continue
		}
		// compute pointer to &sites[i]
		slot := (*C.Site)(unsafe.Pointer(
			uintptr(unsafe.Pointer(sites)) + uintptr(i)*unsafe.Sizeof(C.Site{}),
		))

		slot.url = C.CString(site.url)
		slot.domain = C.CString(site.domain)
		slot.server = C.CString(site.server)
		slot.protocol = C.CString(site.protocol)
		slot.contentType = C.CString(site.contentType)
		slot.body = C.CString(site.body)
		slot.port = C.int(site.port)
	}
	return sites
}

// C-callable wrapper to scrape a single URL
//
// # Parameters
//
//	cUrl (*C.char): A single URL string
//
// # Returns
//
//	*C.Site: A pointer to a C.Site struct with metadata, or nil on error
//
//export scrape_single_url
func scrape_single_url(cUrl *C.char) *C.Site {
	url := C.GoString(cUrl)      // Convert string back to Go string
	site, err := scrapeSite(url) // Get site data
	if err != nil {
		fmt.Printf("Error scraping %s: %v\n", url, err)
		return nil
	}
	// Convert site data back to C struct

	// Allocate memory for C struct
	cSite := (*C.Site)(C.malloc(C.size_t(unsafe.Sizeof(C.Site{}))))

	// Assign values
	cSite.url = C.CString(site.url)
	cSite.domain = C.CString(site.domain)
	cSite.server = C.CString(site.server)
	cSite.protocol = C.CString(site.protocol)
	cSite.contentType = C.CString(site.contentType)
	cSite.body = C.CString(site.body)
	cSite.port = C.int(site.port)

	return cSite
}

// Releases memory allocated for a single C.Site struct
//
// # Parameters
//
//	site (*C.Site): A pointer to the C.Site struct to free
//
//export free_site
func free_site(site *C.Site) {
	if site == nil {
		return
	}
	C.free(unsafe.Pointer(site.url))
	C.free(unsafe.Pointer(site.domain))
	C.free(unsafe.Pointer(site.server))
	C.free(unsafe.Pointer(site.protocol))
	C.free(unsafe.Pointer(site.contentType))
	C.free(unsafe.Pointer(site.body))
	C.free(unsafe.Pointer(site))
}

// Releases memory allocated for an array of C.Site structs
//
// # Parameters
//
//	sites (*C.Site): A pointer to the first element in a C.Site array
//	count (C.int): The number of elements in the array
//
//export free_sites
func free_sites(sites *C.Site, count C.int) {
	for i := range int(count) {
		sitePointer := (*C.Site)(unsafe.Pointer(uintptr(unsafe.Pointer(sites)) + uintptr(i)*unsafe.Sizeof(C.Site{})))

		C.free(unsafe.Pointer(sitePointer.url))
		C.free(unsafe.Pointer(sitePointer.domain))
		C.free(unsafe.Pointer(sitePointer.server))
		C.free(unsafe.Pointer(sitePointer.protocol))
		C.free(unsafe.Pointer(sitePointer.contentType))
		C.free(unsafe.Pointer(sitePointer.body))
	}
	C.free(unsafe.Pointer(sites))
}

func main() {

}
