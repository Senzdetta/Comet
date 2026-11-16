// https://github.com/Senzdetta/Comet

package portscan

import (
    "sync"
    _ "embed"
)

//go:embed services/nmap-services.txt
var nmapServices string

var (
    serviceMap  map[string]Service
    serviceOnce sync.Once

    printMu sync.Mutex
    sctpWarn sync.Once
)

// Copyright (c) 2026 Senzdetta