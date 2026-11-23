// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "strings"
)

func isAssetContentType(contentType string) bool {
    assetPrefixes := []string{
        "text/css",
        "application/javascript",
        "text/javascript",
        "image/",
        "audio/",
        "video/",
        "font/",
        "application/octet-stream",
        "application/pdf",
        "application/zip",
    }

    for _, p := range assetPrefixes {
        if strings.HasPrefix(contentType, p) {
            return true
        }
    }

    return false
}

// Copyright (c) 2026 Senzdetta