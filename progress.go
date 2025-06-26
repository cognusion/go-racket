package racket

import (
	"fmt"
	"log"
)

// ProgressError is a ProgressType when the Data is an error.
// ProgressUpdate is a ProgressType when the Data is a numeric update (ala progress bar +/- math).
// ProgressEsimate is a ProgressType when the Data is a numeric [re]evaluation of how much work is to be performed.
// ProgressMessage is a ProgressType when the Data is a string message.
// ProgressOther is a ProgressType when Data is to be consumed elsewhere, and should not be interpretted outside of that elsewhere.
const (
	ProgressError ProgressType = iota
	ProgressUpdate
	ProgressEstimate
	ProgressMessage
	ProgressOther
)

type (
	// ProgressType is one of the constant types of Progress.
	ProgressType int
	// ProgressErrorFunc is a function that consumes an error.
	ProgressErrorFunc func(error)
	// Progress is a tuple of a ProgressType and Data. It is also an error and a string.
	Progress struct {
		Type ProgressType
		Data any
	}
)

// String returns the stringified version of the type name
func (p ProgressType) String() string {
	switch p {
	case ProgressError:
		return "ProgressError"
	case ProgressUpdate:
		return "ProgressUpdate"
	case ProgressEstimate:
		return "ProgressEstimate"
	case ProgressMessage:
		return "ProgressMessage"
	case ProgressOther:
		return "ProgressOther"
	default:
		return ""
	}
}

// Error returns the Progress Data as an error if Progress is a ProgressError, or nil.
func (p *Progress) Error() error {
	if p.Type == ProgressError {
		return p.Data.(error)
	}
	return nil
}

// String returns a formatted string representation of the ProgressType and the Data.
func (p *Progress) String() string {
	return fmt.Sprintf("%s: %+v", p.Type, p.Data)
}

// ProgressLogger is a helper that can loop over a Progress channel and triage the items generically.
// If non-nil, the supplied ProgressErrorFunc will be called with the error after it is logged or printed:
// Panic'ing or Exit'ing is allowed.
// ProgressBar-related Progress will be sent to the barChan as-is.
func ProgressLogger(outLog *log.Logger, logMessages bool, errf ProgressErrorFunc, progressChan <-chan Progress, barChan chan Progress) {
	for p := range progressChan {
		//outLog.Printf("PROGRESS! %+v\n", p)
		switch p.Type {
		case ProgressError:
			// Always print errors.
			outLog.Printf("[PROGRESS] ERROR: %s\n", p.Data.(error))

			if errf != nil {
				// callback
				errf(p.Data.(error))
			}
		case ProgressMessage:
			if logMessages {
				// Always print if we're logging.
				outLog.Printf("[PROGRESS] %s\n", p.Data.(string))
			}
		case ProgressUpdate, ProgressEstimate:
			if logMessages {
				outLog.Printf("[PROGRESS] %s: %d\n", p.Type.String(), p.Data.(int64))
			}
			if barChan != nil {
				barChan <- p
			}
		default:
			// Always print weird shit.
			outLog.Printf("[PROGRESS] ??: %+v\n", p)
		}
	}
}

// PErrorf returns a ProgressError with a formatted error.
func PErrorf(format string, a ...any) Progress {
	return Progress{
		Type: ProgressError,
		Data: fmt.Errorf(format, a...),
	}
}

// PMessagef returns a ProgressMessage with a formatted string.
func PMessagef(format string, a ...any) Progress {
	return Progress{
		Type: ProgressMessage,
		Data: fmt.Sprintf(format, a...),
	}
}

// PUpdate returns a ProgressUpdate with the specified count.
func PUpdate(count int64) Progress {
	return Progress{
		Type: ProgressUpdate,
		Data: count,
	}
}

// PEstimate returns a ProgressEstimate with the specified estimate.
func PEstimate(estimate int64) Progress {
	return Progress{
		Type: ProgressEstimate,
		Data: estimate,
	}
}
