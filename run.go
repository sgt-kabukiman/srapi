package srapi

import (
	"encoding/json"
	"net/url"
)

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
	}
	Players []struct {
		Relation string `json:"rel"`
		Id       string
		URI      string
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

type runResponse struct {
	Data Run
}

func RunById(id string) (*Run, *Error) {
	return fetchRun(request{"GET", "/runs/" + id, nil, nil, nil})
}

func fetchRun(request request) (*Run, *Error) {
	result := &runResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Run) Game() *Game {
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

func (self *Run) Platform() *Platform {
	if self.PlatformData == nil {
		if len(self.System.Platform) > 0 {
			platform, _ := PlatformById(self.System.Platform)
			return platform
		}

		return nil
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.PlatformData)

	// ... and try to turn it back into something meaningful
	dest := Platform{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest
	}

	return nil
}

func (self *Run) Region() *Region {
	if self.RegionData == nil {
		if len(self.System.Region) > 0 {
			region, _ := RegionById(self.System.Region)
			return region
		}

		return nil
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.RegionData)

	// ... and try to turn it back into something meaningful
	dest := Region{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest
	}

	return nil
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

	return fetchRuns(next.request())
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
