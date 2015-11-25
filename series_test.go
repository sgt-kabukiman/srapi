// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSeries(t *testing.T) {
	Convey("Fetching series by valid IDs", t, func() {
		id := "9v7og6n0" // GTA

		series, err := SeriesByID(id)

		So(err, ShouldBeNil)
		So(series.ID, ShouldEqual, id)
		So(series.Names.International, ShouldEqual, "Grand Theft Auto")
		So(series.Abbreviation, ShouldEqual, "gta")
		So(series.Weblink, ShouldNotBeEmpty)
		So(series.Links, ShouldNotBeEmpty)
		So(series.ModeratorMap(), ShouldNotBeEmpty)

		mods, err := series.Moderators()
		So(err, ShouldBeNil)
		So(mods, ShouldNotBeEmpty)
	})

	Convey("Fetching series by valid abbreviation", t, func() {
		series, err := SeriesByAbbreviation("gta")

		So(err, ShouldBeNil)
		So(series.ID, ShouldEqual, "9v7og6n0")
		So(series.Names.International, ShouldEqual, "Grand Theft Auto")
		So(series.Abbreviation, ShouldEqual, "gta")
		So(series.Weblink, ShouldNotBeEmpty)
		So(series.Links, ShouldNotBeEmpty)
		So(series.ModeratorMap(), ShouldNotBeEmpty)

		mods, err := series.Moderators()
		So(err, ShouldBeNil)
		So(mods, ShouldNotBeEmpty)
	})

	Convey("Fetching series by invalid IDs", t, func() {
		series, err := SeriesByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(series, ShouldBeNil)
	})

	Convey("Fetching series by invalid abbrevitation", t, func() {
		series, err := SeriesByAbbreviation("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(series, ShouldBeNil)
	})

	Convey("Fetching multiple series", t, func() {
		Convey("starting from the beginning", func() {
			seriesList, err := ManySeries(nil, nil, nil)
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldNotBeEmpty)
			So(seriesList.Pagination.Offset, ShouldEqual, 0)

			series := seriesList.Data[0]
			So(series.ID, ShouldNotBeBlank)
			So(series.Names.International, ShouldNotBeBlank)
			So(series.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			seriesList, err := ManySeries(nil, nil, &Cursor{2, 0})
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldNotBeEmpty)
			So(seriesList.Pagination.Offset, ShouldEqual, 2)
			So(seriesList.Pagination.Links, ShouldNotBeEmpty)

			series := seriesList.Data[0]
			So(series.ID, ShouldNotBeBlank)
			So(series.Names.International, ShouldNotBeBlank)
			So(series.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			seriesList, err := ManySeries(nil, nil, &Cursor{0, 3})
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 3)
			So(seriesList.Pagination.Offset, ShouldEqual, 0)
			So(seriesList.Pagination.Max, ShouldEqual, 3)
			So(seriesList.Pagination.Links, ShouldNotBeEmpty)

			series := seriesList.Data[0]
			So(series.ID, ShouldNotBeBlank)
			So(series.Names.International, ShouldNotBeBlank)
			So(series.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the seriesList", func() {
			seriesList, err := ManySeries(nil, nil, &Cursor{0, 1})
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)
			So(seriesList.Pagination.Offset, ShouldEqual, 0)
			So(seriesList.Pagination.Max, ShouldEqual, 1)

			seriesList, err = seriesList.NextPage()
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)
			So(seriesList.Pagination.Offset, ShouldEqual, 1)
			So(seriesList.Pagination.Max, ShouldEqual, 1)

			seriesList, err = seriesList.NextPage()
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)
			So(seriesList.Pagination.Offset, ShouldEqual, 2)
			So(seriesList.Pagination.Max, ShouldEqual, 1)

			seriesList, err = seriesList.PrevPage()
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)
			So(seriesList.Pagination.Offset, ShouldEqual, 1)
			So(seriesList.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			seriesList, err := ManySeries(nil, nil, nil)

			seriesList, err = seriesList.PrevPage()
			So(err, ShouldNotBeNil)
			So(seriesList, ShouldNotBeNil)
		})

		Convey("test the SeriesFilter", func() {
			// check abbrevitation
			filter := SeriesFilter{Abbreviation: "gta"}

			seriesList, err := ManySeries(&filter, nil, nil)
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)

			// check name
			filter = SeriesFilter{Name: "mario"}
			cursor := Cursor{Max: 5}

			seriesList, err = ManySeries(&filter, nil, &cursor)
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 5)

			// check moderator
			filter = SeriesFilter{Moderator: "r5j52gjv"}

			seriesList, err = ManySeries(&filter, nil, nil)
			So(err, ShouldBeNil)
			So(len(seriesList.Data), ShouldBeBetween, 2, 5) // Sorry Josh, but I don't assume it's gonna be more than 5 #Kappa
		})
	})

	Convey("Fetching games of a series", t, func() {
		series, err := SeriesByID("9v7og6n0") // Gameboy
		So(err, ShouldBeNil)

		games, err := series.Games(nil, nil)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.Data[0].ID
		})

		games, err = series.Games(nil, &Sorting{Direction: Descending})
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.Data[0].ID)
		})
	})
}
