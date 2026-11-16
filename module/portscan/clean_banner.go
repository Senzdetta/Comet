// https://github.com/Senzdetta/Comet

package portscan

import (
    "strings"
)

func cleanBanner(b []byte) string {
    s := strings.TrimSpace(string(b))
    if i := strings.IndexAny(s, "\r\n"); i >= 0 {
        s = s[:i]
    }

    s = strings.Map(
        func(r rune) rune {
            if r < 32 || r > 126 {
                return -1
            }
            return r
        },
        s,
    )

    if len(s) > 80 {
        s = s[:80]
    }

    if s == "" {
        return "unknown"
    }

    return s
}

// Copyright (c) 2026 Senzdetta