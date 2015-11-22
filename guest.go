package srapi

import "net/url"

type Guest struct {
	Name  string
	Links []Link
}

type guestResponse struct {
	Data Guest
}

func GuestById(name string) (*Guest, *Error) {
	return fetchGuest(request{"GET", "/guests/" + url.QueryEscape(name), nil, nil, nil})
}

func fetchGuest(request request) (*Guest, *Error) {
	result := &guestResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Guest) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	link := firstLink(self, "runs")
	if link == nil {
		return nil
	}

	runs, _ := fetchRuns(link.request(filter, sort))
	return runs
}

// for the 'hasLinks' interface
func (self *Guest) links() []Link {
	return self.Links
}
