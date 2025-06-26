package racket

import (
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ProgressLogger(t *testing.T) {
	defer leaktest.Check(t)()

	const ProgressCrap ProgressType = 1024
	disco := log.New(io.Discard, "", 0)
	pchan := make(chan Progress)
	bchan := make(chan Progress)
	defer close(pchan) // closing this is the signal for ProgressLogger to peace out.
	defer close(bchan)

	Convey("When a ProgressLogger is set up, it behahves as expected.", t, func() {
		errorCount := 0
		errf := func(e error) {
			errorCount++
		}
		go ProgressLogger(disco, true, errf, pchan, bchan)

		// The easy
		pchan <- PMessagef("Hello")
		pchan <- PErrorf("Error!")

		// Make sure the bar is notified
		pchan <- PEstimate(42)
		So(<-bchan, ShouldEqual, PEstimate(42))

		// Make sure the bar is notified
		pchan <- PUpdate(-1)
		So(<-bchan, ShouldEqual, PUpdate(-1))

		// Make sure weird stuff doesn't blow up
		pchan <- Progress{
			Type: ProgressCrap,
			Data: "CRAP!",
		}

		// Make sure errorCount was eventually incremented
		So(errorCount, ShouldEqual, 1)
	})

}

func Test_ProgressType(t *testing.T) {
	Convey("Undefined ProgressTypes behave and resolve properly", t, func() {
		const ProgressCrap ProgressType = 1024
		pe := Progress{
			Type: ProgressCrap,
			Data: "CRAP!",
		}
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressCrap)
		So(pe.Type.String(), ShouldEqual, "")
		So(pe.Data, ShouldHaveSameTypeAs, "Hello World!")
		So(pe.Error(), ShouldBeNil)
		So(pe.String(), ShouldEqual, ": CRAP!")
	})

	Convey("ProgressError and shortcuts, behave and resolve properly", t, func() {
		pe := PErrorf("an ERROR")
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressError)
		So(pe.Type.String(), ShouldEqual, "ProgressError")
		So(pe.Data, ShouldBeError)
		So(pe.Error(), ShouldEqual, fmt.Errorf("an ERROR"))
		So(pe.String(), ShouldEqual, "ProgressError: an ERROR")
	})

	Convey("ProgressMessage and shortcuts, behave and resolve properly", t, func() {
		pe := PMessagef("MESSAGE!")
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressMessage)
		So(pe.Type.String(), ShouldEqual, "ProgressMessage")
		So(pe.Data, ShouldHaveSameTypeAs, "Hello World")
		So(pe.Error(), ShouldBeNil)
		So(pe.String(), ShouldEqual, "ProgressMessage: MESSAGE!")
	})

	Convey("ProgressUpdate and shortcuts, behave and resolve properly", t, func() {
		pe := PUpdate(42)
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressUpdate)
		So(pe.Type.String(), ShouldEqual, "ProgressUpdate")
		So(pe.Data, ShouldHaveSameTypeAs, int64(1024))
		So(pe.Error(), ShouldBeNil)
		So(pe.String(), ShouldEqual, "ProgressUpdate: 42")
	})

	Convey("ProgressEstimate and shortcuts, behave and resolve properly", t, func() {
		pe := PEstimate(4026)
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressEstimate)
		So(pe.Type.String(), ShouldEqual, "ProgressEstimate")
		So(pe.Data, ShouldHaveSameTypeAs, int64(1024))
		So(pe.Error(), ShouldBeNil)
		So(pe.String(), ShouldEqual, "ProgressEstimate: 4026")
	})

	Convey("ProgressOther behaves and resolve properly", t, func() {
		pe := Progress{
			Type: ProgressOther,
			Data: io.Discard,
		}
		So(pe, ShouldHaveSameTypeAs, Progress{})
		So(pe.Type, ShouldEqual, ProgressOther)
		So(pe.Type.String(), ShouldEqual, "ProgressOther")
		So(pe.Data, ShouldHaveSameTypeAs, io.Discard)
		So(pe.Error(), ShouldBeNil)
		So(pe.String(), ShouldEqual, "ProgressOther: {}")
	})
}
