// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGuests(t *testing.T) {
	countRequests = true

	Convey("Fetching valid guest names should succeed.", t, func() {
		name := "SgtRockworth"

		guest, err := GuestByName(name)

		So(err, ShouldBeNil)
		So(guest.Name, ShouldEqual, name)
		So(guest.Links, ShouldNotBeEmpty)
	})

	Convey("Fetching unknown guests should fail.", t, func() {
		guest, err := GuestByName("i9r29fgwiurh2igtrw89fw7f")
		So(err, ShouldNotBeNil)
		So(guest, ShouldBeNil)
	})

	Convey("Each guest should have a non-empty collection of runs.", t, func() {
		guest, err := GuestByName("SgtRockworth")
		So(err, ShouldBeNil)

		runs, err := guest.Runs(nil, nil, NoEmbeds)
		So(err, ShouldBeNil)
		So(runs, ShouldNotBeEmpty)
	})
}
