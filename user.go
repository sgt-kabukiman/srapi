// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"time"
)

// User represents a user.
type User struct {
	ID    string
	Names struct {
		International string
		Japanese      string
	}
	Weblink   string
	NameStyle struct {
		Style     string
		Color     *NameColor
		ColorFrom *NameColor `json:"color-from"`
		ColorTo   *NameColor `json:"color-to"`
	} `json:"name-style"`
	Role     string
	Signup   *time.Time
	Location struct {
		Country Location
		Region  *Location
	}
	Twitch        *SocialLink
	Hitbox        *SocialLink
	YouTube       *SocialLink
	Twitter       *SocialLink
	SpeedRunsLive *SocialLink
	Links         []Link
}

// SocialLink is a minimal link that points to an external website.
type SocialLink struct {
	URI string
}

// Location is a country/region code with names.
type Location struct {
	Code  string
	Names struct {
		International string
		Japanese      string
	}
}

// NameColor contains hex color codes for light and dark backgrounds, used
// to display usernames on speedrun.com.
type NameColor struct {
	Light string
	Dark  string
}

// toUser transforms a data blob to a User struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toUser(data interface{}, isResponse bool) *User {
	if data == nil {
		return nil
	}

	if isResponse {
		dest := userResponse{}

		if recast(data, &dest) == nil {
			return &dest.Data
		}
	} else {
		dest := User{}

		if recast(data, &dest) == nil {
			return &dest
		}
	}

	return nil
}

// toUserCollection transforms a data blob to a UserCollection.
// If data is nil or casting was unsuccessful, an empty UserCollection
// is returned.
func toUserCollection(data interface{}) *UserCollection {
	tmp := &UserCollection{}
	recast(data, tmp)

	return tmp
}

// userResponse models the actual API response from the server
type userResponse struct {
	// the one user contained in the response
	Data User
}

// UserByID tries to fetch a single user, identified by their ID.
// When an error is returned, the returned user is nil.
func UserByID(id string) (*User, *Error) {
	return fetchUser(request{"GET", "/users/" + id, nil, nil, nil, ""})
}

// Runs fetches a list of runs done by the user, optionally filtered
// and sorted. This function always returns a RunCollection.
func (u *User) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(u, "runs"), filter, sort, embeds)
}

// ModeratedGames fetches a list of games moderated by the user, optionally
// filtered and sorted. This function always returns a GameCollection.
func (u *User) ModeratedGames(filter *GameFilter, sort *Sorting, embeds string) (*GameCollection, *Error) {
	return fetchGamesLink(firstLink(u, "games"), filter, sort, embeds)
}

// PersonalBests fetches a list of PBs by the user, optionally filtered and
// sorted.
func (u *User) PersonalBests(filter *PersonalBestFilter, embeds string) (*PersonalBestCollection, *Error) {
	return fetchPersonalBestsLink(firstLink(u, "personal-bests"), filter, embeds)
}

// for the 'hasLinks' interface
func (u *User) links() []Link {
	return u.Links
}

// UserFilter represents the possible filtering options when fetching a list
// of users.
type UserFilter struct {
	Lookup        string
	Name          string
	Twitch        string
	Hitbox        string
	Twitter       string
	SpeedRunsLive string
}

// applyToURL merged the filter into a URL.
func (uf *UserFilter) applyToURL(u *url.URL) {
	if uf == nil {
		return
	}

	values := u.Query()

	if len(uf.Lookup) > 0 {
		values.Set("lookup", uf.Lookup)
	}

	if len(uf.Name) > 0 {
		values.Set("name", uf.Name)
	}

	if len(uf.Twitch) > 0 {
		values.Set("twitch", uf.Twitch)
	}

	if len(uf.Hitbox) > 0 {
		values.Set("hitbox", uf.Hitbox)
	}

	if len(uf.Twitter) > 0 {
		values.Set("twitter", uf.Twitter)
	}

	if len(uf.SpeedRunsLive) > 0 {
		values.Set("speedrunslive", uf.SpeedRunsLive)
	}

	u.RawQuery = values.Encode()
}

// Users retrieves a collection of users from  speedrun.com. In most cases, you
// will filter the game, as paging through *all* users takes A LOT of requests.
func Users(f *UserFilter, s *Sorting, c *Cursor) (*UserCollection, *Error) {
	return fetchUsers(request{"GET", "/users", f, s, c, ""})
}

// fetchUser fetches a single user from the network. If the request failed,
// the returned user is nil. Otherwise, the error is nil.
func fetchUser(request request) (*User, *Error) {
	result := &userResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchUserLink tries to fetch a given link and interpret the response as
// a single user. If the link is nil or the user could not be fetched,
// nil is returned.
func fetchUserLink(link requestable) (*User, *Error) {
	if !link.exists() {
		return nil, nil
	}

	return fetchUser(link.request(nil, nil, ""))
}

// fetchUsers fetches a list of users from the network. It always
// returns a collection, even when an error is returned.
func fetchUsers(request request) (*UserCollection, *Error) {
	result := &UserCollection{}
	err := httpClient.do(request, result)

	return result, err
}
