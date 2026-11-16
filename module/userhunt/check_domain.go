// https://github.com/Senzdetta/Comet

package userhunt

import (
    "io"
    "strings"
    "sync"
    "sync/atomic"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func checkDomain(
    counter *int64,
    totalVariants int,
    client *http.Client,
    username, domain string,
    checkBody bool,
    results *[]Result,
    mu *sync.Mutex,
) {
    domain = strings.TrimSpace(domain)
    if domain == "" || strings.HasPrefix(domain, "#") {
        return
    }

    base := strings.ReplaceAll(domain, "{}", username)
    variants := generateVariants(base)

    for _, url := range variants {
        current := atomic.AddInt64(counter, 1)

        mu.Lock()
        pcat.Printf(
            "%s[*] %sChecking: %s(%s%d%s/%s%d%s)%s: %s%s%s",
            color.B, color.N, color.DG,
            color.GG, current, color.DG,
            color.GG, totalVariants, color.DG,
            color.N, color.GG, url, color.N,
        )
        mu.Unlock()

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            continue
        }

        req.Header.Set(
            "User-Agent",
            "https://github.com/Senzdetta/Comet",
        )

        resp, err := client.Do(req)
        if err != nil {
            continue
        }

        status := resp.StatusCode
        if checkBody {
            bodyBytes, err := io.ReadAll(
                io.LimitReader(
                    resp.Body,
                    1024*64,
                ),
            )

            resp.Body.Close()
            if err != nil {
                continue
            }

            bodyStr := string(bodyBytes)
            if status == 200 {
                lowerBody := strings.ToLower(bodyStr)
                if strings.Contains(
                    lowerBody, "username not found",
                ) || strings.Contains(
                    lowerBody, "page not found",
                ) || strings.Contains(
                    lowerBody, "user doesn't exist",
                ) || strings.Contains(
                    lowerBody, "no profile",
                ) || strings.Contains(
                    lowerBody, "invalid username",
                ) {
                    continue
                }

                mu.Lock()
                *results = append(
                    *results, Result{
                        URL: url,
                        Status: status,
                    },
                )
                mu.Unlock()
                return
            }
        } else {
            resp.Body.Close()
            if (status == 200 ||
                status == 301 ||
                status == 302 ||
                status == 403) {
                    mu.Lock()
                    *results = append(
                        *results, Result{
                            URL: url,
                            Status: status,
                        },
                    )
                    mu.Unlock()
                    return
            }
        }
    }
}

// Copyright (c) 2026 Senzdetta