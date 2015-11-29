// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPersonalBests(t *testing.T) {
	countRequests = true

	pac, _ := UserByID("wzx7q875")

	Convey("Test fetching related resources", t, func() {
		Convey("Game", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				game, err := pbs.First().Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game.ID, ShouldEqual, "om1m3625")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "game")
				So(err, ShouldBeNil)

				before := requestCount
				game, err := pbs.First().Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game.ID, ShouldEqual, "om1m3625")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Category", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				category, err := pbs.First().Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category.ID, ShouldEqual, "w20p0zkn")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "category")
				So(err, ShouldBeNil)

				before := requestCount
				category, err := pbs.First().Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category.ID, ShouldEqual, "w20p0zkn")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Level", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				level, err := pbs.Get(1).Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level.ID, ShouldEqual, "krdn5dm2")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "level")
				So(err, ShouldBeNil)

				before := requestCount
				level, err := pbs.Get(1).Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level.ID, ShouldEqual, "krdn5dm2")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Platform", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				platform, err := pbs.First().Platform()
				So(err, ShouldBeNil)
				So(platform.ID, ShouldEqual, "rdjq4vwe")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "platform")
				So(err, ShouldBeNil)

				before := requestCount
				platform, err := pbs.First().Platform()
				So(err, ShouldBeNil)
				So(platform.ID, ShouldEqual, "rdjq4vwe")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Region", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				region, err := pbs.First().Region()
				So(err, ShouldBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "region")
				So(err, ShouldBeNil)

				before := requestCount
				region, err := pbs.First().Region()
				So(err, ShouldBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Players", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)

				players, err := pbs.First().Players()
				So(err, ShouldBeNil)
				So(players.Size(), ShouldEqual, 1)
				So(players.First().User.ID, ShouldEqual, "wzx7q875")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "players")
				So(err, ShouldBeNil)

				before := requestCount
				players, err := pbs.First().Players()
				So(err, ShouldBeNil)
				So(players.Size(), ShouldEqual, 1)
				So(players.First().User.ID, ShouldEqual, "wzx7q875")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Examiner", func() {
			pbs, err := pac.PersonalBests(nil, NoEmbeds)
			So(err, ShouldBeNil)

			examiner, err := pbs.Get(3).Examiner()
			So(err, ShouldBeNil)
			So(examiner, ShouldNotBeNil)
			So(examiner.ID, ShouldEqual, "y8d4yl86")
		})
	})
}
