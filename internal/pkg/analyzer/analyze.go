package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/fetcher"
)

type WebPageMeta struct {
	Url           string         `json:"url"`
	Title         string         `json:"title"`
	HtmlVersion   string         `json:"html_version"`
	HeadTags      map[string]int `json:"head_tags"`
	HasLoginForm  bool           `json:"has_login_form"`
	InternalLinks int            `json:"internal_links"`
	ExternalLinks int            `json:"external_links"`
	BrokenLinks   int            `json:"broken_links"`
}

func AnalyzeWebPage(in string) (*WebPageMeta, error) {
	// Validate the input URL
	parsedUrl, err := url.ParseRequestURI(in)
	if err != nil {
		return nil, err
	}

	// Fetch the webpage content
	res, err := http.Get(parsedUrl.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Check if the response status code is OK
	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	out := &WebPageMeta{
		Url:         parsedUrl.String(),
		Title:       d.Find("title").Text(),
		HtmlVersion: GetHtmlVersion(d),
	}

	// Count HTML head tags
	hTags := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		hTags[tag] = d.Find(tag).Length()
	}

	baseUrl := parsedUrl.Host
	hyperLinks := d.Find("a[href]")
	out.HeadTags = hTags

	var wg sync.WaitGroup
	var mu sync.Mutex

	internalLinks, externalLinks, brokenLinks := 0, 0, 0

	hyperLinks.Each(func(i int, x *goquery.Selection) {
		href, _ := x.Attr("href")
		wg.Add(1)

		go func(link string) {

			defer wg.Done()

			absoluteUrl, err := url.ParseRequestURI(link)

			// If the link is invalid or a mailto link, skip it
			if err != nil || absoluteUrl.Scheme == "" || absoluteUrl.Scheme == "mailto" {
				return
			}

			fullUrl := link

			if !absoluteUrl.IsAbs() {
				fullUrl = parsedUrl.ResolveReference(absoluteUrl).String()
			}

			isInternalLink := strings.Contains(absoluteUrl.Host, baseUrl)

			// check if this link is accessible or not
			if !fetcher.IsLinkReachable(fullUrl) {
				mu.Lock()
				brokenLinks++
				mu.Unlock()
			}

			if isInternalLink {
				mu.Lock()
				internalLinks++
				mu.Unlock()
			} else {
				mu.Lock()
				externalLinks++
				mu.Unlock()
			}
		}(href)
	})

	wg.Wait()

	out.InternalLinks = internalLinks
	out.ExternalLinks = externalLinks
	out.BrokenLinks = brokenLinks
	out.HasLoginForm = CheckForLoginForm(d)

	return out, nil
}

func GetHtmlVersion(doc *goquery.Document) string {
	documentNode := doc.Nodes[0].FirstChild
	if documentNode != nil && documentNode.Type == 10 {
		return "HTML5"
	}
	return "Unknown"
}

func CheckForLoginForm(doc *goquery.Document) bool {
	found := false
	doc.Find("form").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find("input[type='password']").Length() > 0 {
			found = true
			return false
		}
		return true
	})
	return found
}
