package racket

import (
	"github.com/spf13/cast"
)

// Work is a representation of specification to pass to a Worker doing a Job.
type Work struct {
	config map[string]any
}

// NewWork takes a map and returns a specified unit of Work.
func NewWork(config map[string]any) Work {
	return Work{
		config: config,
	}
}

// Get returns the value associated with the key, or nil.
func (w *Work) Get(key string) any {
	return w.config[key]
}

// GetString returns the string-ified value associated with the key.
func (w *Work) GetString(key string) string {
	return cast.ToString(w.config[key])
}

// GetBool returns the bool-ified value associated with the key.
func (w *Work) GetBool(key string) bool {
	return cast.ToBool(w.config[key])
}

// GetInt returns the int-ifiied value associated with the key.
func (w *Work) GetInt(key string) int {
	return cast.ToInt(w.config[key])
}
