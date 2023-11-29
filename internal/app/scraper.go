package app

import (
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/maksimulitin/internal/util"
	"github.com/maksimulitin/pkg/model"
)

func ExtractSitemapURLs(startURL string) []string {
	worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func() { worklist <- []string{startURL} }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			n++
			go func(link string) {
				response, err := MakeRequest(link)
				if err != nil {
					log.Printf("Error retrieving URL: %s", link)
				}
				urls, _ := ExtractUrls(response)
				if err != nil {
					log.Printf("Error extracting document from response, URL: %s", link)
				}
				sitemapFiles, pages := util.IsSitemap(urls)
				if sitemapFiles != nil {
					worklist <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}
	return toCrawl
}

func MakeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", util.RandomUserAgent())
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ScrapeUrls(urls []string, parser model.Parser, concurrency int) []model.SeoData {
	tokens := make(chan struct{}, concurrency)
	var n int
	n++
	worklist := make(chan []string)
	results := []model.SeoData{}
	go func() { worklist <- urls }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, token chan struct{}) {
					log.Printf("Requesting URL: %s", url)
					res, err := ScrapePage(url, tokens, parser)
					if err != nil {
						log.Printf("Encountered error, URL: %s", url)
					} else {
						results = append(results, res)
					}
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}
	return results
}

func ExtractUrls(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")
	for i := range sel.Nodes {
		loc := sel.Eq(i)
		result := loc.Text()
		results = append(results, result)
	}
	return results, nil
}

func ScrapePage(url string, token chan struct{}, parser model.Parser) (model.SeoData, error) {
	res, err := CrawlPage(url, token)
	if err != nil {
		return model.SeoData{}, err
	}
	data, err := parser.GetSeoData(res)
	if err != nil {
		return model.SeoData{}, err
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

func ScrapeSitemap(url string, parser model.Parser, concurrency int) []model.SeoData {
	results := ExtractSitemapURLs(url)
	res := ScrapeUrls(results, parser, concurrency)
	return res
}
