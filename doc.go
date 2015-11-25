// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

// Package srapi implements a wrapper around the speedrun.com REST API.
//
// This package supports the anonymous parts of the API, meaning that notifications
// and the profile (for which an API key must be used) are not included.
//
// There is no concept of a client struct or interface; package users will just
// call global functions which rely on the internal http client. Observe this
// simple example:
//
//     import "github.com/sgt-kabukiman/srapi"
//
//     game, err := srapi.GameByAbbreviation("smw", srapi.NoEmbeds)
//     if err == nil {
//         categories := game.Categories(nil, nil, srapi.NoEmbeds)
//     }
//
// Usually, there are two functions per resource; one to get a single object
// (like Game(string)) and one to fetch a collection of objects (like Games()).
// For collections, it's usually possible to specify a filter, sorting options
// as well as a cursor (collections are paginated). All three are optional.
//
// Where applicable, embeds can be used to fetch multiple related resources in
// one request. The package does its best to handle related resources transparently,
// i.e. it will use embedded data when available and otherwise fall back to
// performing more requests as needed. When filtering/sorting options are
// available (e.g. in game.Categories()), those are only effective when the
// resources are *not* embedded. A future improvement to this package could be to
// apply the sorting manually to the embedded data, but for now, that's not being
// done.
//
// Note that due to some limitations in the API, embedding resources can sometimes
// make information unavailable: Embedding moderators in a game resource loses
// the moderator levels. This is not a bug in this package, but a limitation of
// the actual API.
//
// For simplicity reasons, the exported structs often contain *Data properties,
// like ModeratorsData, PlatformsData etc. -- these are used for handling embeds
// and should never be touched by anyone outside of this package. They are a quirk
// and the cost for not wrapping everything in unexported structs and having large
// interfaces all over the place.
//
// This package does not handle any throttling itself. Note that the speedrun.com
// API only allows for a certain number of requests per minute, so make sure to
// slow your calls down yourself if needed.
//
// Due to the usage of net.http, this package is safe for use in concurrent
// goroutines.
package srapi
