// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlatforms(t *testing.T) {
	countRequests = true

	gameboy := "o232q83p"

	Convey("Fetching platforms by valid IDs", t, func() {
		platform, err := PlatformByID(gameboy)
		So(err, ShouldBeNil)
		So(platform.ID, ShouldEqual, gameboy)
		So(platform.Name, ShouldEqual, "Game Boy")
		So(platform.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching platforms by invalid IDs", t, func() {
		platform, err := PlatformByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(platform, ShouldBeNil)
	})

	Convey("Fetching multiple platforms", t, func() {
		platforms, err := Platforms(nil, &Cursor{0, 1})
		So(err, ShouldBeNil)
		So(platforms.Pagination.Offset, ShouldEqual, 0)
		So(platforms.Pagination.Max, ShouldEqual, 1)

		num := 0

		// read a few pages, 7 is arbitrary
		platforms.Walk(func(p *Platform) bool {
			So(p.ID, ShouldNotBeBlank)

			num++
			return num < 7
		})
	})

	Convey("Fetching runs of a platform", t, func() {
		platform, err := PlatformByID(gameboy)
		So(err, ShouldBeNil)

		runs, err := platform.Runs(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of runs should be fine", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)

			firstID = runs.Data[0].ID
		})

		runs, err = platform.Runs(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, runs.Data[0].ID)
		})
	})

	Convey("Fetching games of a platform", t, func() {
		platform, err := PlatformByID(gameboy)
		So(err, ShouldBeNil)

		games, err := platform.Games(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.Data[0].ID
		})

		games, err = platform.Games(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.Data[0].ID)
		})
	})
}
