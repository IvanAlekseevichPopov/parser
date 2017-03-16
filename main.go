// Echo4 выводит аргументы командной строки
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
			return
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					fmt.Println("Found href:", a.Val)
					break
				}
			}
			// isAnchor := t.Data == "a"
			// if isAnchor {
			// 	fmt.Println("We found a link!")
			// }
		}
	}
	// bytes, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println("HTML:\n\n", string(bytes))

	resp.Body.Close()
	// fmt.Println(site)
	// doc, err := goquery.NewDocument(site)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// doc.Find("a").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	// band := s.Find("a").Text()
	// 	// title := s.Find("i").Text()
	// 	// for _, b := range s.Attr("href") {
	// 	// 	fmt.Printf("%s \n", b)
	// 	// }
	// 	fmt.Printf(s.Attr("href"))
	// 	// os.Exit(0)
	// })
	// fmt.Println(strings.TrimSpace(resp.Find("a").Href()))
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
