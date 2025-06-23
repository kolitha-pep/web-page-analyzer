package analyzer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/fetcher"
	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/utils"
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
	QueryTime     float64        `json:"query_time"`
}

func AnalyzeWebPage(in string) (*WebPageMeta, error) {
	startedAt := time.Now()

	// Validate the input URL
	in = ensureHTTPS(in)
	parsedUrl, err := url.ParseRequestURI(in)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Fetch the webpage content
	res, err := http.Get(parsedUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the web page content: %w", err)
	}

	defer res.Body.Close()

	// Check if the response status code is OK
	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Determine HTML version of the document
	htmlVersion, err := getHtmlVersion(res)
	if err != nil {
		return nil, fmt.Errorf("failed to determine HTML version: %w", err)
	}

	// goquery to parse the HTML document
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document from response: %w", err)
	}

	// creating output object
	out := &WebPageMeta{
		Url:          parsedUrl.String(),
		Title:        d.Find("title").Text(),
		HtmlVersion:  htmlVersion,
		HeadTags:     countHtmlHeadTags(d),
		HasLoginForm: checkForLoginForm(d),
	}

	out.InternalLinks, out.ExternalLinks, out.BrokenLinks = analyzeHyperlinks(d, parsedUrl)
	out.QueryTime = utils.RoundFloat(time.Since(startedAt).Seconds(), 2)

	return out, nil
}

func analyzeHyperlinks(doc *goquery.Document, parsedUrl *url.URL) (int, int, int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	internalLinks, externalLinks, brokenLinks := 0, 0, 0

	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")
		wg.Add(1)

		go func(link string) {
			defer wg.Done()
			absoluteUrl, err := url.Parse(link)

			// If the link is invalid or a mailto link, skip it
			if err != nil || absoluteUrl.Scheme == "mailto" {
				return
			}

			fullUrl := link

			// If the link is relative, resolve it against the base URL
			if !absoluteUrl.IsAbs() {
				fullUrl = parsedUrl.ResolveReference(absoluteUrl).String()
				absoluteUrl.Host = parsedUrl.Host
			}

			isInternalLink := strings.Contains(absoluteUrl.Host, parsedUrl.Host)

			// check if this link is accessible or not
			if !fetcher.IsLinkReachable(fullUrl) {
				mu.Lock()
				brokenLinks++
				mu.Unlock()
			}

			// If the link is internal, increment internalLinks, otherwise externalLinks
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

	return internalLinks, externalLinks, brokenLinks
}

func countHtmlHeadTags(doc *goquery.Document) map[string]int {
	hTags := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		hTags[tag] = doc.Find(tag).Length()
	}

	return hTags
}

func getHtmlVersion(resp *http.Response) (string, error) {
	// reading only the first 8KB of the response body
	const maxRead = 8192

	// Use a buffer to read the response body without consuming it
	var buf bytes.Buffer
	tee := io.TeeReader(io.LimitReader(resp.Body, maxRead), &buf)

	body, err := io.ReadAll(tee)
	if err != nil {
		return "Unknown", fmt.Errorf("failed to read body: %w", err)
	}

	// Reset the body to the beginning for further processing
	resp.Body = io.NopCloser(io.MultiReader(&buf, resp.Body))

	content := strings.ToLower(string(body))
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

func checkForLoginForm(doc *goquery.Document) bool {
	return doc.Find("form").FilterFunction(func(_ int, form *goquery.Selection) bool {
		hasPassword := form.Find("input[type='password']").Length() > 0
		if !hasPassword {
			return false
		}

		// Check for common username/email fields
		hasUserField := false
		form.Find("input").EachWithBreak(func(i int, input *goquery.Selection) bool {
			typ, _ := input.Attr("type")
			name, _ := input.Attr("name")
			id, _ := input.Attr("id")

			if typ == "email" {
				hasUserField = true
				return false
			}

			if typ == "text" && (containsLoginKeyword(name) || containsLoginKeyword(id)) {
				hasUserField = true
				return false
			}
			return true
		})

		return hasUserField
	}).Length() > 0
}

func containsLoginKeyword(attr string) bool {
	attr = strings.ToLower(attr)
	return strings.Contains(attr, "user") || strings.Contains(attr, "login") || strings.Contains(attr, "email")
}

func ensureHTTPS(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "https://" + url
}
