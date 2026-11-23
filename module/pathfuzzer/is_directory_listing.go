// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "strings"
)

func isDirectoryListing(bodySample string) bool {
    return strings.Contains(
        bodySample, "index of /",
    ) || strings.Contains(
        bodySample, "parent directory",
    )
}

// Copyright (c) 2026 Senzdetta