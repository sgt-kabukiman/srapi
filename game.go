// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"strconv"
)

// Game represents a single game or romhack.
type Game struct {
	// unique ID of this game
	ID string

	// contains the japanese and international names; japanese is relatively
	// rare, international names are always present
	Names struct {
		International string
		Japanese      string
	}

	// unique abbreviation of the game, e.g. "smw" for Super Mario World
	Abbreviation string

	// link to the game page on speedrun.com
	Weblink string

	// year in which the game was released
	Released int

	// ruleset for the game
	Ruleset struct {
		ShowMilliseconds    bool           `json:"show-milliseconds"`
		RequireVerification bool           `json:"require-verification"`
		RequireVideo        bool           `json:"require-video"`
		RunTimes            []TimingMethod `json:"run-times"`
		DefaultTime         TimingMethod   `json:"default-time"`
		EmulatorsAllowed    bool           `json:"emulators-allowed"`
	}

	// whether or not this is a romhack
	Romhack bool

	// date and time when the game was added on speedrun.com; can be an empty
	// string for old games
	Created string

	// list of assets (images) for the game page design on speedrun.com, like
	// icons for trophies, background images etc.
	Assets map[string]*AssetLink

	// API links to related resources
	Links []Link

	// do not use this field directly, use the available methods
	PlatformsData interface{} `json:"platforms"`

	// do not use this field directly, use the available methods
	RegionsData interface{} `json:"regions"`

	// do not use this field directly, use the available methods
	ModeratorsData interface{} `json:"moderators"`

	// do not use this field directly, use the available methods
	CategoriesData interface{} `json:"categories"`

	// do not use this field directly, use the available methods
	LevelsData interface{} `json:"levels"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

// toGame transforms a data blob to a Game struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toGame(data interface{}) *Game {
	dest := Game{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// gameResponse models the actual API response from the server
type gameResponse struct {
	// the one game contained in the response
	Data Game
}

// GameByID tries to fetch a single game or romhack, identified by its ID.
// When an error is returned, the returned game is nil.
func GameByID(id string, embeds string) (*Game, *Error) {
	return fetchGame(request{"GET", "/games/" + id, nil, nil, nil, embeds})
}

// GameByAbbreviation tries to fetch a single game or romhack, identified by its
// abbreviation. This is convenient for resolving abbreviations, but as they can
// change (in constrast to the ID, which is fixed), it should be used with
// caution.
// When an error is returned, the returned game is nil.
func GameByAbbreviation(abbrev string, embeds string) (*Game, *Error) {
	return GameByID(abbrev, embeds)
}

// Series fetches the series the game belongs to. This returns only nil if there
// is broken data on speedrun.com.
func (g *Game) Series(embeds string) (*Series, *Error) {
	return fetchOneSeriesLink(firstLink(g, "series"), embeds)
}

// PlatformIDs returns a list of platform IDs this game is assigned to. This is
// always available; when the platforms are embedded, the IDs are collected from
// the respective objects.
func (g *Game) PlatformIDs() ([]string, *Error) {
	var result []string

	switch asserted := g.PlatformsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, something := range asserted {
			id, okay := something.(string)
			if okay {
				result = append(result, id)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	// TODO: skip the conversion back and forth and just assert our way through the available data
	case map[string]interface{}:
		platforms, err := g.Platforms()
		if err != nil {
			return result, err
		}

		for _, platform := range platforms {
			result = append(result, platform.ID)
		}
	}

	return result, nil
}

// Platforms returns a list of pointers to platform structs. If platforms were
// not embedded, they are fetched from the network, causing one request per
// platform.
func (g *Game) Platforms() ([]*Platform, *Error) {
	var result []*Platform

	switch asserted := g.PlatformsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		ids, err := g.PlatformIDs()
		if err != nil {
			return result, err
		}

		for _, id := range ids {
			platform, err := PlatformByID(id)
			if err != nil {
				return result, err
			}

			result = append(result, platform)
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		result = toPlatformCollection(asserted).platforms()
	}

	return result, nil
}

// RegionIDs returns a list of region IDs this game is assigned to. This is
// always available; when the regions are embedded, the IDs are collected from
// the respective objects.
func (g *Game) RegionIDs() ([]string, *Error) {
	var result []string

	switch asserted := g.RegionsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, something := range asserted {
			id, okay := something.(string)
			if okay {
				result = append(result, id)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	// TODO: skip the conversion back and forth and just assert our way through the available data
	case map[string]interface{}:
		regions, err := g.Regions()
		if err != nil {
			return result, err
		}

		for _, region := range regions {
			result = append(result, region.ID)
		}
	}

	return result, nil
}

// Regions returns a list of pointers to region structs. If regions were
// not embedded, they are fetched from the network, causing one request per
// region.
func (g *Game) Regions() ([]*Region, *Error) {
	var result []*Region

	switch asserted := g.RegionsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		ids, err := g.RegionIDs()
		if err != nil {
			return result, err
		}

		for _, id := range ids {
			region, err := RegionByID(id)
			if err != nil {
				return result, err
			}

			result = append(result, region)
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		result = toRegionCollection(asserted).regions()
	}

	return result, nil
}

// Categories returns the list of categories for this game. If they were not
// embedded, one additional request is performed and only then are filter and
// sort taken into account.
func (g *Game) Categories(filter *CategoryFilter, sort *Sorting, embeds string) ([]*Category, *Error) {
	var collection *CategoryCollection
	var err *Error

	if g.VariablesData == nil {
		collection, err = fetchCategoriesLink(firstLink(g, "categories"), filter, sort, embeds)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toCategoryCollection(g.CategoriesData)
	}

	return collection.categories(), nil
}

// Levels returns the list of levels for this game. If they were not embedded,
// one additional request is performed and only then is sort taken into account.
func (g *Game) Levels(sort *Sorting, embeds string) ([]*Level, *Error) {
	var collection *LevelCollection
	var err *Error

	if g.VariablesData == nil {
		collection, err = fetchLevelsLink(firstLink(g, "levels"), nil, sort, embeds)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toLevelCollection(g.CategoriesData)
	}

	return collection.levels(), nil
}

// Variables returns the list of variables for this game. If they were not
// embedded, one additional request is performed and only then is sort taken
// into account.
func (g *Game) Variables(sort *Sorting) ([]*Variable, *Error) {
	var collection *VariableCollection
	var err *Error

	if g.VariablesData == nil {
		collection, err = fetchVariablesLink(firstLink(g, "variables"), nil, sort)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toVariableCollection(g.VariablesData)
	}

	return collection.variables(), nil
}

// Romhacks returns a game collection containing the romhacks for the game.
// It always returns a collection, even when there are no romhacks or the game
// is itself a romhack.
func (g *Game) Romhacks(embeds string) (*GameCollection, *Error) {
	return fetchGamesLink(firstLink(g, "romhacks"), nil, nil, embeds)
}

// ModeratorMap returns a map of user IDs to their respective moderation levels.
// Note that due to limitations of the speedrun.com API, the mod levels are not
// available when moderators have been embedded. In this case, the resulting
// map containts UnknownModLevel for every user. If you need both, there is no
// other way than to perform two requests.
func (g *Game) ModeratorMap() map[string]GameModLevel {
	return recastToModeratorMap(g.ModeratorsData)
}

// Moderators returns a list of users that are moderators of the game. If
// moderators were not embedded, they will be fetched individually from the
// network.
func (g *Game) Moderators() ([]*User, *Error) {
	return recastToModerators(g.ModeratorsData)
}

// PrimaryLeaderboard fetches the primary leaderboard, if any, for the game.
// The result can be nil.
func (g *Game) PrimaryLeaderboard(options *LeaderboardOptions, embeds string) (*Leaderboard, *Error) {
	return fetchLeaderboardLink(firstLink(g, "leaderboard"), options, embeds)
}

// Records fetches a list of leaderboards for the game. This includes (by default)
// full-game and per-level leaderboards and is therefore paginated as a collection.
// This function always returns a LeaderboardCollection.
func (g *Game) Records(filter *LeaderboardFilter, embeds string) (*LeaderboardCollection, *Error) {
	return fetchLeaderboardsLink(firstLink(g, "records"), filter, nil, embeds)
}

// Runs fetches a list of runs done in the given game, optionally filtered
// and sorted. This function always returns a RunCollection.
func (g *Game) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(g, "runs"), filter, sort, embeds)
}

// for the 'hasLinks' interface
func (g *Game) links() []Link {
	return g.Links
}

// GameFilter represents the possible filtering options when fetching a list
// of games.
type GameFilter struct {
	Name         string
	Abbreviation string
	Released     int
	Platform     string
	Region       string
	Moderator    string
	Romhack      OptionalFlag
}

// applyToURL merged the filter into a URL.
func (gf *GameFilter) applyToURL(u *url.URL) {
	if gf == nil {
		return
	}

	values := u.Query()

	if len(gf.Name) > 0 {
		values.Set("name", gf.Name)
	}

	if len(gf.Abbreviation) > 0 {
		values.Set("abbreviation", gf.Abbreviation)
	}

	if gf.Released > 0 {
		values.Set("released", strconv.Itoa(gf.Released))
	}

	if len(gf.Platform) > 0 {
		values.Set("platform", gf.Platform)
	}

	if len(gf.Region) > 0 {
		values.Set("region", gf.Region)
	}

	if len(gf.Moderator) > 0 {
		values.Set("moderator", gf.Moderator)
	}

	gf.Romhack.applyToQuery("romhack", &values)

	u.RawQuery = values.Encode()
}

// GameCollection is one page of the entire game list. It consists of the
// games as well as some pagination information (like links to the next or
// previous page).
type GameCollection struct {
	Data       []Game
	Pagination Pagination
}

// Games retrieves a collection of games from the entire set of games on
// speedrun.com. In most cases, you will filter the game, as paging through
// *all* games takes A LOT of requests. For this, you should use BulkMode, which
// is not yet supported by this API.
func Games(f *GameFilter, s *Sorting, c *Cursor, embeds string) (*GameCollection, *Error) {
	return fetchGames(request{"GET", "/games", f, s, c, embeds})
}

// games returns a list of pointers to the games; used for cases where there is
// no pagination and the caller wants to return a flat slice of games instead of
// a collection (which would be misleading, as collections imply pagination).
func (gc *GameCollection) games() []*Game {
	var result []*Game

	for idx := range gc.Data {
		result = append(result, &gc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// games. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (gc *GameCollection) NextPage() (*GameCollection, *Error) {
	return gc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// games. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (gc *GameCollection) PrevPage() (*GameCollection, *Error) {
	return gc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (gc *GameCollection) fetchLink(name string) (*GameCollection, *Error) {
	next := firstLink(&gc.Pagination, name)
	if next == nil {
		return &GameCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchGamesLink(next, nil, nil, "")
}

// fetchGame fetches a single game from the network. If the request failed,
// the returned game is nil. Otherwise, the error is nil.
func fetchGame(request request) (*Game, *Error) {
	result := &gameResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchGameLink tries to fetch a given link and interpret the response as
// a single game. If the link is nil or the game could not be fetched,
// nil is returned.
func fetchGameLink(link requestable, embeds string) (*Game, *Error) {
	if !link.exists() {
		return nil, nil
	}

	return fetchGame(link.request(nil, nil, embeds))
}

// fetchGames fetches a list of games from the network. It always
// returns a collection, even when an error is returned.
func fetchGames(request request) (*GameCollection, *Error) {
	result := &GameCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchGamesLink tries to fetch a given link and interpret the response as
// a list of games. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchGamesLink(link requestable, filter filter, sort *Sorting, embeds string) (*GameCollection, *Error) {
	if !link.exists() {
		return &GameCollection{}, nil
	}

	return fetchGames(link.request(filter, sort, embeds))
}
