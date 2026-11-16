// https://github.com/Senzdetta/Comet

package dnslookup

import (
    "fmt"
    "context"
    "time"
    "net"
    "strings"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func runNativeLookup(cfg Config) []DNSResult {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        time.Duration(cfg.Timeout)*time.Second,
    )
    defer cancel()

    var results []DNSResult

    customResolver := &net.Resolver{
        PreferGo: true,
        Dial: func(
            ctx context.Context,
            network, address string,
        ) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Duration(
                    cfg.Timeout,
                ) * time.Second,
            }

            return d.DialContext(
                ctx, network,
                cfg.DNSResolv,
            )
        },
    }

    if ips, err := customResolver.LookupIPAddr(
        ctx, cfg.Domain,
    ); err == nil {
        for _, ip := range ips {
            key := "A"
            if ip.IP.To4() == nil {
                key = "AAAA"
            }

            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%s%s: %s%s%s)%s",
                color.B, color.N, color.GG, cfg.Domain, color.DG,
                color.N, strings.ToLower(key), color.BB, ip.IP.String(), color.DG,
                color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: key,
                    Value: ip.IP.String(),
                },
            )
        }
    }

    if mxRecords, err := customResolver.LookupMX(
        ctx, cfg.Domain,
    ); err == nil {
        for _, mx := range mxRecords {
            val := fmt.Sprintf(
                "%d %s",
                mx.Pref, mx.Host,
            )

            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%smx pref: %s%d%s, %smx host: %s%s%s)%s",
                color.B, color.N, color.GG, cfg.Domain, color.DG,
                color.N, color.BB, mx.Pref, color.DG,
                color.N, color.BB, mx.Host, color.DG,
                color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: "MX",
                    Value: val,
                },
            )
        }
    }

    if nsRecords, err := customResolver.LookupNS(
        ctx, cfg.Domain,
    ); err == nil {
        for _, ns := range nsRecords {
            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%sns: %s%s%s)%s",
                color.B, color.N, color.GG, cfg.Domain, color.DG,
                color.N, color.BB, ns.Host, color.DG, color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: "NS",
                    Value: ns.Host,
                },
            )
        }
    }

    if txtRecords, err := customResolver.LookupTXT(
        ctx, cfg.Domain,
    ); err == nil {
        for _, txt := range txtRecords {
            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%stxt: %s%s%s)%s",
                color.B, color.N, color.GG, cfg.Domain, color.DG,
                color.N, color.BB, txt, color.DG, color.N,
            )

            results = append(
                results,
                DNSResult{
                    Type: "TXT",
                    Value: txt,
                },
            )
        }
    }

    if cname, err := customResolver.LookupCNAME(
        ctx, cfg.Domain,
    ); err == nil && cname != cfg.Domain && cname != "" {
        pcat.Printf(
            "%s[*] %sLookup: %s%s %s(%scname: %s%s%s)%s",
            color.B, color.N, color.GG, cfg.Domain, color.DG,
            color.N, color.BB, cname, color.DG, color.N,
        )

        results = append(
            results,
            DNSResult{
                Type: "CNAME",
                Value: cname,
            },
        )
    }

    if soaVal := fetchSOARecord(
        cfg.Domain, cfg.DNSResolv, cfg.Timeout,
    ); soaVal != "" {
        pcat.Printf(
            "%s[*] %sLookup: %s%s %s(%ssoa: %s%s%s)%s",
            color.B, color.N, color.GG, cfg.Domain, color.DG,
            color.N, color.BB, soaVal, color.DG, color.N,
        )

        results = append(
            results,
            DNSResult{
                Type: "SOA",
                Value: soaVal,
            },
        )
    }

    fmt.Printf("\r\x1b[K")
    return results
}

// Copyright (c) 2026 Senzdetta