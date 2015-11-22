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

// recastToModeratorMap returns a map of user IDs to their respective moderation
// levels. Note that due to limitations of the speedrun.com API, the mod levels
// are not available when moderators have been embedded. In this case, the
// resulting map containts UnknownModLevel for every user. If you need both,
// there is no other way than to perform two requests.
func recastToModeratorMap(data interface{}) map[string]GameModLevel {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := data.(map[string]GameModLevel)
	if okay {
		return assertedMap
	}

	// maybe we got a list of embedded users
	result := make(map[string]GameModLevel, 0)
	tmp := UserCollection{}

	if recast(data, &tmp) == nil {
		for _, user := range tmp.users() {
			result[user.ID] = UnknownModLevel
		}
	}

	return result
}

// recastToModerators returns a list of users that are moderators of the series.
// If moderators were not embedded, they will be fetched individually from the
// network.
func recastToModerators(data interface{}) []*User {
	// we have a simple map between user IDs and mod levels
	assertedMap, okay := data.(map[string]GameModLevel)
	if okay {
		var result []*User

		for userID := range assertedMap {
			user, err := UserByID(userID)
			if err == nil {
				result = append(result, user)
			}
		}

		return result
	}

	return toUserCollection(data).users()
}
