// https://github.com/Senzdetta/Comet

package help

import (
    "embed"
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
)

//go:embed banner/banner.txt
var asciiFS embed.FS

func banner() {
    data, err := asciiFS.ReadFile("banner/banner.txt")
    if err != nil {
        fmt.Printf(
            "%s[!] %sError loading banner!\n",
            color.R, color.N,
        )
        return
    }

    fmt.Printf(
        "%s%s%s\n",
        color.B, string(data), color.N,
    )
}

// Copyright (c) 2026 Senzdetta