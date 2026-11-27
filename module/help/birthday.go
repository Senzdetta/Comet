// https://github.com/Senzdetta/Comet

package help

import (
    "fmt"
    "time"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/vcore"
)

func birthday() {
    birthDate := "03-22"
    now := time.Now().Format("01-02")
    if now == birthDate {
        fmt.Printf(
            "%s› %sHappy birthday for %s%s %s🎉\n",
            color.R, color.N, color.GG, vcore.Name, color.N,
        )
        fmt.Println()
    }
}

// Copyright (c) 2026 Senzdetta