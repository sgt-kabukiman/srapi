// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

// Category is a structure representing a game category, either per-game or per-level.
type Category struct {
	// unique category ID
	ID string

	// category name, for example "Any%"
	Name string

	// link to this category on speedrun.com
	Weblink string

	// either "per-game" or "per-level"
	Type string

	// ruleset for the category, arbitrary text
	Rules string

	// definition on how many players are needed/allowed for runs in this category
	Players struct {
		Type  string
		Value int
	}

	// whether or not this is a misc (fun) category
	Miscellaneous bool

	// API links to related resources
	Links []Link

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

// toCategory transforms a data blob to a Category struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toCategory(data interface{}, isResponse bool) *Category {
	if data == nil {
		return nil
	}

	if isResponse {
		dest := categoryResponse{}

		if recast(data, &dest) == nil {
			return &dest.Data
		}
	} else {
		dest := Category{}

		if recast(data, &dest) == nil {
			return &dest
		}
	}

	return nil
}

// toCategoryCollection transforms a data blob to a CategoryCollection.
// If data is nil or casting was unsuccessful, an empty CategoryCollection
// is returned.
func toCategoryCollection(data interface{}) *CategoryCollection {
	tmp := &CategoryCollection{}
	recast(data, tmp)

	return tmp
}

// categoryResponse models the actual API response from the server
type categoryResponse struct {
	// the one category contained in the response
	Data Category
}

// CategoryByID tries to fetch a single category, identified by its ID.
// When an error is returned, the returned category is nil.
func CategoryByID(id string, embeds string) (*Category, *Error) {
	return fetchCategory(request{"GET", "/categories/" + id, nil, nil, nil, embeds})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (c *Category) Game(embeds string) (*Game, *Error) {
	if c.GameData == nil {
		return fetchGameLink(firstLink(c, "game"), embeds)
	}

	return toGame(c.GameData, true), nil
}

// Variables extracts the embedded variables, if possible, otherwise it will
// fetch them by doing one additional request. sort is only relevant when the
// variables are not already embedded.
func (c *Category) Variables(sort *Sorting) (*VariableCollection, *Error) {
	var collection *VariableCollection
	var err *Error

	if c.VariablesData == nil {
		collection, err = fetchVariablesLink(firstLink(c, "variables"), nil, sort)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toVariableCollection(c.VariablesData)
	}

	return collection, nil
}

// PrimaryLeaderboard fetches the primary leaderboard, if any, for the category.
// The result can be nil.
func (c *Category) PrimaryLeaderboard(options *LeaderboardOptions, embeds string) (*Leaderboard, *Error) {
	return fetchLeaderboardLink(firstLink(c, "leaderboard"), options, embeds)
}

// Records fetches a list of leaderboards for the category. For full-game
// categories, the list will contain one leaderboard, otherwise it will have one
// per level. This function always returns a LeaderboardCollection.
func (c *Category) Records(filter *LeaderboardFilter, embeds string) (*LeaderboardCollection, *Error) {
	return fetchLeaderboardsLink(firstLink(c, "records"), filter, nil, embeds)
}

// Runs fetches a list of runs done in the given category, optionally filtered
// and sorted. This function always returns a RunCollection.
func (c *Category) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(c, "records"), filter, sort, embeds)
}

// for the 'hasLinks' interface
func (c *Category) links() []Link {
	return c.Links
}

// CategoryFilter represents the possible filtering options when fetching a list
// of categories.
type CategoryFilter struct {
	Miscellaneous OptionalFlag
}

// applyToURL merged the filter into a URL.
func (cf *CategoryFilter) applyToURL(u *url.URL) {
	if cf == nil {
		return
	}

	values := u.Query()
	cf.Miscellaneous.applyToQuery("miscellaneous", &values)
	u.RawQuery = values.Encode()
}

// fetchCategory fetches a single category from the network. If the request failed,
// the returned category is nil. Otherwise, the error is nil.
func fetchCategory(request request) (*Category, *Error) {
	result := &categoryResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchCategoryLink tries to fetch a given link and interpret the response as
// a single category. If the link is nil or the category could not be fetched,
// nil is returned.
func fetchCategoryLink(link requestable, embeds string) (*Category, *Error) {
	if !link.exists() {
		return nil, nil
	}

	return fetchCategory(link.request(nil, nil, embeds))
}

// fetchCategories fetches a list of categories from the network. It always
// returns a collection, even when an error is returned.
func fetchCategories(request request) (*CategoryCollection, *Error) {
	result := &CategoryCollection{}
	err := httpClient.do(request, result)

	return result, err
}

// fetchCategoriesLink tries to fetch a given link and interpret the response as
// a list of categories. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchCategoriesLink(link requestable, filter filter, sort *Sorting, embeds string) (*CategoryCollection, *Error) {
	if !link.exists() {
		return &CategoryCollection{}, nil
	}

	return fetchCategories(link.request(filter, sort, embeds))
}
