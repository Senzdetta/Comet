// https://github.com/Senzdetta/Comet

package utils

import (
    "os"
    "path/filepath"
)

func Root() string {
    exePath, err := os.Executable()
    if err == nil {
        exePath, err = filepath.EvalSymlinks(exePath)
    }

    if err != nil {
        return ""
    }

    return filepath.Clean(filepath.Dir(exePath))
}

// Copyright (c) 2026 Senzdetta