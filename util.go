// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "encoding/json"

// hasLinks describes a struct that has API links attached to it
type hasLinks interface {
	links() []Link
}

// firstLink returns the first link with a matching relation attribute or nil
// if there is no such link.
func firstLink(linked hasLinks, name string) *Link {
	for _, link := range linked.links() {
		if link.Relation == name {
			return &link
		}
	}

	return nil
}

// recast takes some data blob, turns it into JSON and unmarshals that JSON
// into dest. It's a very ugly hack to type-assert structures which are not
// type-safe during compile time.
func recast(data interface{}, dest interface{}) error {
	// convert generic mess into JSON
	encoded, _ := json.Marshal(data)

	// ... and try to turn it back into something meaningful
	return json.Unmarshal(encoded, dest)
}
