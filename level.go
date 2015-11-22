// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

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

func toLevel(data interface{}) *Level {
	dest := Level{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

func toLevelCollection(data interface{}) *LevelCollection {
	tmp := &LevelCollection{}
	recast(data, tmp)

	return tmp
}

type levelResponse struct {
	Data Level
}

func LevelById(id string) (*Level, *Error) {
	return fetchLevel(request{"GET", "/levels/" + id, nil, nil, nil})
}

func (self *Level) Game() *Game {
	return fetchGameLink(firstLink(self, "game"))
}

func (self *Level) Categories(filter *CategoryFilter, sort *Sorting) []*Category {
	var collection *CategoryCollection

	if self.CategoriesData == nil {
		collection = fetchCategoriesLink(firstLink(self, "categories"), filter, sort)
	} else {
		collection = toCategoryCollection(self.CategoriesData)
	}

	return collection.categories()
}

func (self *Level) Variables(sort *Sorting) []*Variable {
	var collection *VariableCollection

	if self.VariablesData == nil {
		collection = fetchVariablesLink(firstLink(self, "variables"), nil, sort)
	} else {
		collection = toVariableCollection(self.VariablesData)
	}

	return collection.variables()
}

func (self *Level) PrimaryLeaderboard(options *LeaderboardOptions) *Leaderboard {
	return fetchLeaderboardLink(firstLink(self, "leaderboard"), options)
}

func (self *Level) Records(filter *LeaderboardFilter) *LeaderboardCollection {
	return fetchLeaderboardsLink(firstLink(self, "records"), filter, nil)
}

func (self *Level) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "runs"), filter, sort)
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

func fetchLevel(request request) (*Level, *Error) {
	result := &levelResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchLevelLink(link *Link) *Level {
	if link == nil {
		return nil
	}

	level, _ := fetchLevel(link.request(nil, nil))
	return level
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

func fetchLevelsLink(link *Link, filter filter, sort *Sorting) *LevelCollection {
	if link == nil {
		return &LevelCollection{}
	}

	collection, _ := fetchLevels(link.request(filter, sort))
	return collection
}
