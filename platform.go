package srapi

type Platform struct {
	Id       string
	Name     string
	Released int
	Links    []Link
}

type platformResponse struct {
	Data Platform
}

func PlatformById(id string) (*Platform, *Error) {
	request := request{"GET", "/platforms/" + id, nil, nil, nil}
	result := &platformResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
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

	return fetchPlatforms(next.request())
}

func fetchPlatforms(request request) (*PlatformCollection, *Error) {
	result := &PlatformCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
