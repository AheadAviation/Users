package users

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddressesHATEOAS(t *testing.T) {

	Convey("Given a new address", t, func() {
		domain = "example.com"
		a := Address{ID: "test"}

		Convey("When adding links", func() {
			a.AddLinks()
			h := Href{"http://example.com/addresses/test"}

			Convey("Then link should equal thet test link", func() {
				So(a.Links["address"], ShouldResemble, h)
			})

		})

	})
}
