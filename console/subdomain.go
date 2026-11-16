// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "strings"
    "github.com/Senzdetta/Comet/module/subdomain"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Subdomain struct{}
func (c Subdomain) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    opts := subdomain.Options{
        Mode: "native",
        Threads: 50,
        Timeout: 15,
        Retry: 5,
        Api: "https://crt.sh/?q=%.{target}&output=json",
        LookupIP: true,
        Wordlist: "",
        DNSResolv: "8.8.8.8:53",
    }

    for i := 2; i < len(args); i++ {
        switch args[i] {
            case "--mode":
                if i+1 < len(args) {
                    opts.Mode = strings.ToLower(args[i+1])
                    i++
                }
            case "--threads":
                if i+1 < len(args) {
                    if val, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        opts.Threads = val
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    if val, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        opts.Timeout = val
                    }
                    i++
                }
            case "--retry":
                if i+1 < len(args) {
                    if val, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        opts.Retry = val
                    }
                    i++
                }
            case "--api":
                if i+1 < len(args) {
                    opts.Api = args[i+1]
                    i++
                }
            case "--lookup-ip":
                if i+1 < len(args) {
                    val := strings.ToLower(args[i+1])
                    opts.LookupIP = (val == "true" || val == "1")
                    i++
                }
            case "--wordlist":
                if i+1 < len(args) {
                    opts.Wordlist = args[i+1]
                    i++
                }
            case "--dnsresolv":
                if i+1 < len(args) {
                    opts.DNSResolv = args[i+1]
                    i++
                }
            default:
                if opts.Domain == "" && args[i][0] != '-' {
                    opts.Domain = args[i]
                }
        }
    }

    if opts.Domain == "" {
        invinput.MissingArgument()
        os.Exit(1)
    }

    subdomain.Execute(opts)
}

// Copyright (c) 2026 Senzdetta