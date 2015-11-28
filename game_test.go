// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGames(t *testing.T) {
	countRequests = true

	superMarioSunshine := "v1pxjz68"
	gtavc := "29d30dlp"

	Convey("Fetching games by valid IDs", t, func() {
		game, err := GameByID(superMarioSunshine, NoEmbeds)

		So(err, ShouldBeNil)
		So(game.ID, ShouldEqual, superMarioSunshine)
		So(game.Names.International, ShouldEqual, "Super Mario Sunshine")
		So(game.Names.Japanese, ShouldEqual, "スーパーマリオサンシャイン")
		So(game.Abbreviation, ShouldEqual, "sms")
		So(game.Weblink, ShouldNotBeEmpty)
		So(game.Released, ShouldEqual, 2002)
		So(game.Ruleset.ShowMilliseconds, ShouldBeFalse)
		So(game.Ruleset.RequireVerification, ShouldBeTrue)
		So(game.Ruleset.RequireVideo, ShouldBeFalse)
		So(game.Ruleset.RunTimes, ShouldHaveLength, 2)
		So(game.Ruleset.DefaultTime, ShouldEqual, TimingRealtime)
		So(game.Ruleset.EmulatorsAllowed, ShouldBeTrue)
		So(game.Romhack, ShouldBeFalse)
		So(game.Created, ShouldNotBeEmpty)
		So(game.Assets, ShouldNotBeEmpty)
		So(game.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching games by invalid IDs", t, func() {
		game, err := GameByID("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(game, ShouldBeNil)
	})

	Convey("Fetching games by abbreviation", t, func() {
		game, err := GameByAbbreviation("gtavc", NoEmbeds)
		So(err, ShouldBeNil)
		So(game, ShouldNotBeNil)
		So(game.Abbreviation, ShouldEqual, "gtavc")
	})

	Convey("Fetching multiple games", t, func() {
		Convey("starting from the beginning", func() {
			games, err := Games(nil, nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			game := games.Data[0]
			So(game.ID, ShouldNotBeBlank)
			So(game.Names.International, ShouldNotBeBlank)
			So(game.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			games, err := Games(nil, nil, &Cursor{2, 0}, NoEmbeds)
			So(err, ShouldBeNil)
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 2)
			So(games.Pagination.Links, ShouldNotBeEmpty)

			game := games.Data[0]
			So(game.ID, ShouldNotBeBlank)
			So(game.Names.International, ShouldNotBeBlank)
			So(game.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			games, err := Games(nil, nil, &Cursor{0, 3}, NoEmbeds)
			So(err, ShouldBeNil)
			So(games.Data, ShouldHaveLength, 3)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(games.Pagination.Max, ShouldEqual, 3)
			So(games.Pagination.Links, ShouldNotBeEmpty)

			game := games.Data[0]
			So(game.ID, ShouldNotBeBlank)
			So(game.Names.International, ShouldNotBeBlank)
			So(game.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the games", func() {
			games, err := Games(nil, nil, &Cursor{0, 1}, NoEmbeds)
			So(err, ShouldBeNil)
			So(games.Data, ShouldHaveLength, 1)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(games.Pagination.Max, ShouldEqual, 1)

			games, err = games.NextPage()
			So(err, ShouldBeNil)
			So(games.Data, ShouldHaveLength, 1)
			So(games.Pagination.Offset, ShouldEqual, 1)
			So(games.Pagination.Max, ShouldEqual, 1)

			games, err = games.NextPage()
			So(err, ShouldBeNil)
			So(games.Data, ShouldHaveLength, 1)
			So(games.Pagination.Offset, ShouldEqual, 2)
			So(games.Pagination.Max, ShouldEqual, 1)

			games, err = games.PrevPage()
			So(err, ShouldBeNil)
			So(games.Data, ShouldHaveLength, 1)
			So(games.Pagination.Offset, ShouldEqual, 1)
			So(games.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			games, err := Games(nil, nil, nil, NoEmbeds)

			games, err = games.PrevPage()
			So(err, ShouldNotBeNil)
			So(games, ShouldNotBeNil)
		})
	})

	Convey("Get the series from a game", t, func() {
		game, err := GameByAbbreviation("gtavc", NoEmbeds)

		series, err := game.Series(NoEmbeds)
		So(err, ShouldBeNil)
		So(series, ShouldNotBeNil)
		So(series.Abbreviation, ShouldEqual, "gta")
	})

	Convey("Fetching related resources", t, func() {
		Convey("Platforms", func() {
			Convey("IDs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				ids, err := game.PlatformIDs()
				So(err, ShouldBeNil)
				So(ids, ShouldHaveLength, 2)
				So(ids[0], ShouldEqual, "1rjz039w")
				So(ids[1], ShouldEqual, "4nv59gjk")
			})

			Convey("IDs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "platforms")

				before := requestCount
				ids, err := game.PlatformIDs()
				So(err, ShouldBeNil)
				So(ids, ShouldHaveLength, 2)
				So(ids[0], ShouldEqual, "1rjz039w")
				So(ids[1], ShouldEqual, "4nv59gjk")
				So(requestCount, ShouldEqual, before)
			})

			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				platforms, err := game.Platforms()
				So(err, ShouldBeNil)
				So(platforms, ShouldHaveLength, 2)
				So(platforms[0].ID, ShouldEqual, "1rjz039w")
				So(platforms[1].ID, ShouldEqual, "4nv59gjk")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "platforms")

				before := requestCount
				platforms, err := game.Platforms()
				So(err, ShouldBeNil)
				So(platforms, ShouldHaveLength, 2)
				So(platforms[0].ID, ShouldEqual, "1rjz039w")
				So(platforms[1].ID, ShouldEqual, "4nv59gjk")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Regions", func() {
			Convey("IDs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				ids, err := game.RegionIDs()
				So(err, ShouldBeNil)
				So(ids, ShouldHaveLength, 4)
				So(ids[0], ShouldEqual, "pr184lqn")
				So(ids[1], ShouldEqual, "e6lxy1dz")
				So(ids[2], ShouldEqual, "o316x197")
				So(ids[3], ShouldEqual, "p2g50lnk")
			})

			Convey("IDs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "regions")

				before := requestCount
				ids, err := game.RegionIDs()
				So(err, ShouldBeNil)
				So(ids, ShouldHaveLength, 4)
				So(ids[0], ShouldEqual, "pr184lqn")
				So(ids[1], ShouldEqual, "e6lxy1dz")
				So(ids[2], ShouldEqual, "o316x197")
				So(ids[3], ShouldEqual, "p2g50lnk")
				So(requestCount, ShouldEqual, before)
			})

			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				regions, err := game.Regions()
				So(err, ShouldBeNil)
				So(regions, ShouldHaveLength, 4)
				So(regions[0].ID, ShouldEqual, "pr184lqn")
				So(regions[1].ID, ShouldEqual, "e6lxy1dz")
				So(regions[2].ID, ShouldEqual, "o316x197")
				So(regions[3].ID, ShouldEqual, "p2g50lnk")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "regions")

				before := requestCount
				regions, err := game.Regions()
				So(err, ShouldBeNil)
				So(regions, ShouldHaveLength, 4)
				So(regions[0].ID, ShouldEqual, "pr184lqn")
				So(regions[1].ID, ShouldEqual, "e6lxy1dz")
				So(regions[2].ID, ShouldEqual, "o316x197")
				So(regions[3].ID, ShouldEqual, "p2g50lnk")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Categories", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				categories, err := game.Categories(nil, nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(categories, ShouldHaveLength, 22)
				So(categories[0].ID, ShouldEqual, "n2y3r8do")
				So(categories[1].ID, ShouldEqual, "7kjqlxd3")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "categories")

				before := requestCount
				categories, err := game.Categories(nil, nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(categories, ShouldHaveLength, 22)
				So(categories[0].ID, ShouldEqual, "n2y3r8do")
				So(categories[1].ID, ShouldEqual, "7kjqlxd3")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Levels", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				levels, err := game.Levels(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(levels, ShouldHaveLength, 14)
				So(levels[0].ID, ShouldEqual, "xd4e80wm")
				So(levels[1].ID, ShouldEqual, "nwlzepdv")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "levels")

				before := requestCount
				levels, err := game.Levels(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(levels, ShouldHaveLength, 14)
				So(levels[0].ID, ShouldEqual, "xd4e80wm")
				So(levels[1].ID, ShouldEqual, "nwlzepdv")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Variables", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				variables, err := game.Variables(nil)
				So(err, ShouldBeNil)
				So(variables, ShouldHaveLength, 2)
				So(variables[0].ID, ShouldEqual, "38dz6zn0")
				So(variables[1].ID, ShouldEqual, "r8r157ne")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "variables")

				before := requestCount
				variables, err := game.Variables(nil)
				So(err, ShouldBeNil)
				So(variables, ShouldHaveLength, 2)
				So(variables[0].ID, ShouldEqual, "38dz6zn0")
				So(variables[1].ID, ShouldEqual, "r8r157ne")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Moderators", func() {
			modIDs := []string{"vqxkmj07", "3qjn18m1", "gpj064jw"}

			Convey("Map", func() {
				game, _ := GameByID(gtavc, NoEmbeds)

				mods := game.ModeratorMap()
				So(mods, ShouldHaveLength, 3)
				So(mods, ShouldContainKey, "vqxkmj07")
				So(mods, ShouldContainKey, "3qjn18m1")
				So(mods, ShouldContainKey, "gpj064jw")
				So(mods["vqxkmj07"], ShouldEqual, NormalModerator)
				So(mods["3qjn18m1"], ShouldEqual, SuperModerator)
				So(mods["gpj064jw"], ShouldEqual, NormalModerator)
			})

			Convey("Map with embedding", func() {
				game, _ := GameByID(gtavc, "moderators")

				before := requestCount
				mods := game.ModeratorMap()
				So(mods, ShouldHaveLength, 3)
				So(mods, ShouldContainKey, "vqxkmj07")
				So(mods, ShouldContainKey, "3qjn18m1")
				So(mods, ShouldContainKey, "gpj064jw")
				So(mods["vqxkmj07"], ShouldEqual, UnknownModLevel)
				So(mods["3qjn18m1"], ShouldEqual, UnknownModLevel)
				So(mods["gpj064jw"], ShouldEqual, UnknownModLevel)
				So(requestCount, ShouldEqual, before)
			})

			Convey("Users", func() {
				game, err := GameByID(gtavc, NoEmbeds)

				mods, err := game.Moderators()
				So(err, ShouldBeNil)
				So(mods, ShouldHaveLength, 3)
				So(mods[0].ID, ShouldBeIn, modIDs)
				So(mods[1].ID, ShouldBeIn, modIDs)
				So(mods[2].ID, ShouldBeIn, modIDs)
			})

			Convey("Users with embedding", func() {
				game, err := GameByID(gtavc, "moderators")

				before := requestCount
				mods, err := game.Moderators()
				So(err, ShouldBeNil)
				So(mods, ShouldHaveLength, 3)
				So(mods[0].ID, ShouldBeIn, modIDs)
				So(mods[1].ID, ShouldBeIn, modIDs)
				So(mods[2].ID, ShouldBeIn, modIDs)
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Romhacks", func() {
			game, err := GameByID(gtavc, NoEmbeds)

			romhacks, err := game.Romhacks(NoEmbeds)
			So(err, ShouldBeNil)
			So(romhacks.Data, ShouldHaveLength, 1)
			So(romhacks.Data[0].ID, ShouldEqual, "4pdv9k1w")
			So(romhacks.Data[0].Abbreviation, ShouldEqual, "gtavc_chaos")
		})

		Convey("Primary leaderboard", func() {
			game, err := GameByID(gtavc, NoEmbeds)

			leaderboard, err := game.PrimaryLeaderboard(&LeaderboardOptions{Top: 5}, NoEmbeds)
			So(err, ShouldBeNil)
			So(leaderboard, ShouldNotBeNil)
			So(leaderboard.Runs, ShouldHaveLength, 5)
		})

		Convey("Records", func() {
			game, err := GameByID(gtavc, NoEmbeds)

			leaderboards, err := game.Records(nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(leaderboards, ShouldNotBeNil)
			So(leaderboards.Data, ShouldHaveLength, 8)
			So(leaderboards.Data[0].Runs, ShouldNotBeEmpty)
		})

		Convey("Runs", func() {
			game, err := GameByID(gtavc, NoEmbeds)

			runs, err := game.Runs(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(runs, ShouldNotBeNil)
			So(runs.Data, ShouldHaveLength, 20)
		})
	})
}
