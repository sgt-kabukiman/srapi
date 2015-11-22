// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"strconv"
)

// requestable describes anything that can turn itself into a request.
type requestable interface {
	request(filter, *Sorting) request
}

// Link represent a generic link, most often used for API links.
type Link struct {
	Relation string `json:"rel"`
	URI      string `json:"uri"`
}

// request turns a link into a GET request.
func (l *Link) request(filter filter, sort *Sorting) request {
	relURL := l.URI[len(BaseURL):]

	return request{"GET", relURL, filter, sort, nil}
}

// AssetLink is a link pointing to an image, having width and height values.
type AssetLink struct {
	Link

	Width  int
	Height int
}

// Pagination contains information on how to navigate through multiple pages
// of results.
type Pagination struct {
	Offset int
	Max    int
	Size   int
	Links  []Link
}

// for the 'hasLinks' interface
func (p *Pagination) links() []Link {
	return p.Links
}

// filter describes anything that can apply itself to a URL.
type filter interface {
	applyToURL(*url.URL)
}

// Cursor represents the current position in a collection.
type Cursor struct {
	Offset int
	Max    int
}

// applyToURL merged the filter into a URL.
func (c *Cursor) applyToURL(u *url.URL) {
	values := u.Query()

	values.Set("offset", strconv.Itoa(c.Offset))
	values.Set("max", strconv.Itoa(c.Max))

	u.RawQuery = values.Encode()
}

// Direction is a sorting order
type Direction int

const (
	// Ascending sorts a...z
	Ascending Direction = iota

	// Descending sorts z...a
	Descending
)

// Sorting represents the sorting options when requesting a list of items from the API.
type Sorting struct {
	OrderBy   string
	Direction Direction
}

// applyToURL merged the filter into a URL.
func (s *Sorting) applyToURL(u *url.URL) {
	values := u.Query()
	dir := "asc"

	if s.Direction == Descending {
		dir = "desc"
	}

	values.Set("orderby", s.OrderBy)
	values.Set("direction", dir)

	u.RawQuery = values.Encode()
}

// TimingMethod specifies what time was measured for a run.
type TimingMethod string

const (
	// TimingRealtime is realtime with loading times.
	TimingRealtime TimingMethod = "realtime"

	// TimingRealtimeWithoutLoads is realtime without loads.
	TimingRealtimeWithoutLoads = "realtime_noloads"

	// TimingIngameTime is using the in-game timer.
	TimingIngameTime = "ingame"
)

// GameModLevel determines the power level of a moderator.
type GameModLevel string

const (
	// NormalModerator users can do game-related things.
	NormalModerator GameModLevel = "moderator"

	// SuperModerator users can appoint other moderators.
	SuperModerator = "super-moderator"

	// UnknownModLevel is used for when moderators have been embedded and there
	// is no information available about their actual level.
	UnknownModLevel = "unknown"
)
