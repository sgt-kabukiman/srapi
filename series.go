// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

// Series represents a series.
type Series struct {
	// unique ID
	ID string

	// names, international is always set, japanese sometimes
	Names struct {
		International string
		Japanese      string
	}

	// unique abbreviation of the game, e.g. "gta" for Grand Theft Auto
	Abbreviation string

	// link to the series page on speedrun.com
	Weblink string

	// list of assets (images) for the series page design on speedrun.com, like
	// icons for trophies, background images etc.
	Assets map[string]*AssetLink

	// API links to related resources
	Links []Link

	// do not use this field directly, use the available methods
	ModeratorsData interface{} `json:"moderators"`
}

// seriesResponse models the actual API response from the server
type seriesResponse struct {
	// the one series contained in the response
	Data Series
}

// SeriesByID tries to fetch a single series, identified by its ID.
// When an error is returned, the returned series is nil.
func SeriesByID(id string) (*Series, *Error) {
	return fetchOneSeries(request{"GET", "/series/" + id, nil, nil, nil})
}

// SeriesByAbbreviation tries to fetch a single series, identified by its
// abbreviation. This is convenient for resolving abbreviations, but as they can
// change (in constrast to the ID, which is fixed), it should be used with
// caution.
// When an error is returned, the returned series is nil.
func SeriesByAbbreviation(abbrev string) (*Series, *Error) {
	return SeriesByID(abbrev)
}

// Games fetches the list of games for the series, optionally filtering it.
func (s *Series) Games(filter *GameFilter, sort *Sorting) *GameCollection {
	return fetchGamesLink(firstLink(s, "games"), filter, sort)
}

// ModeratorMap returns a map of user IDs to their respective moderation levels.
// Note that due to limitations of the speedrun.com API, the mod levels are not
// available when moderators have been embedded. In this case, the resulting
// map containts UnknownModLevel for every user. If you need both, there is no
// other way than to perform two requests.
func (s *Series) ModeratorMap() map[string]GameModLevel {
	return recastToModeratorMap(&s.ModeratorsData)
}

// Moderators returns a list of users that are moderators of the series. If
// moderators were not embedded, they will be fetched individually from the
// network.
func (s *Series) Moderators() []*User {
	return recastToModerators(&s.ModeratorsData)
}

// for the 'hasLinks' interface
func (s *Series) links() []Link {
	return s.Links
}

// SeriesFilter represents the possible filtering options when fetching a list
// of series.
type SeriesFilter struct {
	Name         string
	Abbreviation string
	Moderator    string
}

// applyToURL merged the filter into a URL.
func (sf *SeriesFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(sf.Name) > 0 {
		values.Set("name", sf.Name)
	}

	if len(sf.Abbreviation) > 0 {
		values.Set("abbreviation", sf.Abbreviation)
	}

	if len(sf.Moderator) > 0 {
		values.Set("moderator", sf.Moderator)
	}

	u.RawQuery = values.Encode()
}

// SeriesCollection is one page of the entire series list. It consists of the
// series as well as some pagination information (like links to the next or
// previous page).
type SeriesCollection struct {
	Data       []Series
	Pagination Pagination
}

// ManySeries retrieves a collection of series.
func ManySeries(f *SeriesFilter, s *Sorting, c *Cursor) (*SeriesCollection, *Error) {
	return fetchManySeries(request{"GET", "/series", f, s, c})
}

// series returns a list of pointers to the series; used for cases where there is
// no pagination and the caller wants to return a flat slice of series instead of
// a collection (which would be misleading, as collections imply pagination).
func (sc *SeriesCollection) series() []*Series {
	var result []*Series

	for idx := range sc.Data {
		result = append(result, &sc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// series. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (sc *SeriesCollection) NextPage() (*SeriesCollection, *Error) {
	return sc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// series. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (sc *SeriesCollection) PrevPage() (*SeriesCollection, *Error) {
	return sc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (sc *SeriesCollection) fetchLink(name string) (*SeriesCollection, *Error) {
	next := firstLink(&sc.Pagination, name)
	if next == nil {
		return &SeriesCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchManySeries(next.request(nil, nil))
}

// fetchOneSeries fetches a single series from the network. If the request failed,
// the returned series is nil. Otherwise, the error is nil.
func fetchOneSeries(request request) (*Series, *Error) {
	result := &seriesResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchOneSeriesLink tries to fetch a given link and interpret the response as
// a single series. If the link is nil or the series could not be fetched,
// nil is returned.
func fetchOneSeriesLink(link requestable) *Series {
	if link == nil {
		return nil
	}

	series, _ := fetchOneSeries(link.request(nil, nil))
	return series
}

// fetchManySeries fetches a list of series from the network. It always
// returns a collection, even when an error is returned.
func fetchManySeries(request request) (*SeriesCollection, *Error) {
	result := &SeriesCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
