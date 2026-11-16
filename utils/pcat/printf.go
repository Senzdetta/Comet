// https://github.com/Senzdetta/Comet

package pcat

import (
    "os"
    "fmt"
    "golang.org/x/term"
)

func Printf(format string, a ...interface{}) {
    msg := fmt.Sprintf(format, a...)
    fd := int(os.Stdout.Fd())

    if term.IsTerminal(fd) {
        width := getWidth()
        msg = truncate(msg, width)
        fmt.Printf(
            "\r%s\x1b[K",
            msg,
        )
    } else {
        fmt.Printf(
            "%s\n",
            msg,
        )
    }
}

// Copyright (c) 2026 Senzdetta