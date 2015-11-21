package srapi

import "encoding/json"

type hasLinks interface {
	links() []Link
}

func firstLink(linked hasLinks, name string) *Link {
	for _, link := range linked.links() {
		if link.Relation == name {
			return &link
		}
	}

	return nil
}

func recast(data interface{}, dest interface{}) error {
	// convert generic mess into JSON
	encoded, _ := json.Marshal(data)

	// ... and try to turn it back into something meaningful
	return json.Unmarshal(encoded, dest)
}
