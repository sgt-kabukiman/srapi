// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Variable represents a variable.
type Variable struct {
	// `category` is not mapped on purpose, so we can have a Category()
	// method. Little would be gained from polluting the variable struct
	// with a field for just the category ID.

	ID    string
	Name  string
	Scope struct {
		Type  string
		Level string
	}
	Mandatory   bool
	UserDefined bool `json:"user-defined"`
	Obsoletes   bool
	Values      struct {
		Choices map[string]string
		Default string
	}
	Links []Link
}

// toVariableCollection transforms a data blob to a VariableCollection.
// If data is nil or casting was unsuccessful, an empty VariableCollection
// is returned.
func toVariableCollection(data interface{}) *VariableCollection {
	tmp := &VariableCollection{}
	recast(data, tmp)

	return tmp
}

// variableResponse models the actual API response from the server
type variableResponse struct {
	// the one variable contained in the response
	Data Variable
}

// VariableByID tries to fetch a single variable, identified by its ID.
// When an error is returned, the returned game is nil.
func VariableByID(id string) (*Variable, *Error) {
	return fetchVariable(request{"GET", "/variables/" + id, nil, nil, nil})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (v *Variable) Game() *Game {
	return fetchGameLink(firstLink(v, "game"))
}

// Category extracts the embedded category, if possible, otherwise it will fetch
// the category by doing one additional request. This can return nil.
func (v *Variable) Category() *Category {
	return fetchCategoryLink(firstLink(v, "category"))
}

// for the 'hasLinks' interface
func (v *Variable) links() []Link {
	return v.Links
}

// VariableCollection is a list of variables. It consists of the variables as
// well as some pagination information (like links to the next or previous page).
type VariableCollection struct {
	Data       []Variable
	Pagination Pagination
}

// variables returns a list of pointers to the variables; used for cases where
// there is no pagination and the caller wants to return a flat slice of variables
// instead of a collection (which would be misleading, as collections imply
// pagination).
func (vc *VariableCollection) variables() []*Variable {
	var result []*Variable

	for idx := range vc.Data {
		result = append(result, &vc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// variables. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (vc *VariableCollection) NextPage() (*VariableCollection, *Error) {
	return vc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// variables. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (vc *VariableCollection) PrevPage() (*VariableCollection, *Error) {
	return vc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (vc *VariableCollection) fetchLink(name string) (*VariableCollection, *Error) {
	next := firstLink(&vc.Pagination, name)
	if next == nil {
		return &VariableCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchVariables(next.request(nil, nil))
}

// fetchVariable fetches a single variable from the network. If the request
// failed, the returned variable is nil. Otherwise, the error is nil.
func fetchVariable(request request) (*Variable, *Error) {
	result := &variableResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchVariables fetches a list of variables from the network. It always
// returns a collection, even when an error is returned.
func fetchVariables(request request) (*VariableCollection, *Error) {
	result := &VariableCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchVariablesLink tries to fetch a given link and interpret the response as
// a list of variables. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchVariablesLink(link requestable, filter filter, sort *Sorting) *VariableCollection {
	if link == nil {
		return &VariableCollection{}
	}

	collection, _ := fetchVariables(link.request(filter, sort))
	return collection
}
