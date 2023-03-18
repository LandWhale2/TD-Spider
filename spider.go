package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"os"

	"golang.org/x/net/html"
	"github.com/orcaman/concurrent-map"
)

var requestCount uint64
var parseCount uint64
var isVisit = cmap.New()
var workerCheck = cmap.New()
var keyword string

var requestWorkerCount = 200
var parseWorkerCount = 50
var exitSignal = make(chan bool)

var httpClient *http.Client

func init() {
    httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
    transport := &http.Transport{
        MaxIdleConnsPerHost:   20,
        IdleConnTimeout:       60 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }

    return &http.Client{
        Timeout:   10 * time.Second,
        Transport: transport,
    }
}

type HTMLData struct {
    HTMLText string
    URL      string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <URL>")
		os.Exit(1)
	}

	startURL := os.Args[1]
	keyword = os.Args[2]

	var wg sync.WaitGroup
	urlChannel := make(chan string, 100000)
	HTMLChannel := make(chan HTMLData, 1000)

	wg.Add(parseWorkerCount)
	urlChannel <- startURL

	for i := 0; i < requestWorkerCount; i++ {
		go doRequest(&wg, urlChannel, HTMLChannel)
	}

	for i := 0; i < parseWorkerCount; i++ {
		go doParse(&wg, urlChannel, HTMLChannel, i)
	}

	go printCount()

	fmt.Println("********** START **********")
	wg.Wait()
	close(urlChannel)
    close(HTMLChannel)
	exitSignal <- true
}

func printCount() {
	i := 1
	for {
		select {
		case <-exitSignal:
			return
		default:
			time.Sleep(time.Second)

			fmt.Println("[",i, "SECOND]", "REQUEST_COUNT : ", requestCount, "PARSING_COUNT : ", parseCount)
			i++
		}
	}
}

func doRequest(wg *sync.WaitGroup, urlChannel chan string, HTMLChannel chan HTMLData) {
    for url := range urlChannel {
        atomic.AddUint64(&requestCount, 1)
        urlString := url

        response := getHttpResponse(urlString)
        if response != "" {
            HTMLChannel <- HTMLData{HTMLText: response, URL: urlString}
        }
    }
    fmt.Println("[DONE] doRequest")
    wg.Done()
}


func doParse(wg *sync.WaitGroup, urlChannel chan string, HTMLChannel chan HTMLData, number int) {
	defer wg.Done()

    for htmlData := range HTMLChannel {
        currentURL := htmlData.URL
        node, err := html.Parse(strings.NewReader(htmlData.HTMLText))
        if err == nil {
            findHTML(node, urlChannel, currentURL)
            atomic.AddUint64(&parseCount, 1)
            workerCheck.Set(string(number), true)
        }
    }
    fmt.Println("[DONE] doParse")
}



func getHttpResponse(url string) string {
    resp, err := httpClient.Get(url)
    if err != nil {
        return ""
    }

    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return ""
    }
    return string(data)
}

func findHTML(node *html.Node, urlChannel chan string, currentURL string) {
	if node.Data == "script" || node.Data == "style"{
		return	
	}
	
	getDensityTextScore := getDensityTextScore(node)
	
	if node.Type == html.ElementNode{
		if getDensityTextScore >= 0{
			if node.Data=="a" && len(node.Attr)>0{
				for _, at := range node.Attr {
					if at.Key == "href" && len(strings.TrimSpace(at.Val)) > 0{
						if strings.HasPrefix(at.Val, "http") {
							url := at.Val
							_, isvisited := isVisit.Get(url)

							if isvisited == false {
								isVisit.Set(url, true)
								urlChannel <- url
							}
						}
						break
					}
				}
			}
		}

	}else if node.Type == html.TextNode {
		if getDensityTextScore >= 0{
			if len(strings.TrimSpace(node.Data)) > 0 {
				if strings.Contains(node.Data, keyword) {
					// fmt.Println("[URL]", currentURL, " [TEXT]", node.Data)
				}
				
			}
		}
	}
	
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		findHTML(c, urlChannel, currentURL)
	}
	
}

func chlidNodeLength(node *html.Node, count int) int {
	if node == nil {
		return 0
	}

	if node.Type == html.ElementNode{
		count++
	}
	
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		count = Max(chlidNodeLength(c, count), count)
	}

	return count
}

func chlidNodeTextLengthSum(node *html.Node, count int) int {
	if node == nil {
		return 0
	}

	if node.Type == html.TextNode{
		if len(strings.TrimSpace(node.Data)) != 0{
			count += len(node.Data)
		}
	}
	
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		count = Max(chlidNodeTextLengthSum(c, count), count)
	}

	return count
}

func getDensityTextScore(node *html.Node) float32 {
	childNodeTagLengthNumber := chlidNodeLength(node, 0)
	chlidNodeTextSumNumber := chlidNodeTextLengthSum(node, 0)

	if childNodeTagLengthNumber == 0{
		childNodeTagLengthNumber = 1
	}

	return float32(chlidNodeTextSumNumber) / float32(childNodeTagLengthNumber)
}

func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}