// https://github.com/Senzdetta/Comet

package invinput

import (
    "fmt"
    "os"
    "strings"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/vcore"
)

func InvalidOption(input string) {
    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sInvalid option: %s%s%s\n",
        color.R, color.N, color.GG, input, color.N,
    )

    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sTry: %s%s --help%s\n",
        color.R, color.N, color.GG, strings.ToLower(vcore.Name), color.N,
    )
}

// Copyright (c) 2026 Senzdetta