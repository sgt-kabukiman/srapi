// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// Player is either a User or a Guest, i.e. only one of the two will ever be
// non-nil.
type Player struct {
	User  *User
	Guest *Guest
}

// Name returns the international name for the player.
func (p *Player) Name() string {
	if p.User != nil {
		return p.User.Names.International
	} else if p.Guest != nil {
		return p.Guest.Name
	} else {
		return "(neither guest nor user)"
	}
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

// checks if the link exists
func (pl *PlayerLink) exists() bool {
	return pl != nil
}

// request turns a link into a request
func (pl *PlayerLink) request(filter filter, sort *Sorting, embeds string) request {
	relURL := pl.URI[len(BaseURL):]

	return request{"GET", relURL, filter, sort, nil, embeds}
}

// fetch retrieves the user or guest the link points to
func (pl *PlayerLink) fetch() (*Player, *Error) {
	player := &Player{}

	switch pl.Relation {
	case "user":
		user, err := fetchUserLink(pl)
		if err != nil {
			return player, err
		}

		player.User = user

	case "guest":
		guest, err := fetchGuestLink(pl)
		if err != nil {
			return player, err
		}

		player.Guest = guest
	}

	return player, nil
}
