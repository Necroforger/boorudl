package downloader

// SearchResultsByScore is used to sort search results by score
// It will place the highest scores first
type SearchResultsByScore []*Post

// Less implements the sortable interface
func (s SearchResultsByScore) Less(a, b int) bool {
	return s[a].Score > s[b].Score
}

// Swap implements the sortable interface
func (s SearchResultsByScore) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

// Len implements the sortable interface
func (s SearchResultsByScore) Len() int {
	return len(s)
}
