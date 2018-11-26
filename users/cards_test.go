package users

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCardsHATEOAS(t *testing.T) {
	Convey("Given a new card", t, func() {
		domain = "example.com"
		c := Card{ID: "test"}

		Convey("When adding links", func() {
			c.AddLinks()
			h := Href{"http://example.com/cards/test"}

			Convey("Then link should equal thet test link", func() {
				So(c.Links["card"], ShouldResemble, h)
			})

		})

	})
}

func TestMaskCCs(t *testing.T) {
	Convey("Given a new card", t, func() {
		n := "1234567812345678"
		c := Card{LongNum: n}

		Convey("When the card is masked", func() {
			c.MaskCC()
			m := "************5678"

			Convey("Then the card number should be masked", func() {
				So(c.LongNum, ShouldEqual, m)
			})
		})

	})
}
