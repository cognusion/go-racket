

# racket
`import "github.com/cognusion/go-racket"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package racket is a manager for Jobs, Work, and Progress.

A Job is a repetitive task that uses a common Supervisor to ensure Work is properly distributed,
that the correct number of workers are available to do the Work, and that those workers can
send Progress along as-needed.

Work is a map of parameters that a worker can take to do its Job. It is important that Work contains
at least all of the parameters a worker expects.

Progress is a typed communication construct that allows workers to send a predictable and actionable
set of information. Some Progress may be ignored sometimes (e.g. if one is not tracking units of work
(ala a progress bar) then the ProgressUpdate and ProgressEstimate types are simply discarded).

The interfaces are written such that one may use the existing Job system via NewJob() or implement their
own. Likewise one might use Progress and Work and ignore the Job altogether.


##### Example :
``` go
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
```



## <a name="pkg-index">Index</a>
* [func ProgressLogger(outLog *log.Logger, logMessages bool, errf ProgressErrorFunc, progressChan &lt;-chan Progress, barChan chan Progress)](#ProgressLogger)
* [type Job](#Job)
  * [func NewJob(workerFunc WorkerFunc) Job](#NewJob)
* [type Progress](#Progress)
  * [func PErrorf(format string, a ...any) Progress](#PErrorf)
  * [func PEstimate(estimate int64) Progress](#PEstimate)
  * [func PMessagef(format string, a ...any) Progress](#PMessagef)
  * [func PUpdate(count int64) Progress](#PUpdate)
  * [func (p *Progress) Error() error](#Progress.Error)
  * [func (p *Progress) String() string](#Progress.String)
* [type ProgressErrorFunc](#ProgressErrorFunc)
* [type ProgressType](#ProgressType)
  * [func (p ProgressType) String() string](#ProgressType.String)
* [type Work](#Work)
  * [func NewWork(config map[string]any) Work](#NewWork)
  * [func (w *Work) Get(key string) any](#Work.Get)
  * [func (w *Work) GetBool(key string) bool](#Work.GetBool)
  * [func (w *Work) GetInt(key string) int](#Work.GetInt)
  * [func (w *Work) GetString(key string) string](#Work.GetString)
* [type WorkerFunc](#WorkerFunc)

#### <a name="pkg-examples">Examples</a>
* [Package](#example-)

#### <a name="pkg-files">Package files</a>
[job.go](https://github.com/cognusion/go-racket/tree/master/job.go) [progress.go](https://github.com/cognusion/go-racket/tree/master/progress.go) [work.go](https://github.com/cognusion/go-racket/tree/master/work.go)





## <a name="ProgressLogger">func</a> [ProgressLogger](https://github.com/cognusion/go-racket/tree/master/progress.go?s=2013:2147#L68)
``` go
func ProgressLogger(outLog *log.Logger, logMessages bool, errf ProgressErrorFunc, progressChan <-chan Progress, barChan chan Progress)
```
ProgressLogger is a helper that can loop over a Progress channel and triage the items generically.
If non-nil, the supplied ProgressErrorFunc will be called with the error after it is logged or printed:
Panic'ing or Exit'ing is allowed.
ProgressBar-related Progress will be sent to the barChan as-is.




## <a name="Job">type</a> [Job](https://github.com/cognusion/go-racket/tree/master/job.go?s=1251:2032#L27)
``` go
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
```
Job is a repetitive task that uses a common Supervisor to ensure Work is properly distributed,
that the correct number of workers are available to do the Work, and that those workers can
send Progress along as-needed.







### <a name="NewJob">func</a> [NewJob](https://github.com/cognusion/go-racket/tree/master/job.go?s=2685:2723#L57)
``` go
func NewJob(workerFunc WorkerFunc) Job
```
NewJob consumes a WorkerFunc to accomplish Work, and returns a Job.





## <a name="Progress">type</a> [Progress](https://github.com/cognusion/go-racket/tree/master/progress.go?s=897:948#L27)
``` go
type Progress struct {
    Type ProgressType
    Data any
}

```
Progress is a tuple of a ProgressType and Data. It is also an error and a string.







### <a name="PErrorf">func</a> [PErrorf](https://github.com/cognusion/go-racket/tree/master/progress.go?s=2895:2941#L100)
``` go
func PErrorf(format string, a ...any) Progress
```
PErrorf returns a ProgressError with a formatted error.


### <a name="PEstimate">func</a> [PEstimate](https://github.com/cognusion/go-racket/tree/master/progress.go?s=3455:3494#L124)
``` go
func PEstimate(estimate int64) Progress
```
PEstimate returns a ProgressEstimate with the specified estimate.


### <a name="PMessagef">func</a> [PMessagef](https://github.com/cognusion/go-racket/tree/master/progress.go?s=3089:3137#L108)
``` go
func PMessagef(format string, a ...any) Progress
```
PMessagef returns a ProgressMessage with a formatted string.


### <a name="PUpdate">func</a> [PUpdate](https://github.com/cognusion/go-racket/tree/master/progress.go?s=3286:3320#L116)
``` go
func PUpdate(count int64) Progress
```
PUpdate returns a ProgressUpdate with the specified count.





### <a name="Progress.Error">func</a> (\*Progress) [Error](https://github.com/cognusion/go-racket/tree/master/progress.go?s=1420:1452#L52)
``` go
func (p *Progress) Error() error
```
Error returns the Progress Data as an error if Progress is a ProgressError, or nil.




### <a name="Progress.String">func</a> (\*Progress) [String](https://github.com/cognusion/go-racket/tree/master/progress.go?s=1613:1647#L60)
``` go
func (p *Progress) String() string
```
String returns a formatted string representation of the ProgressType and the Data.




## <a name="ProgressErrorFunc">type</a> [ProgressErrorFunc](https://github.com/cognusion/go-racket/tree/master/progress.go?s=780:809#L25)
``` go
type ProgressErrorFunc func(error)
```
ProgressErrorFunc is a function that consumes an error.










## <a name="ProgressType">type</a> [ProgressType](https://github.com/cognusion/go-racket/tree/master/progress.go?s=702:718#L23)
``` go
type ProgressType int
```
ProgressType is one of the constant types of Progress.


``` go
const (
    ProgressError ProgressType = iota
    ProgressUpdate
    ProgressEstimate
    ProgressMessage
    ProgressOther
)
```
ProgressError is a ProgressType when the Data is an error.
ProgressUpdate is a ProgressType when the Data is a numeric update (ala progress bar +/- math).
ProgressEsimate is a ProgressType when the Data is a numeric [re]evaluation of how much work is to be performed.
ProgressMessage is a ProgressType when the Data is a string message.
ProgressOther is a ProgressType when Data is to be consumed elsewhere, and should not be interpretted outside of that elsewhere.










### <a name="ProgressType.String">func</a> (ProgressType) [String](https://github.com/cognusion/go-racket/tree/master/progress.go?s=1011:1048#L34)
``` go
func (p ProgressType) String() string
```
String returns the stringified version of the type name




## <a name="Work">type</a> [Work](https://github.com/cognusion/go-racket/tree/master/work.go?s=131:174#L8)
``` go
type Work struct {
    // contains filtered or unexported fields
}

```
Work is a representation of specification to pass to a Worker doing a Job.







### <a name="NewWork">func</a> [NewWork](https://github.com/cognusion/go-racket/tree/master/work.go?s=237:277#L13)
``` go
func NewWork(config map[string]any) Work
```
NewWork takes a map and returns a specified unit of Work.





### <a name="Work.Get">func</a> (\*Work) [Get](https://github.com/cognusion/go-racket/tree/master/work.go?s=376:410#L20)
``` go
func (w *Work) Get(key string) any
```
Get returns the value associated with the key, or nil.




### <a name="Work.GetBool">func</a> (\*Work) [GetBool](https://github.com/cognusion/go-racket/tree/master/work.go?s=658:697#L30)
``` go
func (w *Work) GetBool(key string) bool
```
GetBool returns the bool-ified value associated with the key.




### <a name="Work.GetInt">func</a> (\*Work) [GetInt](https://github.com/cognusion/go-racket/tree/master/work.go?s=802:839#L35)
``` go
func (w *Work) GetInt(key string) int
```
GetInt returns the int-ifiied value associated with the key.




### <a name="Work.GetString">func</a> (\*Work) [GetString](https://github.com/cognusion/go-racket/tree/master/work.go?s=507:550#L25)
``` go
func (w *Work) GetString(key string) string
```
GetString returns the string-ified value associated with the key.




## <a name="WorkerFunc">type</a> [WorkerFunc](https://github.com/cognusion/go-racket/tree/master/job.go?s=2251:2320#L44)
``` go
type WorkerFunc func(id any, work Work, progressChan chan<- Progress)
```
WorkerFunc is a definition for how to accomplish Work!
Each invocation can assume it has been giving a unique ID, has it's own unique Work, and it can send
various Progress updates over the supplied channel.














- - -
Generated by [godoc2md](http://github.com/cognusion/godoc2md)
