package crawler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/spike01/webcrawler/client"
	"github.com/spike01/webcrawler/page"
)

type MockClient struct {
	mockResponses map[string]string
}

func (m *MockClient) Do(r *http.Request) (*http.Response, error) {
	body := bytes.NewBufferString(m.mockResponses[r.URL.String()])
	return &http.Response{
		Body: ioutil.NopCloser(body),
	}, nil
}

func TestCrawl(t *testing.T) {
	mockResponses := make(map[string]string)
	mockResponses["http://example.com"] = templateLinks(`<a href="http://example.com/page1">`)
	mockResponses["http://example.com/page1"] = templateLinks(`<a href="http://example.com/page2"> <a href="http://example.com/page3">`)
	mockResponses["http://example.com/page2"] = templateLinks(`<a href="http://example.com"> <a href="http://example.com/page1>`)
	mockResponses["http://example.com/page3"] = templateLinks(`<a href="http://example.com"> <a href="http://example.com/page2>`)

	mockClient := client.NewClient(&MockClient{mockResponses})
	crawler := NewCrawler(mockClient, false)

	pages, err := crawler.Crawl("http://example.com")
	if err != nil {
		t.Fail()
	}

	expected := buildExpected(mockResponses, t)

	if len(expected) != len(pages) {
		t.Fatalf("Failed - expected length %s, actual length %s", len(expected), len(pages))
	}

	// Both slices should have same elements, not necessarily in same order
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].URL.String() < expected[j].URL.String()
	})
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].URL.String() < pages[j].URL.String()
	})

	for i := range expected {
		if !reflect.DeepEqual(expected[i].Links, pages[i].Links) {
			t.Fatalf("Failed - expected %s, actual %s", expected[i].Links, pages[i].Links)
		}
	}
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

func buildExpected(mockResponses map[string]string, t *testing.T) []*page.Page {
	var expected []*page.Page
	for k, v := range mockResponses {
		p, err := page.NewPage(k, v)
		if err != nil {
			t.Fail()
		}
		expected = append(expected, p)
	}
	return expected
}
