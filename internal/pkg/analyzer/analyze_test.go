package analyzer

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestAnalyze_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockHTML := `
		<!DOCTYPE html>
		<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Main Title</h1>
				<h2>Subtitle</h2>
				<h3>Section Title</h3>
				<h4>Subsection Title</h4>
				<h5>Minor Title</h5>
				<h6>Least Important Title</h6>
				<h6>Least Important Title</h6>
				<h6>Least Important Title</h6>
				<a href="http://example.com/internal">Internal</a>
				<a href="https://external.com">External</a>
				<form><input type="password"/></form>
			</body>
		</html>`

	httpmock.RegisterResponder("GET", "http://example.com",
		httpmock.NewStringResponder(200, mockHTML))

	httpmock.RegisterResponder("HEAD", "http://example.com/internal",
		httpmock.NewStringResponder(200, ""))

	httpmock.RegisterResponder("HEAD", "https://external.com",
		httpmock.NewStringResponder(404, ""))

	result, err := AnalyzeWebPage("http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "HTML5", result.HtmlVersion)
	assert.Equal(t, "Test Page", result.Title)
	assert.Equal(t, 1, result.HeadTags["h1"])
	assert.Equal(t, 1, result.HeadTags["h2"])
	assert.Equal(t, 1, result.HeadTags["h3"])
	assert.Equal(t, 1, result.HeadTags["h4"])
	assert.Equal(t, 1, result.HeadTags["h5"])
	assert.Equal(t, 3, result.HeadTags["h6"])
	assert.Equal(t, 1, result.InternalLinks)
	assert.Equal(t, 1, result.ExternalLinks)
	assert.Equal(t, 1, result.BrokenLinks)
	assert.True(t, result.HasLoginForm)
}
