package chanx

import (
	"errors"
	"reflect"
)

var (
	errSrcMustBeChan         = errors.New("src must be a channel")
	errSrcMustBeReadable     = errors.New("src channel must be a readable:[<-chan] or [chan]")
	errDstMustBeChannel      = errors.New("all dst must be channels")
	errSrcDstTypeMustBeMatch = errors.New("dst channel element type must match src channel")
	errDstMustBeWriteable    = errors.New("dst channels must be writable:[chan<-] or [chan]")
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

	// Validate and prepare destination channels
	if len(dst) == 0 {
		return nil
	}

	dstVals := make([]reflect.Value, 0, len(dst))
	for _, d := range dst {
		dstVal := reflect.ValueOf(d)
		if dstVal.Kind() != reflect.Chan {
			return errDstMustBeChannel
		}
		// Check if channel element types match
		if dstVal.Type().Elem() != srcVal.Type().Elem() {
			return errSrcDstTypeMustBeMatch
		}
		// Check if channel is writable
		if dir := dstVal.Type().ChanDir(); dir != reflect.SendDir && dir != reflect.BothDir {
			return errDstMustBeWriteable
		}
		dstVals = append(dstVals, dstVal)
	}

	// Start broadcasting in a separate goroutine
	go func() {
		defer func() {
			// Close all destination channels when done
			if r := recover(); r != nil {
				//XXX: do something?
			}
			for _, d := range dstVals {
				d.Close()
			}
		}()

		for {
			// Receive from source channel
			v, ok := srcVal.Recv()
			if !ok {
				// Source channel is closed
				return
			}

			// Send to all destination channels
			for _, d := range dstVals {
				d.Send(v)
			}
		}
	}()

	return nil
}
