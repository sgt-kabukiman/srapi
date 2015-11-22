package srapi

type Region struct {
	Id    string
	Name  string
	Links []Link
}

func toRegion(data interface{}) *Region {
	dest := Region{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// TODO: Maybe wrap this "data" element away in the HTTP client when it knows
// that we fetch one single object.
type regionResponse struct {
	Data Region
}

func RegionById(id string) (*Region, *Error) {
	request := request{"GET", "/regions/" + id, nil, nil, nil}
	result := &regionResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Region) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	link := firstLink(self, "runs")
	if link == nil {
		return nil
	}

	runs, _ := fetchRuns(link.request(filter, sort))
	return runs
}

func (self *Region) Games(filter *GameFilter, sort *Sorting) *GameCollection {
	link := firstLink(self, "games")
	if link == nil {
		return nil
	}

	games, _ := fetchGames(link.request(filter, sort))
	return games
}

// for the 'hasLinks' interface
func (self *Region) links() []Link {
	return self.Links
}

type RegionCollection struct {
	Data       []Region
	Pagination Pagination
}

func Regions(s *Sorting, c *Cursor) (*RegionCollection, *Error) {
	return fetchRegions(request{"GET", "/regions", nil, s, c})
}

func (self *RegionCollection) regions() []*Region {
	result := make([]*Region, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

func (self *RegionCollection) NextPage() (*RegionCollection, *Error) {
	return self.fetchLink("next")
}

func (self *RegionCollection) PrevPage() (*RegionCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *RegionCollection) fetchLink(name string) (*RegionCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchRegions(next.request(nil, nil))
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchRegions(request request) (*RegionCollection, *Error) {
	result := &RegionCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
