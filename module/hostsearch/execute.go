// https://github.com/Senzdetta/Comet

package hostsearch

import (
    "os"
    "bufio"
    "fmt"
    "net"
    "strings"
    "time"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(domain, apiURL string, timeout int, retries int) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sWarning: this is a deprecated feature, use %s--subdomain %sinstead.\n",
        color.YY, color.N, color.GG, color.N,
    )

    fmt.Println()

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(domain),
    )

    pcat.Printstat(
        "API Endpoint",
        pcat.Disconf(apiURL),
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

    pcat.Printstat(
        "Retries", pcat.Disconf(
            fmt.Sprintf(
                "%d",
                retries,
            ),
        ),
    )

    fmt.Println()

    targetURL := apiURL
    if !strings.Contains(targetURL, "=") {
        targetURL = fmt.Sprintf(
            "%s?q=%s",
            strings.TrimSuffix(apiURL, "/"),
            domain,
        )
    } else if !strings.Contains(targetURL, domain) {
        targetURL = fmt.Sprintf(
            "%s%s",
            apiURL, domain,
        )
    }

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
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

    for i := 0; i < retries; i++ {
        req, _ := http.NewRequest("GET", targetURL, nil)
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
            color.CC, retries, color.DG, color.N,
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
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sServer returned error code: %s%d%s\n",
            color.R, color.N, color.GG, resp.StatusCode, color.N,
        )
        return
    }

    var list []HostResult

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "API count exceeded") {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sRate limit reached!\n",
                color.R, color.N,
            )
            return
        }

        parts := strings.Split(line, ",")
        if len(parts) == 2 {
            sub := strings.TrimSpace(parts[0])
            ip := strings.TrimSpace(parts[1])

            list = append(
                list, HostResult{
                    Subdomain: sub,
                    IPAddress: ip,
                },
            )
        }
    }

    if len(list) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sNo data found or target invalid!\n",
            color.R, color.N,
        )
        return
    }

    for _, item := range list {
        pcat.Mklog(
            "hostsearch",
            ":",
            fmt.Sprintf(
                "Found: %s:%s",
                item.Subdomain, item.IPAddress,
            ),
        )
    }

    pcat.Printres(
        list, func(r HostResult) string {
            return fmt.Sprintf(
                "%s%s %s(%s%s%s)%s",
                color.GG, r.Subdomain, color.DG,
                color.CC, r.IPAddress, color.DG, color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta