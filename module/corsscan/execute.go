// https://github.com/Senzdetta/Comet

package corsscan

import (
    "os"
    "fmt"
    "net"
    "strings"
    "sync"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    targetInput, origin, customHeader string,
    threads, timeout int,
) {
    targetURL := targetInput
    if !strings.HasPrefix(
        targetURL, "http://",
    ) && !strings.HasPrefix(
        targetURL, "https://",
    ) {
        targetURL = "http://" + targetURL
    }

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

    pcat.Printstat(
        "Origin",
        pcat.Disconf(origin),
    )

    pcat.Printstat(
        "Header",
        pcat.Disconf(customHeader),
    )

    pcat.Printstat(
        "Threads", pcat.Disconf(
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
                Timeout: 5 * time.Second,
                KeepAlive: 30 * time.Second,
            }).DialContext,
            TLSHandshakeTimeout: 5 * time.Second,
            ResponseHeaderTimeout: 5 * time.Second,
        },
    }

    _, err := client.Head(targetURL)
    if err != nil {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sTarget: %s%s %snot responding!\n",
            color.R, color.N, color.GG, targetURL, color.N,
        )
        return
    }

    var mu sync.Mutex
    var vulns []VulnResult

    targets := []string{targetURL}

    var wg sync.WaitGroup
    sem := make(chan struct{}, threads)

    for _, tURL := range targets {
        wg.Add(1)
        sem <- struct{}{}

        go func(t string) {
            defer wg.Done()
            defer func() { <-sem }()

            mu.Lock()
            pcat.Printf(
                "%s[*] %sScanning: %s%s%s",
                color.B, color.N, color.GG, t, color.N,
            )
            mu.Unlock()

            req, err := http.NewRequest("GET", t, nil)
            if err != nil {
                return
            }

            req.Header.Set(
                "User-Agent",
                "https://github.com/Senzdetta/Comet",
            )
            req.Header.Set("Origin", origin)

            if customHeader != "" {
                parts := strings.SplitN(customHeader, ":", 2)
                if len(parts) == 2 {
                    req.Header.Set(
                        strings.TrimSpace(parts[0]),
                        strings.TrimSpace(parts[1]),
                    )
                }
            }

            resp, err := client.Do(req)
            if err != nil {
                return
            }
            resp.Body.Close()

            acao := resp.Header.Get(
                "Access-Control-Allow-Origin",
            )

            acac := resp.Header.Get(
                "Access-Control-Allow-Credentials",
            )

            if acao == origin || acao == "*" {
                vulnSeverity := "low/mid"
                if acac == "true" && acao != "*" {
                    vulnSeverity = "high/critical"
                }

                mu.Lock()
                vulns = append(
                    vulns, VulnResult{
                        Target: t,
                        Severity: vulnSeverity,
                        Acao: acao,
                        Acac: acac,
                    },
                )
                mu.Unlock()
            }
        }(tURL)
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    pcat.Printres(
        vulns, func(v VulnResult) string {
            return fmt.Sprintf(
                "%s%s %s(%sseverity: %s%s%s, %sacao: %s%s%s, %sacac: %s%s%s)%s",
                color.GG, v.Target,
                color.DG, color.N, color.GG, v.Severity,
                color.DG, color.N, color.GG, v.Acao,
                color.DG, color.N, color.GG, v.Acac,
                color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta