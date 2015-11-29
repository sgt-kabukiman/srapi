// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Platform represents a platform.
type Platform struct {
	// the unique ID
	ID string

	// the name of the platform
	Name string

	// the year the platform was released
	Released int

	// API links to related resources
	Links []Link
}

// toPlatform transforms a data blob to a Platform struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toPlatform(data interface{}, isResponse bool) *Platform {
	if data == nil {
		return nil
	}

	if isResponse {
		dest := platformResponse{}

		if recast(data, &dest) == nil {
			return &dest.Data
		}
	} else {
		dest := Platform{}

		if recast(data, &dest) == nil {
			return &dest
		}
	}

	return nil
}

// toPlatformCollection transforms a data blob to a PlatformCollection.
// If data is nil or casting was unsuccessful, an empty PlatformCollection
// is returned.
func toPlatformCollection(data interface{}) *PlatformCollection {
	tmp := &PlatformCollection{}
	recast(data, tmp)

	return tmp
}

// platformResponse models the actual API response from the server
type platformResponse struct {
	// the one platform contained in the response
	Data Platform
}

// PlatformByID tries to fetch a single platform, identified by its ID.
// When an error is returned, the returned platform is nil.
func PlatformByID(id string) (*Platform, *Error) {
	return fetchPlatform(request{"GET", "/platforms/" + id, nil, nil, nil, ""})
}

// Runs fetches a list of runs done on the platform, optionally filtered and
// sorted. This function always returns a RunCollection.
func (p *Platform) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(p, "runs"), filter, sort, embeds)
}

// Games fetches a list of games available on the platform, optionally filtered
// and sorted. This function always returns a GameCollection.
func (p *Platform) Games(filter *GameFilter, sort *Sorting, embeds string) (*GameCollection, *Error) {
	return fetchGamesLink(firstLink(p, "games"), filter, sort, embeds)
}

// for the 'hasLinks' interface
func (p *Platform) links() []Link {
	return p.Links
}

// Platforms retrieves a collection of platforms
func Platforms(s *Sorting, c *Cursor) (*PlatformCollection, *Error) {
	return fetchPlatforms(request{"GET", "/platforms", nil, s, c, ""})
}

// fetchPlatform fetches a single platform from the network. If the request failed,
// the returned platform is nil. Otherwise, the error is nil.
func fetchPlatform(request request) (*Platform, *Error) {
	result := &platformResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchPlatforms fetches a list of platforms from the network. It always
// returns a collection, even when an error is returned.
func fetchPlatforms(request request) (*PlatformCollection, *Error) {
	result := &PlatformCollection{}
	err := httpClient.do(request, result)

	return result, err
}
