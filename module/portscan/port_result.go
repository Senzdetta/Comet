// https://github.com/Senzdetta/Comet

package portscan

type PortResult struct {
    IP string
    Port int
    Proto string
    Status string
    Banner string
    Svc Service
}

// Copyright (c) 2026 Senzdetta