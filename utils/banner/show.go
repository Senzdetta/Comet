// https://github.com/Senzdetta/Comet

package banner

import (
    "embed"
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
)

//go:embed banner.txt
var asciiFS embed.FS

func Show() {
    data, err := asciiFS.ReadFile("banner.txt")
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