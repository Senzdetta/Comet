// https://github.com/Senzdetta/Comet

package dnslookup

import (
    "os"
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func Execute(cfg Config) {
    fmt.Printf("\x1b[?25l")
    defer fmt.Printf("\x1b[?25h")

    fmt.Printf(
        "%s[*] %sConfigurations:\n",
        color.B, color.N,
    )

    pcat.Printstat(
        "Target",
        pcat.Disconf(cfg.Domain),
    )

    pcat.Printstat(
        "Threads",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                cfg.Threads,
            ),
        ),
    )

    pcat.Printstat(
        "Timeout",
        pcat.Disconf(
            fmt.Sprintf(
                "%d",
                cfg.Timeout,
            ),
        ),
    )

    if cfg.Mode == "native" {
        pcat.Printstat(
            "DNS Resolver",
            pcat.Disconf(cfg.DNSResolv),
        )
    }

    pcat.Printstat(
        "Mode",
        pcat.Disconf(cfg.Mode),
    )

    if cfg.Mode == "api" {
        pcat.Printstat(
            "API Endpoint",
            pcat.Disconf(cfg.APIURL),
        )

        pcat.Printstat(
            "Retries",
            pcat.Disconf(
                fmt.Sprintf(
                    "%d",
                    cfg.Retries,
                ),
            ),
        )
    }

    fmt.Println()

    var results []DNSResult

    if cfg.Mode == "native" {
        results = runNativeLookup(cfg)
    } else if cfg.Mode == "api" {
        results = runAPILookup(cfg)
    } else {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid mode: %s%s%s\n",
            color.R, color.N, color.GG, cfg.Mode, color.N,
        )
        return
    }

    if len(results) > 0 {
        for _, r := range results {
            pcat.Mklog(
                "dnslookup",
                ":",
                fmt.Sprintf(
                    "%s: %s",
                    r.Type, r.Value,
                ),
            )
        }
    }

    pcat.Printres(
        results, func(r DNSResult) string {
            return fmt.Sprintf(
                "%s%s %s(%s%s: %s%s%s)%s",
                color.GG, cfg.Domain, color.DG,
                color.N, r.Type, color.GG, r.Value, color.DG,
                color.N,
            )
        },
    )
}

// Copyright (c) 2026 Senzdetta