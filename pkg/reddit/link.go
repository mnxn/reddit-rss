package reddit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Link contains information about a link.
type Link struct {
	ApprovedBy          string        `json:"approved_by"`
	Archived            bool          `json:"archived"`
	Author              string        `json:"author"`
	AuthorFlairCSSClass string        `json:"author_flair_css_class"`
	AuthorFlairText     string        `json:"author_flair_text"`
	BannedBy            string        `json:"banned_by"`
	BodyHTML            string        `json:"body_html"`
	Clicked             bool          `json:"clicked"`
	ContestMode         bool          `json:"contest_mode"`
	Created             float64       `json:"created"`
	CreatedUtc          float64       `json:"created_utc"`
	Distinguished       string        `json:"distinguished"`
	Domain              string        `json:"domain"`
	Downs               int           `json:"downs"`
	Gilded              int           `json:"gilded"`
	Hidden              bool          `json:"hidden"`
	HideScore           bool          `json:"hide_score"`
	ID                  string        `json:"id"`
	IsSelf              bool          `json:"is_self"`
	Likes               bool          `json:"likes"`
	LinkFlairCSSClass   string        `json:"link_flair_css_class"`
	LinkFlairText       string        `json:"link_flair_text"`
	Locked              bool          `json:"locked"`
	Media               Media         `json:"media"`
	MediaEmbed          interface{}   `json:"media_embed"`
	ModReports          []interface{} `json:"mod_reports"`
	Name                string        `json:"name"`
	NumComments         int           `json:"num_comments"`
	NumReports          int           `json:"num_reports"`
	Over18              bool          `json:"over_18"`
	Permalink           string        `json:"permalink"`
	Quarantine          bool          `json:"quarantine"`
	RemovalReason       interface{}   `json:"removal_reason"`
	ReportReasons       []interface{} `json:"report_reasons"`
	Saved               bool          `json:"saved"`
	Score               int           `json:"score"`
	SecureMedia         SecureMedia   `json:"secure_media"`
	SecureMediaEmbed    interface{}   `json:"secure_media_embed"`
	SelftextHTML        string        `json:"selftext_html"`
	Selftext            string        `json:"selftext"`
	Stickied            bool          `json:"stickied"`
	Subreddit           string        `json:"subreddit"`
	SubredditID         string        `json:"subreddit_id"`
	SuggestedSort       string        `json:"suggested_sort"`
	Thumbnail           string        `json:"thumbnail"`
	Title               string        `json:"title"`
	URL                 string        `json:"url"`
	Ups                 int           `json:"ups"`
	UserReports         []interface{} `json:"user_reports"`
	Visited             bool          `json:"visited"`

	MediaMetadata       map[string]MediaMetadata `json:"media_metadata,omitempty"`
	GalleryData         GalleryData              `json:"gallery_data,omitempty"`
	CrossPostParentList []Link                   `json:"crosspost_parent_list"`
}

type MediaMetadataS struct {
	Width  int    `json:"x"`
	Height int    `json:"y"`
	U      string `json:"u"`
	Mp4    string `json:"mp4"`
	Gif    string `json:"gif"`
}

type MediaMetadata struct {
	S MediaMetadataS `json:"s"`
}

type GalleryDataItem struct {
	MediaID string `json:"media_id"`
	ID      int    `json:"id"`
}

type GalleryData struct {
	Items []GalleryDataItem `json:"items"`
}

const linkType = "t3"

type linkListing struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Children []struct {
			Kind string `json:"kind"`
			Data Link   `json:"data"`
		} `json:"children"`
		After  string      `json:"after"`
		Before interface{} `json:"before"`
	} `json:"data"`
}

// CommentOnLink posts a top-level comment to the given link. Requires the 'submit' OAuth scope.
func (c *Client) CommentOnLink(linkID string, text string) error {
	return c.commentOnThing(fmt.Sprintf("%s_%s", linkType, linkID), text)
}

// DeleteLink deletes a link submitted by the currently authenticated user. Requires the 'edit' OAuth scope.
func (c *Client) DeleteLink(linkID string) error {
	return c.deleteThing(fmt.Sprintf("%s_%s", linkType, linkID))
}

// EditLinkText edits the text of a self post by the currently authenticated user. Requires the 'edit' OAuth scope.
func (c *Client) EditLinkText(linkID string, text string) error {
	return c.editThingText(fmt.Sprintf("%s_%s", linkType, linkID), text)
}

// GetHotLinks retrieves a listing of hot links.
func (c *Client) GetHotLinks(subreddit string) ([]*Link, error) {
	return c.getLinks(subreddit, "hot")
}

// GetNewLinks retrieves a listing of new links.
func (c *Client) GetNewLinks(subreddit string) ([]*Link, error) {
	return c.getLinks(subreddit, "new")
}

// GetTopLinks retrieves a listing of top links.
func (c *Client) GetTopLinks(subreddit string) ([]*Link, error) {
	return c.getLinks(subreddit, "top")
}

// HideLink removes the given link from the user's default view of subreddit listings. Requires the 'report' OAuth scope.
func (c *Client) HideLink(linkID string) error {
	data := url.Values{}
	data.Set("id", fmt.Sprintf("%s_%s", linkType, linkID))
	url := fmt.Sprintf("%s/api/hide", baseAuthURL)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("HTTP Status Code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) getLinks(subreddit string, sort string) ([]*Link, error) {
	url := fmt.Sprintf("%s/r/%s/%s.json", baseURL, subreddit, sort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result linkListing
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var links []*Link
	for _, link := range result.Data.Children {
		links = append(links, &link.Data)
	}

	return links, nil
}
