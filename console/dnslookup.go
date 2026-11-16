// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "strings"
    "github.com/Senzdetta/Comet/module/dnslookup"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Dnslookup struct{}
func (c Dnslookup) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    timeout := 10
    retries := 5
    threads := 100
    mode := "native"
    apiURL := "https://api.hackertarget.com/dnslookup/?q="
    dnsResolv := "8.8.8.8:53"
    targetDomain := ""

    for i := 2; i < len(args); i++ {
        switch args[i] {
            case "--timeout":
                if i+1 < len(args) {
                    if t, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        timeout = t
                    }
                    i++
                }
            case "--retry":
                if i+1 < len(args) {
                    if r, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        retries = r
                    }
                    i++
                }
            case "--threads":
                if i+1 < len(args) {
                    if th, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        threads = th
                    }
                    i++
                }
            case "--mode":
                if i+1 < len(args) {
                    mode = args[i+1]
                    i++
                }
            case "--api":
                if i+1 < len(args) {
                    apiURL = args[i+1]
                    i++
                }
            case "--dnsresolv":
                if i+1 < len(args) {
                    dnsResolv = args[i+1]
                    if !strings.Contains(dnsResolv, ":") {
                        dnsResolv = dnsResolv + ":53"
                    }
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

    dnslookup.Execute(
        dnslookup.Config{
            Domain: targetDomain,
            Timeout: timeout,
            Retries: retries,
            Threads: threads,
            Mode: mode,
            APIURL: apiURL,
            DNSResolv: dnsResolv,
        },
    )
}

// Copyright (c) 2026 Senzdetta