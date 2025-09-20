// Package racket is a manager for Jobs, Work, and Progress.
//
// A Job is a repetitive task that uses a common Supervisor to ensure Work is properly distributed,
// that the correct number of workers are available to do the Work, and that those workers can
// send Progress along as-needed.
//
// Work is a map of parameters that a worker can take to do its Job. It is important that Work contains
// at least all of the parameters a worker expects.
//
// Progress is a typed communication construct that allows workers to send a predictable and actionable
// set of information. Some Progress may be ignored sometimes (e.g. if one is not tracking units of work
// (ala a progress bar) then the ProgressUpdate and ProgressEstimate types are simply discarded).
//
// The interfaces are written such that one may use the existing Job system via NewJob() or implement their
// own. Likewise one might use Progress and Work and ignore the Job altogether.
package racket

import (
	"sync/atomic"
	"time"

	"github.com/cognusion/semaphore"
)

// Job is a repetitive task that uses a common Supervisor to ensure Work is properly distributed,
// that the correct number of workers are available to do the Work, and that those workers can
// send Progress along as-needed.
type Job interface {
	// Supervisor will ensure there are workers to do the Work, and a channel to receive that Work on,
	// while also supplying a means to receive progress reports and how to report back when there is no
	// more work to do.
	Supervisor(maxWorkers int, workChan chan Work) (progressChan chan Progress, doneFunc func())
	// NewWorker will ready a worker to do some Work, giving it an ID to reference it by. Calling this directly
	// is generally unnecessary as Supervisor will handle it.
	NewWorker(id any)
	// IsDone will wait until all of the doled-out Work had been completed, and all of the workers have left.
	// It's flexible enough to be used as a blocking inline "wait" or in a select{} so other things can occur whilst
	// waiting.
	IsDone() <-chan bool
}

// WorkerFunc is a definition for how to accomplish Work!
// Each invocation can assume it has been giving a unique ID, has it's own unique Work, and it can send
// various Progress updates over the supplied channel.
type WorkerFunc func(id any, work Work, progressChan chan<- Progress)

// defaultJob is a Job that takes a dynamic worker definition to accomplish varied Work using the same
// Supervisor system.
type defaultJob struct {
	workerFunc   WorkerFunc
	workChan     chan Work
	workerCount  atomic.Int64
	progressChan chan Progress
	doneChan     chan struct{}
	lock         semaphore.Semaphore
}

// NewJob consumes a WorkerFunc to accomplish Work, and returns a Job.
func NewJob(workerFunc WorkerFunc) Job {
	return &defaultJob{
		workerFunc: workerFunc,
	}
}

// NewWorker spins up a workerFunc to accomplish Work,
// blocking until Work has been accomplished, or there is
// no more to do.
func (j *defaultJob) NewWorker(id any) {
	defer j.lock.Unlock()
	defer j.workerCount.Add(-1)

	select {
	case w := <-j.workChan:
		j.workerFunc(id, w, j.progressChan)
	case <-j.doneChan:
	}
}

// IsDone waits until all of the workers have completed, kind of.
// After done() has been called, if there are zero workers 4 consecutive 10ms polls,
// we assume we are done.
func (j *defaultJob) IsDone() <-chan bool {
	b := make(chan bool)

	go func() {
		var count int
		<-j.doneChan // if doneChan isn't closed, we are definitely not done

		for {
			if j.workerCount.Load() > 0 {
				count = 0
			} else {
				count++
			}
			if count > 4 {
				break
			}
			<-time.After(10 * time.Millisecond)
		}
		b <- true
	}()

	return b
}

// Supervisor spins up maxWorkers, who will wait for Work via workChan, and returns a channel for
// progress reciepts and func to signal when there is no new Work to be added to workChan.
func (j *defaultJob) Supervisor(maxWorkers int, workChan chan Work) (progressChan chan Progress, doneFunc func()) {
	j.doneChan = make(chan struct{})
	j.progressChan = make(chan Progress)
	j.workChan = workChan
	j.lock = semaphore.NewSemaphore(maxWorkers)

	go func() {
		c := 0
		for {
			c++
			select {
			case <-j.lock.Until():
				// woo! make a worker!
				j.workerCount.Add(1)
				go j.NewWorker(c)
			case <-j.doneChan:
				// That's all folks!
				return
			}
		}
	}()

	return j.progressChan, func() { close(j.doneChan) }
}
