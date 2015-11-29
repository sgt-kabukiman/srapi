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
func toRegion(data interface{}, isResponse bool) *Region {
	if data == nil {
		return nil
	}

	if isResponse {
		dest := regionResponse{}

		if recast(data, &dest) == nil {
			return &dest.Data
		}
	} else {
		dest := Region{}

		if recast(data, &dest) == nil {
			return &dest
		}
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
	return fetchRegion(request{"GET", "/regions/" + id, nil, nil, nil, ""})
}

// Runs fetches a list of runs done in the region, optionally filtered and
// sorted. This function always returns a RunCollection.
func (r *Region) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(r, "runs"), filter, sort, embeds)
}

// Games fetches a list of games available in the region, optionally filtered
// and sorted. This function always returns a GameCollection.
func (r *Region) Games(filter *GameFilter, sort *Sorting, embeds string) (*GameCollection, *Error) {
	return fetchGamesLink(firstLink(r, "games"), filter, sort, embeds)
}

// for the 'hasLinks' interface
func (r *Region) links() []Link {
	return r.Links
}

// Regions retrieves a collection of regions
func Regions(s *Sorting, c *Cursor) (*RegionCollection, *Error) {
	return fetchRegions(request{"GET", "/regions", nil, s, c, ""})
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

	return result, err
}
