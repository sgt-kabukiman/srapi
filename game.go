// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"strconv"
)

type Game struct {
	Id    string
	Names struct {
		International string
		Japanese      string
	}
	Abbreviation string
	Weblink      string
	Released     int
	Ruleset      struct {
		ShowMilliseconds    bool           `json:"show-milliseconds"`
		RequireVerification bool           `json:"require-verification"`
		RequireVideo        bool           `json:"require-video"`
		RunTimes            []TimingMethod `json:"run-times"`
		DefaultTime         TimingMethod   `json:"default-time"`
		EmulatorsAllowed    bool           `json:"emulators-allowed"`
	}
	Romhack bool
	Created string
	Assets  map[string]*AssetLink
	Links   []Link

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

func toGame(data interface{}) *Game {
	dest := Game{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

func toGameCollection(data interface{}) *GameCollection {
	tmp := &GameCollection{}
	recast(data, tmp)

	return tmp
}

type gameResponse struct {
	Data Game
}

func GameById(id string) (*Game, *Error) {
	return fetchGame(request{"GET", "/games/" + id, nil, nil, nil})
}

func GameByAbbreviation(abbrev string) (*Game, *Error) {
	return GameById(abbrev)
}

func (self *Game) Series() *Series {
	return fetchOneSeriesLink(firstLink(self, "series"))
}

func (self *Game) PlatformIds() []string {
	result := make([]string, 0)

	switch asserted := self.PlatformsData.(type) {
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
		for _, platform := range self.Platforms() {
			result = append(result, platform.Id)
		}
	}

	return result
}

func (self *Game) Platforms() []*Platform {
	result := make([]*Platform, 0)

	switch asserted := self.PlatformsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, id := range self.PlatformIds() {
			platform, err := PlatformById(id)
			if err == nil {
				result = append(result, platform)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		return toPlatformCollection(asserted).platforms()
	}

	return result
}

func (self *Game) RegionIds() []string {
	result := make([]string, 0)

	switch asserted := self.RegionsData.(type) {
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
		for _, region := range self.Regions() {
			result = append(result, region.Id)
		}
	}

	return result
}

func (self *Game) Regions() []*Region {
	result := make([]*Region, 0)

	switch asserted := self.RegionsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, id := range self.RegionIds() {
			region, err := RegionById(id)
			if err == nil {
				result = append(result, region)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		return toRegionCollection(asserted).regions()
	}

	return result
}

func (self *Game) Categories(filter *CategoryFilter, sort *Sorting) []*Category {
	var collection *CategoryCollection

	if self.VariablesData == nil {
		collection = fetchCategoriesLink(firstLink(self, "categories"), filter, sort)
	} else {
		collection = toCategoryCollection(self.CategoriesData)
	}

	return collection.categories()
}

func (self *Game) Levels(sort *Sorting) []*Level {
	var collection *LevelCollection

	if self.VariablesData == nil {
		collection = fetchLevelsLink(firstLink(self, "levels"), nil, sort)
	} else {
		collection = toLevelCollection(self.CategoriesData)
	}

	return collection.levels()
}

func (self *Game) Variables(sort *Sorting) []*Variable {
	var collection *VariableCollection

	if self.VariablesData == nil {
		collection = fetchVariablesLink(firstLink(self, "variables"), nil, sort)
	} else {
		collection = toVariableCollection(self.VariablesData)
	}

	return collection.variables()
}

func (self *Game) Romhacks() *GameCollection {
	return fetchGamesLink(firstLink(self, "romhacks"), nil, nil)
}

func (self *Game) ModeratorMap() map[string]GameModLevel {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := self.ModeratorsData.(map[string]GameModLevel)
	if okay {
		return assertedMap
	}

	// maybe we got a list of embedded users
	result := make(map[string]GameModLevel, 0)
	tmp := UserCollection{}

	if recast(self.ModeratorsData, &tmp) == nil {
		for _, user := range tmp.users() {
			result[user.Id] = UnknownModLevel
		}
	}

	return result
}

func (self *Game) Moderators() []*User {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := self.ModeratorsData.(map[string]GameModLevel)
	if okay {
		result := make([]*User, 0)

		for userId := range assertedMap {
			user, err := UserById(userId)
			if err == nil {
				result = append(result, user)
			}
		}

		return result
	}

	// maybe we got a list of embedded users
	return toUserCollection(self.ModeratorsData).users()
}

func (self *Game) PrimaryLeaderboard(options *LeaderboardOptions) *Leaderboard {
	return fetchLeaderboardLink(firstLink(self, "leaderboard"), options)
}

func (self *Game) Records(filter *LeaderboardFilter) *LeaderboardCollection {
	return fetchLeaderboardsLink(firstLink(self, "records"), filter, nil)
}

func (self *Game) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "runs"), filter, sort)
}

// for the 'hasLinks' interface
func (self *Game) links() []Link {
	return self.Links
}

type GameCollection struct {
	Data       []Game
	Pagination Pagination
}

func (self *GameCollection) games() []*Game {
	result := make([]*Game, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type GameFilter struct {
	Name         string
	Abbreviation string
	Released     int
	Platform     string
	Region       string
	Moderator    string
	Romhack      *bool
}

func (self *GameFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.Name) > 0 {
		values.Set("name", self.Name)
	}

	if len(self.Abbreviation) > 0 {
		values.Set("abbreviation", self.Abbreviation)
	}

	if self.Released > 0 {
		values.Set("released", strconv.Itoa(self.Released))
	}

	if len(self.Platform) > 0 {
		values.Set("platform", self.Platform)
	}

	if len(self.Region) > 0 {
		values.Set("region", self.Region)
	}

	if len(self.Moderator) > 0 {
		values.Set("moderator", self.Moderator)
	}

	if self.Romhack != nil {
		if *self.Romhack {
			values.Set("romhack", "yes")
		} else {
			values.Set("romhack", "no")
		}
	}

	u.RawQuery = values.Encode()
}

func Games(f *GameFilter, s *Sorting, c *Cursor) (*GameCollection, *Error) {
	return fetchGames(request{"GET", "/games", f, s, c})
}

func (self *GameCollection) NextPage() (*GameCollection, *Error) {
	return self.fetchLink("next")
}

func (self *GameCollection) PrevPage() (*GameCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *GameCollection) fetchLink(name string) (*GameCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchGames(next.request(nil, nil))
}

func fetchGame(request request) (*Game, *Error) {
	result := &gameResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchGameLink(link *Link) *Game {
	if link == nil {
		return nil
	}

	game, _ := fetchGame(link.request(nil, nil))
	return game
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchGames(request request) (*GameCollection, *Error) {
	result := &GameCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func fetchGamesLink(link *Link, filter filter, sort *Sorting) *GameCollection {
	if link == nil {
		return &GameCollection{}
	}

	collection, _ := fetchGames(link.request(filter, sort))
	return collection
}
