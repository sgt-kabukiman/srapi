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
	Romhack    bool
	Platforms  []string
	Regions    []string
	Moderators map[string]GameModLevel
	Created    string
	Assets     map[string]*AssetLink
	Links      []Link
}

type AssetLink struct {
	URI    string
	Width  int
	Height int
}

type gameResponse struct {
	Data Game
}

func GameById(id string) (*Game, *Error) {
	request := request{"GET", "/games/" + id, nil, nil, nil}
	result := &gameResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func GameByAbbreviation(abbrev string) (*Game, *Error) {
	return GameById(abbrev)
}

// for the 'hasLinks' interface
func (self *Game) links() []Link {
	return self.Links
}

type GameCollection struct {
	Data       []Game
	Pagination Pagination
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

	return fetchGames(next.request())
}

func fetchGames(request request) (*GameCollection, *Error) {
	result := &GameCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
