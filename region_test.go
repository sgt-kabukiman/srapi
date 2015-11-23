// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegions(t *testing.T) {
	Convey("Fetching regions by valid IDs", t, func() {
		id := "mol4z19n"

		region, err := RegionByID(id)

		So(err, ShouldBeNil)
		So(region.ID, ShouldEqual, id)
		So(region.Name, ShouldNotBeBlank)
		So(region.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching regions by invalid IDs", t, func() {
		region, err := RegionByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(region, ShouldBeNil)
	})

	Convey("Fetching multiple regions", t, func() {
		Convey("starting from the beginning", func() {
			regions, err := Regions(nil, nil)

			So(err, ShouldBeNil)
			So(regions.Data, ShouldNotBeEmpty)
			So(regions.Pagination.Offset, ShouldEqual, 0)

			region := regions.Data[0]
			So(region.ID, ShouldNotBeBlank)
			So(region.Name, ShouldNotBeBlank)
			So(region.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			regions, err := Regions(nil, &Cursor{2, 0})

			So(err, ShouldBeNil)
			So(regions.Data, ShouldNotBeEmpty)
			So(regions.Pagination.Offset, ShouldEqual, 2)
			So(regions.Pagination.Links, ShouldNotBeEmpty)

			region := regions.Data[0]
			So(region.ID, ShouldNotBeBlank)
			So(region.Name, ShouldNotBeBlank)
			So(region.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			regions, err := Regions(nil, &Cursor{0, 3})

			So(err, ShouldBeNil)
			So(regions.Data, ShouldHaveLength, 3)
			So(regions.Pagination.Offset, ShouldEqual, 0)
			So(regions.Pagination.Max, ShouldEqual, 3)
			So(regions.Pagination.Links, ShouldNotBeEmpty)

			region := regions.Data[0]
			So(region.ID, ShouldNotBeBlank)
			So(region.Name, ShouldNotBeBlank)
			So(region.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the regions", func() {
			regions, err := Regions(nil, &Cursor{0, 1})

			So(err, ShouldBeNil)
			So(regions.Data, ShouldHaveLength, 1)
			So(regions.Pagination.Offset, ShouldEqual, 0)
			So(regions.Pagination.Max, ShouldEqual, 1)

			regions, err = regions.NextPage()

			So(err, ShouldBeNil)
			So(regions.Data, ShouldHaveLength, 1)
			So(regions.Pagination.Offset, ShouldEqual, 1)
			So(regions.Pagination.Max, ShouldEqual, 1)

			regions, err = regions.NextPage()

			So(err, ShouldBeNil)
			So(regions.Data, ShouldHaveLength, 1)
			So(regions.Pagination.Offset, ShouldEqual, 2)
			So(regions.Pagination.Max, ShouldEqual, 1)

			regions, err = regions.PrevPage()

			So(err, ShouldBeNil)
			So(regions.Data, ShouldHaveLength, 1)
			So(regions.Pagination.Offset, ShouldEqual, 1)
			So(regions.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			regions, err := Regions(nil, nil)

			regions, err = regions.PrevPage()

			So(err, ShouldNotBeNil)
			So(regions, ShouldNotBeNil)
		})
	})

	Convey("Fetching runs of a region", t, func() {
		region, err := RegionByID("e6lxy1dz") // Europe/PAL
		So(err, ShouldBeNil)

		runs := region.Runs(nil, nil)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of runs should be fine", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)

			firstID = runs.Data[0].ID
		})

		runs = region.Runs(nil, &Sorting{Direction: Descending})
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, runs.Data[0].ID)
		})
	})

	Convey("Fetching games of a region", t, func() {
		region, err := RegionByID("e6lxy1dz") // Europe/PAL
		So(err, ShouldBeNil)

		games := region.Games(nil, nil)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.Data[0].ID
		})

		games = region.Games(nil, &Sorting{Direction: Descending})
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.Data[0].ID)
		})
	})
}
