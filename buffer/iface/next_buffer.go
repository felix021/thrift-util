package iface

import "errors"

var (
	ErrOutOfBound = errors.New("out of bound")
)

type NextBuffer interface {
	Next(n int) (p []byte, err error)
	Skip(n int) (err error)
}
