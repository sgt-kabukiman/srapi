// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSeries(t *testing.T) {
	countRequests = true

	gta := "9v7og6n0"

	Convey("Fetching series by valid IDs", t, func() {
		series, err := SeriesByID(gta, NoEmbeds)

		So(err, ShouldBeNil)
		So(series.ID, ShouldEqual, gta)
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
		series, err := SeriesByAbbreviation("gta", NoEmbeds)

		So(err, ShouldBeNil)
		So(series.ID, ShouldEqual, gta)
		So(series.Names.International, ShouldEqual, "Grand Theft Auto")
		So(series.Abbreviation, ShouldEqual, "gta")
		So(series.Weblink, ShouldNotBeEmpty)
		So(series.Links, ShouldNotBeEmpty)

		m := series.ModeratorMap()
		So(m, ShouldNotBeEmpty)

		for _, level := range m {
			So(level, ShouldNotEqual, UnknownModLevel)
		}

		mods, err := series.Moderators()
		So(err, ShouldBeNil)
		So(mods, ShouldNotBeEmpty)
	})

	Convey("Fetching series by invalid IDs", t, func() {
		series, err := SeriesByID("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(series, ShouldBeNil)
	})

	Convey("Fetching series by invalid abbrevitation", t, func() {
		series, err := SeriesByAbbreviation("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(series, ShouldBeNil)
	})

	Convey("embed moderators in series", t, func() {
		series, err := SeriesByID(gta, "moderators")
		So(err, ShouldBeNil)

		before := requestCount
		m := series.ModeratorMap()
		So(m, ShouldNotBeEmpty)
		So(requestCount, ShouldEqual, before)

		for _, level := range m {
			So(level, ShouldEqual, UnknownModLevel)
		}

		mods, err := series.Moderators()
		So(err, ShouldBeNil)
		So(mods.Size(false), ShouldBeBetween, 3, 100)
		So(mods.First().Names.International, ShouldNotBeEmpty)
	})

	Convey("Fetching multiple series", t, func() {
		seriesList, err := ManySeries(nil, nil, &Cursor{0, 1}, NoEmbeds)
		So(err, ShouldBeNil)
		So(seriesList.Pagination.Offset, ShouldEqual, 0)
		So(seriesList.Pagination.Max, ShouldEqual, 1)

		num := 0

		// read a few pages, 7 is arbitrary
		seriesList.Walk(func(s *Series) bool {
			So(s.ID, ShouldNotBeBlank)

			num++
			return num < 7
		})

		Convey("test the SeriesFilter", func() {
			// check abbrevitation
			filter := SeriesFilter{Abbreviation: "gta"}

			seriesList, err := ManySeries(&filter, nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 1)

			// check name
			filter = SeriesFilter{Name: "mario"}
			cursor := Cursor{Max: 5}

			seriesList, err = ManySeries(&filter, nil, &cursor, NoEmbeds)
			So(err, ShouldBeNil)
			So(seriesList.Data, ShouldHaveLength, 5)

			// check moderator
			filter = SeriesFilter{Moderator: "r5j52gjv"}

			seriesList, err = ManySeries(&filter, nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(len(seriesList.Data), ShouldBeBetween, 2, 5) // Sorry Josh, but I don't assume it's gonna be more than 5 #Kappa
		})
	})

	Convey("Fetching games of a series", t, func() {
		series, err := SeriesByID(gta, NoEmbeds)
		So(err, ShouldBeNil)

		games, err := series.Games(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.Data[0].ID
		})

		games, err = series.Games(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.Data[0].ID)
		})
	})
}
