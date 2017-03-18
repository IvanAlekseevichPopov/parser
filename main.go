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
	links := make(map[int]string)

	flag.StringVar(&site, "site", "", "http://example.com")
	flag.IntVar(&queryLimit, "l", -1, "Limit of queries. Unlimited by default")
	flag.IntVar(&threads, "t", 20, "Quantity of threads")
	flag.Parse()

	if !isVaildSiteName(site) {
		fmt.Println("Yum must specify correct site name to parse")
		os.Exit(0)
	}

	fmt.Printf("Start parsing %s\n", site)

	links = getLinks(site)
	// for {
	// fmt.Printf("%s\n", links[0])
	// delete(links, 0)
	// fmt.Printf("%s\n", links[1])
	// delete(links, 1)
	// }

	for _, link := range links {
		fmt.Printf("%s\n", link)
	}
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
	userAgentCollection := [15]string{
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0",
		"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G920V Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; SM-G928X Build/LMY47X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 6P Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/51.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/50.0",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
	}

	rand.Seed(time.Now().UnixNano())
	return userAgentCollection[rand.Intn(15)]
}

func getLinks(url string) map[int]string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", randomUserAgent())
	resp, _ := client.Do(req)

	z := html.NewTokenizer(resp.Body)
	links := make(map[int]string)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			resp.Body.Close()
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					if -1 != strings.Index(a.Val, "catalog") {
						links[len(links)] = a.Val
					}
					break
				}
			}
		}
	}
}
