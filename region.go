package srapi

type Region struct {
	Id    string
	Name  string
	Links []Link
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

	return fetchRegions(next.request())
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
