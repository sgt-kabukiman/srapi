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
		users, err := Users(nil, nil, &Cursor{0, 1})
		So(err, ShouldBeNil)
		So(users.Pagination.Offset, ShouldEqual, 0)
		So(users.Pagination.Max, ShouldEqual, 1)

		num := 0

		// read a few pages, 7 is arbitrary
		users.Walk(func(u *User) bool {
			So(u.ID, ShouldNotBeBlank)

			num++
			return num < 7
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
			So(runs.Pagination.Offset, ShouldEqual, 0)

			firstID = runs.First().ID
		})

		runs, err = user.Runs(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(runs.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, runs.First().ID)
		})
	})

	Convey("Fetching the games a user moderates", t, func() {
		user, err := UserByID(odyssic)
		So(err, ShouldBeNil)

		games, err := user.ModeratedGames(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)

		firstID := ""

		Convey("first page of games should be fine", func() {
			So(games.Pagination.Offset, ShouldEqual, 0)

			firstID = games.First().ID
		})

		games, err = user.ModeratedGames(nil, &Sorting{Direction: Descending}, NoEmbeds)
		So(err, ShouldBeNil)

		Convey("sorting order should be taken into account", func() {
			So(games.Pagination.Offset, ShouldEqual, 0)
			So(firstID, ShouldNotEqual, games.First().ID)
		})
	})

	Convey("Fetching PBs of a user", t, func() {
		user, err := UserByID(pac)
		So(err, ShouldBeNil)

		Convey("unfiltered", func() {
			pbs, err := user.PersonalBests(nil, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs.First().Rank, ShouldBeGreaterThanOrEqualTo, 1)
			So(pbs.First().Run.ID, ShouldNotBeEmpty)
		})

		Convey("only first place", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Top: 1}, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs.Size(false), ShouldBeLessThan, 5)
		})

		Convey("only in one series", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Top: 1, Series: "049rqr4v"}, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs.Size(false), ShouldEqual, 1)
		})

		Convey("only in one game", func() {
			pbs, err := user.PersonalBests(&PersonalBestFilter{Game: "om1m3625"}, NoEmbeds)
			So(err, ShouldBeNil)
			So(pbs.Size(false), ShouldEqual, 1)
		})
	})
}
