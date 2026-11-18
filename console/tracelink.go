// https://github.com/Senzdetta/Comet

package console

import (
    "fmt"
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/tracelink"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Tracelink struct{}
func (c Tracelink) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    target := args[2]
    threads := 50
    timeout := 15
    recursive := false

    for i := 3; i < len(args); i++ {
        switch args[i] {
            case "--threads":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil {
                        threads = t
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil {
                        timeout = t
                    }
                    i++
                }
            case "--recursive":
                recursive = true
        }
    }

    if threads <= 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid threads value!\n",
            color.R, color.N,
        )
        os.Exit(1)
    }

    if timeout <= 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid timeout value!\n",
            color.R, color.N,
        )
        os.Exit(1)
    }

    tracelink.Execute(
        target,
        threads,
        timeout,
        recursive,
    )
}

// Copyright (c) 2026 Senzdetta