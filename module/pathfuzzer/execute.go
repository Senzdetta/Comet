// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "fmt"
    "net"
    "sort"
    "strings"
    "sync"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    target string,
    wordlist string,
    timeout int,
    recursive bool,
    threads int,
) {
    target = strings.TrimRight(target, "/")

    openWords := wordlist
    if openWords == "" {
        openWords = "wordlist/filedirs.txt"
    }

    scanner, closeFn, total, err := openWordlist(
        WordlistFS, openWords,
    )

    if err != nil {
        fmt.Printf(
            "%s[!] %sWordlist: %s%s %snot found!\n",
            color.R, color.N, color.GG, wordlist, color.N,
        )
        return
    }
    defer closeFn()

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
        Transport: &http.Transport{
            DialContext: (&net.Dialer{
                Timeout: 2 * time.Second,
                KeepAlive: 2 * time.Second,
            }).DialContext,
            TLSHandshakeTimeout: 2 * time.Second,
            ResponseHeaderTimeout: 3 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
        },
    }

    _, err = client.Head(target)
    if err != nil {
        fmt.Printf(
            "%s[!] %sTarget: %s%s %snot responding!\n",
            color.R, color.N, color.GG, target, color.N,
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
        pcat.Disconf(target),
    )

    displayWordlist := wordlist
    if displayWordlist == "" {
        displayWordlist = "embedded:wordlist/filedirs.txt"
    }

    pcat.Printstat(
        "Wordlist",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                displayWordlist, total,
            ),
        ),
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
                "%d",
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

    var wg sync.WaitGroup
    sem := make(chan struct{}, threads)

    var mu sync.Mutex
    resultsMap := make(map[int][]string)
    var dirs []string

    idx := 1
    for scanner.Scan() {
        dir := scanner.Text()

        wg.Add(1)
        sem <- struct{}{}

        go func(i int, d string) {
            defer wg.Done()
            fuzz(
                client, target, d,
                i, recursive, total,
                &resultsMap, &dirs, &mu,
            )
            <-sem
        }(idx, dir)
        idx++
    }

    wg.Wait()

    if recursive {
        for _, dir := range dirs {
            subScanner, subClose, _, err := openWordlist(
                WordlistFS, wordlist,
            )

            if err != nil {
                continue
            }

            idx := 1
            for subScanner.Scan() {
                sub := strings.TrimSpace(subScanner.Text())

                wg.Add(1)
                sem <- struct{}{}

                go func(i int, path string) {
                    defer wg.Done()
                    fuzz(
                        client, target, path,
                        i, false, total,
                        &resultsMap, &dirs, &mu,
                    )
                    <-sem
                }(idx, dir+"/"+sub)
                idx++
            }

            wg.Wait()
            subClose()
        }
    }

    fmt.Printf("\r\x1b[K")

    var fuzzResults []FuzzResult
    statuses := []int{200, 301, 302}

    for _, status := range statuses {
        urls, ok := resultsMap[status]
        if !ok || len(urls) == 0 {
            continue
        }

        sort.Strings(urls)
        for _, u := range urls {
            fuzzResults = append(
                fuzzResults, FuzzResult{
                    URL: u,
                    Status: status,
                },
            )

            pcat.Mklog(
                "pathfuzzer",
                ":",
                fmt.Sprintf(
                    "%s (resp: %d)",
                    u, status,
                ),
            )
        }
    }

    pcat.Printres(
        fuzzResults, func(r FuzzResult) string {
            return fmt.Sprintf(
                "%s%s %s(%sresp: %s%d%s)%s",
                color.GG, r.URL, color.DG, color.N,
                color.GG, r.Status, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta