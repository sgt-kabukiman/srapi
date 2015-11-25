// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlatforms(t *testing.T) {
	Convey("Fetching platforms by valid IDs", t, func() {
		id := "o232q83p"

		platform, err := PlatformByID(id)

		So(err, ShouldBeNil)
		So(platform.ID, ShouldEqual, id)
		So(platform.Name, ShouldNotBeBlank)
		So(platform.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching platforms by invalid IDs", t, func() {
		platform, err := PlatformByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(platform, ShouldBeNil)
	})

	Convey("Fetching multiple platforms", t, func() {
		Convey("starting from the beginning", func() {
			platforms, err := Platforms(nil, nil)
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldNotBeEmpty)
			So(platforms.Pagination.Offset, ShouldEqual, 0)

			platform := platforms.Data[0]
			So(platform.ID, ShouldNotBeBlank)
			So(platform.Name, ShouldNotBeBlank)
			So(platform.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			platforms, err := Platforms(nil, &Cursor{2, 0})
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldNotBeEmpty)
			So(platforms.Pagination.Offset, ShouldEqual, 2)
			So(platforms.Pagination.Links, ShouldNotBeEmpty)

			platform := platforms.Data[0]
			So(platform.ID, ShouldNotBeBlank)
			So(platform.Name, ShouldNotBeBlank)
			So(platform.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			platforms, err := Platforms(nil, &Cursor{0, 3})
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldHaveLength, 3)
			So(platforms.Pagination.Offset, ShouldEqual, 0)
			So(platforms.Pagination.Max, ShouldEqual, 3)
			So(platforms.Pagination.Links, ShouldNotBeEmpty)

			platform := platforms.Data[0]
			So(platform.ID, ShouldNotBeBlank)
			So(platform.Name, ShouldNotBeBlank)
			So(platform.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the platforms", func() {
			platforms, err := Platforms(nil, &Cursor{0, 1})
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldHaveLength, 1)
			So(platforms.Pagination.Offset, ShouldEqual, 0)
			So(platforms.Pagination.Max, ShouldEqual, 1)

			platforms, err = platforms.NextPage()
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldHaveLength, 1)
			So(platforms.Pagination.Offset, ShouldEqual, 1)
			So(platforms.Pagination.Max, ShouldEqual, 1)

			platforms, err = platforms.NextPage()
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldHaveLength, 1)
			So(platforms.Pagination.Offset, ShouldEqual, 2)
			So(platforms.Pagination.Max, ShouldEqual, 1)

			platforms, err = platforms.PrevPage()
			So(err, ShouldBeNil)
			So(platforms.Data, ShouldHaveLength, 1)
			So(platforms.Pagination.Offset, ShouldEqual, 1)
			So(platforms.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			platforms, err := Platforms(nil, nil)

			platforms, err = platforms.PrevPage()
			So(err, ShouldNotBeNil)
			So(platforms, ShouldNotBeNil)
		})
	})

	Convey("Fetching runs of a platform", t, func() {
		platform, err := PlatformByID("o232q83p") // Gameboy
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
		platform, err := PlatformByID("o232q83p") // Gameboy
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
