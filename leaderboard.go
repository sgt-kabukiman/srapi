// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"strconv"
)

// Leaderboard represents a leaderboard, i.e. a collection of ranked runs for a
// certain configuration of game, category, level and a few others.
type Leaderboard struct {
	// a link to the leaderboard on speedrun.com
	Weblink string

	// whether or not emulators are allowed
	Emulators bool

	// what platform, if any (otherwise this is empty), is the leaderboard limited to
	Platform string

	// what region, if any (otherwise this is empty), is the leaderboard limited to
	Region string

	// whether or not to only take runs with videos into account
	VideoOnly bool `json:"video-only"`

	// the timing method used to compare runs against each other
	Timing TimingMethod

	// the chosen variables (keys) and values (values) for the leaderboard, both
	// given as their respective IDs
	Values map[string]string

	// the runs, sorted from best to worst
	Runs []RankedRun

	// API links to related resources
	Links []Link

	// do not use this field directly, use the available methods
	PlatformsData interface{} `json:"platforms"`

	// do not use this field directly, use the available methods
	RegionsData interface{} `json:"regions"`

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	CategoryData interface{} `json:"category"`

	// do not use this field directly, use the available methods
	LevelData interface{} `json:"level"`

	// do not use this field directly, use the available methods
	PlayersData interface{} `json:"players"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

// RankedRun is a run with an assigned rank. As the rank only makes sense when
// a specific ruleset (video-only? realtime or ingame time? etc.) is applied,
// normal runs do not have a rank; only those in leaderboards have.
type RankedRun struct {
	// the embedded run
	Run

	// the rank, starting at 1
	Rank int
}

// leaderboardResponse models the actual API response from the server
type leaderboardResponse struct {
	// the one leaderboard contained in the response
	Data Leaderboard
}

// FullGameLeaderboard retrieves a the leaderboard for a specific game and one of
// its full-game categories. An error is returned if no category is given or if
// a per-level category is given. If no game is given, it is fetched automatically,
// but if you have it already at hand, you can save one request by specifying it.
func FullGameLeaderboard(game *Game, cat *Category, options *LeaderboardOptions, embeds string) (*Leaderboard, *Error) {
	if cat == nil {
		return nil, &Error{"", "", ErrorBadLogic, "No category given."}
	}

	if cat.Type != "per-game" {
		return nil, &Error{"", "", ErrorBadLogic, "The given category is not a full-game category."}
	}

	if game == nil {
		var err *Error

		game, err = cat.Game("")
		if err != nil {
			return nil, err
		}
	}

	return fetchLeaderboard(request{"GET", "/leaderboards/" + game.ID + "/category/" + cat.ID, options, nil, nil, embeds})
}

// LevelLeaderboard retrieves a the leaderboard for a specific game and one of
// its levels in a specific category. An error is returned if no category or
// level is given or if a full-game category is given. If no game is given, it
// is fetched automatically, but if you have it already at hand, you can save
// one request by specifying it.
func LevelLeaderboard(game *Game, cat *Category, level *Level, options *LeaderboardOptions, embeds string) (*Leaderboard, *Error) {
	if cat == nil {
		return nil, &Error{"", "", ErrorBadLogic, "No category given."}
	}

	if level == nil {
		return nil, &Error{"", "", ErrorBadLogic, "No level given."}
	}

	if cat.Type != "per-level" {
		return nil, &Error{"", "", ErrorBadLogic, "The given category is not a individual-level category."}
	}

	if game == nil {
		var err *Error

		game, err = level.Game("")
		if err != nil {
			return nil, err
		}
	}

	return fetchLeaderboard(request{"GET", "/leaderboards/" + game.ID + "/level/" + level.ID + "/" + cat.ID, options, nil, nil, embeds})
}

// Game returns the game that the leaderboard is for. If it was not embedded, it
// is fetched from the network. Except for broken data on speedrun.com, this
// should never return nil.
func (lb *Leaderboard) Game(embeds string) (*Game, *Error) {
	// we only have the game ID at hand
	asserted, okay := lb.GameData.(string)
	if okay {
		return GameByID(asserted, embeds)
	}

	return toGame(lb.GameData, true), nil
}

// Category returns the category that the leaderboard is for. If it was not
// embedded, it is fetched from the network. Except for broken data on
// speedrun.com, this should never return nil.
func (lb *Leaderboard) Category(embeds string) (*Category, *Error) {
	// we only have the category ID at hand
	asserted, okay := lb.CategoryData.(string)
	if okay {
		return CategoryByID(asserted, embeds)
	}

	return toCategory(lb.CategoryData, true), nil
}

// Level returns the level that the leaderboard is for. If it's a full-game
// leaderboard, nil is returned. If the level was not embedded, it is fetched
// from the network.
func (lb *Leaderboard) Level(embeds string) (*Level, *Error) {
	if lb.LevelData == nil {
		return nil, nil
	}

	// we only have the level ID at hand
	asserted, okay := lb.LevelData.(string)
	if okay {
		return LevelByID(asserted, embeds)
	}

	return toLevel(lb.LevelData, true), nil
}

// Platforms returns a list of all platforms that are used in the leaderboard.
// If they have not been embedded, an empty slice is returned.
func (lb *Leaderboard) Platforms() []*Platform {
	return toPlatformCollection(lb.PlatformsData).platforms()
}

// Regions returns a list of all regions that are used in the leaderboard.
// If they have not been embedded, an empty slice is returned.
func (lb *Leaderboard) Regions() []*Region {
	return toRegionCollection(lb.RegionsData).regions()
}

// Players returns a list of all players that are present in the leaderboard.
// If they have not been embedded, an empty slice is returned.
func (lb *Leaderboard) Players() []*Player {
	var result []*Player

	// players have not been embedded
	if lb.PlayersData == nil {
		return result
	}

	return recastToPlayerList(lb.PlayersData)
}

// Variables returns a list of all variables that are present in the leaderboard.
// If they have not been embedded, an empty slice is returned.
func (lb *Leaderboard) Variables() []*Variable {
	return toVariableCollection(lb.VariablesData).variables()
}

// for the 'hasLinks' interface
func (lb *Leaderboard) links() []Link {
	return lb.Links
}

// LeaderboardOptions are the options that can be used to further narrow down a
// leaderboard to only a subset of runs.
type LeaderboardOptions struct {
	// If set to a value >0, only this many places are returned. Note that there
	// can be multiple runs with the same rank, so you can end up with
	// len(runs) > Top. This value is ignored when set to anything else.
	Top int

	// The platform ID to restrict the leaderboard to.
	Platform string

	// The platform ID to restrict the leaderboard to.
	Region string

	// When set, can control if all or no runs are done on emulators.
	Emulators OptionalFlag

	// When set, can control if all or no runs are required to have a video.
	VideoOnly OptionalFlag

	// the timing method that should be used to compare runs; not all are
	// allowed for all games, a server-side error will be returned if an invalid
	// choice was made.
	Timing TimingMethod

	// ISO 8601 date; when given, only runs done before this date will be considerd
	Date string

	// map of variable IDs to value IDs
	Values map[string]string
}

// applyToURL merged the filter into a URL.
func (lo *LeaderboardOptions) applyToURL(u *url.URL) {
	if lo == nil {
		return
	}

	values := u.Query()

	if lo.Top > 0 {
		values.Set("top", strconv.Itoa(lo.Top))
	}

	if len(lo.Platform) > 0 {
		values.Set("platform", lo.Platform)
	}

	if len(lo.Region) > 0 {
		values.Set("region", lo.Region)
	}

	if len(lo.Timing) > 0 {
		values.Set("timing", string(lo.Timing))
	}

	if len(lo.Date) > 0 {
		values.Set("date", lo.Date)
	}

	lo.Emulators.applyToQuery("emulators", &values)
	lo.VideoOnly.applyToQuery("video-only", &values)

	for varID, valueID := range lo.Values {
		values.Set("var-"+varID, valueID)
	}

	u.RawQuery = values.Encode()
}

// LeaderboardFilter represents the possible filtering options when fetching a
// list of leaderboards.
type LeaderboardFilter struct {
	// If set to a value >0, only this many places are returned. Note that there
	// can be multiple runs with the same rank, so you can end up with
	// len(runs) > Top. This value is ignored when set to anything else.
	Top int

	// If set, can be used to skip returning empty leaderboards.
	SkipEmpty OptionalFlag
}

// applyToURL merged the filter into a URL.
func (lf *LeaderboardFilter) applyToURL(u *url.URL) {
	if lf == nil {
		return
	}

	values := u.Query()

	if lf.Top > 0 {
		values.Set("top", strconv.Itoa(lf.Top))
	}

	lf.SkipEmpty.applyToQuery("skip-empty", &values)

	u.RawQuery = values.Encode()
}

// LeaderboardCollection is one page of a paginated list of leaderboards. It
// consists of the leaderboards as well as some pagination information (like
// links to the next or previous page).
type LeaderboardCollection struct {
	Data       []Leaderboard
	Pagination Pagination
}

// leaderboards returns a list of pointers to the leaderboards; used for cases
// where there is no pagination and the caller wants to return a flat slice of
// leaderboards instead of a collection (which would be misleading, as
// collections imply pagination).
func (lc *LeaderboardCollection) leaderboards() []*Leaderboard {
	var result []*Leaderboard

	for idx := range lc.Data {
		result = append(result, &lc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// leaderboards. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (lc *LeaderboardCollection) NextPage() (*LeaderboardCollection, *Error) {
	return lc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// leaderboards. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (lc *LeaderboardCollection) PrevPage() (*LeaderboardCollection, *Error) {
	return lc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (lc *LeaderboardCollection) fetchLink(name string) (*LeaderboardCollection, *Error) {
	next := firstLink(&lc.Pagination, name)
	if next == nil {
		return &LeaderboardCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchLeaderboardsLink(next, nil, nil, "")
}

// fetchLeaderboard fetches a single leaderboard from the network. If the request
// failed, the returned leaderboard is nil. Otherwise, the error is nil.
func fetchLeaderboard(request request) (*Leaderboard, *Error) {
	result := &leaderboardResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchLeaderboardLink tries to fetch a given link and interpret the response as
// a single leaderboard. If the link is nil or the leaderboard could not be fetched,
// nil is returned.
func fetchLeaderboardLink(link requestable, options *LeaderboardOptions, embeds string) (*Leaderboard, *Error) {
	if !link.exists() {
		return nil, nil
	}

	return fetchLeaderboard(link.request(options, nil, embeds))
}

// fetchLeaderboards fetches a list of leaderboards from the network. It always
// returns a collection, even when an error is returned.
func fetchLeaderboards(request request) (*LeaderboardCollection, *Error) {
	result := &LeaderboardCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchLeaderboardsLink tries to fetch a given link and interpret the response as
// a list of leaderboards. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchLeaderboardsLink(link requestable, filter filter, sort *Sorting, embeds string) (*LeaderboardCollection, *Error) {
	if !link.exists() {
		return &LeaderboardCollection{}, nil
	}

	return fetchLeaderboards(link.request(filter, sort, embeds))
}
