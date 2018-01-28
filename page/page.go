package page

import (
	"bytes"
	"net/url"

	"golang.org/x/net/html"
)

type Page struct {
	URL   *url.URL
	HTML  *html.Node
	Links []*url.URL
}

func NewPage(rawURL string, rawHTML string) (*Page, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBufferString(rawHTML)
	h, err := html.Parse(b)
	if err != nil {
		return nil, err
	}

	var links []*url.URL
	extract(&links, h)
	filtered := filter(&links, u)
	duplicatesRemoved := removeDuplicates(&filtered)

	return &Page{u, h, duplicatesRemoved}, nil
}

func (p Page) String() string {
	var b bytes.Buffer

	b.WriteString(p.URL.String())
	links := p.Links
	for i, link := range links {
		b.WriteString("\n")
		if i == len(links)-1 {
			b.WriteString("└──")
			b.WriteString(link.String())
			break
		}
		b.WriteString("├──")
		b.WriteString(link.String())
	}

	return b.String()
}

func extract(l *[]*url.URL, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				u, err := url.Parse(a.Val)
				if err != nil {
					break
				}
				*l = append(*l, u)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extract(l, c)
	}
}

func filter(links *[]*url.URL, u *url.URL) []*url.URL {
	var filtered []*url.URL

	for _, l := range *links {
		if !(l.Scheme == "http" || l.Scheme == "https" || !l.IsAbs()) {
			continue
		}
		if l.Fragment != "" {
			l.Fragment = ""
		}
		if l.RawQuery != "" {
			l.RawQuery = ""
		}
		if l.Host == "" {
			l.Host = u.Host
			l.Scheme = u.Scheme
		}
		if l.Host == u.Host {
			filtered = append(filtered, l)
		}
	}

	return filtered
}

func removeDuplicates(links *[]*url.URL) []*url.URL {
	found := map[string]bool{}
	result := []*url.URL{}

	for i := range *links {
		if found[(*links)[i].String()] {
			continue
		}
		found[(*links)[i].String()] = true
		result = append(result, (*links)[i])
	}

	return result
}
