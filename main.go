// parser - scope links of site
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	resp, _ := http.Get(site)
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break //not working return - exit
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					// fmt.Println("Found href:", a.Val)
					links[len(links)] = a.Val
					break
				}
			}
			// isAnchor := t.Data == "a"
			// if isAnchor {
			// 	fmt.Println("We found a link!")
			// }
		}
	}
	resp.Body.Close()
	fmt.Println(links)

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
