// https://github.com/Senzdetta/Comet

package version

import (
    "fmt"
    "github.com/Senzdetta/Comet/utils/vcore"
    "github.com/Senzdetta/Comet/utils/color"
)

func Execute() {
    fmt.Printf(
        "%s- %s%s %s-%s\n",
        color.DG, color.GG, vcore.Name, color.DG, color.N,
    )

    fmt.Printf(
        "%sVersion: %s%s%s\n",
        color.N, color.GG, vcore.Version, color.N,
    )

    fmt.Printf(
        "%sDeveloper: %s%s%s\n",
        color.N, color.GG, vcore.Developer, color.N,
    )

    fmt.Printf(
        "%sHomepage: %s%s%s\n",
        color.N, color.GG, vcore.Homepage, color.N,
    )

    fmt.Printf(
        "%sLanguage: %s%s%s\n",
        color.N, color.GG, vcore.Langver(), color.N,
    )

    fmt.Printf(
        "%sInstall path: %s%s%s\n",
        color.N, color.GG, vcore.Root(), color.N,
    )

    fmt.Printf(
        "%sExecutable path: %s%s%s\n",
        color.N, color.GG, vcore.Exepath(), color.N,
    )

    fmt.Printf(
        "%sOS System: %s%s%s\n",
        color.N, color.GG, vcore.OS(), color.N,
    )

    fmt.Printf(
        "%sOS Arch: %s%s%s\n",
        color.N, color.GG, vcore.Arch(), color.N,
    )
}

// Copyright (c) 2026 Senzdetta