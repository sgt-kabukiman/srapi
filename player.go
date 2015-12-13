// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

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

// toLink returns a link pointing to this player.
func (p *Player) toLink() PlayerLink {
	var link PlayerLink

	if p.User != nil {
		link = PlayerLink{
			Link: Link{
				Relation: "user",
				URI:      BaseURL + "/users/" + p.User.ID,
			},
			ID:   p.User.ID,
			Name: "",
		}

		// try our best to use the provided link
		l := firstLink(p.User, "self")
		if l != nil {
			link.URI = l.URI
		}
	} else {
		link = PlayerLink{
			Link: Link{
				Relation: "guest",
				URI:      BaseURL + "/guests/" + url.QueryEscape(p.Guest.Name),
			},
			ID:   "",
			Name: p.Guest.Name,
		}

		// try our best to use the provided link
		l := firstLink(p.Guest, "self")
		if l != nil {
			link.URI = l.URI
		}
	}

	return link
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
