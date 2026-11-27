// https://github.com/Senzdetta/Comet

package vcore

import (
    "os"
    "os/exec"
    "path/filepath"
)

func Exepath() string {
    if len(os.Args) == 0 || os.Args[0] == "" {
        return "unknown"
    }

    if filepath.Base(os.Args[0]) != os.Args[0] {
        abs, err := filepath.Abs(os.Args[0])
        if err == nil {
            return filepath.Clean(abs)
        }
        return filepath.Clean(os.Args[0])
    }

    if p, err := exec.LookPath(os.Args[0]); err == nil {
        abs, err := filepath.Abs(p)
        if err == nil {
            return filepath.Clean(abs)
        }
        return filepath.Clean(p)
    }

    return os.Args[0]
}

// Copyright (c) 2026 Senzdetta