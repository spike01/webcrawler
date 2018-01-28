package page

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestLinks(t *testing.T) {
	var tests = []struct {
		description   string
		domain        string
		actualLinks   string
		expectedLinks []string
	}{
		{
			"Ignore other domains",
			"http://example.com",
			`<a href="http://example.com/something"> <a href="http://notexample.com">`,
			[]string{"http://example.com/something"},
		},
		{
			"Ignore duplicates",
			"http://example.com",
			`<a href="http://example.com/something"> <a href="http://example.com/something">`,
			[]string{"http://example.com/something"},
		},
		{
			"Handle non-absolute links",
			"http://example.com",
			`<a href="/something"> <a href="/something-else">`,
			[]string{"http://example.com/something", "http://example.com/something-else"},
		},
		{
			"Ignore subdomains",
			"http://example.com",
			`<a href="http://example.com/something"> <a href="http://sub.example.com/something-else">`,
			[]string{"http://example.com/something"},
		},
		{
			"Ignore non-http schemes",
			"http://example.com",
			`<a href="http://example.com/something"> <a href="mailto:contact@example.com"> <a href="tel:02038178870">`,
			[]string{"http://example.com/something"},
		},
		{
			"Ignore fragments",
			"http://example.com",
			`<a href="http://example.com/something"> <a href="http://example.com/something#fragment">`,
			[]string{"http://example.com/something"},
		},
		{
			"Ignore query strings",
			"http://example.com",
			`<a href="http://example.com/something?query=this"> <a href="http://example.com/something?query=that">`,
			[]string{"http://example.com/something"},
		},
	}

	for _, tt := range tests {
		rawHTML := templateLinks(tt.actualLinks)
		page, err := NewPage(tt.domain, rawHTML)
		if err != nil {
			t.Fail()
		}

		var links = tt.expectedLinks

		if !reflect.DeepEqual(page.Links, convertToUrls(links, t)) {
			t.Errorf("%s failed - expected %s, actual %s", tt.description, tt.expectedLinks, page.Links)
		}
	}
}

func convertToUrls(stringLinks []string, t *testing.T) []*url.URL {
	var expectedLinks []*url.URL

	for _, link := range stringLinks {
		u, err := url.Parse(link)
		if err != nil {
			t.Errorf("Malformed URL: %s", link)
		}

		expectedLinks = append(expectedLinks, u)
	}

	return expectedLinks
}

func templateLinks(links string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <title>Test</title>
</head>

<body>
    <header>
        <h1>Test</h1>
    </header>
		%s
</body>
</html>
	`, links)
}
