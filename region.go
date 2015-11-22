// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

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

func toRegionCollection(data interface{}) *RegionCollection {
	tmp := &RegionCollection{}
	recast(data, tmp)

	return tmp
}

// TODO: Maybe wrap this "data" element away in the HTTP client when it knows
// that we fetch one single object.
type regionResponse struct {
	Data Region
}

func RegionById(id string) (*Region, *Error) {
	return fetchRegion(request{"GET", "/regions/" + id, nil, nil, nil})
}

func (self *Region) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "runs"), filter, sort)
}

func (self *Region) Games(filter *GameFilter, sort *Sorting) *GameCollection {
	return fetchGamesLink(firstLink(self, "games"), filter, sort)
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

func fetchRegion(request request) (*Region, *Error) {
	result := &regionResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchRegionLink(link *Link) *Region {
	if link == nil {
		return nil
	}

	region, _ := fetchRegion(link.request(nil, nil))
	return region
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

func fetchRegionsLink(link *Link, filter filter, sort *Sorting) *RegionCollection {
	if link == nil {
		return &RegionCollection{}
	}

	collection, _ := fetchRegions(link.request(filter, sort))
	return collection
}
