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
	result := make(map[string]GameModLevel)

	// both embedded and non-embedded moderators look at least like this
	// and we cannot type assert to map[string]string directly.
	assertedMap, okay := data.(map[string]interface{})
	if okay {
		if isResponseLike(assertedMap) {
			tmp := UserCollection{}

			if recast(data, &tmp) == nil {
				for _, user := range tmp.users() {
					result[user.ID] = UnknownModLevel
				}
			}
		} else {
			for userID, something := range assertedMap {
				level, okay := something.(string)
				if okay {
					result[userID] = GameModLevel(level)
				}
			}
		}
	}

	return result
}

// recastToModerators returns a list of users that are moderators of the series.
// If moderators were not embedded, they will be fetched individually from the
// network.
func recastToModerators(data interface{}) ([]*User, *Error) {
	var result []*User

	// both embedded and non-embedded moderators look at least like this
	// and we cannot type assert to map[string]string directly.
	assertedMap, okay := data.(map[string]interface{})
	if okay {
		if isResponseLike(assertedMap) {
			return toUserCollection(data).users(), nil
		}

		for userID := range assertedMap {
			user, err := UserByID(userID)
			if err != nil {
				return result, err
			}

			result = append(result, user)
		}
	}

	return result, nil
}

// recastToPlayerList casts a player blob into a list of players.
func recastToPlayerList(data interface{}) []*Player {
	var result []*Player

	tmp := playerCollection{}
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
					result = append(result, &player)
				}
			}
		}
	}

	return result
}

func isResponseLike(data map[string]interface{}) bool {
	isResponse := len(data) == 1

	for key := range data {
		if key != "data" {
			isResponse = false
		}
	}

	return isResponse
}
