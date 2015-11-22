package srapi

import "encoding/json"

type Level struct {
	Id      string
	Name    string
	Weblink string
	Rules   string
	Links   []Link

	// do not use this field directly, use the available methods
	CategoriesData interface{} `json:"categories"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

type levelResponse struct {
	Data Level
}

func LevelById(id string) (*Level, *Error) {
	return fetchLevel(request{"GET", "/levels/" + id, nil, nil, nil})
}

func fetchLevel(request request) (*Level, *Error) {
	result := &levelResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Level) Game() *Game {
	link := firstLink(self, "game")
	if link == nil {
		return nil
	}

	game, _ := fetchGame(link.request(nil, nil))
	return game
}

func (self *Level) Categories(filter *CategoryFilter, sort *Sorting) []*Category {
	if self.CategoriesData == nil {
		link := firstLink(self, "categories")
		if link == nil {
			return nil
		}

		collection, _ := fetchCategories(link.request(filter, sort))

		return collection.categories()
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.CategoriesData)

	// ... and try to turn it back into something meaningful
	dest := CategoryCollection{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return dest.categories()
	}

	return make([]*Category, 0)
}

func (self *Level) Variables(sort *Sorting) []*Variable {
	if self.VariablesData == nil {
		link := firstLink(self, "variables")
		if link == nil {
			return nil
		}

		collection, _ := fetchVariables(link.request(nil, sort))

		return collection.variables()
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.VariablesData)

	// ... and try to turn it back into something meaningful
	dest := VariableCollection{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return dest.variables()
	}

	return make([]*Variable, 0)
}

func (self *Level) PrimaryLeaderboard(options *LeaderboardOptions) *Leaderboard {
	link := firstLink(self, "leaderboard")
	if link == nil {
		return nil
	}

	leaderboard, _ := fetchLeaderboard(link.request(options, nil))
	return leaderboard
}

func (self *Level) Records(filter *LeaderboardFilter) *LeaderboardCollection {
	link := firstLink(self, "records")
	if link == nil {
		return nil
	}

	leaderboards, _ := fetchLeaderboards(link.request(filter, nil))
	return leaderboards
}

func (self *Level) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	link := firstLink(self, "runs")
	if link == nil {
		return nil
	}

	runs, _ := fetchRuns(link.request(filter, sort))
	return runs
}

// for the 'hasLinks' interface
func (self *Level) links() []Link {
	return self.Links
}

type LevelCollection struct {
	Data       []Level
	Pagination Pagination
}

func (self *LevelCollection) levels() []*Level {
	result := make([]*Level, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

func (self *LevelCollection) NextPage() (*LevelCollection, *Error) {
	return self.fetchLink("next")
}

func (self *LevelCollection) PrevPage() (*LevelCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *LevelCollection) fetchLink(name string) (*LevelCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchLevels(next.request(nil, nil))
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchLevels(request request) (*LevelCollection, *Error) {
	result := &LevelCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
