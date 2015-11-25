// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Level represents a level.
type Level struct {
	// the unique ID
	ID string

	// the name of the level
	Name string

	// a link to the leaderboard for this level on speedrun.com
	Weblink string

	// rules for the level; arbitrary text
	Rules string

	// API links to related resources
	Links []Link

	// do not use this field directly, use the available methods
	CategoriesData interface{} `json:"categories"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

// toLevel transforms a data blob to a Level struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toLevel(data interface{}) *Level {
	dest := Level{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// toLevelCollection transforms a data blob to a LevelCollection.
// If data is nil or casting was unsuccessful, an empty LevelCollection
// is returned.
func toLevelCollection(data interface{}) *LevelCollection {
	tmp := &LevelCollection{}
	recast(data, tmp)

	return tmp
}

// levelResponse models the actual API response from the server
type levelResponse struct {
	// the one level contained in the response
	Data Level
}

// LevelByID tries to fetch a single level, identified by its ID.
// When an error is returned, the returned level is nil.
func LevelByID(id string) (*Level, *Error) {
	return fetchLevel(request{"GET", "/levels/" + id, nil, nil, nil})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (l *Level) Game() (*Game, *Error) {
	return fetchGameLink(firstLink(l, "game"))
}

// Categories extracts the embedded categories, if possible, otherwise it will
// fetch them by doing one additional request. filter and sort are only relevant
// when the categories are not already embedded.
func (l *Level) Categories(filter *CategoryFilter, sort *Sorting) ([]*Category, *Error) {
	var collection *CategoryCollection
	var err *Error

	if l.CategoriesData == nil {
		collection, err = fetchCategoriesLink(firstLink(l, "categories"), filter, sort)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toCategoryCollection(l.CategoriesData)
	}

	return collection.categories(), nil
}

// Variables extracts the embedded variables, if possible, otherwise it will
// fetch them by doing one additional request. sort is only relevant when the
// variables are not already embedded.
func (l *Level) Variables(sort *Sorting) ([]*Variable, *Error) {
	var collection *VariableCollection
	var err *Error

	if l.VariablesData == nil {
		collection, err = fetchVariablesLink(firstLink(l, "variables"), nil, sort)
		if err != nil {
			return nil, err
		}
	} else {
		collection = toVariableCollection(l.VariablesData)
	}

	return collection.variables(), nil
}

// PrimaryLeaderboard fetches the primary leaderboard, if any, for the level.
// The result can be nil.
func (l *Level) PrimaryLeaderboard(options *LeaderboardOptions) (*Leaderboard, *Error) {
	return fetchLeaderboardLink(firstLink(l, "leaderboard"), options)
}

// Records fetches a list of leaderboards for the level, assuming the default
// category. This function always returns a LeaderboardCollection.
func (l *Level) Records(filter *LeaderboardFilter) (*LeaderboardCollection, *Error) {
	return fetchLeaderboardsLink(firstLink(l, "records"), filter, nil)
}

// Runs fetches a list of runs done in the given level and its default category,
// optionally filtered and sorted. This function always returns a RunCollection.
func (l *Level) Runs(filter *RunFilter, sort *Sorting) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(l, "runs"), filter, sort)
}

// for the 'hasLinks' interface
func (l *Level) links() []Link {
	return l.Links
}

// LevelCollection is one page of a level list. It consists of the levels as
// well as some pagination information (like links to the next or previous page).
type LevelCollection struct {
	Data       []Level
	Pagination Pagination
}

// levels returns a list of pointers to the levels; used for cases where
// there is no pagination and the caller wants to return a flat slice of levels
// instead of a collection (which would be misleading, as collections imply
// pagination).
func (lc *LevelCollection) levels() []*Level {
	var result []*Level

	for idx := range lc.Data {
		result = append(result, &lc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// levels. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (lc *LevelCollection) NextPage() (*LevelCollection, *Error) {
	return lc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// levels. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (lc *LevelCollection) PrevPage() (*LevelCollection, *Error) {
	return lc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (lc *LevelCollection) fetchLink(name string) (*LevelCollection, *Error) {
	next := firstLink(&lc.Pagination, name)
	if next == nil {
		return &LevelCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchLevelsLink(next, nil, nil)
}

// fetchLevel fetches a single level from the network. If the request failed,
// the returned level is nil. Otherwise, the error is nil.
func fetchLevel(request request) (*Level, *Error) {
	result := &levelResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchLevels fetches a list of levels from the network. It always
// returns a collection, even when an error is returned.
func fetchLevels(request request) (*LevelCollection, *Error) {
	result := &LevelCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchLevelsLink tries to fetch a given link and interpret the response as
// a list of levels. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchLevelsLink(link requestable, filter filter, sort *Sorting) (*LevelCollection, *Error) {
	if !link.exists() {
		return &LevelCollection{}, nil
	}

	return fetchLevels(link.request(filter, sort))
}
