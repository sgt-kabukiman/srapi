// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

type Guest struct {
	Name  string
	Links []Link
}

func toGuest(data interface{}) *Guest {
	dest := Guest{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

type guestResponse struct {
	Data Guest
}

func GuestById(name string) (*Guest, *Error) {
	return fetchGuest(request{"GET", "/guests/" + url.QueryEscape(name), nil, nil, nil})
}

func (self *Guest) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(self, "runs"), filter, sort)
}

// for the 'hasLinks' interface
func (self *Guest) links() []Link {
	return self.Links
}

func fetchGuest(request request) (*Guest, *Error) {
	result := &guestResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func fetchGuestLink(link *Link) *Guest {
	if link == nil {
		return nil
	}

	guest, _ := fetchGuest(link.request(nil, nil))
	return guest
}
