// https://github.com/Senzdetta/Comet

package header

import (
    "os"
    "fmt"
    "strings"
    "time"
    "net"
    "crypto/tls"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(
    targetURL string,
    timeout time.Duration,
    followRedirect bool,
) {
    if !strings.HasPrefix(
        targetURL, "http://",
    ) && !strings.HasPrefix(
        targetURL, "https://",
    ) {
        targetURL = "https://" + targetURL
    }

    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(targetURL),
    )

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(
            fmt.Sprintf(
                "%v",
                timeout,
            ),
        ),
    )

    pcat.Printstat(
        "Redirect",
        pcat.Disconf(
            fmt.Sprintf(
                "%t",
                followRedirect,
            ),
        ),
    )

    fmt.Println()

    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout: timeout,
        }).DialContext,
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }

    client := &http.Client{
        Timeout: timeout,
        Transport: transport,
    }

    if !followRedirect {
        client.CheckRedirect = func(
            req *http.Request,
            via []*http.Request,
        ) error {
            return http.ErrUseLastResponse
        }
    }

    startTime := time.Now()
    req, _ := http.NewRequest("HEAD", targetURL, nil)
    req.Header.Set(
        "User-Agent",
        "https://github.com/Senzdetta/Comet",
    )

    resp, err := client.Do(req)
    if err != nil {
        req.Method = "GET"
        resp, err = client.Do(req)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sError connecting to target: %s%v%s\n",
                color.R, color.N, color.GG, err, color.N,
            )
            return
        }
        defer resp.Body.Close()
    } else {
        defer resp.Body.Close()
    }

    latency := time.Since(startTime)

    ipAddr := "Unknown"
    if addr, err := net.LookupHost(
        resp.Request.URL.Hostname(),
    ); err == nil && len(addr) > 0 {
        ipAddr = addr[0]
    }

    var results []HeaderResult

    results = append(
        results,
        HeaderResult{
            Section: "Common info",
            Item: "Status",
            Value: fmt.Sprintf(
                "%d (%s)",
                resp.StatusCode,
                resp.Status,
            ),
        },
    )

    results = append(
        results,
        HeaderResult{
            Section: "Common info",
            Item: "Target IP",
            Value: ipAddr,
        },
    )

    results = append(
        results,
        HeaderResult{
            Section: "Common info",
            Item: "Latency",
            Value: latency.Round(
                time.Millisecond,
            ).String(),
        },
    )

    if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
        cert := resp.TLS.PeerCertificates[0]
        daysLeft := int(
            time.Until(
                cert.NotAfter,
            ).Hours() / 24,
        )

        results = append(
            results,
            HeaderResult{
                Section: "SSL/TLS cert info",
                Item: "SSL/TLS issuer",
                Value: cert.Issuer.CommonName,
            },
        )

        results = append(
            results,
            HeaderResult{
                Section: "SSL/TLS cert info",
                Item: "Cipher",
                Value: tls.CipherSuiteName(
                    resp.TLS.CipherSuite,
                ),
            },
        )

        if daysLeft < 0 {
            results = append(
                results,
                HeaderResult{
                    Section: "SSL/TLS cert info",
                    Item: "SSL expiry",
                    Value: "EXPIRED",
                },
            )
        } else {
            results = append(
                results,
                HeaderResult{
                    Section: "SSL/TLS cert info",
                    Item: "SSL expiry",
                    Value: fmt.Sprintf(
                        "%d days remaining (%s)",
                        daysLeft,
                        cert.NotAfter.Format("2006-01-02"),
                    ),
                },
            )
        }
    } else {
        results = append(
            results,
            HeaderResult{
                Section: "SSL/TLS cert info",
                Item: "SSL/TLS",
                Value: "None (HTTP Plain)",
            },
        )
    }

    securityHeaders := []string{
        "Strict-Transport-Security",
        "X-Frame-Options",
        "X-Content-Type-Options",
        "Content-Security-Policy",
        "X-XSS-Protection",
        "Referrer-Policy",
        "Permissions-Policy",
    }

    for _, sh := range securityHeaders {
        val := resp.Header.Get(sh)
        if val != "" {
            results = append(
                results,
                HeaderResult{
                    Section: "Security headers info",
                    Item: sh,
                    Value: val,
                },
            )
        } else {
            results = append(
                results,
                HeaderResult{
                    Section: "Security headers info",
                    Item: sh,
                    Value: "Unknown",
                },
            )
        }
    }

    cookies := resp.Cookies()
    for _, cookie := range cookies {
        var flags []string
        if cookie.HttpOnly {
            flags = append(flags, "HttpOnly")
        }

        if cookie.Secure {
            flags = append(flags, "Secure")
        }

        if cookie.SameSite != http.SameSiteDefaultMode {
            flags = append(
                flags,
                fmt.Sprintf(
                    "SameSite=%d",
                    cookie.SameSite,
                ),
            )
        }

        flagStr := "None"
        if len(flags) > 0 {
            flagStr = strings.Join(flags, ", ")
        }

        results = append(
            results,
            HeaderResult{
                Section: fmt.Sprintf(
                    "Detected Cookies (%d)",
                    len(cookies),
                ),
                Item: cookie.Name,
                Value: fmt.Sprintf(
                    "[Flags: %s]",
                    flagStr,
                ),
            },
        )
    }

    for key, values := range resp.Header {
        valueStr := strings.Join(values, ", ")
        results = append(
            results, HeaderResult{
                Section: "Raw headers",
                Item: key,
                Value: valueStr,
            },
        )
    }

    currentSection := ""
    var sectionItems []HeaderResult

    printSection := func(sec string, items []HeaderResult) {
        if len(items) == 0 {
            return
        }

        fmt.Printf(
            "%s--------------- %s%s %s---------------%s\n",
            color.DG, color.GG, sec, color.DG, color.N,
        )

        pcat.Printres(
            items, func(r HeaderResult) string {
                return fmt.Sprintf(
                    "%s%s: %s%s%s",
                    color.N, r.Item,
                    color.GG, r.Value, color.N,
                )
            },
        )

        fmt.Println()
    }

    for _, r := range results {
        if r.Section != currentSection {
            if currentSection != "" {
                printSection(currentSection, sectionItems)
                sectionItems = nil
            }
            currentSection = r.Section
        }
        sectionItems = append(sectionItems, r)
    }

    printSection(currentSection, sectionItems)
}

// Copyright (c) 2026 Senzdetta