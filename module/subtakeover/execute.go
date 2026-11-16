// https://github.com/Senzdetta/Comet

package subtakeover

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "os"
    "strings"
    "sync"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/module/subdomain"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Takeover(
    domain string,
    threads, timeout int,
    dnsResolv string,
    wordlistPath string,
    fingerprintPath string,
) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    var words []string
    dispWordlist := wordlistPath

    if wordlistPath == "" {
        dispWordlist = "embedded:wordlist/seclists/subdomains-top1million-20000.txt"
        words = parseWordlist(dataWordlist)
    } else {
        f, err := os.Open(wordlistPath)
        if err != nil {
            fmt.Printf(
                "%s[!] %sError opening wordlist: %s%v%s\n",
                color.R, color.N, color.GG, err, color.N,
            )
            return
        }
        defer f.Close()

        var lines []string
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            lines = append(lines, scanner.Text())
        }
        words = parseWordlist(strings.Join(lines, "\n"))
    }

    fingerprints := parseFingerprint(dataFingerprint)
    dispFingerprint := fingerprintPath

    if fingerprintPath == "" {
        dispFingerprint = "embedded:wordlist/fingerprint.hcl"
    } else {
        contentBytes, err := os.ReadFile(fingerprintPath)
        if err == nil {
            fingerprints = parseFingerprint(
                string(contentBytes),
            )
        }
    }

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(domain),
    )

    pcat.Printstat(
        "Wordlist",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                dispWordlist, len(words),
            ),
        ),
    )

    pcat.Printstat(
        "Fingerprints",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                dispFingerprint, len(fingerprints),
            ),
        ),
    )

    pcat.Printstat(
        "DNS Server",
        pcat.Disconf(dnsResolv),
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

    fmt.Println()

    subs := subdomain.BruteForceDNS(
        domain, words, threads,
        dnsResolv, timeout,
    )

    if len(subs) == 0 {
        fmt.Printf(
            "%s[!] %sNo subdomains found!\n",
            color.R, color.N,
        )
        return
    }

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

    var wg sync.WaitGroup
    sem := make(chan struct{}, threads)

    var mu sync.Mutex
    var takeoverResults []TakeoverResult

    totalSubs := len(subs)

    completed := 0
    for _, sub := range subs {
        wg.Add(1)
        sem <- struct{}{}

        go func(target string) {
            defer wg.Done()
            defer func() {
                <-sem
            }()

            mu.Lock()
            completed++
            pcat.Printf(
                "%s[*] %sScanning %s(%s%d%s/%s%d%s)%s: %s%s%s",
                color.B, color.N, color.DG,
                color.GG, completed, color.DG,
                color.GG, totalSubs, color.DG, color.N,
                color.GG, target, color.N,
            )
            mu.Unlock()

            var resp *http.Response
            var reqErr error

            testURL := "http://" + target
            req, _ := http.NewRequest("GET", testURL, nil)
            req.Header.Set(
                "User-Agent",
                "https://github.com/Senzdetta/Comet",
            )

            resp, reqErr = client.Do(req)
            if reqErr != nil {
                testURL = "https://" + target
                req, _ = http.NewRequest("GET", testURL, nil)
                req.Header.Set(
                    "User-Agent",
                    "https://github.com/Senzdetta/Comet",
                )
                resp, reqErr = client.Do(req)
            }

            if reqErr != nil || resp == nil {
                return
            }
            defer resp.Body.Close()

            bodyBytes, _ := io.ReadAll(resp.Body)
            bodyStr := string(bodyBytes)

            for prov, sig := range fingerprints {
                if strings.Contains(bodyStr, sig) {
                    mu.Lock()
                    takeoverResults = append(
                        takeoverResults, TakeoverResult{
                            Target: target,
                            Provider: prov,
                            Signature: sig,
                        },
                    )
                    mu.Unlock()
                    break
                }
            }
        }(sub)
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    if len(takeoverResults) == 0 {
        fmt.Printf(
            "%s[!] %sNo takeover vulnerability found!\n",
            color.R, color.N,
        )
        return
    }

    for _, res := range takeoverResults {
        pcat.Mklog(
            "subtakeover",
            ":",
            fmt.Sprintf(
                "Takeover: %s (%s: %s)",
                res.Target, res.Provider, res.Signature,
            ),
        )
    }

    pcat.Printres(
        takeoverResults, func(r TakeoverResult) string {
            return fmt.Sprintf(
                "%s%s %s(%s%s: %s%s%s)%s",
                color.GG, r.Target, color.DG,
                color.N, r.Provider, color.GG, r.Signature, color.DG,
                color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta