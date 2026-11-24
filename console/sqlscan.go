// https://github.com/Senzdetta/Comet

package console

import (
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/sqlscan"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Sqlscan struct{}
func (c Sqlscan) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    targetURL := ""
    wordlist := ""
    signature := ""
    threads := 50
    timeout := 15

    for i := 2; i < len(args); i++ {
        switch args[i] {
            case "--payload":
                if i+1 < len(args) {
                    wordlist = args[i+1]
                    i++
                }
            case "--signature":
                if i+1 < len(args) {
                    signature = args[i+1]
                    i++
                }
            case "--threads":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil {
                        threads = t
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    t, err := strconv.Atoi(args[i+1])
                    if err == nil {
                        timeout = t
                    }
                    i++
                }
            default:
                if targetURL == "" && args[i][0] != '-' {
                    targetURL = args[i]
                }
        }
    }

    if targetURL == "" {
        invinput.MissingArgument()
        os.Exit(1)
    }

    sqlscan.Execute(
        targetURL,
        wordlist,
        signature,
        threads,
        timeout,
    )
}

// Copyright (c) 2026 Senzdetta