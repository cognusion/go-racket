package racket

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Work(t *testing.T) {

	Convey("When Work is created, values are accessible as expected", t, func() {
		w := NewWork(map[string]any{
			"Hello":      "World",
			"Truth":      true,
			"The Answer": 42,
		})

		So(w.Get("Does not exist"), ShouldBeNil)
		So(w.GetString("Hello"), ShouldEqual, "World")
		So(w.GetBool("Truth"), ShouldBeTrue)
		So(w.GetInt("The Answer"), ShouldEqual, 42)

		Convey("... When Getters are type mismatched, nothing blows up and the values are casted.", func() {
			So(w.GetBool("Hello"), ShouldBeFalse)
			So(w.GetString("The Answer"), ShouldEqual, "42") // the string of the int
			So(w.GetInt("Truth"), ShouldEqual, 1)            // int true
		})

	})
}
