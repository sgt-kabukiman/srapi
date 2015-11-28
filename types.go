// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// package version
const Version = "1.0"

// requestable describes anything that can turn itself into a request.
type requestable interface {
	exists() bool
	request(filter, *Sorting, string) request
}

// Link represent a generic link, most often used for API links.
type Link struct {
	Relation string `json:"rel"`
	URI      string `json:"uri"`
}

// checks if the link exists
func (l *Link) exists() bool {
	return l != nil
}

// request turns a link into a GET request.
func (l *Link) request(filter filter, sort *Sorting, embeds string) request {
	relURL := l.URI[len(BaseURL):]

	return request{"GET", relURL, filter, sort, nil, embeds}
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
	if c == nil {
		return
	}

	values := u.Query()

	if c.Offset > 0 {
		values.Set("offset", strconv.Itoa(c.Offset))
	}

	if c.Max > 0 {
		values.Set("max", strconv.Itoa(c.Max))
	}

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

// use this to denote no embeds
const NoEmbeds = ""

// Sorting represents the sorting options when requesting a list of items from the API.
type Sorting struct {
	OrderBy   string
	Direction Direction
}

// applyToURL merged the filter into a URL.
func (s *Sorting) applyToURL(u *url.URL) {
	if s == nil {
		return
	}

	values := u.Query()
	dir := "asc"

	if s.Direction == Descending {
		dir = "desc"
	}

	values.Set("orderby", s.OrderBy)
	values.Set("direction", dir)

	u.RawQuery = values.Encode()
}

// OptionalFlag represents a tri-state of true, false and undefined and is used for
// flags in collection filters.
type OptionalFlag int

const (
	// Undefined represents an unset flag
	Undefined OptionalFlag = iota

	// Yes is true.
	Yes

	// No is false.
	No
)

// applyToURL sets the flag in a query string if it's not Undefined.
func (f OptionalFlag) applyToQuery(name string, values *url.Values) {
	if f == Yes || f == No {
		values.Set(name, f.String())
	}
}

// String returns a string representation.
func (f OptionalFlag) String() string {
	switch f {
	case Yes:
		return "yes"

	case No:
		return "no"

	default:
		return ""
	}
}

// TimingMethod specifies what time was measured for a run.
type TimingMethod string

const (
	// TimingRealtime is realtime with loading times.
	TimingRealtime TimingMethod = "realtime"

	// TimingRealtimeWithoutLoads is realtime without loads.
	TimingRealtimeWithoutLoads TimingMethod = "realtime_noloads"

	// TimingIngameTime is using the in-game timer.
	TimingIngameTime TimingMethod = "ingame"
)

// GameModLevel determines the power level of a moderator.
type GameModLevel string

const (
	// NormalModerator users can do game-related things.
	NormalModerator GameModLevel = "moderator"

	// SuperModerator users can appoint other moderators.
	SuperModerator GameModLevel = "super-moderator"

	// UnknownModLevel is used for when moderators have been embedded and there
	// is no information available about their actual level.
	UnknownModLevel GameModLevel = "unknown"
)

// dateLayout describes the format for ISO 8601 dates
var dateLayout = "2006-01-02"

// DateParseError is an error that occurs when a JSON string is not a valid date
var DateParseError = errors.New(`DateParseError: should be a string formatted as "2006-01-02"`)

// Date is a custom time.Time wrapper that allows dates without times in JSON
// documents.
type Date struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format(dateLayout) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) != len(`"2006-01-02"`) {
		return DateParseError
	}

	ret, err := time.Parse(dateLayout, s[1:11])
	if err != nil {
		return err
	}

	d.Time = ret

	return nil
}

// DurationParseError is an error that occurs when a JSON value is not a valid float value
var DurationParseError = errors.New(`DurationParseError: value should be a valid float`)

// Duration is a custom time.Time wrapper that allows dates without times in JSON
// documents.
type Duration struct {
	time.Duration
}

// MarshalJSON implements the json.Marshaler interface
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.3f", d.Seconds())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Duration) UnmarshalJSON(b []byte) error {
	parsed, err := strconv.ParseFloat(string(b), 32)

	if err != nil {
		return DurationParseError
	}

	d.Duration = time.Duration(parsed * float64(time.Second))

	return nil
}

// Format returns a human readable time in the form of "[[HH:]MM:]SS[.MS]".
func (d *Duration) Format() string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % (3600)
	seconds := int(d.Seconds()) % 60
	milli := (d.Seconds() - float64(int(d.Seconds())))

	list := make([]string, 0)

	if hours > 0 {
		list = append(list, fmt.Sprintf("%02d", hours))
	}

	if len(list) > 0 || minutes > 0 {
		list = append(list, fmt.Sprintf("%02d", minutes))
	}

	if len(list) > 0 || seconds > 0 {
		list = append(list, fmt.Sprintf("%02d", seconds))
	}

	formatted := strings.TrimPrefix(strings.Join(list, ":"), "0")

	if milli >= 0.001 {
		formatted += fmt.Sprintf(".%02d", int(milli*1000 + 0.5)) // +0.5 for easy rounding
	}

	return formatted
}
