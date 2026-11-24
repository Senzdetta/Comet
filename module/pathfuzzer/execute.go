// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "fmt"
    "net"
    "os"
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

    words, err := loadWordlistSlice(openWords)
    if err != nil {
        displayErrWordlist := wordlist
        if displayErrWordlist == "" {
            displayErrWordlist = openWords
        }

        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sWordlist: %s%s %snot found!\n",
            color.R, color.N, color.GG, displayErrWordlist, color.N,
        )

        return
    }

    total := len(words)

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
        fmt.Fprintf(
            os.Stderr,
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
    var fuzzResults []FuzzResult
    var dirs []string
    seenDirs := make(map[string]bool)

    runBatch := func(prefix string, allowRecurse bool) {
        for i, w := range words {
            relPath := w
            if prefix != "" {
                relPath = prefix + "/" + w
            }

            wg.Add(1)
            sem <- struct{}{}

            go func(idx int, p string) {
                defer wg.Done()
                fuzz(
                    client, target, p,
                    idx, allowRecurse, total,
                    &fuzzResults, &dirs, &mu,
                )
                <-sem
            }(i+1, relPath)
        }
        wg.Wait()
    }

    runBatch("", recursive)
    if recursive {
        for dIdx := 0; dIdx < len(dirs); dIdx++ {
            dir := dirs[dIdx]
            if seenDirs[dir] {
                continue
            }

            seenDirs[dir] = true
            runBatch(dir, true)
        }
    }

    fmt.Printf("\r\x1b[K")

    sort.Slice(
        fuzzResults, func(i, j int) bool {
            if fuzzResults[i].Status != fuzzResults[j].Status {
                return fuzzResults[i].Status < fuzzResults[j].Status
            }
            return fuzzResults[i].URL < fuzzResults[j].URL
        },
    )

    pcat.Printres(
        fuzzResults, func(r FuzzResult) string {
            return fmt.Sprintf(
                "%s%s %s(%sresp: %s%d%s, %stype %s%s%s, %slenght: %s%d%s)%s",
                color.GG, r.URL, color.DG, color.N,
                color.GG, r.Status, color.DG,
                color.N, color.GG, r.Type, color.DG,
                color.N, color.GG, r.Length, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta