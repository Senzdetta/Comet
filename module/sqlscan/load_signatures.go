// https://github.com/Senzdetta/Comet

package sqlscan

import (
    "os"
    "fmt"
    "io/fs"
)

func LoadSignatures(
    path string,
    embeddedFS fs.FS,
) ([]string, error) {
    if path == "" {
        file, err := embeddedFS.Open(
            "signature/signatures.txt",
        )

        if err != nil {
            return nil, err
        }

        defer file.Close()

        return loadLine(file)
    }

    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if info.IsDir() {
        return nil, fmt.Errorf(
            "signature path %s is a directory",
            path,
        )
    }

    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    return loadLine(file)
}

// Copyright (c) 2026 Senzdetta