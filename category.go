package srapi

import "net/url"

type Category struct {
	Id      string
	Name    string
	Weblink string
	Type    string
	Rules   string
	Players struct {
		Type  string
		Value int
	}
	Miscellaneous bool
	Links         []Link

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

func toCategory(data interface{}) *Category {
	dest := Category{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

func toCategoryCollection(data interface{}) *CategoryCollection {
	tmp := &CategoryCollection{}
	recast(data, tmp)

	return tmp
}

type categoryResponse struct {
	Data Category
}

func CategoryById(id string) (*Category, *Error) {
	return fetchCategory(request{"GET", "/categories/" + id, nil, nil, nil})
}

func (self *Category) Game() *Game {
	if self.GameData == nil {
		return fetchGameLink(firstLink(self, "game"))
	}

	return toGame(self.GameData)
}

func (self *Category) Variables(sort *Sorting) []*Variable {
	var collection *VariableCollection

	if self.VariablesData == nil {
		collection = fetchVariablesLink(firstLink(self, "variables"), nil, sort)
	} else {
		collection = toVariableCollection(self.VariablesData)
	}

	return collection.variables()
}

func (self *Category) PrimaryLeaderboard(options *LeaderboardOptions) *Leaderboard {
	return fetchLeaderboardLink(firstLink(self, "leaderboard"), options)
}

func (self *Category) Records(filter *LeaderboardFilter) *LeaderboardCollection {
	return fetchLeaderboardsLink(firstLink(self, "records"), filter, nil)
}

func (self *Category) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "records"), filter, sort)
}

// for the 'hasLinks' interface
func (self *Category) links() []Link {
	return self.Links
}

type CategoryCollection struct {
	Data       []Category
	Pagination Pagination
}

func (self *CategoryCollection) categories() []*Category {
	result := make([]*Category, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type CategoryFilter struct {
	Miscellaneous *bool
}

func (self *CategoryFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if self.Miscellaneous != nil {
		if *self.Miscellaneous {
			values.Set("miscellaneous", "yes")
		} else {
			values.Set("miscellaneous", "no")
		}
	}

	u.RawQuery = values.Encode()
}

func (self *CategoryCollection) NextPage() (*CategoryCollection, *Error) {
	return self.fetchLink("next")
}

func (self *CategoryCollection) PrevPage() (*CategoryCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *CategoryCollection) fetchLink(name string) (*CategoryCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchCategories(next.request(nil, nil))
}

func fetchCategory(request request) (*Category, *Error) {
	result := &categoryResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchCategoryLink(link *Link) *Category {
	if link == nil {
		return nil
	}

	category, _ := fetchCategory(link.request(nil, nil))
	return category
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchCategories(request request) (*CategoryCollection, *Error) {
	result := &CategoryCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func fetchCategoriesLink(link *Link, filter filter, sort *Sorting) *CategoryCollection {
	if link == nil {
		return &CategoryCollection{}
	}

	collection, _ := fetchCategories(link.request(filter, sort))
	return collection
}
