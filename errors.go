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
)
