// https://github.com/Senzdetta/Comet

package selfuninstall

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/vcore"
)

func Execute() {
    scriptPath := filepath.Join(
        vcore.Root(), "uninstall.sh",
    )

    info, err := os.Stat(scriptPath)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Fprintf(
                os.Stderr,
                "%s[!] %sFile: %s%s %snot found!\n",
                color.R, color.N, color.GG, scriptPath, color.N,
            )
            return
        }

        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sError accessing uninstall script: %s%s%s\n",
            color.R, color.N, color.GG, err, color.N,
        )

        return
    }

    if info.IsDir() {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sPath: %s%s %sis a directory!\n",
            color.R, color.N, color.GG, scriptPath, color.N,
        )
        return
    }

    cmd := exec.Command("bash", scriptPath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Fprintf(
            os.Stderr,
            "%s[!] %sUninstall failed: %s%s%s\n",
            color.R, color.N, color.GG, err, color.N,
        )
        return
    }
}

// Copyright (c) 2026 Senzdetta