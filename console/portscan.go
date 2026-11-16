// https://github.com/Senzdetta/Comet

package console

import (
    "fmt"
    "os"
    "strings"
    "github.com/Senzdetta/Comet/module/portscan"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Portscan struct{}
func (c Portscan) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    ip := args[2]
    start := 1
    end := 65535
    protocol := "tcp,udp,sctp"
    only := ""
    threads := "100"
    timeout := "500"
    servicePath := ""
    skipHostCheck := false

    rest := args[3:]
    if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
        st, en, err := portscan.ParseRange(rest[0])
        if err != nil {
            fmt.Printf(
                "%s[!] %sError parsing range: %s%v%s\n",
                color.R, color.N, color.GG, err, color.N,
            )
            os.Exit(1)
        }

        if st > 0 {
            start = st
        }

        if en > 0 {
            end = en
        }

        rest = rest[1:]
    }

    for i := 0; i < len(rest); i++ {
        switch rest[i] {
            case "--skip-host-check":
                skipHostCheck = true
            case "--protocol":
                if i+1 >= len(rest) {
                    invinput.MissingArgument()
                    os.Exit(1)
                }
                i++
                protocol = rest[i]
            case "--only":
                if i+1 >= len(rest) {
                    invinput.MissingArgument()
                    os.Exit(1)
                }
                i++
                only = rest[i]
            case "--threads":
                if i+1 >= len(rest) {
                    invinput.MissingArgument()
                    os.Exit(1)
                }
                i++
                threads = rest[i]
            case "--timeout":
                if i+1 >= len(rest) {
                    invinput.MissingArgument()
                    os.Exit(1)
                }
                i++
                timeout = rest[i]
            case "--service":
                if i+1 >= len(rest) {
                    invinput.MissingArgument()
                    os.Exit(1)
                }
                i++
                servicePath = rest[i]
            default:
                invinput.InvalidOption(
                    fmt.Sprintf(
                        "--portscan %s",
                        rest[i],
                    ),
                )
                os.Exit(1)
        }
    }

    portscan.Execute(
        ip,
        start,
        end,
        protocol,
        only,
        threads,
        timeout,
        servicePath,
        skipHostCheck,
    )
}

// Copyright (c) 2026 Senzdetta