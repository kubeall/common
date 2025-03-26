// Copyright 2022 The efucloud.com Authors.
//go:build linux || darwin
// +build linux darwin

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
