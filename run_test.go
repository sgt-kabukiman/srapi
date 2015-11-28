// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRuns(t *testing.T) {
	countRequests = true

	destinyWR := "dy4285nm"

	Convey("Fetching runs by valid IDs", t, func() {
		run, err := RunByID(destinyWR, NoEmbeds)

		So(err, ShouldBeNil)
		So(run.ID, ShouldEqual, destinyWR)
		So(run.Weblink, ShouldNotBeEmpty)
		So(run.Videos.Links, ShouldNotBeEmpty)
		So(run.Comment, ShouldBeEmpty)
		So(run.Status.Status, ShouldEqual, "verified")
		So(run.Status.Examiner, ShouldEqual, "e8el4pj0")
		So(run.Status.VerifyDate.Format("2006-01-02T15:04:05Z"), ShouldEqual, "2015-08-23T02:33:17Z")
		So(run.Date.Format("2006-01-02"), ShouldEqual, "2015-07-23")
		So(run.Submitted.Format("2006-01-02T15:04:05Z"), ShouldEqual, "2015-08-23T02:13:18Z")
		So(run.Times.Primary.String(), ShouldEqual, "5m55s")
		So(run.Times.Primary.Format(), ShouldEqual, "5:55")
		So(run.Times.Realtime.String(), ShouldEqual, "5m55s")
		So(run.System.Platform, ShouldEqual, "lk3gl4jd")
		So(run.System.Emulated, ShouldBeFalse)
		So(run.System.Region, ShouldBeEmpty)
		So(run.Values, ShouldNotBeEmpty)
		So(run.Values, ShouldContainKey, "5lyjpkl4")
		So(run.Links, ShouldNotBeEmpty)

		Convey("Check a run with a comment", func() {
			run, _ := RunByID("wzp1d7rz", NoEmbeds)
			So(run.Comment, ShouldNotBeEmpty)
		})

		Convey("Check a run with a region", func() {
			run, _ := RunByID("6yj1pwoy", NoEmbeds)
			So(run.System.Region, ShouldEqual, "pr184lqn")
			So(run.Times.IngameTime.String(), ShouldEqual, "5m15s")
			So(run.Times.IngameTime.Format(), ShouldEqual, "5:15")
		})

		Convey("Check a run with splits", func() {
			run, _ := RunByID("x7z0ooz5", NoEmbeds)
			So(run.Splits, ShouldNotBeNil)
			So(run.Splits.URI, ShouldNotBeEmpty)
		})

		Convey("Check a run with milliseconds", func() {
			run, _ := RunByID("dy43g2zl", NoEmbeds)
			So(run.Times.Primary.Seconds(), ShouldAlmostEqual, 2164.890, 0.001)
			So(run.Times.Primary.Format(), ShouldEqual, "36:04.890")
			So(run.Times.Realtime.Seconds(), ShouldAlmostEqual, 2164.890, 0.001)
			So(run.Times.Realtime.Format(), ShouldEqual, "36:04.890")
			So(run.Times.RealtimeWithoutLoads.Seconds(), ShouldAlmostEqual, 1492.240, 0.001)
			So(run.Times.RealtimeWithoutLoads.Format(), ShouldEqual, "24:52.240")
			So(run.Times.IngameTime.Seconds(), ShouldAlmostEqual, 1492, 0.001)
			So(run.Times.IngameTime.Format(), ShouldEqual, "24:52")
		})
	})

	Convey("Fetching runs by invalid IDs", t, func() {
		run, err := RunByID("i_do_not_exist", NoEmbeds)
		So(err, ShouldNotBeNil)
		So(run, ShouldBeNil)
	})

	Convey("Fetching multiple runs", t, func() {
		Convey("starting from the beginning", func() {
			runs, err := Runs(nil, nil, nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)

			run := runs.Data[0]
			So(run.ID, ShouldNotBeBlank)
			So(run.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			runs, err := Runs(nil, nil, &Cursor{2, 0}, NoEmbeds)
			So(err, ShouldBeNil)
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 2)
			So(runs.Pagination.Links, ShouldNotBeEmpty)

			run := runs.Data[0]
			So(run.ID, ShouldNotBeBlank)
			So(run.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			runs, err := Runs(nil, nil, &Cursor{0, 3}, NoEmbeds)
			So(err, ShouldBeNil)
			So(runs.Data, ShouldHaveLength, 3)
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(runs.Pagination.Max, ShouldEqual, 3)
			So(runs.Pagination.Links, ShouldNotBeEmpty)

			run := runs.Data[0]
			So(run.ID, ShouldNotBeBlank)
			So(run.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the runs", func() {
			runs, err := Runs(nil, nil, &Cursor{0, 1}, NoEmbeds)
			So(err, ShouldBeNil)
			So(runs.Data, ShouldHaveLength, 1)
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(runs.Pagination.Max, ShouldEqual, 1)

			runs, err = runs.NextPage()
			So(err, ShouldBeNil)
			So(runs.Data, ShouldHaveLength, 1)
			So(runs.Pagination.Offset, ShouldEqual, 1)
			So(runs.Pagination.Max, ShouldEqual, 1)

			runs, err = runs.NextPage()
			So(err, ShouldBeNil)
			So(runs.Data, ShouldHaveLength, 1)
			So(runs.Pagination.Offset, ShouldEqual, 2)
			So(runs.Pagination.Max, ShouldEqual, 1)

			runs, err = runs.PrevPage()
			So(err, ShouldBeNil)
			So(runs.Data, ShouldHaveLength, 1)
			So(runs.Pagination.Offset, ShouldEqual, 1)
			So(runs.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			runs, err := Runs(nil, nil, nil, NoEmbeds)

			runs, err = runs.PrevPage()
			So(err, ShouldNotBeNil)
			So(runs, ShouldNotBeNil)
		})
	})

	Convey("Fetching related resources", t, func() {
		Convey("Game", func() {
			Convey("Without embedding", func() {
				run, err := RunByID(destinyWR, NoEmbeds)

				game, err := run.Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game, ShouldNotBeNil)
				So(game.ID, ShouldEqual, "y65r341e")
			})

			Convey("With embedding", func() {
				run, err := RunByID(destinyWR, "game")

				before := requestCount
				game, err := run.Game(NoEmbeds)
				So(err, ShouldBeNil)
				So(game, ShouldNotBeNil)
				So(game.ID, ShouldEqual, "y65r341e")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Category", func() {
			Convey("Without embedding", func() {
				run, err := RunByID(destinyWR, NoEmbeds)

				category, err := run.Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category, ShouldNotBeNil)
				So(category.ID, ShouldEqual, "mkey4926")
			})

			Convey("With embedding", func() {
				run, err := RunByID(destinyWR, "category")

				before := requestCount
				category, err := run.Category(NoEmbeds)
				So(err, ShouldBeNil)
				So(category, ShouldNotBeNil)
				So(category.ID, ShouldEqual, "mkey4926")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Level", func() {
			Convey("Without embedding", func() {
				run, err := RunByID(destinyWR, NoEmbeds)

				level, err := run.Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level, ShouldNotBeNil)
				So(level.ID, ShouldEqual, "ldy5j7w3")
			})

			Convey("With embedding", func() {
				run, err := RunByID(destinyWR, "level")

				before := requestCount
				level, err := run.Level(NoEmbeds)
				So(err, ShouldBeNil)
				So(level, ShouldNotBeNil)
				So(level.ID, ShouldEqual, "ldy5j7w3")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Platform", func() {
			Convey("Without embedding", func() {
				run, err := RunByID(destinyWR, NoEmbeds)

				platform, err := run.Platform()
				So(err, ShouldBeNil)
				So(platform, ShouldNotBeNil)
				So(platform.ID, ShouldEqual, "lk3gl4jd")
			})

			Convey("With embedding", func() {
				run, err := RunByID(destinyWR, "platform")

				before := requestCount
				platform, err := run.Platform()
				So(err, ShouldBeNil)
				So(platform, ShouldNotBeNil)
				So(platform.ID, ShouldEqual, "lk3gl4jd")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Region", func() {
			Convey("Without embedding", func() {
				run, err := RunByID("68m7g4m0", NoEmbeds)

				region, err := run.Region()
				So(err, ShouldBeNil)
				So(region, ShouldNotBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
			})

			Convey("With embedding", func() {
				run, err := RunByID("68m7g4m0", "region")

				before := requestCount
				region, err := run.Region()
				So(err, ShouldBeNil)
				So(region, ShouldNotBeNil)
				So(region.ID, ShouldEqual, "pr184lqn")
				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Players", func() {
			guests := []string{"Ehroar", "Snead"}

			Convey("Without embedding", func() {
				run, err := RunByID(destinyWR, NoEmbeds)

				players, err := run.Players()
				So(err, ShouldBeNil)
				So(players, ShouldHaveLength, 3)

				for _, player := range players {
					if player.Guest != nil {
						So(player.Guest.Name, ShouldBeIn, guests)
					} else {
						So(player.User.ID, ShouldEqual, "y8dnlgj6")
					}
				}
			})

			Convey("With embedding", func() {
				run, err := RunByID(destinyWR, "players")

				before := requestCount
				players, err := run.Players()
				So(err, ShouldBeNil)
				So(players, ShouldHaveLength, 3)

				for _, player := range players {
					if player.Guest != nil {
						So(player.Guest.Name, ShouldBeIn, guests)
					} else {
						So(player.User.ID, ShouldEqual, "y8dnlgj6")
					}
				}

				So(requestCount, ShouldEqual, before)
			})
		})

		Convey("Examiner", func() {
			run, err := RunByID(destinyWR, NoEmbeds)

			examiner, err := run.Examiner()
			So(err, ShouldBeNil)
			So(examiner, ShouldNotBeNil)
			So(examiner.ID, ShouldEqual, "e8el4pj0")
		})
	})
}
