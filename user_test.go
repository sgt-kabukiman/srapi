// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUsers(t *testing.T) {
	countRequests = true

	pac := "wzx7q875"
	odyssic := "gpj064jw"

	Convey("Fetching users by valid IDs", t, func() {
		user, err := UserByID(pac)

		So(err, ShouldBeNil)
		So(user.ID, ShouldEqual, pac)
		So(user.Names.International, ShouldEqual, "Pac")
		So(user.Names.Japanese, ShouldEqual, "パック")
		So(user.Weblink, ShouldNotBeEmpty)
		So(user.NameStyle.Style, ShouldEqual, "gradient")
		So(user.NameStyle.Color, ShouldBeNil)
		So(user.NameStyle.ColorFrom, ShouldNotBeNil)
		So(user.NameStyle.ColorFrom.Light, ShouldNotBeNil)
		So(user.NameStyle.ColorFrom.Dark, ShouldNotBeNil)
		So(user.NameStyle.ColorTo, ShouldNotBeNil)
		So(user.NameStyle.ColorTo.Light, ShouldNotBeNil)
		So(user.NameStyle.ColorTo.Dark, ShouldNotBeNil)
		So(user.Role, ShouldEqual, "programmer")
		So(user.Signup.Format("2006-01-02T15:04:05Z"), ShouldEqual, "2013-12-09T12:03:01Z")
		So(user.Location.Country.Code, ShouldNotBeEmpty)
		So(user.Twitch, ShouldNotBeNil)
		So(user.Hitbox, ShouldNotBeNil)
		So(user.YouTube, ShouldNotBeNil)
		So(user.Twitter, ShouldNotBeNil)
		So(user.SpeedRunsLive, ShouldNotBeNil)
		So(user.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching users by invalid IDs", t, func() {
		user, err := UserByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(user, ShouldBeNil)
	})

	Convey("Fetching multiple users", t, func() {
		Convey("starting from the beginning", func() {
			users, err := Users(nil, nil, nil)
			So(err, ShouldBeNil)
			So(users.Data, ShouldNotBeEmpty)
			So(users.Pagination.Offset, ShouldEqual, 0)

			user := users.Data[0]
			So(user.ID, ShouldNotBeBlank)
			So(user.Names.International, ShouldNotBeBlank)
			So(user.Links, ShouldNotBeEmpty)
		})

		Convey("skipping the first few", func() {
			users, err := Users(nil, nil, &Cursor{2, 0})
			So(err, ShouldBeNil)
			So(users.Data, ShouldNotBeEmpty)
			So(users.Pagination.Offset, ShouldEqual, 2)
			So(users.Pagination.Links, ShouldNotBeEmpty)

			user := users.Data[0]
			So(user.ID, ShouldNotBeBlank)
			So(user.Names.International, ShouldNotBeBlank)
			So(user.Links, ShouldNotBeEmpty)
		})

		Convey("limited to just a few", func() {
			users, err := Users(nil, nil, &Cursor{0, 3})
			So(err, ShouldBeNil)
			So(users.Data, ShouldHaveLength, 3)
			So(users.Pagination.Offset, ShouldEqual, 0)
			So(users.Pagination.Max, ShouldEqual, 3)
			So(users.Pagination.Links, ShouldNotBeEmpty)

			user := users.Data[0]
			So(user.ID, ShouldNotBeBlank)
			So(user.Names.International, ShouldNotBeBlank)
			So(user.Links, ShouldNotBeEmpty)
		})

		Convey("paging through the users", func() {
			users, err := Users(nil, nil, &Cursor{0, 1})
			So(err, ShouldBeNil)
			So(users.Data, ShouldHaveLength, 1)
			So(users.Pagination.Offset, ShouldEqual, 0)
			So(users.Pagination.Max, ShouldEqual, 1)

			users, err = users.NextPage()
			So(err, ShouldBeNil)
			So(users.Data, ShouldHaveLength, 1)
			So(users.Pagination.Offset, ShouldEqual, 1)
			So(users.Pagination.Max, ShouldEqual, 1)

			users, err = users.NextPage()
			So(err, ShouldBeNil)
			So(users.Data, ShouldHaveLength, 1)
			So(users.Pagination.Offset, ShouldEqual, 2)
			So(users.Pagination.Max, ShouldEqual, 1)

			users, err = users.PrevPage()
			So(err, ShouldBeNil)
			So(users.Data, ShouldHaveLength, 1)
			So(users.Pagination.Offset, ShouldEqual, 1)
			So(users.Pagination.Max, ShouldEqual, 1)
		})

		Convey("the prev page from the beginning should yield an error", func() {
			users, err := Users(nil, nil, nil)

			users, err = users.PrevPage()
			So(err, ShouldNotBeNil)
			So(users, ShouldNotBeNil)
		})

		Convey("check the lookup filter", func() {
			users, _ := Users(&UserFilter{Lookup: "Pac"}, nil, nil)
			So(users.Data, ShouldHaveLength, 1)
		})

		Convey("check the name filter", func() {
			users, _ := Users(nil, nil, nil)
			first := users.Data[0].ID

			users, _ = Users(&UserFilter{Name: "Pac"}, nil, nil)
			So(users.Data[0].ID, ShouldNotEqual, first)
		})

		Convey("check the twitch filter", func() {
			users, _ := Users(&UserFilter{Twitch: "Pac__"}, nil, nil)
			So(users.Data, ShouldHaveLength, 1)
		})

		Convey("check the hitbox filter", func() {
			users, _ := Users(&UserFilter{Hitbox: "Pac"}, nil, nil)
			So(users.Data, ShouldHaveLength, 1)
		})

		Convey("check the twitter filter", func() {
			users, _ := Users(&UserFilter{Twitter: "pac____"}, nil, nil)
			So(users.Data, ShouldHaveLength, 1)
		})

		Convey("check the speedrunslive filter", func() {
			users, _ := Users(&UserFilter{SpeedRunsLive: "Pac"}, nil, nil)
			So(users.Data, ShouldHaveLength, 1)
		})
	})

	Convey("Fetching runs of a user", t, func() {
		user, err := UserByID(pac)
		So(err, ShouldBeNil)

		runs, err := user.Runs(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of runs should be fine", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)

			firstID = runs.Data[0].ID
		})

		runs, err = user.Runs(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(runs.Data, ShouldNotBeEmpty)
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, runs.Data[0].ID)
		})
	})

	Convey("Fetching the games a user moderates", t, func() {
		user, err := UserByID(odyssic)
		So(err, ShouldBeNil)

		games, err := user.ModeratedGames(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.Data[0].ID
		})

		games, err = user.ModeratedGames(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Data, ShouldNotBeEmpty)
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.Data[0].ID)
		})
	})

	Convey("Fetching PBs of a user", t, func() {
		user, err := UserByID(pac)
		So(err, ShouldBeNil)

		Convey("unfiltered", func() {
			pbs, err := user.PersonalBests(nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs, ShouldNotBeEmpty)
			So(pbs[0].Rank, ShouldBeGreaterThanOrEqualTo, 1)
			So(pbs[0].Run.ID, ShouldNotBeEmpty)
		})

		Convey("only first place", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Top: 1}, NoEmbeds)
			So(err, ShouldBeNil)
			So(len(pbs), ShouldBeLessThan, 5)
		})

		Convey("only in one series", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Top: 1, Series: "049rqr4v"}, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs, ShouldHaveLength, 1)
		})

		Convey("only in one game", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Game: "om1m3625"}, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs, ShouldHaveLength, 1)
		})
	})
}
