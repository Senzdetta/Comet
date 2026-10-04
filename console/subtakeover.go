// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/subtakeover"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type SubTakeover struct{}
func (c SubTakeover) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    targetDomain := ""
    threads := 100
    timeout := 10
    retries := 5

    for i := 2; i < len(args); i++ {
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
            case "--retry":
                if i+1 < len(args) {
                    r, err := strconv.Atoi(args[i+1])
                    if err == nil {
                        retries = r
                    }
                    i++
                }
            default:
                if (targetDomain == "" &&
                    args[i][0] != '-') {
                        targetDomain = args[i]
                }
        }
    }

    if targetDomain == "" {
        invinput.MissingArgument()
        os.Exit(1)
    }

    subtakeover.Takeover(
        targetDomain,
        threads,
        timeout,
        retries,
    )
}

// Copyright (c) 2026 Senzdetta