// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/subtakeover"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Subtakeover struct{}
func (c Subtakeover) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    targetDomain := ""
    threads := 50
    timeout := 15
    dnsResolv := "8.8.8.8:53"
    wordlistPath := ""
    fingerprintPath := ""

    for i := 2; i < len(args); i++ {
        switch args[i] {
            case "--threads":
                if i+1 < len(args) {
                    if t, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        threads = t
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    if t, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        timeout = t
                    }
                    i++
                }
            case "--dnsresolv":
                if i+1 < len(args) {
                    dnsResolv = args[i+1]
                    i++
                }
            case "--wordlist":
                if i+1 < len(args) {
                    wordlistPath = args[i+1]
                    i++
                }
            case "--fingerprint":
                if i+1 < len(args) {
                    fingerprintPath = args[i+1]
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

    subtakeover.Takeover(
        targetDomain,
        threads,
        timeout,
        dnsResolv,
        wordlistPath,
        fingerprintPath,
    )
}

// Copyright (c) 2026 Senzdetta