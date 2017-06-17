package extractor

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Paheal ...
type Paheal struct {
	client *http.Client
}

// NewPaheal ...
func NewPaheal(client *http.Client) *Paheal {
	return &Paheal{client: client}
}

// Search ...
func (p *Paheal) Search(q SearchQuery) (Posts, error) {
	posts := Posts{}
	q.Page++
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://rule34.paheal.net/post/list/%s/%d", url.QueryEscape(q.Tags), q.Page), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("ui-tnc-agreed", "true")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Accept-Language", "en-US,en;q=0.8")
	request.Header.Add("Accept", "*/*")

	resp, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".shm-thumb.thumb").Each(func(n int, s *goquery.Selection) {
		posts = append(posts, Post{
			Tags:         s.AttrOr("data-tags", ""),
			ThumbnailURL: s.Find("img").AttrOr("src", ""),
			ImageURL:     s.Find("a:contains('Image Only')").AttrOr("href", ""),
		})
	})

	if len(posts) == 0 {
		return nil, ErrNoPosts
	}

	return posts, nil
}
