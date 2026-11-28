// https://github.com/Senzdetta/Comet

package portscan

import (
    "os"
    "fmt"
    "strconv"
    "strings"
    "sync"
    "time"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    ip string,
    start, end int,
    protocol, only, threads, timeout, servicePath string,
    skipHostCheck bool,
) {
    if servicePath != "" {
        _, err := os.Stat(servicePath)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sFailed to open service: %s%s%s\n",
                color.R, color.N, color.GG, servicePath, color.N,
            )
            return
        }
    }

    var protos []string
    seenProto := map[string]bool{}

    for _, p := range strings.Split(
        strings.ToLower(protocol), ",",
    ) {
        p = strings.TrimSpace(p)
        if (p == "tcp" || p == "udp" || p == "sctp") && !seenProto[p] {
            seenProto[p] = true
            protos = append(protos, p)
        }
    }

    if len(protos) == 0 {
        protos = []string{
            "tcp",
            "udp",
            "sctp",
        }
    }

    var onlyList []string
    onlySet := map[string]bool{}

    for _, o := range strings.Split(
        strings.ToLower(only), ",",
    ) {
        o = strings.TrimSpace(o)
        switch o {
            case "open":
            case "filtered", "open|filtered":
                o = "filtered"
            case "closed", "close":
                o = "closed"
            default:
                continue
        }

        if !onlySet[o] {
            onlySet[o] = true
            onlyList = append(onlyList, o)
        }
    }

    nThreads, err := strconv.Atoi(threads)
    if err != nil || nThreads < 1 || nThreads > 10000 {
        nThreads = 0
    }

    nTimeout := 0
    if ms, err := strconv.Atoi(timeout); err == nil && ms > 0 {
        nTimeout = ms
    }

    if !skipHostCheck && !checkHost(ip) {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sHost: %s%s %sis down!\n",
            color.R, color.N, color.GG, ip, color.N,
        )
        return
    }

    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(ip),
    )

    pcat.Printstat(
        "Port Range",
        pcat.Disconf(
            fmt.Sprintf(
                "%d:%d",
                start, end,
            ),
        ),
    )

    pcat.Printstat(
        "Protocol",
        pcat.Disconf(
            strings.Join(protos, ","),
        ),
    )

    dispService := servicePath
    if dispService == "" {
        dispService = "embedded:services/nmap-services.txt"
    }

    pcat.Printstat(
        "Services",
        pcat.Disconf(dispService),
    )

    pcat.Printstat(
        "Only",
        pcat.Disconf(
            strings.Join(onlyList, ","),
        ),
    )

    pcat.Printstat(
        "Threads",
        pcat.Disconf(threads),
    )

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(timeout),
    )

    fmt.Println()

    var wg sync.WaitGroup
    var mu sync.Mutex
    var portResults []PortResult

    for _, pr := range protos {
        t := 500 * time.Millisecond
        n := 200

        switch pr {
            case "udp":
                t = 1500 * time.Millisecond
                n = 50
            case "sctp":
                t = 1000 * time.Millisecond
                n = 50
        }

        if nTimeout > 0 {
            t = time.Duration(nTimeout) * time.Millisecond
        }

        if nThreads > 0 {
            n = nThreads
        }

        sem := make(chan struct{}, n)

        for port := start; port <= end; port++ {
            wg.Add(1)
            sem <- struct{}{}

            go func(p int, pr string, t time.Duration) {
                defer wg.Done()
                defer func() {
                    <-sem
                }()

                mu.Lock()
                pcat.Printf(
                    "%s[*] %sScanning %s(%s%s%s, %s%d%s/%s%d%s)%s: %s%s%s:%s%d%s",
                    color.B, color.N, color.DG,
                    color.GG, pr, color.DG,
                    color.GG, p, color.DG, color.GG, end, color.DG,
                    color.N,
                    color.GG, ip, color.DG, color.BB, p, color.N,
                )
                mu.Unlock()

                r := scanRun(ip, p, pr, t)
                if allow(onlySet, r.Status) {
                    svc := lookupService(
                        r.Port, r.Proto,
                        servicePath,
                    )

                    banner := r.Banner
                    if banner == "" {
                        banner = "unknown"
                    }

                    mu.Lock()
                    portResults = append(
                        portResults, PortResult{
                            IP: ip,
                            Port: r.Port,
                            Proto: r.Proto,
                            Status: r.Status,
                            Banner: banner,
                            Svc: svc,
                        },
                    )
                    mu.Unlock()
                }
            }(port, pr, t)
        }
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    if len(portResults) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo port matched on ip: %s%s%s\n",
            color.R, color.N, color.GG, ip, color.N,
        )
        return
    }

    pcat.Printres(
        portResults, func(r PortResult) string {
            return fmt.Sprintf(
                "%s%s:%d%s (%sprotocol: %s%s, %sservice: %s%s, %sstatus: %s%s%s, %sprob: %s%s, %sbanner: %s%s, %sdesc: %s%s%s)%s",
                color.GG, r.IP, r.Port, color.DG,
                color.N, color.GG, r.Proto,
                color.N, color.GG, r.Svc.Name,
                color.N, color.GG, r.Status, color.DG,
                color.N, color.GG, r.Svc.Prob,
                color.N, color.GG, r.Banner,
                color.N, color.GG, r.Svc.Desc, color.DG,
                color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta