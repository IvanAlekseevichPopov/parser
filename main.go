// parser - scope links of site
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	var site string
	var queryLimit int
	var threads int
	// links := make(map[int]string)

	flag.StringVar(&site, "site", "", "http://example.com")
	flag.IntVar(&queryLimit, "l", -1, "Limit of queries. Unlimited by default")
	flag.IntVar(&threads, "t", 20, "Quantity of threads")
	flag.Parse()

	if !isVaildSiteName(site) {
		fmt.Println("Yum must specify correct site name to parse")
		os.Exit(0)
	}

	fmt.Printf("Start parsing %s\n", site)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", site, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	resp, _ := client.Do(req)

	// fmt.Printf("%T\n", resp)
	links := getLinks(resp)

	for _, link := range links {
		fmt.Printf("%s\n", link)
	}

	resp.Body.Close()
}

func isVaildSiteName(name string) bool {
	if name == "" {
		return false
	}

	if -1 == strings.Index(name, "http://") {
		return false
	}

	return true
}

func randomUserAgent() string {
	userAgentCollection := [3]string{
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
	}

	rand.Seed(time.Now().UnixNano())
	return userAgentCollection[rand.Intn(3)]
}

func getLinks(resp *http.Response) map[int]string {
	z := html.NewTokenizer(resp.Body)
	links := make(map[int]string)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					// fmt.Println("Found href:", a.Val)
					links[len(links)] = a.Val
					break
				}
			}
		}
	}
}
