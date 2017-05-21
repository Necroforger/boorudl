package downloader

import "sort"

/////////////////////////////////////////
//              POSTS
////////////////////////////////////////

// Posts supplies methods for filtering and managing a slice of search results.
type Posts []Post

// SortByScore sorts the slice in order of score from highest score to lowest.
func (p Posts) SortByScore() Posts {
	sort.Sort(postsByScore(p))
	return p
}

// SortByID sorts the slice in order if the post IDS from lowest to highest.
func (p Posts) SortByID() Posts {
	sort.Sort(postsByID(p))
	return p
}

// SortByRating sorts the slice in the alphabetical order of the post's rating
func (p Posts) SortByRating() Posts {
	sort.Sort(postsByRating(p))
	return p
}

// RemoveByRating removes all posts with the specified rating from the slice
// The available ratings are stored in constants in downloader.go
// It will return the an array of posts that were removed.
func (p Posts) RemoveByRating(rating string) (removed Posts) {
	removed = Posts{}

	for i, post := range p {
		if post.Rating == rating {
			p = append(p[i:], p[i+1:]...)
		}
	}

	return
}

/////////////////////////////////////////
//              postsByScore
////////////////////////////////////////

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
//              postsByRating
////////////////////////////////////////

type postsByRating []Post

func (s postsByRating) Less(a, b int) bool {
	return s[a].Rating < s[b].Rating
}

// Swap implements the sortable interface
func (s postsByRating) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

// Len implements the sortable interface
func (s postsByRating) Len() int {
	return len(s)
}

/////////////////////////////////////////
//              postsById
////////////////////////////////////////

type postsByID []Post

func (s postsByID) Less(a, b int) bool {
	return s[a].ID < s[b].ID
}

// Swap implements the sortable interface
func (s postsByID) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

// Len implements the sortable interface
func (s postsByID) Len() int {
	return len(s)
}
