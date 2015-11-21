package srapi

type Variable struct {
	// `category` is not mapped on purpose, so we can have a Category()
	// method. Little would be gained from polluting the variable struct
	// with a field for just the category ID.

	Id    string
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

type variableResponse struct {
	Data Variable
}

func VariableById(id string) (*Variable, *Error) {
	return fetchVariable(request{"GET", "/variables/" + id, nil, nil, nil})
}

func fetchVariable(request request) (*Variable, *Error) {
	result := &variableResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Variable) Game() *Game {
	link := firstLink(self, "game")
	if link == nil {
		return nil
	}

	game, _ := fetchGame(link.request())
	return game
}

func (self *Variable) Category() *Category {
	link := firstLink(self, "category")
	if link == nil {
		return nil
	}

	category, _ := fetchCategory(link.request())
	return category
}

// for the 'hasLinks' interface
func (self *Variable) links() []Link {
	return self.Links
}

type VariableCollection struct {
	Data       []Variable
	Pagination Pagination
}

func (self *VariableCollection) variables() []*Variable {
	result := make([]*Variable, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

func (self *VariableCollection) NextPage() (*VariableCollection, *Error) {
	return self.fetchLink("next")
}

func (self *VariableCollection) PrevPage() (*VariableCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *VariableCollection) fetchLink(name string) (*VariableCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchVariables(next.request())
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchVariables(request request) (*VariableCollection, *Error) {
	result := &VariableCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
