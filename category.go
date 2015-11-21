package srapi

import (
	"encoding/json"
	"net/url"
)

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

type categoryResponse struct {
	Data Category
}

func CategoryById(id string) (*Category, *Error) {
	return fetchCategory(request{"GET", "/categories/" + id, nil, nil, nil})
}

func fetchCategory(request request) (*Category, *Error) {
	result := &categoryResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Category) Game() *Game {
	if self.GameData == nil {
		link := firstLink(self, "game")
		if link == nil {
			return nil
		}

		game, _ := fetchGame(link.request())
		return game
	}

	// convert generic mess into JSON
	encoded, _ := json.Marshal(self.GameData)

	// ... and try to turn it back into something meaningful
	dest := gameResponse{}
	err := json.Unmarshal(encoded, &dest)
	if err == nil {
		return &dest.Data
	}

	return nil
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

	return fetchCategories(next.request())
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
