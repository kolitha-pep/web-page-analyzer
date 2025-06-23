package fetcher

import (
	"net/http"
	"time"
)

// IsLinkReachable checks the give n link is reachable or not by sending a HEAD request
func IsLinkReachable(link string) bool {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	// Send a HEAD request to the link
	resp, err := client.Head(link)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		return false
	}

	return true
}

func HttpGet(link string) (*http.Response, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Send a GET request to the link
	return client.Get(link)
}
