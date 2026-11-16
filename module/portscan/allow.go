// https://github.com/Senzdetta/Comet

package portscan

func allow(only map[string]bool, status string) bool {
    if status == "unsupported" {
        return false
    }

    if len(only) == 0 {
        return status == "open"
    }

    if status == "open|filtered" {
        status = "filtered"
    }

    return only[status]
}

// Copyright (c) 2026 Senzdetta