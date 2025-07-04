package chanx

import (
	"errors"
	"reflect"
)

var (
	errSrcMustBeChan     = errors.New("src must be a channel")
	errSrcMustBeReadable = errors.New("src channel must be a readable:[<-chan] or [chan]")
)

// Broadcast
// Read from src, send to the dst which could be a list
func Broadcast(src any, dst ...any) error {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Chan {
		return errSrcMustBeChan
	}

	//src channel must be readable
	if dir := srcVal.Type().ChanDir(); dir != reflect.RecvDir && dir != reflect.BothDir {
		return errSrcMustBeReadable
	}

	// TODO

	return nil
}
