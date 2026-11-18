// https://github.com/Senzdetta/Comet

package invinput

import (
    "fmt"
    "os"
    "github.com/Senzdetta/Comet/utils/color"
)

func MissingArgument() {
    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sMissing argument!\n",
        color.R, color.N,
    )

    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sTry: %scomet --help%s\n",
        color.R, color.N, color.GG, color.N,
    )
}

// Copyright (c) 2026 Senzdetta