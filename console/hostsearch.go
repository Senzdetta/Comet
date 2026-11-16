// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/hostsearch"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Hostsearch struct{}
func (c Hostsearch) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    timeout := 10
    retries := 5
    apiURL := "https://api.hackertarget.com/hostsearch/?q="
    targetDomain := ""

    for i := 2; i < len(args); i++ {
        switch args[i] {
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
            case "--api":
                if i+1 < len(args) {
                    apiURL = args[i+1]
                    i++
                }
            default:
                if targetDomain == "" && args[i][0] != '-' {
                    targetDomain = args[i]
                }
        }
    }

    if targetDomain == "" {
        invinput.MissingArgument()
        os.Exit(1)
    }

    hostsearch.Execute(
        targetDomain,
        apiURL,
        timeout,
        retries,
    )
}

// Copyright (c) 2026 Senzdetta