package srapi

import (
	"net/url"
	"strconv"
)

type User struct {
	Id    string
	Names struct {
		International string
		Japanese      string
	}
	Weblink   string
	NameStyle struct {
		Style     string
		Color     *NameColor
		ColorFrom *NameColor `json:"color-from"`
		ColorTo   *NameColor `json:"color-to"`
	} `json:"name-style"`
	Role     string
	Signup   string
	Location struct {
		Country Location
		Region  *Location
	}
	Twitch        *SocialLink
	Hitbox        *SocialLink
	YouTube       *SocialLink
	Twitter       *SocialLink
	SpeedRunsLive *SocialLink
	Links         []Link
}

type SocialLink struct {
	URI string
}

type Location struct {
	Code  string
	Names struct {
		International string
		Japanese      string
	}
}

type NameColor struct {
	Light string
	Dark  string
}

type userResponse struct {
	Data User
}

func UserById(id string) (*User, *Error) {
	return fetchUser(request{"GET", "/users/" + id, nil, nil, nil})
}

func fetchUser(request request) (*User, *Error) {
	result := &userResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *User) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	link := firstLink(self, "runs")
	if link == nil {
		return nil
	}

	runs, _ := fetchRuns(link.request(filter, sort))
	return runs
}

func (self *User) ModeratedGames(filter *GameFilter, sort *Sorting) *GameCollection {
	link := firstLink(self, "games")
	if link == nil {
		return nil
	}

	games, _ := fetchGames(link.request(filter, sort))
	return games
}

// for the 'hasLinks' interface
func (self *User) links() []Link {
	return self.Links
}

type PersonalBest struct {
	Rank int
	Run  Run

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

func (self *PersonalBest) Game() *Game {
	if self.GameData == nil {
		return nil
	}

	// convert generic mess into JSON
	dest := Game{}

	if recast(self.GameData, &dest) == nil {
		return &dest
	}

	return nil
}

func (self *PersonalBest) Category() *Category {
	if self.CategoryData == nil {
		return nil
	}

	// convert generic mess into JSON
	dest := Category{}

	if recast(self.CategoryData, &dest) == nil {
		return &dest
	}

	return nil
}

func (self *PersonalBest) Level() *Level {
	if self.LevelData == nil {
		return nil
	}

	// convert generic mess into JSON
	dest := Level{}

	if recast(self.LevelData, &dest) == nil {
		return &dest
	}

	return nil
}

func (self *PersonalBest) Platform() *Platform {
	if self.PlatformData == nil {
		if len(self.Run.System.Platform) > 0 {
			platform, _ := PlatformById(self.Run.System.Platform)
			return platform
		}

		return nil
	}

	// convert generic mess into JSON
	dest := Platform{}

	if recast(self.PlatformData, &dest) == nil {
		return &dest
	}

	return nil
}

func (self *PersonalBest) Region() *Region {
	if self.RegionData == nil {
		if len(self.Run.System.Region) > 0 {
			region, _ := RegionById(self.Run.System.Region)
			return region
		}

		return nil
	}

	// convert generic mess into JSON
	dest := Region{}

	if recast(self.RegionData, &dest) == nil {
		return &dest
	}

	return nil
}

func (self *PersonalBest) Players() []*Player {
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
	}

	return result
}

func (self *PersonalBest) Examiner() *User {
	link := firstLink(&self.Run, "examiner")
	if link == nil {
		return nil
	}

	examiner, _ := fetchUser(link.request(nil, nil))
	return examiner
}

type personalBestResponse struct {
	Data []PersonalBest
}

func (self *personalBestResponse) personalBests() []*PersonalBest {
	result := make([]*PersonalBest, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type PersonalBestFilter struct {
	Top    int
	Series string
	Game   string
}

func (self *PersonalBestFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if self.Top > 0 {
		values.Set("top", strconv.Itoa(self.Top))
	}

	if len(self.Series) > 0 {
		values.Set("series", self.Series)
	}

	if len(self.Game) > 0 {
		values.Set("game", self.Game)
	}

	u.RawQuery = values.Encode()
}

func (self *User) PersonalBests(filter *PersonalBestFilter) []*PersonalBest {
	link := firstLink(self, "personal-bests")
	if link == nil {
		return make([]*PersonalBest, 0)
	}

	tmp := personalBestResponse{}
	err := httpClient.do(link.request(filter, nil), &tmp)
	if err != nil {
		return make([]*PersonalBest, 0)
	}

	return tmp.personalBests()
}

type UserCollection struct {
	Data       []User
	Pagination Pagination
}

func (self *UserCollection) users() []*User {
	result := make([]*User, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type UserFilter struct {
	Lookup        string
	Name          string
	Twitch        string
	Hitbox        string
	Twitter       string
	SpeedRunsLive string
}

func (self *UserFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.Lookup) > 0 {
		values.Set("lookup", self.Lookup)
	}

	if len(self.Name) > 0 {
		values.Set("name", self.Name)
	}

	if len(self.Twitch) > 0 {
		values.Set("twitch", self.Twitch)
	}

	if len(self.Hitbox) > 0 {
		values.Set("hitbox", self.Hitbox)
	}

	if len(self.Twitter) > 0 {
		values.Set("twitter", self.Twitter)
	}

	if len(self.SpeedRunsLive) > 0 {
		values.Set("speedrunslive", self.SpeedRunsLive)
	}

	u.RawQuery = values.Encode()
}

func Users(f *UserFilter, s *Sorting, c *Cursor) (*UserCollection, *Error) {
	return fetchUsers(request{"GET", "/users", f, s, c})
}

func (self *UserCollection) NextPage() (*UserCollection, *Error) {
	return self.fetchLink("next")
}

func (self *UserCollection) PrevPage() (*UserCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *UserCollection) fetchLink(name string) (*UserCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchUsers(next.request(nil, nil))
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchUsers(request request) (*UserCollection, *Error) {
	result := &UserCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
