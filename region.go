// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Region represents a geographic region.
type Region struct {
	// the unique ID
	ID string

	// the name of the region
	Name string

	// API links to related resources
	Links []Link
}

// toRegion transforms a data blob to a Region struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toRegion(data interface{}) *Region {
	dest := Region{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// toRegionCollection transforms a data blob to a RegionCollection.
// If data is nil or casting was unsuccessful, an empty RegionCollection
// is returned.
func toRegionCollection(data interface{}) *RegionCollection {
	tmp := &RegionCollection{}
	recast(data, tmp)

	return tmp
}

// regionResponse models the actual API response from the server
type regionResponse struct {
	// the one region contained in the response
	Data Region
}

// RegionByID tries to fetch a single region, identified by its ID.
// When an error is returned, the returned region is nil.
func RegionByID(id string) (*Region, *Error) {
	return fetchRegion(request{"GET", "/regions/" + id, nil, nil, nil})
}

// Runs fetches a list of runs done in the region, optionally filtered and
// sorted. This function always returns a RunCollection.
func (r *Region) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(r, "runs"), filter, sort)
}

// Games fetches a list of games available in the region, optionally filtered
// and sorted. This function always returns a GameCollection.
func (r *Region) Games(filter *GameFilter, sort *Sorting) *GameCollection {
	return fetchGamesLink(firstLink(r, "games"), filter, sort)
}

// for the 'hasLinks' interface
func (r *Region) links() []Link {
	return r.Links
}

// RegionCollection is one page of a region list. It consists of the regions
// as well as some pagination information (like links to the next or previous page).
type RegionCollection struct {
	Data       []Region
	Pagination Pagination
}

// Regions retrieves a collection of regions
func Regions(s *Sorting, c *Cursor) (*RegionCollection, *Error) {
	return fetchRegions(request{"GET", "/regions", nil, s, c})
}

// regions returns a list of pointers to the regions; used for cases where
// there is no pagination and the caller wants to return a flat slice of
// regions instead of a collection (which would be misleading, as collections
// imply pagination).
func (rc *RegionCollection) regions() []*Region {
	var result []*Region

	for idx := range rc.Data {
		result = append(result, &rc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// regions. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (rc *RegionCollection) NextPage() (*RegionCollection, *Error) {
	return rc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// regions. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (rc *RegionCollection) PrevPage() (*RegionCollection, *Error) {
	return rc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (rc *RegionCollection) fetchLink(name string) (*RegionCollection, *Error) {
	next := firstLink(&rc.Pagination, name)
	if next == nil {
		return &RegionCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchRegions(next.request(nil, nil))
}

// fetchRegion fetches a single region from the network. If the request failed,
// the returned region is nil. Otherwise, the error is nil.
func fetchRegion(request request) (*Region, *Error) {
	result := &regionResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchRegions fetches a list of regions from the network. It always
// returns a collection, even when an error is returned.
func fetchRegions(request request) (*RegionCollection, *Error) {
	result := &RegionCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
