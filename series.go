package srapi

import "net/url"

type Series struct {
	Id    string
	Names struct {
		International string
		Japanese      string
	}
	Abbreviation string
	Weblink      string
	Assets       map[string]*AssetLink
	Links        []Link

	// do not use this field directly, use the available methods
	ModeratorsData interface{} `json:"moderators"`
}

type seriesResponse struct {
	Data Series
}

func SeriesById(id string) (*Series, *Error) {
	return fetchOneSeries(request{"GET", "/series/" + id, nil, nil, nil})
}

func fetchOneSeries(request request) (*Series, *Error) {
	result := &seriesResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func SeriesByAbbreviation(abbrev string) (*Series, *Error) {
	return SeriesById(abbrev)
}

func (self *Series) Games() *GameCollection {
	link := firstLink(self, "games")
	if link == nil {
		return nil
	}

	collection, _ := fetchGames(link.request())

	return collection
}

func (self *Series) ModeratorMap() map[string]GameModLevel {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := self.ModeratorsData.(map[string]GameModLevel)
	if okay {
		return assertedMap
	}

	// maybe we got a list of embedded users
	result := make(map[string]GameModLevel, 0)
	tmp := UserCollection{}

	if recast(self.ModeratorsData, &tmp) == nil {
		for _, user := range tmp.users() {
			result[user.Id] = UnknownModLevel
		}
	}

	return result
}

func (self *Series) Moderators() []*User {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := self.ModeratorsData.(map[string]GameModLevel)
	if okay {
		result := make([]*User, 0)

		for userId := range assertedMap {
			user, err := UserById(userId)
			if err == nil {
				result = append(result, user)
			}
		}

		return result
	}

	// maybe we got a list of embedded users
	tmp := UserCollection{}

	if recast(self.ModeratorsData, &tmp) == nil {
		return tmp.users()
	}

	return make([]*User, 0)
}

// for the 'hasLinks' interface
func (self *Series) links() []Link {
	return self.Links
}

type SeriesCollection struct {
	Data       []Series
	Pagination Pagination
}

func (self *SeriesCollection) series() []*Series {
	result := make([]*Series, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type SeriesFilter struct {
	Name         string
	Abbreviation string
	Moderator    string
}

func (self *SeriesFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.Name) > 0 {
		values.Set("name", self.Name)
	}

	if len(self.Abbreviation) > 0 {
		values.Set("abbreviation", self.Abbreviation)
	}

	if len(self.Moderator) > 0 {
		values.Set("moderator", self.Moderator)
	}

	u.RawQuery = values.Encode()
}

func ManySeries(f *SeriesFilter, s *Sorting, c *Cursor) (*SeriesCollection, *Error) {
	return fetchManySeries(request{"GET", "/series", f, s, c})
}

func (self *SeriesCollection) NextPage() (*SeriesCollection, *Error) {
	return self.fetchLink("next")
}

func (self *SeriesCollection) PrevPage() (*SeriesCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *SeriesCollection) fetchLink(name string) (*SeriesCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchManySeries(next.request())
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchManySeries(request request) (*SeriesCollection, *Error) {
	result := &SeriesCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
