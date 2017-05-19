package downloader

// SearchResult represents an Image search result
type SearchResult struct {
	FileExtension string
	ImageURL      string
	ThumbnailURL  string
	Tags          string
	Author        string
	ID            int
}

// Searchable represents a searchable booru that returns a
// *SearchResult.
type Searchable interface {
	Search(q SearchQuery) *SearchResult
}

// SearchQuery is a searchquery used to provide optional values in searches
// Tags:   	Space separated tags to search for
// Limit:   The number of results to retrieve
// Page: 	The number of the next page.
// Random: 	Retrieve results randomly
type SearchQuery struct {
	Tags   string
	Limit  int
	Page   int
	Random bool
}
