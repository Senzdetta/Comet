// https://github.com/Senzdetta/Comet

package paramscan

import (
    "time"
    "bufio"
    "strings"
    "fmt"
    "sync"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func runAPIMode(opts Options) {
    targetClean := opts.Target
    targetClean = strings.TrimPrefix(targetClean, "http://")
    targetClean = strings.TrimPrefix(targetClean, "https://")
    targetClean = strings.Split(targetClean, "/")[0]

    apiURL := strings.ReplaceAll(
        opts.CustomAPI, "{target}", targetClean,
    )

    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(targetClean),
    )

    pcat.Printstat(
        "Mode",
        pcat.Disconf("api"),
    )

    pcat.Printstat(
        "API Endpoint",
        pcat.Disconf(opts.CustomAPI),
    )

    pcat.Printstat(
        "Fuzz",
        pcat.Disconf(
            fmt.Sprintf(
                "%t",
                opts.Fuzz,
            ),
        ),
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

    pcat.Printstat(
        "Threads",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                opts.Threads,
            ),
        ),
    )

    fmt.Println()

    client := &http.Client{
        Timeout: 60 * time.Second,
    }

    var resp *http.Response
    var err error

    for i := 0; i < opts.Retry; i++ {
        pcat.Printf(
            "%s[!] %sConnecting to API, retrying %s(%s%d%s/%s%d%s)%s",
            color.R, color.N,
            color.DG, color.GG, i+1, color.DG,
            color.CC, opts.Retry, color.DG, color.N,
        )

        req, _ := http.NewRequest("GET", apiURL, nil)
        req.Header.Set(
            "User-Agent",
            "https://github.com/Senzdetta/Comet",
        )

        resp, err = client.Do(req)
        if err == nil && resp.StatusCode == 200 {
            break
        }

        if resp != nil {
            resp.Body.Close()
        }

        time.Sleep(
            time.Duration(i+2) * time.Second,
        )
    }

    fmt.Printf("\r\x1b[K")

    if err != nil || resp == nil {
        fmt.Printf(
            "%s[!] %sFailed to connect to API: %s%v%s\n",
            color.R, color.N, color.GG, err, color.N,
        )
        return
    }
    defer resp.Body.Close()

    var patterns []string
    seen := make(map[string]bool)
    scanner := bufio.NewScanner(resp.Body)

    for scanner.Scan() {
        link := scanner.Text()
        if strings.Contains(link, "?") {
            var item string
            if !opts.Fuzz {
                item = link
            } else {
                item = formatToFuzz(link)
            }

            if !seen[item] {
                patterns = append(patterns, item)
                seen[item] = true
            }
        }
    }

    if len(patterns) == 0 {
        fmt.Printf(
            "%s[!] %sNo parameterized urls found!\n",
            color.R, color.N,
        )
        return
    }

    var paramResults []ParamResult
    if !opts.Fuzz {
        for _, p := range patterns {
            paramResults = append(
                paramResults, ParamResult{
                    URL: p,
                    Status: 200,
                },
            )
        }
    } else {
        var wg sync.WaitGroup
        sem := make(chan struct{}, opts.Threads)
        var mu sync.Mutex

        fuzzClient := &http.Client{
            Timeout: time.Duration(opts.Timeout) * time.Second,
            Transport: &http.Transport{
                MaxIdleConns: opts.Threads,
                MaxIdleConnsPerHost: opts.Threads,
            },
        }

        for _, p := range patterns {
            for val := 1; val <= 5; val++ {
                wg.Add(1)
                sem <- struct{}{}

                go func(pattern string, v int) {
                    defer wg.Done()
                    defer func() {
                        <-sem
                    }()

                    testURL := strings.ReplaceAll(
                        pattern, "FUZZ",
                        fmt.Sprintf("%d", v),
                    )

                    mu.Lock()
                    pcat.Printf(
                        "%s[*] %sFuzzing: %s%s%s",
                        color.B, color.N, color.GG, testURL, color.N,
                    )
                    mu.Unlock()

                    for retryCount := 0; retryCount < opts.Retry; retryCount++ {
                        req, _ := http.NewRequest(
                            "GET", testURL, nil,
                        )

                        req.Header.Set(
                            "User-Agent",
                            "https://github.com/Senzdetta/Comet",
                        )

                        res, err := fuzzClient.Do(req)
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
                            break
                        }
                        time.Sleep(200 * time.Millisecond)
                    }
                }(p, val)
            }
        }
        wg.Wait()
        fmt.Printf("\r\x1b[K")
    }

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
                "%s%s %(%sstatus: %s%d%s)%s",
                color.GG, r.URL, color.DG, color.N,
                color.GG, r.Status, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta