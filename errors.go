package onet

import (
	"github.com/libs4go/errors"
)

// ScopeOfAPIError .
const errVendor = "onet"

// errors
var (
	ErrParams        = errors.New("params error", errors.WithVendor(errVendor), errors.WithCode(-1))
	ErrExists        = errors.New("transport load error", errors.WithVendor(errVendor), errors.WithCode(-2))
	ErrNotFound      = errors.New("resource not found", errors.WithVendor(errVendor), errors.WithCode(-3))
	ErrProtocolValue = errors.New("protocol value check error", errors.WithVendor(errVendor), errors.WithCode(-4))
	ErrClosed        = errors.New("the socket closed", errors.WithVendor(errVendor), errors.WithCode(-5))
	ErrMuxNotFound   = errors.New("mux session not found", errors.WithVendor(errVendor), errors.WithCode(-6))
)
