// https://github.com/Senzdetta/Comet

package subdomain

import (
    "sync"
    "context"
    "fmt"
    "time"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func BruteForceDNS(
    domain string,
    words []string,
    threads int,
    dnsAddr string,
    timeoutSec int,
) []string {
    jobs := make(chan string, len(words))
    results := make(chan string, len(words))
    var wg sync.WaitGroup

    ctx, cancel := context.WithTimeout(
        context.Background(),
        time.Duration(timeoutSec)*time.Second,
    )
    defer cancel()

    total := len(words)
    var mu sync.Mutex

    completed := 0
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for w := range jobs {
                target := fmt.Sprintf(
                    "%s.%s", w, domain,
                )

                mu.Lock()
                completed++
                pcat.Printf(
                    "%s[*] %sBruteforcing %s(%s%d%s/%s%d%s)%s: %s%s%s",
                    color.B, color.N, color.DG,
                    color.GG, completed, color.DG,
                    color.GG, total, color.DG,
                    color.N, color.GG, target, color.N,
                )
                mu.Unlock()

                ips := resolveIP(ctx, target, dnsAddr)
                if len(ips) > 0 {
                    results <- target
                }
            }
        }()
    }

    for _, w := range words {
        jobs <- w
    }
    close(jobs)

    wg.Wait()
    fmt.Printf("\r\x1b[K")
    close(results)

    var found []string
    for res := range results {
        found = append(found, res)
    }

    return found
}

// Copyright (c) 2026 Senzdetta