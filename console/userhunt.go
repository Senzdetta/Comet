// https://github.com/Senzdetta/Comet

package console

import (
    "strconv"
    "fmt"
    "os"
    "time"
    "github.com/Senzdetta/Comet/module/userhunt"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Userhunt struct{}
func (c Userhunt) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    username := args[2]
    maxThreads := 100
    timeoutSec := 50 * time.Second
    checkBody := true
    wordlistPath := ""

    for i := 3; i < len(args); i++ {
        switch args[i] {
            case "--threads":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil && t > 0 {
                        maxThreads = t
                    } else {
                        fmt.Printf(
                            "%s[!] %sInvalid thread value: %s%s%s\n",
                            color.R, color.N, color.GG, args[i+1], color.N,
                        )
                        os.Exit(1)
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil && t > 0 {
                        timeoutSec = time.Duration(t) * time.Second
                    } else {
                        fmt.Printf(
                            "%s[!] %sInvalid timeout value: %s%s%s\n",
                            color.R, color.N, color.GG, args[i+1], color.N,
                        )
                        os.Exit(1)
                    }
                    i++
                }
            case "--check-body":
                if i+1 < len(args) {
                    val := args[i+1]
                    if val == "true" {
                        checkBody = true
                    } else if val == "false" {
                        checkBody = false
                    } else {
                        fmt.Printf(
                            "%s[!] %sInvalid check-body value: %s%s%s\n",
                            color.R, color.N, color.GG, val, color.N,
                        )
                        os.Exit(1)
                    }
                    i++
                }
            case "--site":
                if i+1 < len(args) {
                    wordlistPath = args[i+1]
                    i++
                }
        }
    }

    userhunt.Execute(
        username,
        wordlistPath,
        maxThreads,
        timeoutSec,
        checkBody,
    )
}

// Copyright (c) 2026 Senzdetta