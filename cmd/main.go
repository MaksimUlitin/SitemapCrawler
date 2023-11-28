package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SeoData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type DefaultParser struct {
}

type Parser interface {
}

var UserAgentes = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func main() {
	p := DefaultParser{}
	results := ScrapeSiteMap("", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}

func RandomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNam := rand.Int() % len(UserAgentes)
	return UserAgentes[randNam]
}

func isSitemap(urls []string) ([]string, []string) {
	sitemapFile := []string{}
	pages := []string{}

	for _, page := range urls {
		foundSistemFile := strings.Contains(page, "XML")
		if foundSistemFile == true {
			fmt.Println("Found Sistem", page)
			sitemapFile = append(sitemapFile, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFile, pages
}

func ExtractSiteMapURls(startURL string) []string {
	Worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func() { Worklist <- []string{startURL} }()
	for ; n > 0; n-- {
	}
	list := <-Worklist
	for _, link := range list {
		n++
		go func(link string) {
			response, err := MakeRequest(link)
			if err != nil {
				log.Printf("error retrieving URL: %s ", link)
			}
			urls, _ := ExtractSiteMapURls(response)
			if err != nil {
				log.Printf("Error extracting document from response, URL: %s", link)
			}
			sitemapFiles, pages := isSitemap(urls)
			if sitemapFiles != nil {
				Worklist <- sitemapFiles
			}
			for _, page := range pages {
				toCrawl = append(toCrawl, page)
			}

		}(link)
	}
}

func MakeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", RandomUserAgent())
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err

}

func ScrapeURLs() {

}

func ScrapePage(url string, token chan struct{}, parser Parser) (SeoData, error) {

	res, err := CrawlPage(url, token)
	if err != nil {
		return SeoData{}, err
	}
	data, err := parser.GetSEOData(res)
	if err != nil {
		return SeoData{}, err
	}
	return data, nil
}

func CrawlPage(url string, tokens chan struct{}) (*http.Response, error) {
	tokens <- struct{}{}
	resp, err := MakeRequest(url)
	<-tokens
	if err != nil {
		return nil, err
	}

	return resp, err

}

func (d DefaultParser) GetSEOData(resp *http.Response) (SeoData, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return SeoData{}, err
	}
	result := SeoData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.MetaDescription, _ = doc.Find("meta[name^=description]").Attr("content")
	return result, err
}

func ScrapeSiteMap(url string, parser Parser, concurrency int) []SeoData {
	result := ExtractSiteMapURls(url)
	res := ScrapeURLs(result, parser, concurrency)
	return res
}
