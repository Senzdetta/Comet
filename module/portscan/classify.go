// https://github.com/Senzdetta/Comet

package portscan

import (
    "syscall"
    "errors"
)

func classify(err error) string {
    if errors.Is(
        err,
        syscall.ECONNREFUSED,
    ) {
        return "closed"
    }
    return "filtered"
}

// Copyright (c) 2026 Senzdetta