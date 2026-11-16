// https://github.com/Senzdetta/Comet

package pcat

import (
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
)

func Printres[T any](items []T, formatter func(item T) string) {
    if len(items) == 0 {
        fmt.Printf(
            "%s[*] %sTotal found: %s0%s\n",
            color.B, color.N, color.GG, color.N,
        )
    } else {
        fmt.Printf(
            "%s[+] %sTotal found: %s%d%s\n",
            color.B, color.N, color.GG, len(items), color.N,
        )

        for _, item := range items {
            fmt.Printf(
                "%s* %s%s\n",
                color.DG, color.N, formatter(item),
            )
        }
    }
}

// Copyright (c) 2026 Senzdetta