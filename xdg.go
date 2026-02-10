package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"fmt"
	"os"
	"strings"
)

var xdg_runtime_dir string

func DetectXdgRuntimeDir() {
	// get xdg_runtime_dir
	xdg_runtime_dir = os.Getenv("XDG_RUNTIME_DIR")
	if xdg_runtime_dir == "" {
		xdg_runtime_dir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	xdg_runtime_dir = strings.TrimRight(xdg_runtime_dir, "/")
}
