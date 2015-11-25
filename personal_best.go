// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"net/url"
	"strconv"
)

// PersonalBest represents one PB of a given user in a certain game/category/level
type PersonalBest struct {
	// the rank on the default leaderboard (i.e. with no options set) of this run
	Rank int

	// the run itpb
	Run Run

	// do not use this field directly, use the available methods
	PlatformData interface{} `json:"platform"`

	// do not use this field directly, use the available methods
	RegionData interface{} `json:"region"`

	// do not use this field directly, use the available methods
	PlayersData interface{} `json:"players"`

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	CategoryData interface{} `json:"category"`

	// do not use this field directly, use the available methods
	LevelData interface{} `json:"level"`
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (pb *PersonalBest) Game() (*Game, *Error) {
	if pb.GameData == nil {
		return pb.Run.Game()
	}

	return toGame(pb.GameData), nil
}

// Category extracts the embedded category, if possible, otherwise it will fetch
// the category by doing one additional request. If nothing on the server side is
// fubar, then this function should never return nil.
func (pb *PersonalBest) Category() (*Category, *Error) {
	if pb.CategoryData == nil {
		return pb.Run.Category()
	}

	return toCategory(pb.CategoryData), nil
}

// Level extracts the embedded level, if possible, otherwise it will fetch the
// level by doing one additional request. For full-game runs, this returns nil.
func (pb *PersonalBest) Level() (*Level, *Error) {
	if pb.LevelData == nil {
		return pb.Run.Level()
	}

	return toLevel(pb.LevelData), nil
}

// Platform extracts the embedded platform, if possible, otherwise it will fetch
// the platform by doing one additional request. Not all runs have platforms
// attached, so this can return nil.
func (pb *PersonalBest) Platform() (*Platform, *Error) {
	if pb.PlatformData == nil {
		return pb.Run.Platform()
	}

	return toPlatform(pb.PlatformData), nil
}

// Region extracts the embedded region, if possible, otherwise it will fetch
// the region by doing one additional request. Not all runs have regions
// attached, so this can return nil.
func (pb *PersonalBest) Region() (*Region, *Error) {
	if pb.RegionData == nil {
		return pb.Run.Region()
	}

	return toRegion(pb.RegionData), nil
}

// Players returns a list of all players that aparticipated in this PB.
// If they have not been embedded, they are fetched individually from the
// network, one request per player.
func (pb *PersonalBest) Players() ([]*Player, *Error) {
	if pb.PlayersData == nil {
		return pb.Run.Players()
	}

	return recastToPlayerList(pb.PlayersData), nil
}

// Examiner returns the user that examined the run after submission. This can
// be nil.
func (pb *PersonalBest) Examiner() (*User, *Error) {
	return fetchUserLink(firstLink(&pb.Run, "examiner"))
}

// personalBestsResponse models the actual API response from the server
type personalBestsResponse struct {
	// the contained personal best runs
	Data []PersonalBest
}

// personalBests returns a list of pointers to the PBs; used for cases where
// there is no pagination and the caller wants to return a flat slice of
// PBs instead of a collection (which would be misleading, as collections
// imply pagination).
func (pbr *personalBestsResponse) personalBests() []*PersonalBest {
	var result []*PersonalBest

	for idx := range pbr.Data {
		result = append(result, &pbr.Data[idx])
	}

	return result
}

// PersonalBestFilter represents the possible filtering options when fetching a
// list of PBs.
type PersonalBestFilter struct {
	// If set to >0, only return PBs with a rank better or equal to this.
	Top int

	// series ID
	Series string

	// game ID
	Game string
}

// applyToURL merged the filter into a URL.
func (pbf *PersonalBestFilter) applyToURL(u *url.URL) {
	if pbf == nil {
		return
	}

	values := u.Query()

	if pbf.Top > 0 {
		values.Set("top", strconv.Itoa(pbf.Top))
	}

	if len(pbf.Series) > 0 {
		values.Set("series", pbf.Series)
	}

	if len(pbf.Game) > 0 {
		values.Set("game", pbf.Game)
	}

	u.RawQuery = values.Encode()
}

// fetchVariables fetches a list of PBs from the network. It always
// returns a collection, even when an error is returned.
func fetchPersonalBests(request request) ([]*PersonalBest, *Error) {
	result := &personalBestsResponse{}
	err := httpClient.do(request, result)

	return result.personalBests(), err
}

// fetchPersonalBestsLink tries to fetch a given link and interpret the response as
// a list of PBs. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchPersonalBestsLink(link requestable, filter *PersonalBestFilter) ([]*PersonalBest, *Error) {
	if !link.exists() {
		return make([]*PersonalBest, 0), nil
	}

	return fetchPersonalBests(link.request(filter, nil))
}
