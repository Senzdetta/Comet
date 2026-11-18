// https://github.com/Senzdetta/Comet

package subdomain

import (
    "fmt"
    "net"
    "strings"
    "time"
    "net/http"
    "encoding/json"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func fetch(
    domain string,
    apiPattern string,
    timeout int,
    retries int,
) ([]string, error) {
    apiURL := strings.ReplaceAll(
        apiPattern, "{target}", domain,
    )

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: 100,
            MaxIdleConnsPerHost: 100,
            DialContext: (&net.Dialer{
                Timeout: 15 * time.Second,
                KeepAlive: 30 * time.Second,
            }).DialContext,
            TLSHandshakeTimeout: 10 * time.Second,
            ResponseHeaderTimeout: 15 * time.Second,
        },
    }

    var resp *http.Response
    var err error

    for i := 0; i < retries; i++ {
        pcat.Printf(
            "%s[!] %sConnection unstable, retrying api %s(%s%d%s/%s%d%s)%s",
            color.R, color.N,
            color.DG, color.GG, i+1, color.DG,
            color.GG, retries, color.DG, color.N,
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

        time.Sleep(time.Duration(i+2) * time.Second)
    }

    fmt.Printf("\r\x1b[K")

    if err != nil || resp == nil {
        return nil, fmt.Errorf(
            "failed to connect to api: %v",
            err,
        )
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf(
            "server returned status code: %d",
            resp.StatusCode,
        )
    }

    var results []CrtResult
    err = json.NewDecoder(resp.Body).Decode(&results)
    if err != nil {
        return nil, err
    }

    uniqueSubs := make(map[string]bool)
    for _, res := range results {
        subs := strings.Split(res.NameValue, "\n")
        for _, sub := range subs {
            sub = strings.TrimSpace(sub)
            sub = strings.TrimPrefix(sub, "*.")
            if strings.HasSuffix(sub, domain) {
                uniqueSubs[sub] = true
            }
        }
    }

    var sortedSubs []string
    for sub := range uniqueSubs {
        sortedSubs = append(sortedSubs, sub)
    }

    return sortedSubs, nil
}

// Copyright (c) 2026 Senzdetta