package extractor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Endpoint constants for danbooru
const (
	EndpointDanbooru      = "https://danbooru.donmai.us"
	EndpointDanbooruPosts = EndpointDanbooru + "/posts.json"
	DanbooruMaxLimit      = 1000
)

// DanbooruPost stores the JSON information returned from the posts endpoint
type DanbooruPost struct {
	ID                  int         `json:"id"`
	CreatedAt           string      `json:"created_at"`
	UploaderID          int         `json:"uploader_id"`
	Score               int         `json:"score"`
	Source              string      `json:"source"`
	Md5                 string      `json:"md5"`
	LastCommentBumpedAt string      `json:"last_comment_bumped_at"`
	Rating              string      `json:"rating"`
	ImageWidth          int         `json:"image_width"`
	ImageHeight         int         `json:"image_height"`
	TagString           string      `json:"tag_string"`
	IsNoteLocked        bool        `json:"is_note_locked"`
	FavCount            int         `json:"fav_count"`
	FileExt             string      `json:"file_ext"`
	LastNotedAt         interface{} `json:"last_noted_at"`
	IsRatingLocked      bool        `json:"is_rating_locked"`
	ParentID            int         `json:"parent_id"`
	HasChildren         bool        `json:"has_children"`
	ApproverID          interface{} `json:"approver_id"`
	TagCountGeneral     int         `json:"tag_count_general"`
	TagCountArtist      int         `json:"tag_count_artist"`
	TagCountCharacter   int         `json:"tag_count_character"`
	TagCountCopyright   int         `json:"tag_count_copyright"`
	FileSize            int         `json:"file_size"`
	IsStatusLocked      bool        `json:"is_status_locked"`
	FavString           string      `json:"fav_string"`
	PoolString          string      `json:"pool_string"`
	UpScore             int         `json:"up_score"`
	DownScore           int         `json:"down_score"`
	IsPending           bool        `json:"is_pending"`
	IsFlagged           bool        `json:"is_flagged"`
	IsDeleted           bool        `json:"is_deleted"`
	TagCount            int         `json:"tag_count"`
	UpdatedAt           string      `json:"updated_at"`
	IsBanned            bool        `json:"is_banned"`
	PixivID             interface{} `json:"pixiv_id"`
	LastCommentedAt     string      `json:"last_commented_at"`
	HasActiveChildren   bool        `json:"has_active_children"`
	BitFlags            int         `json:"bit_flags"`
	UploaderName        string      `json:"uploader_name"`
	HasLarge            bool        `json:"has_large"`
	TagStringArtist     string      `json:"tag_string_artist"`
	TagStringCharacter  string      `json:"tag_string_character"`
	TagStringCopyright  string      `json:"tag_string_copyright"`
	TagStringGeneral    string      `json:"tag_string_general"`
	HasVisibleChildren  bool        `json:"has_visible_children"`
	ChildrenIds         interface{} `json:"children_ids"`
	FileURL             string      `json:"file_url"`
	LargeFileURL        string      `json:"large_file_url"`
	PreviewFileURL      string      `json:"preview_file_url"`
}

// Danbooru searches images from http://danbooru.donmai.us/
type Danbooru struct {
	client *http.Client
}

// NewDanbooru returns a pointer to new Danbooru struct
func NewDanbooru(client *http.Client) *Danbooru {
	return &Danbooru{
		client,
	}
}

func (d *Danbooru) searchURL(tags string, limit int, page int, random bool) string {

	rand := ""
	if random {
		rand = "&random=true"
	}

	return EndpointDanbooruPosts +
		fmt.Sprintf(
			"?limit=%d&page=%d"+rand+"&tags=%s",
			limit, page, url.QueryEscape(tags),
		)
}

// Search Searches the booru for images
// Limit
// Random
// Page
// Tags
func (d *Danbooru) Search(q SearchQuery) (results Posts, err error) {
	// Danbooru pages start at one, so add one to the page index.
	q.Page++

	if q.Limit > DanbooruMaxLimit {
		q.Limit = DanbooruMaxLimit
	}

	return d.search(q)
}

func (d *Danbooru) search(q SearchQuery) (results Posts, err error) {
	results = []Post{}

	res, err := d.client.Get(d.searchURL(q.Tags, q.Limit, q.Page, q.Random))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var posts []*DanbooruPost
	err = json.NewDecoder(res.Body).Decode(&posts)
	if err != nil {
		return
	}

	if len(posts) == 0 {
		return nil, ErrNoPosts
	}

	for _, v := range posts {

		// If the file does not exist, skip to the next post.
		if v.LargeFileURL == "" || v.PreviewFileURL == "" {
			continue
		}

		var imageURL string
		if strings.HasPrefix(v.LargeFileURL, "http://") ||
			strings.HasPrefix(v.LargeFileURL, "https://") {
			imageURL = v.LargeFileURL
		} else {
			imageURL = path.Join(EndpointDanbooru, v.LargeFileURL)
		}

		results = append(results, Post{
			ImageURL:     imageURL,
			ThumbnailURL: EndpointDanbooru + v.PreviewFileURL,
			Author:       v.UploaderName,
			ID:           v.ID,
			Tags:         v.TagString,
			Rating:       v.Rating,
			Score:        v.Score,
		})

	}

	return
}
