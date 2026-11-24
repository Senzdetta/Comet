// https://github.com/Senzdetta/Comet

package userhunt

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "sync"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    username string,
    wordlistPath string,
    maxThreads int,
    timeoutSec time.Duration,
    checkBody bool,
) {
    var (
        domains []string
        total int
        loadErr error
    )

    if wordlistPath == "" {
        file, err := SitesFS.Open("sites/sites.txt")
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sEmbedded list not found!\n",
                color.R, color.N,
            )
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if line != "" && !strings.HasPrefix(line, "#") {
                domains = append(domains, line)
            }
        }
    } else {
        file, err := os.Open(wordlistPath)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sWordlist file not found: %s%v%s\n",
                color.R, color.N, color.GG, err, color.N,
            )
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if line != "" && !strings.HasPrefix(line, "#") {
                domains = append(domains, line)
            }
        }

        loadErr = scanner.Err()
        if loadErr != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sError reading wordlist: %s%v%s\n",
                color.R, color.N, color.GG, loadErr, color.N,
            )
            return
        }
    }

    total = len(domains)
    if total == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo domains loaded from wordlist!\n",
            color.R, color.N,
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
        pcat.Disconf(username),
    )

    displayWordlist := wordlistPath
    if wordlistPath == "" {
        displayWordlist = "embedded:sites/sites.txt"
    }

    pcat.Printstat(
        "Wordlist",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d x 8)",
                displayWordlist, total,
            ),
        ),
    )

    pcat.Printstat(
        "Threads",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                maxThreads,
            ),
        ),
    )

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(timeoutSec.String()),
    )

    pcat.Printstat(
        "Check Body",
        pcat.Disconf(
            fmt.Sprintf(
                "%t",
                checkBody,
            ),
        ),
    )

    fmt.Println()

    client := &http.Client{
        Timeout: timeoutSec,
        Transport: &http.Transport{
            DialContext: (&net.Dialer{
                Timeout: 5 * time.Second,
                KeepAlive: 5 * time.Second,
            }).DialContext,
            TLSHandshakeTimeout: 5 * time.Second,
            ResponseHeaderTimeout: 6 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
        },
    }

    var wg sync.WaitGroup
    sem := make(chan struct{}, maxThreads)

    var mu sync.Mutex
    var huntResults []Result

    totalVariants := len(domains) * 8
    var atomicCounter int64 = 0

    for _, domain := range domains {
        wg.Add(1)
        sem <- struct{}{}

        go func(d string) {
            defer wg.Done()
            checkDomain(
                &atomicCounter, totalVariants, client,
                username, d, checkBody,
                &huntResults, &mu,
            )
            <-sem
        }(domain)
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    if len(huntResults) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo username found!\n",
            color.R, color.N,
        )
        return
    }

    pcat.Printres(
        huntResults, func(r Result) string {
            return fmt.Sprintf(
                "%s%s %s(%sstatus: %s%d%s)%s",
                color.GG, r.URL, color.DG,
                color.N, color.GG, r.Status, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta