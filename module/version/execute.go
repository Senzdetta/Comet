// https://github.com/Senzdetta/Comet

package version

import (
    "fmt"
    "github.com/Senzdetta/Comet/utils"
    "github.com/Senzdetta/Comet/utils/color"
)

const (
    name = "Comet"
    version = "v0.1.20261125"
    developer = "Senzdetta"
    homepage = "https://github.com/Senzdetta/Comet"
)

func Execute() {
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

    fmt.Printf(
        "%sInstall path: %s%s%s\n",
        color.N, color.GG, utils.Root(), color.N,
    )
}

// Copyright (c) 2026 Senzdetta