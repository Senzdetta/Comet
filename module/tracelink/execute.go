// https://github.com/Senzdetta/Comet

package tracelink

import (
    "os"
    "fmt"
    "sync"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    target string,
    threads int,
    timeout int,
    recursive bool,
) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(target),
    )

    pcat.Printstat(
        "Threads",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                threads,
            ),
        ),
    )

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(
            fmt.Sprintf(
                "%ds",
                timeout,
            ),
        ),
    )

    pcat.Printstat(
        "Recursive",
        pcat.Disconf(
            fmt.Sprintf(
                "%t",
                recursive,
            ),
        ),
    )

    fmt.Println()

    var traceResults []TraceResult
    seen := make(map[string]struct{})
    seen[target] = struct{}{}

    var mu sync.Mutex
    var wg sync.WaitGroup
    var crawl func(string)

    sem := make(chan struct{}, threads)
    crawl = func(u string) {
        defer wg.Done()

        mu.Lock()
        pcat.Printf(
            "%s[*] %sCrawling: %s%s%s",
            color.B, color.N,
            color.GG, u, color.N,
        )
        mu.Unlock()

        links := fetchLinks(u, timeout)

        for _, link := range links {
            mu.Lock()
            if _, ok := seen[link]; ok {
                mu.Unlock()
                continue
            }

            seen[link] = struct{}{}
            traceResults = append(
                traceResults, TraceResult{
                    SourceLink: u,
                    TargetLink: link,
                },
            )
            mu.Unlock()

            if recursive && isSameHost(target, link) {
                wg.Add(1)
                sem <- struct{}{}

                go func(next string) {
                    defer func() {
                        <-sem
                    }()
                    crawl(next)
                }(link)
            }
        }
    }

    wg.Add(1)
    sem <- struct{}{}
    go func() {
        defer func() {
            <-sem
        }()
        crawl(target)
    }()

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    if len(traceResults) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo links found!\n",
            color.R, color.N,
        )
        return
    }

    pcat.Printres(
        traceResults, func(r TraceResult) string {
            return fmt.Sprintf(
                "%s%s %s-> %s%s%s",
                color.GG, r.SourceLink, color.DG,
                color.GG, r.TargetLink, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta