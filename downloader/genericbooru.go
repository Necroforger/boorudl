package downloader

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Generic Endpoint for booru
const (
	EndpointGenericPosts = "/index.php?page=dapi&s=post&q=index"
)

// GenericPosts is are the generic booru posts
type GenericPosts struct {
	Count  string        `xml:"count,attr"`
	Offset string        `xml:"offset,attr"`
	Post   []GenericPost `xml:"post"`
}

// GenericPost is a generic booru post
type GenericPost struct {
	ParentID      string `xml:"parent_id,attr"`
	Score         string `xml:"score,attr"`
	SampleWidth   string `xml:"sample_width,attr"`
	HasChildren   string `xml:"has_children,attr"`
	HasComments   string `xml:"has_comments,attr"`
	PreviewHeight string `xml:"preview_height,attr"`
	Md5           string `xml:"md5,attr"`
	HasNotes      string `xml:"has_notes,attr"`
	Height        string `xml:"height,attr"`
	Source        string `xml:"source,attr"`
	PreviewURL    string `xml:"preview_url,attr"`
	SampleHeight  string `xml:"sample_height,attr"`
	ID            string `xml:"id,attr"`
	Rating        string `xml:"rating,attr"`
	PreviewWidth  string `xml:"preview_width,attr"`
	Status        string `xml:"status,attr"`
	FileURL       string `xml:"file_url,attr"`
	Tags          string `xml:"tags,attr"`
	Width         string `xml:"width,attr"`
	CreatedAt     string `xml:"created_at,attr"`
	SampleURL     string `xml:"sample_url,attr"`
	Change        string `xml:"change,attr"`
	CreatorID     string `xml:"creator_id,attr"`
}

// GenericBooru fetches images from any booru that uses
// /index.php?page=dapi&s=post&q=index
// As its API path
type GenericBooru struct {
	client       *http.Client
	EndpointRoot string
	SearchLimit  int
}

// NewGenericBooru returns a pointer to a new GenericBooru struct
//		client		:		An http client. It can be used to set custom timeout durations for requests
// 		rootEndpoint:		The base URL of the site
//		searchLimit :		The maximum number of posts you can retrieve per request. The default search limit for boorus is 100.
func NewGenericBooru(client *http.Client, rootEndpoint string, searchLimit int) *GenericBooru {
	return &GenericBooru{
		client:       client,
		EndpointRoot: rootEndpoint,
		SearchLimit:  searchLimit,
	}
}

func (g *GenericBooru) searchURL(limit int, pageid int, tags string, postID string) (string, error) {

	u, err := url.Parse(g.EndpointRoot + EndpointGenericPosts)
	if err != nil {
		return "", err
	}

	// Set the URL scheme if it isn't already set.
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	q := u.Query()
	q.Set("limit", fmt.Sprint(limit))
	q.Set("pid", fmt.Sprint(pageid))
	q.Set("tags", tags)
	if postID != "" {
		q.Set("id", fmt.Sprint(postID))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// Search a generic booru
func (g *GenericBooru) Search(q SearchQuery) ([]*SearchResult, error) {
	// results := []*SearchResult{}

	// Split the limit into multiple queries if its beyond the supported range
	// if q.Limit > g.SearchLimit {
	// 	numpages := q.Limit / g.SearchLimit
	// 	fmt.Println(q.Page, " to ", q.Page+numpages)
	// 	numRequests := q.Page + numpages
	// 	for i := q.Page; i < numRequests; i++ {
	// 		q.Page = i
	// 		q.Limit = g.SearchLimit

	// 		res, err := g.search(q)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		// Return if there are no more images to find
	// 		if len(res) == 0 {
	// 			return results, nil
	// 		}

	// 		results = append(results, res...)
	// 	}
	// 	return results, nil
	// }

	if q.Limit > g.SearchLimit {
		q.Limit = g.SearchLimit
	}

	return g.search(q)
}

func (g *GenericBooru) search(q SearchQuery) ([]*SearchResult, error) {

	results := []*SearchResult{}

	searchURL, err := g.searchURL(q.Limit, q.Page, q.Tags, q.PostID)
	if err != nil {
		return nil, err
	}

	// Request XML post list
	res, err := g.client.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode XML data from response
	var booruPosts GenericPosts
	err = xml.NewDecoder(res.Body).Decode(&booruPosts)
	if err != nil {
		return nil, err
	}

	// Check that posts exist
	posts := booruPosts.Post
	if posts == nil {
		return nil, ErrNoPosts
	}

	// Convert the generic posts into an array of SearchResults
	for _, v := range posts {

		// Parse the URLs and set the schemes to http
		furl, err := url.Parse(v.FileURL)
		if err != nil {
			continue
		}
		if furl.Scheme == "" {
			furl.Scheme = "http"
		}
		purl, err := url.Parse(v.PreviewURL)
		if err != nil {
			continue
		}
		if purl.Scheme == "" {
			purl.Scheme = "http"
		}

		// Convert the ID to an integer value
		ID, err := strconv.Atoi(v.ID)
		if err != nil {
			ID = -1
		}

		score, err := strconv.Atoi(v.Score)
		if err != nil {
			score = -1
		}

		results = append(results, &SearchResult{
			Author:       v.CreatorID,
			ID:           ID,
			Tags:         v.Tags,
			ImageURL:     furl.String(),
			ThumbnailURL: purl.String(),
			Rating:       v.Rating,
			Score:        score,
		})
	}

	return results, nil
}
