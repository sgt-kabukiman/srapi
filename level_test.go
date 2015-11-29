// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLevels(t *testing.T) {
	countRequests = true

	crashTwinsanityJungleBungle := "lewp5z9n"
	gta1LibertyCityGangstaBang := "zldypd3y"
	jfgCerulean := "yweon79l"

	Convey("Fetching levels by valid IDs", t, func() {
		level, err := LevelByID(crashTwinsanityJungleBungle, NoEmbeds)
		So(err, ShouldBeNil)
		So(level.ID, ShouldEqual, crashTwinsanityJungleBungle)
		So(level.Name, ShouldEqual, "Jungle Bungle")
		So(level.Weblink, ShouldNotBeEmpty)
		So(level.Rules, ShouldNotBeEmpty)
		So(level.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching levels by invalid IDs", t, func() {
		level, err := LevelByID("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(level, ShouldBeNil)
	})

	Convey("Get a level's game", t, func() {
		level, err := LevelByID(gta1LibertyCityGangstaBang, NoEmbeds)
		game, err := level.Game(NoEmbeds)
		So(err, ShouldBeNil)
		So(game, ShouldNotBeNil)
		So(game.Abbreviation, ShouldEqual, "gta1")
	})

	Convey("Get a level's categories", t, func() {
		level, err := LevelByID(jfgCerulean, NoEmbeds)
		categories, err := level.Categories(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(categories, ShouldNotBeNil)
		So(categories.Data, ShouldHaveLength, 3)
		So(categories.Data[0].Name, ShouldEqual, "All Tribals")
	})

	Convey("Get a level's categories via embedding", t, func() {
		level, err := LevelByID(jfgCerulean, "categories")

		before := requestCount
		categories, err := level.Categories(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(categories, ShouldNotBeNil)
		So(categories.Data, ShouldHaveLength, 3)
		So(categories.Data[0].Name, ShouldEqual, "All Tribals")
		So(requestCount, ShouldEqual, before)
	})

	Convey("Get a level's variables", t, func() {
		level, err := LevelByID(jfgCerulean, NoEmbeds)
		variables, err := level.Variables(nil)
		So(err, ShouldBeNil)
		So(variables, ShouldNotBeNil)
		So(variables.Data, ShouldHaveLength, 3)
		So(variables.Data[0].Name, ShouldEqual, "Region")
	})

	Convey("Get a level's variables via embedding", t, func() {
		level, err := LevelByID(jfgCerulean, "variables")

		before := requestCount
		variables, err := level.Variables(nil)
		So(err, ShouldBeNil)
		So(variables, ShouldNotBeNil)
		So(variables.Data, ShouldHaveLength, 3)
		So(variables.Data[0].Name, ShouldEqual, "Region")
		So(requestCount, ShouldEqual, before)
	})

	Convey("Fetch the primary leaderboard for a level", t, func() {
		level, err := LevelByID(gta1LibertyCityGangstaBang, NoEmbeds)
		leaderboard, err := level.PrimaryLeaderboard(nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(leaderboard, ShouldNotBeNil)
		So(len(leaderboard.Runs), ShouldBeGreaterThanOrEqualTo, 1)
	})

	Convey("Fetch the records for a level", t, func() {
		level, err := LevelByID(gta1LibertyCityGangstaBang, NoEmbeds)
		leaderboards, err := level.Records(nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(leaderboards, ShouldNotBeNil)
		So(len(leaderboards.Data), ShouldBeGreaterThanOrEqualTo, 2)
		So(leaderboards.Data[0].Runs, ShouldNotBeEmpty)
	})

	Convey("Fetch the runs for a level", t, func() {
		level, err := LevelByID(gta1LibertyCityGangstaBang, NoEmbeds)
		runs, err := level.Runs(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(runs, ShouldNotBeNil)
		So(len(runs.Data), ShouldBeGreaterThanOrEqualTo, 2)
	})
}
