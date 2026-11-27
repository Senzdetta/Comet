// https://github.com/Senzdetta/Comet

package invinput

import (
    "fmt"
    "os"
    "strings"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/vcore"
)

func MissingArgument() {
    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sMissing argument!\n",
        color.R, color.N,
    )

    fmt.Fprintf(
        os.Stderr,
        "%s[!] %sTry: %s%s --help%s\n",
        color.R, color.N, color.GG, strings.ToLower(vcore.Name), color.N,
    )
}

// Copyright (c) 2026 Senzdetta