// https://github.com/Senzdetta/Comet

package pcat

import (
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
)

func Printstat(label, value string) {
    fmt.Printf(
        "%s* %s%s: %s%s%s\n",
        color.DG, color.N, label, color.GG, value, color.N,
    )
}

// Copyright (c) 2026 Senzdetta