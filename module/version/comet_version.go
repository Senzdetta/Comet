// https://github.com/Senzdetta/Comet

package version

import (
    "fmt"
    "github.com/Senzdetta/Comet/utils/color"
)

const (
    name = "Comet"
    version = "v0.1.20261007"
    developer = "Senzdetta"
    homepage = "https://github.com/Senzdetta/Comet"
)

func CometVersion() {
    fmt.Printf(
        "%s- %s%s %s-%s\n",
        color.DG, color.GG, name, color.DG, color.N,
    )

    fmt.Printf(
        "%sVersion: %s%s%s\n",
        color.N, color.GG, version, color.N,
    )

    fmt.Printf(
        "%sDeveloper: %s%s%s\n",
        color.N, color.GG, developer, color.N,
    )

    fmt.Printf(
        "%sHomepage: %s%s%s\n",
        color.N, color.GG, homepage, color.N,
    )
}

// Copyright (c) 2026 Senzdetta