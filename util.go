package srapi

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
