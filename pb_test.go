// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPersonalBests(t *testing.T) {
	pac, _ := UserByID("wzx7q875")

	Convey("Test fetching related resources", t, func() {
		Convey("Game", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				game, err := pbs[0].Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game.ID, ShouldEqual, "om1m3625")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "game")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				game, err := pbs[0].Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game.ID, ShouldEqual, "om1m3625")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Category", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				category, err := pbs[0].Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category.ID, ShouldEqual, "w20p0zkn")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "category")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				category, err := pbs[0].Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category.ID, ShouldEqual, "w20p0zkn")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Level", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				level, err := pbs[1].Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level.ID, ShouldEqual, "krdn5dm2")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "level")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				level, err := pbs[1].Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level.ID, ShouldEqual, "krdn5dm2")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Platform", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				platform, err := pbs[0].Platform()
				So(err, ShouldBeNil)
				So(platform.ID, ShouldEqual, "rdjq4vwe")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "platform")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				platform, err := pbs[0].Platform()
				So(err, ShouldBeNil)
				So(platform.ID, ShouldEqual, "rdjq4vwe")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Region", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				region, err := pbs[0].Region()
				So(err, ShouldBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "region")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				region, err := pbs[0].Region()
				So(err, ShouldBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Players", func() {
			Convey("Without embedding", func() {
				pbs, err := pac.PersonalBests(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				players, err := pbs[0].Players()
				So(err, ShouldBeNil)
				So(players, ShouldHaveLength, 1)
				So(players[0].User.ID, ShouldEqual, "wzx7q875")
			})

			Convey("With embedding", func() {
				pbs, err := pac.PersonalBests(nil, "players")
				So(err, ShouldBeNil)
				So(pbs, ShouldNotBeEmpty)

				before := requestCount
				players, err := pbs[0].Players()
				So(err, ShouldBeNil)
				So(players, ShouldHaveLength, 1)
				So(players[0].User.ID, ShouldEqual, "wzx7q875")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Examiner", func() {
			pbs, err := pac.PersonalBests(nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs, ShouldNotBeEmpty)

			examiner, err := pbs[3].Examiner()
			So(err, ShouldBeNil)
			So(examiner, ShouldNotBeNil)
			So(examiner.ID, ShouldEqual, "y8d4yl86")
		})
	})
}
