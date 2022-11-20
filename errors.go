package storage

import "fmt"

var (
	ErrNothingFound = fmt.Errorf("there is no document with this parameters")
)
