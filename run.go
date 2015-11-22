// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

// Run represents a single run.
type Run struct {
	// unique ID
	ID string

	// link to the run on speedrun.com
	Weblink string

	// Videos submitted for the run
	Videos struct {
		// the original submission value for the "videos" form field on speedrun.com
		Text string

		// list of links pointing to videos on external websites, most likely
		// a lot of twitch or youtube links
		Links []Link
	}

	// the runner's comment
	Comment string

	// run status
	Status struct {
		// Status can be "new", "verified" or "rejected"
		Status string

		// user ID of the last examiner
		Examiner string

		// date when the run was verified
		VerifyDate string `json:"verify-date"`

		// If Status is "rejected", then this field possibly contains the reason
		// for the rejection.
		Reason string
	}

	// the date the run was done on
	Date string

	// timing information, not all of them are filled all the time, except
	// the Primary one
	Times struct {
		// primary time (as defined by the game's DefaultTime setting) as a
		// ISO 8601 duration
		Primary string

		// realtime as a ISO 8601 duration
		Realtime string

		// realtime without loading times as a ISO 8601 duration
		RealtimeWithoutLoads string `json:"realtime_noloads"`

		// in-game time a ISO 8601 duration
		IngameTime string `json:"ingame"`
	}

	// the system the run was done on
	System struct {
		// platform ID
		Platform string

		// whether or not the run was done using an emulator
		Emulated bool

		// region ID
		Region string
	}

	// If available, a link pointing to a website containing the splits. As of
	// 2015, if a link is present, it is pointing to splits.io.
	Splits *Link

	// variable values for the run (mapping of variable ID to value ID)
	Values map[string]string

	// API links to related resources
	Links []Link

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

// runResponse models the actual API response from the server
type runResponse struct {
	// the one run contained in the response
	Data Run
}

// RunByID tries to fetch a single run, identified by its ID.
// When an error is returned, the returned run is nil.
func RunByID(id string) (*Run, *Error) {
	return fetchRun(request{"GET", "/runs/" + id, nil, nil, nil})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (r *Run) Game() *Game {
	// we only have the game ID at hand
	asserted, okay := r.GameData.(string)
	if okay {
		game, _ := GameByID(asserted)
		return game
	}

	return toGame(r.GameData)
}

// Category extracts the embedded category, if possible, otherwise it will fetch
// the game by doing one additional request. If nothing on the server side is
// fubar, then this function should never return nil.
func (r *Run) Category() *Category {
	if r.CategoryData == nil { // should never happen
		return nil
	}

	// we only have the category ID at hand
	asserted, okay := r.CategoryData.(string)
	if okay {
		category, _ := CategoryByID(asserted)
		return category
	}

	return toCategory(r.CategoryData)
}

// Level extracts the embedded level, if possible, otherwise it will fetch
// the game by doing one additional request. It's possible for runs to not have
// levels, so this function can return nil for full-game runs.
func (r *Run) Level() *Level {
	if r.LevelData == nil {
		return nil
	}

	// we only have the level ID at hand
	asserted, okay := r.LevelData.(string)
	if okay {
		level, _ := LevelByID(asserted)
		return level
	}

	return toLevel(r.LevelData)
}

// Platform extracts the embedded platform, if possible, otherwise it will fetch
// the game by doing one additional request. Some runs don't have platforms
// attached, so this can return nil.
func (r *Run) Platform() *Platform {
	if r.PlatformData == nil {
		if len(r.System.Platform) > 0 {
			platform, _ := PlatformByID(r.System.Platform)
			return platform
		}

		return nil
	}

	return toPlatform(r.PlatformData)
}

// Region extracts the embedded region, if possible, otherwise it will fetch
// the game by doing one additional request. Some runs don't have regions
// attached, so this can return nil.
func (r *Run) Region() *Region {
	if r.RegionData == nil {
		if len(r.System.Region) > 0 {
			region, _ := RegionByID(r.System.Region)
			return region
		}

		return nil
	}

	return toRegion(r.RegionData)
}

// Players returns a list of all players that aparticipated in this run.
// If they have not been embedded, they are fetched individually from the
// network, one request per player.
func (r *Run) Players() []*Player {
	var result []*Player

	switch asserted := r.PlayersData.(type) {
	// list of simple links to users/guests, e.g. players=[{rel:..,id:...}, {...}]
	case []interface{}:
		var tmp []PlayerLink

		if recast(asserted, &tmp) == nil {
			for _, link := range tmp {
				player := Player{}

				switch link.Relation {
				case "user":
					if user := fetchUserLink(&link); user != nil {
						player.User = user
					}

				case "guest":
					if guest := fetchGuestLink(&link); guest != nil {
						player.Guest = guest
					}
				}

				if player.User != nil || player.Guest != nil {
					result = append(result, &player)
				}
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		tmp := playerCollection{}

		if recast(asserted, &tmp) == nil {
			// each element in tmp.Data has a rel that tells us whether we have a
			// user or a guest
			for _, playerProps := range tmp.Data {
				rel, exists := playerProps["rel"]
				if exists {
					player := Player{}

					switch rel {
					case "user":
						if user := toUser(playerProps); user != nil {
							player.User = user
						}

					case "guest":
						if guest := toGuest(playerProps); guest != nil {
							player.Guest = guest
						}
					}

					if player.User != nil || player.Guest != nil {
						result = append(result, &player)
					}
				}
			}
		}
	}

	return result
}

// Examiner returns the user that examined the run after submission. This can
// be nil, especially for new runs.
func (r *Run) Examiner() *User {
	return fetchUserLink(firstLink(r, "examiner"))
}

// for the 'hasLinks' interface
func (r *Run) links() []Link {
	return r.Links
}

// RunFilter represents the possible filtering options when fetching a list of
// runs.
type RunFilter struct {
	// a user ID
	User string

	// the name of a guest
	Guest string

	// user ID to fetch runs last examined by this user
	Examiner string

	// game ID
	Game string

	// level ID
	Level string

	// category ID
	Category string

	// platform ID
	Platform string

	// region ID
	Region string

	// when set, controls if all or no runs are on emulator
	Emulated OptionalFlag

	// can be set to "new", "verified" or "rejected"
	Status string
}

// applyToURL merged the filter into a URL.
func (rf *RunFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(rf.User) > 0 {
		values.Set("user", rf.User)
	}

	if len(rf.Guest) > 0 {
		values.Set("guest", rf.Guest)
	}

	if len(rf.Examiner) > 0 {
		values.Set("examiner", rf.Examiner)
	}

	if len(rf.Game) > 0 {
		values.Set("game", rf.Game)
	}

	if len(rf.Level) > 0 {
		values.Set("level", rf.Level)
	}

	if len(rf.Category) > 0 {
		values.Set("category", rf.Category)
	}

	if len(rf.Platform) > 0 {
		values.Set("platform", rf.Platform)
	}

	if len(rf.Region) > 0 {
		values.Set("region", rf.Region)
	}

	if len(rf.Status) > 0 {
		values.Set("status", rf.Status)
	}

	rf.Emulated.applyToQuery("emulated", &values)

	u.RawQuery = values.Encode()
}

// RunCollection is one page of a run list. It consists of the runs as well as
// some pagination information (like links to the next or previous page).
type RunCollection struct {
	Data       []Run
	Pagination Pagination
}

// Runs retrieves a collection of runs, most likely filtered and sorted.
func Runs(f *RunFilter, s *Sorting, c *Cursor) (*RunCollection, *Error) {
	return fetchRuns(request{"GET", "/runs", f, s, c})
}

// runs returns a list of pointers to the runs; used for cases where
// there is no pagination and the caller wants to return a flat slice of
// runs instead of a collection (which would be misleading, as collections
// imply pagination).
func (rc *RunCollection) runs() []*Run {
	var result []*Run

	for idx := range rc.Data {
		result = append(result, &rc.Data[idx])
	}

	return result
}

// NextPage tries to follow the "next" link and retrieve the next page of
// runs. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (rc *RunCollection) NextPage() (*RunCollection, *Error) {
	return rc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// runs. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (rc *RunCollection) PrevPage() (*RunCollection, *Error) {
	return rc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (rc *RunCollection) fetchLink(name string) (*RunCollection, *Error) {
	next := firstLink(&rc.Pagination, name)
	if next == nil {
		return &RunCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchRuns(next.request(nil, nil))
}

// fetchRun fetches a single run from the network. If the request failed,
// the returned run is nil. Otherwise, the error is nil.
func fetchRun(request request) (*Run, *Error) {
	result := &runResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchRunLink tries to fetch a given link and interpret the response as
// a single run. If the link is nil or the run could not be fetched,
// nil is returned.
func fetchRuns(request request) (*RunCollection, *Error) {
	result := &RunCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchRunsLink tries to fetch a given link and interpret the response as
// a list of runs. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchRunsLink(link requestable, filter filter, sort *Sorting) *RunCollection {
	if link == nil {
		return &RunCollection{}
	}

	collection, _ := fetchRuns(link.request(filter, sort))
	return collection
}
