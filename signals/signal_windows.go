// Copyright 2022 The efucloud.com Authors.
//go:build windows
// +build windows

package signals

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
