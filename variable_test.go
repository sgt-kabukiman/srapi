// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVariables(t *testing.T) {
	seedID := "0789j6nw" // gta vc chaos% seed
	seedVar, err := VariableByID(seedID)

	Convey("Fetching variables by valid IDs", t, func() {
		So(err, ShouldBeNil)
		So(seedVar.ID, ShouldEqual, seedID)
		So(seedVar.Name, ShouldEqual, "Seed")
		So(seedVar.Scope.Type, ShouldEqual, "full-game")
		So(seedVar.Scope.Level, ShouldBeEmpty)
		So(seedVar.Mandatory, ShouldBeTrue)
		So(seedVar.UserDefined, ShouldBeTrue)
		So(seedVar.Obsoletes, ShouldBeTrue)
		So(seedVar.Values.Choices, ShouldNotBeEmpty)
		So(seedVar.Values.Default, ShouldBeEmpty)
		So(seedVar.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching variables by invalid IDs", t, func() {
		variable, err := VariableByID("i_do_not_exist")
		So(err, ShouldNotBeNil)
		So(variable, ShouldBeNil)
	})

	Convey("Fetch the game the variable belongs to", t, func() {
		game, err := seedVar.Game(NoEmbeds)
		So(err, ShouldBeNil)
		So(game, ShouldNotBeNil)
		So(game.Names.International, ShouldEqual, "Grand Theft Auto: Vice City Chaos%")
	})

	Convey("Fetch the category the variable belongs to", t, func() {
		category, err := seedVar.Category(NoEmbeds)
		So(err, ShouldBeNil)
		So(category, ShouldBeNil)

		variable, err := VariableByID("5lyjm9l4") // "Charatcer" in CTR
		So(err, ShouldBeNil)
		So(variable, ShouldNotBeNil)

		category, err = variable.Category(NoEmbeds)
		So(err, ShouldBeNil)
		So(category, ShouldNotBeNil)
		So(category.Name, ShouldEqual, "Any%")
	})
}
