package srapi

import (
	"encoding/json"
	"net/url"
)

type Leaderboard struct {
	Weblink   string
	Emulators bool
	Platform  string
	Region    string
	VideoOnly bool `json:"video-only"`
	Timing    TimingMethod
	Values    map[string]string
	Runs      []RankedRun
	Links     []Link

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
}

type RankedRun struct {
	Run

	Rank int
}

type leaderboardResponse struct {
	Data Leaderboard
}

func FullGameLeaderboard(game *Game, cat *Category) (*Leaderboard, *Error) {
	if cat == nil {
		return nil, &Error{"", "", ErrorBadLogic, "No category given."}
	}

	if cat.Type != "per-game" {
		return nil, &Error{"", "", ErrorBadLogic, "The given category is not a full-game category."}
	}

	if game == nil {
		game = cat.Game()
	}

	return fetchLeaderboard(request{"GET", "/leaderboards/" + game.Id + "/category/" + cat.Id, nil, nil, nil})
}

func LevelLeaderboard(game *Game, cat *Category, level *Level) (*Leaderboard, *Error) {
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
		game = level.Game()
	}

	return fetchLeaderboard(request{"GET", "/leaderboards/" + game.Id + "/level/" + level.Id + "/" + cat.Id, nil, nil, nil})
}

func fetchLeaderboard(request request) (*Leaderboard, *Error) {
	result := &leaderboardResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Leaderboard) Game() *Game {
	// we only have the game ID at hand
	asserted, okay := self.GameData.(string)
	if okay {
		game, _ := GameById(asserted)
		return game
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.GameData)

	// ... and try to turn it back into something meaningful
	dest := Game{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest
	}

	return nil
}

func (self *Leaderboard) Category() *Category {
	// we only have the category ID at hand
	asserted, okay := self.CategoryData.(string)
	if okay {
		category, _ := CategoryById(asserted)
		return category
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.CategoryData)

	// ... and try to turn it back into something meaningful
	dest := Category{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest
	}

	return nil
}

func (self *Leaderboard) Level() *Level {
	if self.LevelData == nil {
		return nil
	}

	// we only have the level ID at hand
	asserted, okay := self.LevelData.(string)
	if okay {
		level, _ := LevelById(asserted)
		return level
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.LevelData)

	// ... and try to turn it back into something meaningful
	dest := Level{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest
	}

	return nil
}

func (self *Leaderboard) Platforms() []*Platform {
	if self.PlatformsData == nil {
		// TODO: Walk through all runs, collect platform IDs and fetch them
		return make([]*Platform, 0)
	}

	tmp := PlatformCollection{}

	if recast(self.PlatformsData, &tmp) == nil {
		return tmp.platforms()
	}

	return make([]*Platform, 0)
}

func (self *Leaderboard) Regions() []*Region {
	if self.RegionsData == nil {
		// TODO: Walk through all runs, collect region IDs and fetch them
		return make([]*Region, 0)
	}

	tmp := RegionCollection{}

	if recast(self.RegionsData, &tmp) == nil {
		return tmp.regions()
	}

	return make([]*Region, 0)
}

func (self *Leaderboard) Players() []*Player {
	result := make([]*Player, 0)

	// players have not been embedded
	if self.PlayersData == nil {
		return result
	}

	tmp := playerCollection{}

	if recast(self.PlayersData, &tmp) == nil {
		// each element in tmp.Data has a rel that tells us whether we have a
		// user or a guest
		for _, playerProps := range tmp.Data {
			rel, exists := playerProps["rel"]
			if exists {
				player := Player{}

				switch rel {
				case "user":
					user := User{}

					if recast(playerProps, &user) == nil {
						player.User = &user
					}

				case "guest":
					guest := Guest{}

					if recast(playerProps, &guest) == nil {
						player.Guest = &guest
					}
				}

				if player.User != nil || player.Guest != nil {
					result = append(result, &player)
				}
			}
		}
	}

	return result
}

// for the 'hasLinks' interface
func (self *Leaderboard) links() []Link {
	return self.Links
}

type LeaderboardCollection struct {
	Data       []Leaderboard
	Pagination Pagination
}

func (self *LeaderboardCollection) runs() []*Leaderboard {
	result := make([]*Leaderboard, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type LeaderboardFilter struct {
	User     string
	Guest    string
	Examiner string
	Game     string
	Level    string
	Category string
	Platform string
	Region   string
	Emulated *bool
	Status   string
}

func (self *LeaderboardFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.User) > 0 {
		values.Set("user", self.User)
	}

	if len(self.Guest) > 0 {
		values.Set("guest", self.Guest)
	}

	if len(self.Examiner) > 0 {
		values.Set("examiner", self.Examiner)
	}

	if len(self.Game) > 0 {
		values.Set("game", self.Game)
	}

	if len(self.Level) > 0 {
		values.Set("level", self.Level)
	}

	if len(self.Category) > 0 {
		values.Set("category", self.Category)
	}

	if len(self.Platform) > 0 {
		values.Set("platform", self.Platform)
	}

	if len(self.Region) > 0 {
		values.Set("region", self.Region)
	}

	if len(self.Status) > 0 {
		values.Set("status", self.Status)
	}

	if self.Emulated != nil {
		if *self.Emulated {
			values.Set("emulated", "yes")
		} else {
			values.Set("emulated", "no")
		}
	}

	u.RawQuery = values.Encode()
}

func Leaderboards(f *LeaderboardFilter, s *Sorting, c *Cursor) (*LeaderboardCollection, *Error) {
	return fetchLeaderboards(request{"GET", "/runs", f, s, c})
}

func (self *LeaderboardCollection) NextPage() (*LeaderboardCollection, *Error) {
	return self.fetchLink("next")
}

func (self *LeaderboardCollection) PrevPage() (*LeaderboardCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *LeaderboardCollection) fetchLink(name string) (*LeaderboardCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchLeaderboards(next.request())
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchLeaderboards(request request) (*LeaderboardCollection, *Error) {
	result := &LeaderboardCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
