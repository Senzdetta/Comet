// https://github.com/Senzdetta/Comet

package invinput

import (
    "fmt"
    "os"
    "github.com/Senzdetta/Comet/utils/color"
)

func InvalidOption(input string) {
    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sInvalid option: %s%s%s\n",
        color.R, color.N, color.GG, input, color.N,
    )

    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sTry: %scomet --help%s\n",
        color.R, color.N, color.GG, color.N,
    )
}

// Copyright (c) 2026 Senzdetta