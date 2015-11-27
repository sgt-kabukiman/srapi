// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCategories(t *testing.T) {
	gtavcAny := "nxd1rk8q"
	gta1Any := "jzd368dn"
	ctrAny := "w9d846kn"

	Convey("Fetching categories by valid IDs", t, func() {
		category, err := CategoryByID(gtavcAny, NoEmbeds)
		So(err, ShouldBeNil)
		So(category.ID, ShouldEqual, gtavcAny)
		So(category.Name, ShouldEqual, "Any%")
		So(category.Type, ShouldEqual, "per-game")
		So(category.Weblink, ShouldNotBeEmpty)
		So(category.Rules, ShouldNotBeEmpty)
		So(category.Players.Type, ShouldEqual, "exactly")
		So(category.Players.Value, ShouldEqual, 1)
		So(category.Miscellaneous, ShouldBeFalse)
		So(category.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching categories by invalid IDs", t, func() {
		category, err := CategoryByID("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(category, ShouldBeNil)
	})

	Convey("Get a category's game", t, func() {
		category, err := CategoryByID(gtavcAny, NoEmbeds)

		game, err := category.Game(NoEmbeds)
		So(err, ShouldBeNil)
		So(game, ShouldNotBeNil)
		So(game.Abbreviation, ShouldEqual, "gtavc")
	})

	Convey("Get a category's game via embedding", t, func() {
		category, err := CategoryByID(gtavcAny, "game")

		before := requestCount
		game, err := category.Game(NoEmbeds)
		So(err, ShouldBeNil)
		So(game, ShouldNotBeNil)
		So(game.Abbreviation, ShouldEqual, "gtavc")
		So(requestCount, ShouldEqual, before)
	})

	Convey("Get a category's variables", t, func() {
		category, err := CategoryByID(ctrAny, NoEmbeds)
		variables, err := category.Variables(nil)
		So(err, ShouldBeNil)
		So(variables, ShouldNotBeNil)
		So(variables, ShouldHaveLength, 1)
		So(variables[0].Name, ShouldEqual, "Character")
	})

	Convey("Get a category's variables via embedding", t, func() {
		category, err := CategoryByID(ctrAny, "variables")

		before := requestCount
		variables, err := category.Variables(nil)
		So(err, ShouldBeNil)
		So(variables, ShouldNotBeNil)
		So(variables, ShouldHaveLength, 1)
		So(variables[0].Name, ShouldEqual, "Character")
		So(requestCount, ShouldEqual, before)
	})

	Convey("Fetch the primary leaderboard for a category", t, func() {
		category, err := CategoryByID(gtavcAny, NoEmbeds)
		leaderboard, err := category.PrimaryLeaderboard(&LeaderboardOptions{Top: 5}, NoEmbeds)
		So(err, ShouldBeNil)
		So(leaderboard, ShouldNotBeNil)
		So(leaderboard.Runs, ShouldHaveLength, 5)
	})

	Convey("Fetch the records for a per-game category", t, func() {
		category, err := CategoryByID(gtavcAny, NoEmbeds)
		leaderboards, err := category.Records(nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(leaderboards, ShouldNotBeNil)
		So(leaderboards.Data, ShouldHaveLength, 1)
		So(leaderboards.Data[0].Runs, ShouldNotBeEmpty)
	})

	Convey("Fetch the records for a per-level category", t, func() {
		category, err := CategoryByID(gta1Any, NoEmbeds) // GTA1 Any%
		leaderboards, err := category.Records(nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(leaderboards, ShouldNotBeNil)
		So(len(leaderboards.Data), ShouldBeGreaterThanOrEqualTo, 4)
		So(leaderboards.Data[0].Runs, ShouldNotBeEmpty)
	})

	Convey("Fetch the runs for a category", t, func() {
		category, err := CategoryByID(gta1Any, NoEmbeds) // GTA1 Any%
		runs, err := category.Runs(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(runs, ShouldNotBeNil)
		So(len(runs.Data), ShouldBeGreaterThanOrEqualTo, 6)
	})
}
