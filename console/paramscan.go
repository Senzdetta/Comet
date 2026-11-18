// https://github.com/Senzdetta/Comet

package console

import (
    "fmt"
    "os"
    "strconv"
    "github.com/Senzdetta/Comet/module/paramscan"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/invinput"
)

type Paramscan struct{}
func (c Paramscan) Execute(args []string) {
    if len(args) < 3 {
        invinput.MissingArgument()
        os.Exit(1)
    }

    target := args[2]
    timeout := 10
    threads := 100
    retry := 5
    fuzz := false
    mode := "native"
    customAPI := "http://web.archive.org/cdx/search/cdx?url={target}/*&output=txt&fl=original&collapse=urlkey"
    wordlistPath := ""
    valuePath := ""

    for i := 3; i < len(args); i++ {
        switch args[i] {
            case "--mode":
                if i+1 < len(args) {
                    mode = args[i+1]
                    i++
                }
            case "--fuzz":
                fuzz = true
            case "--threads":
                if i+1 < len(args) {
                    if t, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        threads = t
                    }
                    i++
                }
            case "--timeout":
                if i+1 < len(args) {
                    if t, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        timeout = t
                    }
                    i++
                }
            case "--retry":
                if i+1 < len(args) {
                    if r, err := strconv.Atoi(
                        args[i+1],
                    ); err == nil {
                        retry = r
                    }
                    i++
                }
            case "--api":
                if i+1 < len(args) {
                    customAPI = args[i+1]
                    i++
                }
            case "--wordlist":
                if i+1 < len(args) {
                    wordlistPath = args[i+1]
                    i++
                }
            case "--valueword":
                if i+1 < len(args) {
                    valuePath = args[i+1]
                    i++
                }
        }
    }

    if mode != "native" && mode != "api" {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid mode: %s%s%s\n",
            color.R, color.N, color.GG, mode, color.N,
        )
        os.Exit(1)
    }

    if threads <= 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid threads value!\n",
            color.R, color.N,
        )
        os.Exit(1)
    }

    if timeout <= 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid timeout value!\n",
            color.R, color.N,
        )
        os.Exit(1)
    }

    if retry <= 0 {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sInvalid retry value!\n",
            color.R, color.N,
        )
        os.Exit(1)
    }

    paramscan.Execute(
        paramscan.Options{
            Target: target,
            Mode: mode,
            Threads: threads,
            Timeout: timeout,
            Retry: retry,
            Fuzz: fuzz,
            CustomAPI: customAPI,
            WordlistPath: wordlistPath,
            ValuePath: valuePath,
        },
    )
}

// Copyright (c) 2026 Senzdetta