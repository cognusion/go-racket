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
		stdOut = log.New(os.Stdout, "", 0) // a logger to output to the terminal
		wchan  = make(chan Work)           // a channel to put Work into
	)

	// a WorkerFunc to do work. In this case just push a ProgressMessage onto the
	// progress channel, and atomically increment a number.
	workerFunc := func(id any, work Work, pchan chan<- Progress) {
		pchan <- PMessagef("I am %v! The work number is %d!\n", id, work.GetInt("the number"))
		wCount.Add(1)
	}

	// Create a new Job, that will use our workerFunc
	j := NewJob(workerFunc)

	// Hire a Supervisor to oversee 2 workers, who will get work via our Work channel.
	// The Supervisor will return a channel to receive Progress on, and a func to call
	// when we are done putting Work onto the Work channel.
	pchan, done := j.Supervisor(2, wchan)
	defer close(pchan) // if we don't close pchan, the ProcessLogger (below) never exits cleanly.

	// Spin up a ProgressLogger using our stdOut logger, logging messages,
	// not especially handling errors, reading from pchan, not using a progress bar
	go ProgressLogger(stdOut, true, nil, pchan, nil)

	// Put 100 items of Work into the Work channel.
	for i := range 100 {
		wchan <- NewWork(map[string]any{
			"the number": i,
		})
	}
	done() // signal the Supervisor and any idle workers that we're done giving out Work.

	// wait until all outstanding Work is accomplished. i.e. the Job is done.
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
