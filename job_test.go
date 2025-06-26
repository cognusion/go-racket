package racket

import (
	"io"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	var (
		wCount atomic.Int64                // atomic counter
		stdOut = log.New(os.Stdout, "", 0) //
	)

	workerFunc := func(id any, work Work, pchan chan<- Progress) {
		pchan <- PMessagef("I am %v! The work number is %d!\n", id, work.GetInt("the number"))
		wCount.Add(1)
	}

	j := NewJob(workerFunc)
	wchan := make(chan Work)
	pchan, done := j.Supervisor(2, wchan)
	defer close(pchan) // if we don't close pchan, the ProcessLogger never exits cleanly.

	// Spin up a ProgressLogger using our stdOut logger, logging messages,
	// not especially handling errors, reading from pchan, not using a progress bar
	go ProgressLogger(stdOut, true, nil, pchan, nil)

	for i := range 100 {
		wchan <- NewWork(map[string]any{
			"the number": i,
		})
	}
	done() // signal the Supervisor and any idle workers that we're done giving out jobs.

	// wait until all outstanding Work is accomplished.
	<-j.IsDone()

}

func Test_Job(t *testing.T) {
	defer leaktest.Check(t)()

	disco := log.New(io.Discard, "", 0)
	its := 100

	Convey("When a Job is created, and Work is assigned, everything works.", t, func(c C) {
		var wCount atomic.Int64

		wf := func(id any, work Work, pchan chan<- Progress) {
			pchan <- PMessagef("I am %v!\n", id)
			wCount.Add(1)
		}

		j := NewJob(wf)
		wchan := make(chan Work)
		pchan, done := j.Supervisor(2, wchan)
		defer close(pchan)
		go ProgressLogger(disco, false, nil, pchan, nil)

		for range its {
			wchan <- NewWork(nil)
		}
		done()

		<-j.IsDone()

		c.So(wCount.Load(), ShouldEqual, its)
	})
}
