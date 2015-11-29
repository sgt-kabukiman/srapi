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
	return fetchVariable(request{"GET", "/variables/" + id, nil, nil, nil, ""})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (v *Variable) Game(embeds string) (*Game, *Error) {
	return fetchGameLink(firstLink(v, "game"), embeds)
}

// Category extracts the embedded category, if possible, otherwise it will fetch
// the category by doing one additional request. This can return nil.
func (v *Variable) Category(embeds string) (*Category, *Error) {
	return fetchCategoryLink(firstLink(v, "category"), embeds)
}

// for the 'hasLinks' interface
func (v *Variable) links() []Link {
	return v.Links
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

	return result, err
}

// fetchVariablesLink tries to fetch a given link and interpret the response as
// a list of variables. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchVariablesLink(link requestable, filter filter, sort *Sorting) (*VariableCollection, *Error) {
	if !link.exists() {
		return &VariableCollection{}, nil
	}

	return fetchVariables(link.request(filter, sort, ""))
}
