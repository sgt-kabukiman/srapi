// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

// Guest models a guest on speedrun.com, i.e. someone who is not yet registred but
// already part of the leaderboard.
type Guest struct {
	// the guest's name
	Name string

	// API links to related resources
	Links []Link
}

// toGuest transforms a data blob to a Guest struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toGuest(data interface{}) *Guest {
	dest := Guest{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// guestResponse models the actual API response from the server
type guestResponse struct {
	// the one guest contained in the response
	Data Guest
}

// GuestByName tries to fetch a single guest, identified by their name.
// When an error is returned, the returned guest is nil.
func GuestByName(name string) (*Guest, *Error) {
	return fetchGuest(request{"GET", "/guests/" + url.QueryEscape(name), nil, nil, nil, ""})
}

// Runs fetches a list of runs done by the guest, optionally filtered and sorted.
// This function always returns a RunCollection.
func (g *Guest) Runs(filter *RunFilter, sort *Sorting, embeds string) (*RunCollection, *Error) {
	return fetchRunsLink(firstLink(g, "runs"), filter, sort, embeds)
}

// for the 'hasLinks' interface
func (g *Guest) links() []Link {
	return g.Links
}

// fetchGuest fetches a single guest from the network. If the request failed,
// the returned guest is nil. Otherwise, the error is nil.
func fetchGuest(request request) (*Guest, *Error) {
	result := &guestResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchGuestLink tries to fetch a given link and interpret the response as
// a single guest. If the link is nil or the guest could not be fetched,
// nil is returned.
func fetchGuestLink(link requestable) (*Guest, *Error) {
	if !link.exists() {
		return nil, nil
	}

	return fetchGuest(link.request(nil, nil, ""))
}
