// https://github.com/Senzdetta/Comet

package portscan

import (
    "os"
    "errors"
    "fmt"
    "net"
    "strconv"
    "strings"
    "syscall"
    "time"
    "github.com/ishidawataru/sctp"
    "github.com/Senzdetta/Comet/utils/color"
)

func scanSCTP(
    ip string,
    port int,
    timeout time.Duration,
) (string, string) {
    raddr, err := sctp.ResolveSCTPAddr(
        "sctp",
        net.JoinHostPort(
            ip,
            strconv.Itoa(port),
        ),
    )

    if err != nil {
        return "filtered", ""
    }

    type res struct {
        c   *sctp.SCTPConn
        err error
    }

    ch := make(chan res, 1)
    go func() {
        c, err := sctp.DialSCTP(
            "sctp", nil, raddr,
        )
        ch <- res{c, err}
    }()

    select {
        case r := <-ch:
            if r.err != nil {
                if errors.Is(
                    r.err,
                    syscall.EPROTONOSUPPORT,
                ) || errors.Is(
                    r.err,
                    syscall.ESOCKTNOSUPPORT,
                ) || strings.Contains(
                    r.err.Error(),
                    "not supported",
                ) {
                    sctpWarn.Do(
                        func() {
                            fmt.Fprintf(
                                os.Stderr,
                                "%s[!] %sSTCP not supported by kernel %s(%smodprobe sctp%s)%s\n",
                                color.R, color.N, color.DG, color.GG, color.DG, color.N,
                            )
                        },
                    )
                    return "unsupported", ""
                }
                return classify(r.err), ""
            }
            r.c.Close()
            return "open", "unknown"
        case <-time.After(timeout):
            go func() {
                if r := <-ch; r.c != nil {
                    r.c.Close()
                }
            }()
            return "filtered", "unknown"
    }
}

// Copyright (c) 2026 Senzdetta