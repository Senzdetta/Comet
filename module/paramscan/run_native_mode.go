// https://github.com/Senzdetta/Comet

package paramscan

import (
    "sync"
    "time"
    "fmt"
    "strings"
    "net/url"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func runNativeMode(opts Options) {
    baseURL := opts.Target
    if !strings.HasPrefix(
        baseURL, "http://",
    ) && !strings.HasPrefix(
        baseURL, "https://",
    ) {
        baseURL = "http://" + baseURL
    }

    parsedURL, err := url.Parse(baseURL)
    if err != nil {
        fmt.Printf(
            "%s[!] %sInvalid target URL: %v\n",
            color.R, color.N, err,
        )
        return
    }

    if parsedURL.Path == "" {
        parsedURL.Path = "/"
    }

    cleanBaseURL := fmt.Sprintf(
        "%s://%s%s",
        parsedURL.Scheme, parsedURL.Host, parsedURL.Path,
    )

    words := loadLines(
        opts.WordlistPath,
        "wordlist/seclists/burp-parameter-names.txt",
    )

    values := loadLines(
        opts.ValuePath,
        "wordlist/parameter_values.txt",
    )

    if len(words) == 0 {
        fmt.Printf(
            "%s[!] %sWordlist is empty!\n",
            color.R, color.N,
        )
        return
    }

    if len(values) == 0 {
        values = []string{"1"}
    }

    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(cleanBaseURL),
    )

    pcat.Printstat(
        "Mode",
        pcat.Disconf("native"),
    )

    dispWordlist := opts.WordlistPath
    if dispWordlist == "" {
        dispWordlist = "embedded:wordlist/seclists/burp-parameter-names.txt"
    }

    pcat.Printstat(
        "Wordlist",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                dispWordlist, len(words),
            ),
        ),
    )

    dispValue := opts.ValuePath
    if dispValue == "" {
        dispValue = "embedded:wordlist/parameter_values.txt"
    }

    pcat.Printstat(
        "Value Wordlist",
        pcat.Disconf(
            fmt.Sprintf(
                "%s (%d)",
                dispValue, len(values),
            ),
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
        "Timeout",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                opts.Timeout,
            ),
        ),
    )

    fmt.Println()

    var wg sync.WaitGroup
    sem := make(chan struct{}, opts.Threads)

    var mu sync.Mutex
    var paramResults []ParamResult

    client := &http.Client{
        Timeout: time.Duration(opts.Timeout) * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: opts.Threads,
            MaxIdleConnsPerHost: opts.Threads,
        },
    }

    totalWl := len(words)
    totalVlw := len(values)

    for wIdx, param := range words {
        for vIdx, val := range values {
            wg.Add(1)
            sem <- struct{}{}

            go func(
                pw string,
                vw string,
                wi, vi int,
            ) {
                defer wg.Done()
                defer func() {
                    <-sem
                }()

                testURL := fmt.Sprintf(
                    "%s?%s=%s",
                    cleanBaseURL, pw, vw,
                )

                mu.Lock()
                pcat.Printf(
                    "%s[*] %sFuzzing: %s(%s%d%s/%s%d%s, %s%d%s/%s%d%s)%s: %s%s%s?%s%s%s=%s%s%s",
                    color.B, color.N,
                    color.DG, color.GG, wi+1, color.DG, color.GG, totalWl,
                    color.DG, color.GG, vi+1, color.DG, color.GG, totalVlw,
                    color.DG, color.N,
                    color.GG, cleanBaseURL,
                    color.DG, color.BB, pw,
                    color.DG, color.BB, vw,
                    color.N,
                )
                mu.Unlock()

                res, err := client.Get(testURL)
                if err == nil {
                    if (res.StatusCode == 200 ||
                        res.StatusCode == 301 ||
                        res.StatusCode == 302) {
                            mu.Lock()
                            paramResults = append(
                                paramResults, ParamResult{
                                    URL: testURL,
                                    Status: res.StatusCode,
                                },
                            )
                            mu.Unlock()
                    }
                    res.Body.Close()
                }
            }(param, val, wIdx, vIdx)
        }
    }

    wg.Wait()
    fmt.Printf("\r\x1b[K")

    if len(paramResults) == 0 {
        fmt.Printf(
            "%s[!] %sNo active parameters found!\n",
            color.R, color.N,
        )
        return
    }

    for _, r := range paramResults {
        pcat.Mklog(
            "paramscan",
            ":",
            fmt.Sprintf(
                "Found: %s:%d",
                r.URL, r.Status,
            ),
        )
    }

    pcat.Printres(
        paramResults, func(r ParamResult) string {
            return fmt.Sprintf(
                "%s%s %s(%sstatus: %s%d%s)%s",
                color.GG, r.URL, color.DG, color.N,
                color.GG, r.Status, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta