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
		So(game.Created.Format("2006-01-02T15:04:05Z"), ShouldEqual, "2014-12-07T12:50:20Z")
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
			So(games.Pagination.Offset, ShouldEqual, 0)

			game := games.First()
			So(game, ShouldNotBeNil)
			So(game.ID, ShouldNotBeBlank)
			So(game.Names.International, ShouldNotBeBlank)
			So(game.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			games, err := Games(nil, nil, &Cursor{2, 0}, NoEmbeds)
			So(err, ShouldBeNil)
			So(games.Pagination.Offset, ShouldEqual, 2)
			So(games.Pagination.Links, ShouldNotBeEmpty)

			game := games.First()
			So(game, ShouldNotBeNil)
			So(game.ID, ShouldNotBeBlank)
			So(game.Names.International, ShouldNotBeBlank)
			So(game.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			games, err := Games(nil, nil, &Cursor{0, 3}, NoEmbeds)
			So(err, ShouldBeNil)
			So(len(games.Data), ShouldEqual, 3)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(games.Pagination.Max, ShouldEqual, 3)
			So(games.Pagination.Links, ShouldNotBeEmpty)

			game := games.First()
			So(game, ShouldNotBeNil)
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

			num := 0

			// read a few pages, 7 is arbitrary
			games.Walk(func(g *Game) bool {
				So(g.ID, ShouldNotBeBlank)

				num++
				return num < 7
			})
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
				So(platforms.Data, ShouldHaveLength, 2)
				So(platforms.Data[0].ID, ShouldEqual, "1rjz039w")
				So(platforms.Data[1].ID, ShouldEqual, "4nv59gjk")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "platforms")

				before := requestCount
				platforms, err := game.Platforms()
				So(err, ShouldBeNil)
				So(platforms.Data, ShouldHaveLength, 2)
				So(platforms.Data[0].ID, ShouldEqual, "1rjz039w")
				So(platforms.Data[1].ID, ShouldEqual, "4nv59gjk")
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
				So(regions.Data, ShouldHaveLength, 4)
				So(regions.Data[0].ID, ShouldEqual, "pr184lqn")
				So(regions.Data[1].ID, ShouldEqual, "e6lxy1dz")
				So(regions.Data[2].ID, ShouldEqual, "o316x197")
				So(regions.Data[3].ID, ShouldEqual, "p2g50lnk")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "regions")

				before := requestCount
				regions, err := game.Regions()
				So(err, ShouldBeNil)
				So(regions.Data, ShouldHaveLength, 4)
				So(regions.Data[0].ID, ShouldEqual, "pr184lqn")
				So(regions.Data[1].ID, ShouldEqual, "e6lxy1dz")
				So(regions.Data[2].ID, ShouldEqual, "o316x197")
				So(regions.Data[3].ID, ShouldEqual, "p2g50lnk")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Categories", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				categories, err := game.Categories(nil, nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(categories.Data, ShouldHaveLength, 22)
				So(categories.Data[0].ID, ShouldEqual, "n2y3r8do")
				So(categories.Data[1].ID, ShouldEqual, "7kjqlxd3")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "categories")

				before := requestCount
				categories, err := game.Categories(nil, nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(categories.Data, ShouldHaveLength, 22)
				So(categories.Data[0].ID, ShouldEqual, "n2y3r8do")
				So(categories.Data[1].ID, ShouldEqual, "7kjqlxd3")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Levels", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				levels, err := game.Levels(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(levels.Data, ShouldHaveLength, 14)
				So(levels.Data[0].ID, ShouldEqual, "xd4e80wm")
				So(levels.Data[1].ID, ShouldEqual, "nwlzepdv")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "levels")

				before := requestCount
				levels, err := game.Levels(nil, NoEmbeds)
				So(err, ShouldBeNil)
				So(levels.Data, ShouldHaveLength, 14)
				So(levels.Data[0].ID, ShouldEqual, "xd4e80wm")
				So(levels.Data[1].ID, ShouldEqual, "nwlzepdv")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Variables", func() {
			Convey("Structs", func() {
				game, err := GameByID(superMarioSunshine, NoEmbeds)

				variables, err := game.Variables(nil)
				So(err, ShouldBeNil)
				So(variables.Data, ShouldHaveLength, 2)
				So(variables.Data[0].ID, ShouldEqual, "38dz6zn0")
				So(variables.Data[1].ID, ShouldEqual, "r8r157ne")
			})

			Convey("Structs with embedding", func() {
				game, err := GameByID(superMarioSunshine, "variables")

				before := requestCount
				variables, err := game.Variables(nil)
				So(err, ShouldBeNil)
				So(variables.Data, ShouldHaveLength, 2)
				So(variables.Data[0].ID, ShouldEqual, "38dz6zn0")
				So(variables.Data[1].ID, ShouldEqual, "r8r157ne")
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
				So(mods.Data, ShouldHaveLength, 3)
				So(mods.Data[0].ID, ShouldBeIn, modIDs)
				So(mods.Data[1].ID, ShouldBeIn, modIDs)
				So(mods.Data[2].ID, ShouldBeIn, modIDs)
			})

			Convey("Users with embedding", func() {
				game, err := GameByID(gtavc, "moderators")

				before := requestCount
				mods, err := game.Moderators()
				So(err, ShouldBeNil)
				So(mods.Data, ShouldHaveLength, 3)
				So(mods.Data[0].ID, ShouldBeIn, modIDs)
				So(mods.Data[1].ID, ShouldBeIn, modIDs)
				So(mods.Data[2].ID, ShouldBeIn, modIDs)
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

	Convey("Generic test for all collections", t, func() {
		series, err := SeriesByAbbreviation("gta", NoEmbeds)
		So(err, ShouldBeNil)
		So(series, ShouldNotBeNil)

		Convey("get the first item", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)

			first := games.First()
			So(first, ShouldNotBeNil)
			So(first.ID, ShouldNotBeBlank)
		})

		Convey("iterate over the first few elements", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)
			So(len(games.Data), ShouldBeGreaterThanOrEqualTo, 3)

			idx := 0
			ids := []string{
				games.Data[0].ID,
				games.Data[1].ID,
				games.Data[2].ID,
			}

			iterator := games.Iterator()

			for run := range iterator.Output() {
				So(run.ID, ShouldEqual, ids[idx])
				idx++

				if idx == 3 {
					iterator.Stop()
				}
			}

			So(idx, ShouldEqual, 3)
		})

		Convey("walk over the first few elements", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)
			So(len(games.Data), ShouldBeGreaterThanOrEqualTo, 3)

			idx := 0
			ids := []string{
				games.Data[0].ID,
				games.Data[1].ID,
				games.Data[2].ID,
			}

			games.Walk(func(g *Game) bool {
				So(g.ID, ShouldEqual, ids[idx])
				idx++

				return idx < 3
			})

			So(idx, ShouldEqual, 3)
		})

		Convey("scan for an ID", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)

			So(games.ScanForID("m9dowk1p"), ShouldNotBeNil)
			So(games.ScanForID("yo1yv1q5"), ShouldNotBeNil)
			So(games.ScanForID("jy657deo"), ShouldNotBeNil)
		})

		Convey("return all games in the collection", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)

			structs := games.Games()
			expectedLen := len(games.Data)

			So(len(structs), ShouldEqual, expectedLen)
			So(structs[0].ID, ShouldEqual, games.Data[0].ID)
		})

		Convey("limit collection", func() {
			games, err := series.Games(nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(games, ShouldNotBeNil)

			games = games.Limit(3)

			So(games.Size(true), ShouldEqual, 3)
			So(games.Size(false), ShouldEqual, 3)
			So(len(games.Games()), ShouldEqual, 3)
			So(games.ScanForID("ok6qvxdg"), ShouldBeNil)

			idx := 0

			games.Walk(func(g *Game) bool {
				idx++
				return true
			})

			So(idx, ShouldEqual, 3)
		})
	})
}
