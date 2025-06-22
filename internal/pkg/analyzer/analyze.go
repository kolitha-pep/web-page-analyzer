package analyzer

import (
	"fmt"
	"io"
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
	fmt.Println(parsedUrl.String())
	res, err := http.Get(parsedUrl.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Check if the response status code is OK
	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// goquery to parse the HTML document
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Determine HTML version of the document
	htmlVersion, err := GetHtmlVersion(res)
	if err != nil {
		return nil, fmt.Errorf("failed to determine HTML version: %w", err)
	}

	out := &WebPageMeta{
		Url:         parsedUrl.String(),
		Title:       d.Find("title").Text(),
		HtmlVersion: htmlVersion,
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

func GetHtmlVersion(resp *http.Response) (string, error) {
	const maxRead = 8192 // reading only the first 8KB of the response body
	body := make([]byte, maxRead)
	n, err := resp.Body.Read(body)
	if err != nil && err != io.EOF {
		return "Unknown", err
	}

	content := strings.ToLower(string(body[:n]))
	content = strings.TrimSpace(content)

	// Check for known DOCTYPEs
	switch {
	case strings.HasPrefix(content, "<!doctype html>"):
		return "HTML5", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd html 4.01 transitional//en"`):
		return "HTML 4.01 Transitional", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd html 4.01 strict//en"`):
		return "HTML 4.01 Strict", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd html 4.01 frameset//en"`):
		return "HTML 4.01 Frameset", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd xhtml 1.0 strict//en"`):
		return "XHTML 1.0 Strict", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd xhtml 1.0 transitional//en"`):
		return "XHTML 1.0 Transitional", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd xhtml 1.0 frameset//en"`):
		return "XHTML 1.0 Frameset", nil
	case strings.Contains(content, `<!doctype html public "-//w3c//dtd xhtml 1.1//en"`):
		return "XHTML 1.1", nil
	default:
		return "Unknown", nil
	}
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
