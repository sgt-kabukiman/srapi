// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testStruct struct {
	Foo string
}

func TestUtils(t *testing.T) {
	Convey("Check if we can get the first link from a collection of links.", t, func() {
		g := Guest{}
		g.Links = append(g.Links, Link{"test1", "foo.bar"})
		g.Links = append(g.Links, Link{"test2", "foo.bar"})
		g.Links = append(g.Links, Link{"test3", "foo.bar"})

		link := firstLink(&g, "test1")
		So(link, ShouldNotBeNil)
		So(link.Relation, ShouldEqual, "test1")

		link = firstLink(&g, "test2")
		So(link, ShouldNotBeNil)
		So(link.Relation, ShouldEqual, "test2")

		link = firstLink(&g, "test3")
		So(link, ShouldNotBeNil)
		So(link.Relation, ShouldEqual, "test3")

		link = firstLink(&g, "test4")
		So(link, ShouldBeNil)
	})

	Convey("Test recast()", t, func() {
		tmp := testStruct{}
		So(tmp.Foo, ShouldEqual, "")

		src := testStruct{"bar"}
		err := recast(src, &tmp)
		So(err, ShouldBeNil)
		So(tmp.Foo, ShouldEqual, "bar")
	})
}
