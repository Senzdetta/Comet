// https://github.com/Senzdetta/Comet

package dnslookup

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "net"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func runAPILookup(cfg Config) []DNSResult {
    targetURL := cfg.APIURL
    if !strings.Contains(targetURL, "=") {
        targetURL = fmt.Sprintf(
            "%s?q=%s",
            strings.TrimSuffix(cfg.APIURL, "/"), cfg.Domain,
        )
    } else if !strings.Contains(targetURL, cfg.Domain) {
        targetURL = fmt.Sprintf(
            "%s%s",
            cfg.APIURL, cfg.Domain,
        )
    }

    client := &http.Client{
        Timeout: time.Duration(cfg.Timeout) * time.Second,
        Transport: &http.Transport{
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

    for i := 0; i < cfg.Retries; i++ {
        req, _ := http.NewRequest(
            "GET", targetURL, nil,
        )

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

        pcat.Printf(
            "%s[!] %sConnection unstable, retrying api %s(%s%d%s/%s%d%s)%s",
            color.R, color.N,
            color.DG, color.GG, i+1, color.DG,
            color.CC, cfg.Retries, color.DG, color.N,
        )

        time.Sleep(time.Duration(i+2) * time.Second)
        fmt.Printf("\r\x1b[K")
    }

    if err != nil || resp == nil {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sFailed to connect to api: %s%v%s\n",
            color.R, color.N, color.GG, err, color.N,
        )
        return nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sServer returned error code: %s%d%s\n",
            color.R, color.N, color.GG, resp.StatusCode, color.N,
        )
        return nil
    }

    scanner := bufio.NewScanner(resp.Body)
    var results []DNSResult

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "API count exceeded") {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sRate limit reached!\n",
                color.R, color.N,
            )
            return nil
        }

        if strings.TrimSpace(line) == "" {
            continue
        }

        parts := strings.SplitN(line, ":", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])

            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%s%s: %s%s%s)%s",
                color.B, color.N, color.GG, targetURL, color.DG,
                color.N, strings.ToLower(key), color.BB, value, color.DG,
                color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: key,
                    Value: value,
                },
            )
        } else {
            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%sinfo: %s%s%s)%s",
                color.B, color.N, color.GG, targetURL, color.DG,
                color.N, color.BB, line, color.DG, color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: "INFO",
                    Value: line,
                },
            )
        }
    }

    fmt.Printf("\r\x1b[K")
    return results
}

// Copyright (c) 2026 Senzdetta