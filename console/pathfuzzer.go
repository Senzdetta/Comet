// https://github.com/Senzdetta/Comet

package console

import (
    "fmt"
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/pathfuzzer"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Pathfuzzer struct{}
func (c Pathfuzzer) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    target := args[2]
    wordlist := ""
    threads := 100
    timeout := 15
    recursive := false

    for i := 3; i < len(args); i++ {
        switch args[i] {
            case "--recursive":
                recursive = true
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
            case "--wordlist":
                if i+1 < len(args) {
                    wordlist = args[i+1]
                    i++
                }
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

    pathfuzzer.Execute(
        target,
        wordlist,
        timeout,
        recursive,
        threads,
    )
}

// Copyright (c) 2026 Senzdetta