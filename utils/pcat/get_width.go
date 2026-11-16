// https://github.com/Senzdetta/Comet

package pcat

import (
    "os"
    "strconv"
    "golang.org/x/term"
)

func getWidth() int {
    if colStr := os.Getenv("COLUMNS"); colStr != "" {
        if w, err := strconv.Atoi(
            colStr,
        ); err == nil && w > 0 {
            return w
        }
    }

    fd := int(os.Stdout.Fd())
    if term.IsTerminal(fd) {
        if width, _, err := term.GetSize(
            fd,
        ); err == nil && width > 0 {
            return width
        }
    }

    return 80
}

// Copyright (c) 2026 Senzdetta