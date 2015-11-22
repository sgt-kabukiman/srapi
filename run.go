package srapi

import "net/url"

type Run struct {
	Id      string
	Weblink string
	Videos  struct {
		Text  string
		Links []Link
	}
	Comment string
	Status  struct {
		Status     string
		Examiner   string
		VerifyDate string `json:"verify-date"`
		Reason     string
	}
	Date  string
	Times struct {
		Primary              string
		Realtime             string
		RealtimeWithoutLoads string `json:"realtime_noloads"`
		IngameTime           string `json:"ingame"`
	}
	System struct {
		Platform string
		Emulated bool
		Region   string
	}
	Splits *Link
	Values map[string]string
	Links  []Link

	// do not use this field directly, use the available methods
	PlatformData interface{} `json:"platform"`

	// do not use this field directly, use the available methods
	RegionData interface{} `json:"region"`

	// do not use this field directly, use the available methods
	PlayersData interface{} `json:"players"`

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	CategoryData interface{} `json:"category"`

	// do not use this field directly, use the available methods
	LevelData interface{} `json:"level"`
}

func toRunCollection(data interface{}) *RunCollection {
	tmp := &RunCollection{}
	recast(data, tmp)

	return tmp
}

type runResponse struct {
	Data Run
}

func RunById(id string) (*Run, *Error) {
	return fetchRun(request{"GET", "/runs/" + id, nil, nil, nil})
}

func (self *Run) Game() *Game {
	// we only have the game ID at hand
	asserted, okay := self.GameData.(string)
	if okay {
		game, _ := GameById(asserted)
		return game
	}

	return toGame(self.GameData)
}

func (self *Run) Category() *Category {
	if self.CategoryData == nil {
		return nil
	}

	// we only have the category ID at hand
	asserted, okay := self.CategoryData.(string)
	if okay {
		category, _ := CategoryById(asserted)
		return category
	}

	return toCategory(self.CategoryData)
}

func (self *Run) Level() *Level {
	if self.LevelData == nil {
		return nil
	}

	// we only have the level ID at hand
	asserted, okay := self.LevelData.(string)
	if okay {
		level, _ := LevelById(asserted)
		return level
	}

	return toLevel(self.LevelData)
}

func (self *Run) Platform() *Platform {
	if self.PlatformData == nil {
		if len(self.System.Platform) > 0 {
			platform, _ := PlatformById(self.System.Platform)
			return platform
		}

		return nil
	}

	return toPlatform(self.PlatformData)
}

func (self *Run) Region() *Region {
	if self.RegionData == nil {
		if len(self.System.Region) > 0 {
			region, _ := RegionById(self.System.Region)
			return region
		}

		return nil
	}

	return toRegion(self.RegionData)
}

func (self *Run) Players() []*Player {
	result := make([]*Player, 0)

	switch asserted := self.PlayersData.(type) {
	// list of simple links to users/guests, e.g. players=[{rel:..,id:...}, {...}]
	case []interface{}:
		tmp := make([]PlayerLink, 0)

		if recast(asserted, &tmp) == nil {
			for _, link := range tmp {
				player := Player{}

				switch link.Relation {
				case "user":
					user, err := fetchUser(link.request(nil, nil))
					if err == nil {
						player.User = user
					}

				case "guest":
					guest, err := fetchGuest(link.request(nil, nil))
					if err == nil {
						player.Guest = guest
					}
				}

				if player.User != nil || player.Guest != nil {
					result = append(result, &player)
				}
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		tmp := playerCollection{}

		if recast(asserted, &tmp) == nil {
			// each element in tmp.Data has a rel that tells us whether we have a
			// user or a guest
			for _, playerProps := range tmp.Data {
				rel, exists := playerProps["rel"]
				if exists {
					player := Player{}

					switch rel {
					case "user":
						if user := toUser(playerProps); user != nil {
							player.User = user
						}

					case "guest":
						if guest := toGuest(playerProps); guest != nil {
							player.Guest = guest
						}
					}

					if player.User != nil || player.Guest != nil {
						result = append(result, &player)
					}
				}
			}
		}
	}

	return result
}

func (self *Run) Examiner() *User {
	return fetchUserLink(firstLink(self, "examiner"))
}

// for the 'hasLinks' interface
func (self *Run) links() []Link {
	return self.Links
}

type RunCollection struct {
	Data       []Run
	Pagination Pagination
}

func (self *RunCollection) runs() []*Run {
	result := make([]*Run, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type RunFilter struct {
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

func (self *RunFilter) applyToURL(u *url.URL) {
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

func Runs(f *RunFilter, s *Sorting, c *Cursor) (*RunCollection, *Error) {
	return fetchRuns(request{"GET", "/runs", f, s, c})
}

func (self *RunCollection) NextPage() (*RunCollection, *Error) {
	return self.fetchLink("next")
}

func (self *RunCollection) PrevPage() (*RunCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *RunCollection) fetchLink(name string) (*RunCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchRuns(next.request(nil, nil))
}

func fetchRun(request request) (*Run, *Error) {
	result := &runResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchRuns(request request) (*RunCollection, *Error) {
	result := &RunCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func fetchRunsLink(link *Link, filter filter, sort *Sorting) *RunCollection {
	if link == nil {
		return &RunCollection{}
	}

	collection, _ := fetchRuns(link.request(filter, sort))
	return collection
}
