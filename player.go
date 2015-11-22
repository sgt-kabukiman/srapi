// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Player is either a User or a Guest, i.e. only one of the two will ever be
// non-nil.
type Player struct {
	User  *User
	Guest *Guest
}

// PlayerLink is a special link that points to either a user (then ID is given)
// or a guest (then Name is given).
type PlayerLink struct {
	Link

	// user ID
	ID string

	// guest name
	Name string
}

// request turns a link into a request
func (pl *PlayerLink) request(filter filter, sort *Sorting) request {
	relURL := pl.URI[len(BaseURL):]

	return request{"GET", relURL, filter, sort, nil}
}

// playerCollection is a list of players, used inside Run structs
type playerCollection struct {
	Data []map[string]interface{}
}
