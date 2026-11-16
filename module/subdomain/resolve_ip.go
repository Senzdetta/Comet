// https://github.com/Senzdetta/Comet

package subdomain

import (
    "context"
    "strings"
    "net"
    "time"
)

func resolveIP(ctx context.Context, host, dnsAddr string) []string {
    if !strings.Contains(dnsAddr, ":") {
        dnsAddr = dnsAddr + ":53"
    }

    r := &net.Resolver{
        PreferGo: true,
        Dial: func(
            ctx context.Context,
            network, address string,
        ) (net.Conn, error) {
            d := net.Dialer{
                Timeout: 2 * time.Second,
            }
            return d.DialContext(ctx, "udp", dnsAddr)
        },
    }

    ips, err := r.LookupHost(ctx, host)
    if err != nil {
        return nil
    }

    return ips
}

// Copyright (c) 2026 Senzdetta