// parser - scope links of site
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	var site string
	var queryLimit, threads, timeOut int
	links := make(map[int]string)
	// task := make(chan string)
	response := make(chan map[int]string, 2000)

	flag.StringVar(&site, "site", "", "http://example.com")
	flag.IntVar(&queryLimit, "l", -1, "Limit of queries. Unlimited by default")
	flag.IntVar(&threads, "threads", 20, "Quantity of threads")
	flag.IntVar(&timeOut, "timeout", 10, "Quantity of threads")
	flag.Parse()

	if !isVaildSiteName(site) {
		fmt.Println("Yum must specify correct site name to parse")
		os.Exit(0)
	}

	fmt.Printf("Start parsing %s\n", site)

	go getLinks(site, timeOut, response)
	links = <-response

	for i := 1; i < 50; i++ {
		if i%5 == 0 {
			newLinsks := <-response
			for _, link := range newLinsks {
				// fmt.Printf("%s\n", link)
				//TODO проверка на дубликаты
				links[len(links)] = link
			}
			for threads < runtime.NumGoroutine() {
				fmt.Printf("Too many goroutines(%d). Waiting..\n", runtime.NumGoroutine())
				time.Sleep(time.Second * 1)
			}
		}

		fmt.Printf("Goroutines number - %d\n", runtime.NumGoroutine())
		fmt.Printf("Size map - %d\n", len(links))
		fmt.Printf("Setting task - %s\n", links[i])
		go getLinks(site+links[i], timeOut, response)
	}
}

func getLinks(url string, timeOut int, response chan<- map[int]string) {
	client := &http.Client{
		Timeout: time.Duration(time.Duration(timeOut) * time.Second),
		// Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}

	fmt.Println("taking new Task: " + url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", randomUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	z := html.NewTokenizer(resp.Body)
	links := make(map[int]string)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			resp.Body.Close()
			response <- links
			return
		} else if tt == html.StartTagToken {
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" {
					if isValidLink(a.Val) {
						links[len(links)] = a.Val
					}
					break
				}
			}
		}
	}
}

func isValidLink(str string) bool {
	if -1 == strings.Index(str, "catalog") {
		return false
	}

	if -1 != strings.Index(str, "http:") {
		return false
	}

	if -1 != strings.Index(str, "https:") {
		return false
	}

	return true
}

// func getProxy() string {
// 	return "proxy"
// }

func randomUserAgent() string {
	userAgentCollection := [26]string{
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2204.58 Safari/535.33",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0",
		"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G920V Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; SM-G928X Build/LMY47X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; SM-G928X Build/LMY44X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.35",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 6P Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/51.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/50.0",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
		"Mozilla/5.0 (compatible; MSIE 9.0; AOL 9.7; AOLBuild 4343.19; Windows NT 6.1; WOW64; Trident/5.0; FunWebProducts)",
		"Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B367 Safari/531.21.10",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
		"Mozilla/5.0 (iPod; CPU iPhone OS 5_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A405 Safari/7534.48.3",
		"Mozilla/5.0 (Linux; Android 4.4.4; GT-P5100 Build/KTU84Q) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Safari/537.36",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML like Gecko) Chrome/37.0.2062.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Opera/9.80 (Windows NT 6.2; Win64; x64) Presto/2.12 Version/12.16",
	}

	rand.Seed(time.Now().UnixNano())
	return userAgentCollection[rand.Intn(26)]
}

func isVaildSiteName(name string) bool {
	if name == "" {
		return false
	}

	if -1 == strings.Index(name, "http://") && -1 == strings.Index(name, "https://") {
		return false
	}

	return true
}
