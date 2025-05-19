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
	count := int(cCount)

	// Build a Go slice of length n over the C array
	const maxUrls = 1 << 30 // see below
	urlPtrs := (*[maxUrls]*C.char)(unsafe.Pointer(cUrls))[:count:count]

	goURLs := make([]string, 0, count)
	for i := range count {
		goURLs = append(goURLs, C.GoString(urlPtrs[i]))
	}

	sitesData := ParseURLs(goURLs)

	// Allocate one big C array for all Site structs
	size := C.size_t(cCount) * C.size_t(unsafe.Sizeof(C.Site{}))
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
	// urls := []string{
	// 	"https://www.google.com",
	// 	"https://www.facebook.com",
	// 	"https://www.youtube.com",
	// 	"https://www.twitter.com",
	// 	"https://www.instagram.com",
	// 	"https://www.linkedin.com",
	// 	"https://www.wikipedia.org",
	// 	"https://www.reddit.com",
	// 	"https://www.amazon.com",
	// 	"https://www.netflix.com",
	// 	"https://www.apple.com",
	// 	"https://www.microsoft.com",
	// 	"https://www.dropbox.com",
	// 	"https://www.spotify.com",
	// 	"https://www.tumblr.com",
	// 	"https://www.quora.com",
	// 	"https://www.stackoverflow.com",
	// 	"https://www.medium.com",
	// 	"https://www.bing.com",
	// 	"https://www.paypal.com",
	// 	"https://www.ebay.com",
	// 	"https://www.pinterest.com",
	// 	"https://www.tiktok.com",
	// 	"https://www.cnn.com",
	// 	"https://www.bbc.com",
	// 	"https://www.nytimes.com",
	// 	"https://www.theguardian.com",
	// 	"https://www.washingtonpost.com",
	// 	"https://www.forbes.com",
	// 	"https://www.bloomberg.com",
	// 	"https://www.airbnb.com",
	// 	"https://www.udemy.com",
	// 	"https://www.coursera.org",
	// 	"https://www.khanacademy.org",
	// 	"https://www.github.com",
	// 	"https://www.gitlab.com",
	// 	"https://www.codepen.io",
	// 	"https://www.heroku.com",
	// 	"https://www.digitalocean.com",
	// 	"https://www.slack.com",
	// 	"https://www.zoom.us",
	// 	"https://www.skype.com",
	// 	"https://www.trello.com",
	// 	"https://www.notion.so",
	// 	"https://www.canva.com",
	// 	"https://www.wix.com",
	// 	"https://www.shopify.com",
	// 	"https://www.mozilla.org",
	// 	"https://www.icann.org",
	// 	"https://www.cloudflare.com",
	// 	"https://www.openai.com",
	// 	"https://www.deepmind.com",
	// 	"https://www.ibm.com",
	// 	"https://www.oracle.com",
	// 	"https://www.sap.com",
	// 	"https://www.adobe.com",
	// 	"https://www.salesforce.com",
	// 	"https://www.zendesk.com",
	// 	"https://www.asana.com",
	// 	"https://www.bitbucket.org",
	// 	"https://www.bitly.com",
	// 	"https://www.hubspot.com",
	// 	"https://www.mailchimp.com",
	// 	"https://www.figma.com",
	// 	"https://www.behance.net",
	// 	"https://www.dribbble.com",
	// 	"https://www.envato.com",
	// 	"https://www.codeacademy.com",
	// 	"https://www.pluralsight.com",
	// 	"https://www.edx.org",
	// 	"https://www.futurelearn.com",
	// 	"https://www.teachable.com",
	// 	"https://www.skillshare.com",
	// 	"https://www.lynda.com",
	// 	"https://www.x.com",
	// 	"https://www.aliexpress.com",
	// 	"https://www.flipkart.com",
	// 	"https://www.target.com",
	// 	"https://www.homedepot.com",
	// 	"https://www.walmart.com",
	// 	"https://www.bestbuy.com",
	// 	"https://www.nike.com",
	// 	"https://www.adidas.com",
	// 	"https://www.samsung.com",
	// 	"https://www.huawei.com",
	// 	"https://www.sony.com",
	// 	"https://www.lenovo.com",
	// 	"https://www.dell.com",
	// 	"https://www.hp.com",
	// 	"https://www.intel.com",
	// 	"https://www.amd.com",
	// 	"https://www.nvidia.com",
	// 	"https://www.tesla.com",
	// 	"https://www.ford.com",
	// 	"https://www.gm.com",
	// 	"https://www.toyota.com",
	// 	"https://www.honda.com",
	// 	"https://www.bmw.com",
	// 	"https://www.mercedes-benz.com",
	// 	"https://www.ycombinator.com",
	// 	"https://www.producthunt.com",
	// 	"https://www.crunchbase.com",
	// 	"https://www.techcrunch.com",
	// 	"https://www.engadget.com",
	// 	"https://www.theverge.com",
	// 	"https://www.wired.com",
	// 	"https://www.zdnet.com",
	// 	"https://www.cnet.com",
	// 	"https://www.lifehacker.com",
	// 	"https://www.makeuseof.com",
	// 	"https://www.arstechnica.com",
	// 	"https://www.tomshardware.com",
	// 	"https://www.howtogeek.com",
	// 	"https://www.sciencedaily.com",
	// 	"https://www.nature.com",
	// 	"https://www.sciencemag.org",
	// 	"https://www.popsci.com",
	// 	"https://www.space.com",
	// 	"https://www.nasa.gov",
	// 	"https://www.noaa.gov",
	// 	"https://www.who.int",
	// 	"https://www.cdc.gov",
	// 	"https://www.nih.gov",
	// 	"https://www.whitehouse.gov",
	// 	"https://www.congress.gov",
	// 	"https://www.supremecourt.gov",
	// 	"https://www.un.org",
	// 	"https://www.worldbank.org",
	// 	"https://www.imf.org",
	// 	"https://www.oecd.org",
	// 	"https://www.weforum.org",
	// 	"https://www.undp.org",
	// 	"https://www.unesco.org",
	// 	"https://www.ted.com",
	// 	"https://www.brainpickings.org",
	// 	"https://www.goodreads.com",
	// 	"https://www.bookbub.com",
	// 	"https://www.librarything.com",
	// 	"https://www.archlinux.org",
	// 	"https://www.ubuntu.com",
	// 	"https://www.debian.org",
	// 	"https://www.fedoraproject.org",
	// 	"https://www.linuxmint.com",
	// 	"https://www.kali.org",
	// 	"https://www.gentoo.org",
	// 	"https://www.apache.org",
	// 	"https://www.nginx.com",
	// 	"https://www.mysql.com",
	// 	"https://www.postgresql.org",
	// 	"https://www.mongodb.com",
	// 	"https://www.redis.io",
	// 	"https://www.sqlite.org",
	// 	"https://www.rabbitmq.com",
	// 	"https://www.docker.com",
	// 	"https://www.kubernetes.io",
	// 	"https://www.jenkins.io",
	// 	"https://www.travis-ci.com",
	// 	"https://www.circleci.com",
	// 	"https://www.netlify.com",
	// 	"https://www.vercel.com",
	// 	"https://www.render.com",
	// 	"https://www.herokuapp.com",
	// 	"https://www.supabase.com",
	// 	"https://www.prisma.io",
	// 	"https://www.grafana.com",
	// 	"https://www.prometheus.io",
	// 	"https://www.elastic.co",
	// 	"https://www.datadoghq.com",
	// 	"https://www.splunk.com",
	// 	"https://www.cloudflarestatus.com",
	// 	"https://www.cloudflarestatus.com",
	// 	"https://www.wolframalpha.com",
	// 	"https://www.desmos.com",
	// 	"https://www.geogebra.org",
	// 	"https://www.overleaf.com",
	// 	"https://www.latex-project.org",
	// 	"https://www.arxiv.org",
	// 	"https://www.researchgate.net",
	// 	"https://www.jstor.org",
	// 	"https://www.acm.org",
	// 	"https://www.ieee.org",
	// 	"https://www.springer.com",
	// 	"https://www.elsevier.com",
	// 	"https://www.scopus.com",
	// 	"https://www.duolingo.com",
	// 	"https://www.memrise.com",
	// 	"https://www.busuu.com",
	// 	"https://www.babbel.com",
	// 	"https://www.hellotalk.com",
	// 	"https://www.italki.com",
	// 	"https://www.lingq.com",
	// 	"https://www.openstreetmap.org",
	// 	"https://www.mapbox.com",
	// 	"https://www.here.com",
	// 	"https://www.uber.com",
	// 	"https://www.lyft.com",
	// 	"https://www.doordash.com",
	// 	"https://www.ubereats.com",
	// 	"https://www.grubhub.com",
	// }
	// resp := ParseURLs(urls)
	// for i := range len(resp) {
	// 	if resp[i].url != urls[i] {
	// 		fmt.Printf("%d: \n\t%s %s\n\n", i, resp[i].url, urls[i])
	// 	}
	// }
}
