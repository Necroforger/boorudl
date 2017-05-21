package downloader

import (
	"errors"
	"net/http"
	"net/url"
)

// TODO implement a content filter

// Error Constants
var (
	ErrNoPosts = errors.New("err: query returned no posts")
)

// Constants for ratings
const (
	RatingSafe         = "s"
	RatingQuestionable = "q"
	RatingExplicit     = "e"
)

// Post represents an Image search result
type Post struct {
	ImageURL     string
	ThumbnailURL string
	Tags         string
	Author       string
	Rating       string
	Score        int
	ID           int
}

// Searcher represents a searchable booru that returns a
// *SearchResult.
type Searcher interface {
	Search(q SearchQuery) (Posts, error)
}

// SearchQuery is a searchquery used to provide optional values in searches
// 	Tags:   	Space separated tags to search for
// 	PostID:  Search for a specific post ID
// 	Limit:   The number of results to retrieve
// 	Page: 	The number of the next page.
// 	Random: 	Retrieve results randomly
type SearchQuery struct {
	Tags   string
	PostID string
	Limit  int
	Page   int
	Random bool
}

// NewSearchQuery returns the default values for a search query
func NewSearchQuery() SearchQuery {
	return SearchQuery{
		Limit:  1,
		Page:   0,
		PostID: "",
		Random: false,
		Tags:   "",
	}
}

// Search attempts to search images from the given booru link
func Search(URL string, q SearchQuery) (Posts, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	s := SearcherFromURL(u)

	res, err := s.Search(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SearcherFromURL returns a booru searcher from a hostname
func SearcherFromURL(u *url.URL) Searcher {
	var s Searcher

	// Check for custom searchers, if not, use the GenericBooru
	switch u.Host {
	case "danbooru.donmai.us":
		s = NewDanbooru(&http.Client{})
	default:
		s = NewGenericBooru(&http.Client{}, u.String(), 100)
	}

	return s
}
