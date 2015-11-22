// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

type Platform struct {
	Id       string
	Name     string
	Released int
	Links    []Link
}

func toPlatform(data interface{}) *Platform {
	dest := Platform{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

func toPlatformCollection(data interface{}) *PlatformCollection {
	tmp := &PlatformCollection{}
	recast(data, tmp)

	return tmp
}

type platformResponse struct {
	Data Platform
}

func PlatformById(id string) (*Platform, *Error) {
	return fetchPlatform(request{"GET", "/platforms/" + id, nil, nil, nil})
}

func (self *Platform) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "runs"), filter, sort)
}

func (self *Platform) Games(filter *GameFilter, sort *Sorting) *GameCollection {
	return fetchGamesLink(firstLink(self, "games"), filter, sort)
}

// for the 'hasLinks' interface
func (self *Platform) links() []Link {
	return self.Links
}

type PlatformCollection struct {
	Data       []Platform
	Pagination Pagination
}

func Platforms(s *Sorting, c *Cursor) (*PlatformCollection, *Error) {
	return fetchPlatforms(request{"GET", "/platforms", nil, s, c})
}

func (self *PlatformCollection) platforms() []*Platform {
	result := make([]*Platform, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

func (self *PlatformCollection) NextPage() (*PlatformCollection, *Error) {
	return self.fetchLink("next")
}

func (self *PlatformCollection) PrevPage() (*PlatformCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *PlatformCollection) fetchLink(name string) (*PlatformCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchPlatforms(next.request(nil, nil))
}

func fetchPlatform(request request) (*Platform, *Error) {
	result := &platformResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchPlatformLink(link *Link) *Platform {
	if link == nil {
		return nil
	}

	platform, _ := fetchPlatform(link.request(nil, nil))
	return platform
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchPlatforms(request request) (*PlatformCollection, *Error) {
	result := &PlatformCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func fetchPlatformsLink(link *Link, filter filter, sort *Sorting) *PlatformCollection {
	if link == nil {
		return &PlatformCollection{}
	}

	collection, _ := fetchPlatforms(link.request(filter, sort))
	return collection
}
