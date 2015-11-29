// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

// playersResponse is a list of players, used as a helper struct in recastToPlayerList.
type playersResponse struct {
	Data []map[string]interface{}
}

// toPlayerCollection transforms a data blob to a PlayerCollection.
// If data is nil or casting was unsuccessful, an empty PlayerCollection
// is returned.
func toPlayerCollection(data interface{}) *PlayerCollection {
	result := &PlayerCollection{}

	if data == nil {
		return result
	}

	tmp := playersResponse{}
	if recast(data, &tmp) == nil {
		// each element in tmp.Data has a rel that tells us whether we have a
		// user or a guest
		for _, playerProps := range tmp.Data {
			rel, exists := playerProps["rel"]
			if exists {
				player := Player{}

				switch rel {
				case "user":
					if user := toUser(playerProps, false); user != nil {
						player.User = user
					}

				case "guest":
					if guest := toGuest(playerProps, false); guest != nil {
						player.Guest = guest
					}
				}

				if player.User != nil || player.Guest != nil {
					result.Data = append(result.Data, player)
				}
			}
		}
	}

	return result
}

// PlayerCollection is a list of players.
type PlayerCollection struct {
	Data []Player
}

// PlayerWalkerFunc is a function that can be used in Walk(). If it returns
// true, walking continues, else the walk stops.
type PlayerWalkerFunc func(g *Player) bool

// Players returns a list of pointers to the structs; used for cases where
// there is no pagination and the caller wants to return a flat slice of items
// instead of a collection (which would be misleading, as collections imply
// pagination).
func (c *PlayerCollection) Players() []*Player {
	var result []*Player

	c.Walk(func(item *Player) bool {
		result = append(result, item)
		return true
	})

	return result
}

// Walk applies a function to all items in the collection, in order. If the
// function returns false, iterating will be stopped.
func (c *PlayerCollection) Walk(f PlayerWalkerFunc) {
	it := c.Iterator()

	for item := it.Start(); item != nil; item = it.Next() {
		if !f(item) {
			break
		}
	}
}

// Users is like Players(), but only returns users and skips guests.
func (c *PlayerCollection) Users() []*User {
	var result []*User

	c.Walk(func(item *Player) bool {
		if item.User != nil {
			result = append(result, item.User)
		}

		return true
	})

	return result
}

// Guests is like Players(), but only returns users and skips guests.
func (c *PlayerCollection) Guests() []*Guest {
	var result []*Guest

	c.Walk(func(item *Player) bool {
		if item.Guest != nil {
			result = append(result, item.Guest)
		}

		return true
	})

	return result
}

// Size returns the number of elements in the collection.
func (c *PlayerCollection) Size() int {
	return len(c.Data)
}

// Get returns the n-th element (the first one has idx 0) and nil if there is
// no such index.
func (c *PlayerCollection) Get(idx int) *Player {
	// easy, the idx is on this page
	if idx < len(c.Data) {
		return &c.Data[idx]
	}

	return nil
}

// First returns the first element, if any, otherwise nil.
func (c *PlayerCollection) First() *Player {
	if len(c.Data) == 0 {
		return nil
	}

	return &c.Data[0]
}

// Iterator returns an interator for a PlayerCollection. There can be many
// independent iterators starting from the same collection.
func (c *PlayerCollection) Iterator() PlayerIterator {
	return PlayerIterator{
		origin: c,
		cursor: 0,
	}
}

// PlayerIterator represents a list of games.
type PlayerIterator struct {
	origin *PlayerCollection
	page   *PlayerCollection
	cursor int
}

// Start returns the iterator to the start of the original collection page
// and returns the first element if it exists.
func (i *PlayerIterator) Start() *Player {
	i.cursor = 0
	i.page = i.origin

	return i.fetch()
}

// Next advances to the next item. If there is no further item, nil is
// returned. All further calls to Next would return nil as well.
func (i *PlayerIterator) Next() *Player {
	i.cursor++

	return i.fetch()
}

// fetch tries to return the current item. If it doesn't exist, it attempts
// to fetch the next page and return its first item.
func (i *PlayerIterator) fetch() *Player {
	// easy, just get the next item on the current page
	if i.cursor < len(i.page.Data) {
		return &i.page.Data[i.cursor]
	}

	return nil
}
