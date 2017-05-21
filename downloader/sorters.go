package downloader

import "sort"

// postsByScore is used to sort search results by score
// It will place the highest scores first
type postsByScore []Post

// Less implements the sortable interface
func (s postsByScore) Less(a, b int) bool {
	return s[a].Score > s[b].Score
}

// Swap implements the sortable interface
func (s postsByScore) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

// Len implements the sortable interface
func (s postsByScore) Len() int {
	return len(s)
}

/////////////////////////////////////////
//              POSTS
////////////////////////////////////////

// Posts supplies methods for filtering and managing a slice of search results.
type Posts []Post

// SortByScore sorts the slice in order of score from highest score to lowest.
func (p Posts) SortByScore() {
	sort.Sort(postsByScore(p))
}
