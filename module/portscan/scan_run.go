// https://github.com/Senzdetta/Comet

package portscan

import (
    "time"
)

func scanRun(
    ip string,
    port int,
    proto string,
    timeout time.Duration,
) Result {
    r := Result{
        Port: port,
        Proto: proto,
    }

    switch proto {
        case "tcp":
            r.Status, r.Banner = scanTCP(
                ip, port, timeout,
            )
        case "udp":
            r.Status, r.Banner = scanUDP(
                ip, port, timeout,
            )
        case "sctp":
            r.Status, r.Banner = scanSCTP(
                ip, port, timeout,
            )
    }

    return r
}

// Copyright (c) 2026 Senzdetta