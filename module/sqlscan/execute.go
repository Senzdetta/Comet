// https://github.com/Senzdetta/Comet

package sqlscan

import (
    "fmt"
    "io"
    "net"
    "sync"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    targetURL, wordlistPath string,
    threads, timeout int,
) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(targetURL),
    )

    payloads, err := LoadWordlist(wordlistPath, PayloadFS)
    if err != nil {
        fmt.Printf(
            "%s[!] %sError loading wordlist: %s%v%s\n",
            color.R, color.N, color.GG, err, color.N,
        )
        return
    }

    if len(payloads) == 0 {
        fmt.Printf(
            "%s[!] %sNo valid payloads found!%s\n",
            color.R, color.N, color.N,
        )
        return
    }

    displayPayload := wordlistPath
    if displayPayload == "" {
        displayPayload = "embedded:payload/"
    }

    pcat.Printstat(
        "Payloads",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                displayPayload, len(payloads),
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

    fmt.Println()

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: threads,
            MaxIdleConnsPerHost: threads,
            DialContext: (&net.Dialer{
                Timeout: 10 * time.Second,
                KeepAlive: 30 * time.Second,
            }).DialContext,
            ResponseHeaderTimeout: 10 * time.Second,
        },
    }

    _, err = client.Head(targetURL)
    if err != nil {
        fmt.Printf(
            "%s[!] %sTarget: %s%s %snot responding!\n",
            color.R, color.N, color.GG, targetURL, color.N,
        )
        return
    }

    var wg sync.WaitGroup
    sem := make(chan struct{}, threads)

    var mu sync.Mutex
    var sqlResults []SqlResult

    idx := 0
    for _, payload := range payloads {
        testURLs := buildTestURLs(targetURL, payload)
        for _, tURL := range testURLs {
            wg.Add(1)
            sem <- struct{}{}

            go func(testURL, p string) {
                defer wg.Done()
                defer func() {
                    <-sem
                }()

                mu.Lock()
                pcat.Printf(
                    "%s[*] %sScanning %s(%s%d%s/%s%d%s)%s: %s%s%s%s%s",
                    color.B, color.N, color.DG, color.GG, idx, color.DG,
                    color.GG, len(payloads), color.DG, color.N,
                    color.GG, targetURL, color.BB, p, color.N,
                )
                mu.Unlock()

                req, err := http.NewRequest(
                    "GET", testURL, nil,
                )

                if err != nil {
                    return
                }

                req.Header.Set(
                    "User-Agent",
                    "https://github.com/Senzdetta/Comet",
                )

                resp, err := client.Do(req)
                if err != nil || resp == nil {
                    return
                }
                defer resp.Body.Close()

                bodyBytes, err := io.ReadAll(resp.Body)
                if err != nil {
                    return
                }
                bodyStr := string(bodyBytes)

                isVuln, sig := checkVulnerable(
                    bodyStr,
                    resp.StatusCode,
                )

                if isVuln {
                    mu.Lock()
                    sqlResults = append(
                        sqlResults, SqlResult{
                            URL: testURL,
                            Signature: sig,
                            Resp: resp.StatusCode,
                        },
                    )

                    pcat.Mklog(
                        "sqlscan",
                        ":",
                        fmt.Sprintf(
                            "%s (signature: %s, resp: %d)",
                            testURL, sig, resp.StatusCode,
                        ),
                    )
                    mu.Unlock()
                }
            }(tURL, payload)
            idx++
        }
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    pcat.Printres(
        sqlResults, func(r SqlResult) string {
            return fmt.Sprintf(
                "%s%s %s(%ssignature: %s%s%s, %sresp: %s%d%s)%s",
                color.GG, r.URL, color.DG,
                color.N, color.GG, r.Signature, color.DG,
                color.N, color.GG, r.Resp, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta