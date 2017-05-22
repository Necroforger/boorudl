package downloader

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

// GooglePost ...
type GooglePost struct {
	Cb  int    `json:"cb"`
	ID  string `json:"id"`  // Image ID
	Isu string `json:"isu"` // Image source
	Itg int    `json:"itg"`
	Ity string `json:"ity"` // File extension
	Oh  int    `json:"oh"`
	Ou  string `json:"ou"` // Image URL
	Ow  int    `json:"ow"`
	Pt  string `json:"pt"`  // Picture Title
	Rid string `json:"rid"` // ID
	Ru  string `json:"ru"`
	S   string `json:"s"`
	St  string `json:"st"`
	Th  int    `json:"th"`
	Tu  string `json:"tu"` // Thumbnail URL
	Tw  int    `json:"tw"`
}

// GoogleImages ...
type GoogleImages struct {
	client *http.Client
}

// NewGoogleImages returns a pointer to a new google images object
func NewGoogleImages(client *http.Client) *GoogleImages {
	return &GoogleImages{
		client: client,
	}
}

// SearchURL generates a URL to search from
func (g *GoogleImages) SearchURL(query string) string {
	return "https://www.google.com/search?biw=985&bih=661&tbm=isch&sa=1&q=" + url.QueryEscape(query)
}

// Search implements the searcher interface
func (g *GoogleImages) Search(q SearchQuery) (Posts, error) {
	return g.search(q)
}

func (g *GoogleImages) search(q SearchQuery) (Posts, error) {
	if q.Page > 0 {
		return nil, ErrNoPosts
	}

	results := Posts{}

	request, err := http.NewRequest("GET", g.SearchURL(q.Tags), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	request.Header.Add("Postman-Token", "d33e606b-d057-738d-d94b-d37c2d4a4634")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Accept-Language", "en-US,en;q=0.8")
	request.Header.Add("Accept", "*/*")
	// request.Header.Add("Accept-Encoding", "gzip, deflate, sdch")

	resp, err := g.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body string
	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		body = string(b)
	} else {
		return nil, err
	}

	// fmt.Println(body)

	jsonStart := `class="rg_meta">{`
	jsonEnd := `</div>`

	for i := 0; i < q.Limit; i++ {
		var jsonData string

		startIndex := strings.Index(body, jsonStart)
		if startIndex < 0 {
			break
		}

		endIndex := strings.Index(body[startIndex:], jsonEnd)
		if endIndex < 0 {
			break
		}
		endIndex += startIndex

		jsonData = body[startIndex+len(jsonStart)-1 : endIndex]
		body = body[endIndex+len(jsonEnd):]

		var goo GooglePost
		if err := json.Unmarshal([]byte(jsonData), &goo); err == nil {
			// fmt.Println("Unmarshaled", jsonData)
			results = append(results, Post{
				ImageURL:     goo.Ou,
				ThumbnailURL: goo.Tu,
				Rating:       "s",
				Score:        0,
				ID:           rand.Int(), // Generate a random ID, because google images uses strings for its ids
			})
		}
	}

	if len(results) == 0 {
		return nil, ErrNoPosts
	}

	return results, nil
}
