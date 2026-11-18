// https://github.com/Senzdetta/Comet

package subdomain

import (
    "os"
    "context"
    "fmt"
    "sort"
    "strings"
    "sync"
    "time"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(opts Options) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    var words []string
    var errLoad error

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(opts.Domain),
    )

    pcat.Printstat(
        "Mode",
        pcat.Disconf(opts.Mode),
    )

    if opts.Mode == "api" {
        pcat.Printstat(
            "API",
            pcat.Disconf(opts.Api),
        )

        pcat.Printstat(
            "Retry",
            pcat.Disconf(
                fmt.Sprintf(
                    "%d",
                    opts.Retry,
                ),
            ),
        )
    } else {
        words, errLoad = loadWordlist(opts.Wordlist)
        if errLoad != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sWordlist Error: %s%v%s\n",
                color.R, color.N, color.GG, errLoad, color.N,
            )
            return
        }

        dispWordlist := opts.Wordlist
        if dispWordlist == "" {
            dispWordlist = "embedded:wordlist/subdomains-top1million-20000.txt"
        }

        pcat.Printstat(
            "Wordlist",
            fmt.Sprintf(
                "%s (%d)",
                pcat.Disconf(dispWordlist), len(words),
            ),
        )

        pcat.Printstat(
            "Threads",
            pcat.Disconf(
                fmt.Sprintf(
                    "%d",
                    opts.Threads,
                ),
            ),
        )

        pcat.Printstat(
            "DNS Server",
            pcat.Disconf(opts.DNSResolv),
        )
    }

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(
            fmt.Sprintf(
                "%ds",
                opts.Timeout,
            ),
        ),
    )

    pcat.Printstat(
        "Lookup IP",
        pcat.Disconf(
            fmt.Sprintf(
                "%t",
                opts.LookupIP,
            ),
        ),
    )

    fmt.Println()

    var rawSubs []string
    var err error

    if opts.Mode == "api" {
        rawSubs, err = fetch(
            opts.Domain, opts.Api,
            opts.Timeout, opts.Retry,
        )

        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sError: %s%v%s\n",
                color.R, color.N, color.GG, err, color.N,
            )
            return
        }
    } else {
        rawSubs = BruteForceDNS(
            opts.Domain, words,
            opts.Threads, opts.DNSResolv,
            opts.Timeout,
        )
    }

    if len(rawSubs) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo subdomains found for: %s%s%s\n",
            color.R, color.N, color.GG, opts.Domain, color.N,
        )
        return
    }

    var finalResults []SubResult
    var mu sync.Mutex
    var wg sync.WaitGroup

    if opts.LookupIP {
        sem := make(chan struct{}, opts.Threads)
        ctx, cancel := context.WithTimeout(
            context.Background(),
            time.Duration(opts.Timeout)*time.Second,
        )
        defer cancel()

        for _, sub := range rawSubs {
            wg.Add(1)
            sem <- struct{}{}

            go func(s string) {
                defer wg.Done()
                defer func() {
                    <-sem
                }()

                mu.Lock()
                pcat.Printf(
                    "%s[*] %sResolving: %s%s%s",
                    color.B, color.N, color.GG, s, color.N,
                )
                mu.Unlock()

                ips := resolveIP(ctx, s, opts.DNSResolv)
                ipStr := ""
                if len(ips) > 0 {
                    ipStr = strings.Join(ips, ", ")
                }

                mu.Lock()
                finalResults = append(
                    finalResults, SubResult{
                        Subdomain: s,
                        IPs: ipStr,
                    },
                )
                mu.Unlock()
            }(sub)
        }

        wg.Wait()
        fmt.Printf("\r\x1b[K")
    } else {
        for _, sub := range rawSubs {
            finalResults = append(
                finalResults, SubResult{
                    Subdomain: sub,
                    IPs: "",
                },
            )
        }
    }

    if len(finalResults) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo valid subdomains resolved!\n",
            color.R, color.N,
        )
        return
    }

    sort.Slice(
        finalResults, func(i, j int) bool {
            return finalResults[i].Subdomain < finalResults[j].Subdomain
        },
    )

    for _, res := range finalResults {
        ipLogPart := ""
        if res.IPs != "" {
            ipLogPart = fmt.Sprintf(
                " [%s]", res.IPs,
            )
        }

        pcat.Mklog(
            "subdomain",
            ":",
            fmt.Sprintf(
                "Found: %s%s",
                res.Subdomain, ipLogPart,
            ),
        )
    }

    pcat.Printres(
        finalResults, func(r SubResult) string {
            ipPart := ""
            if r.IPs != "" {
                splitIPs := strings.Split(r.IPs, ", ")
                formattedIPs := make([]string, len(splitIPs))
                for i, ip := range splitIPs {
                    formattedIPs[i] = fmt.Sprintf(
                        "%s%s%s",
                        color.GG, ip, color.N,
                    )
                }

                joinedIPs := strings.Join(
                    formattedIPs, fmt.Sprintf(
                        "%s, %s",
                        color.DG, color.GG,
                    ),
                )

                ipPart = fmt.Sprintf(
                    "%s(%sip: %s%s%s)",
                    color.DG, color.N, joinedIPs, color.DG, color.N,
                )
            }

            return fmt.Sprintf(
                "%s%s %s%s",
                color.GG, r.Subdomain, ipPart, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta