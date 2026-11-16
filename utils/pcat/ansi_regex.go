// https://github.com/Senzdetta/Comet

package pcat

import (
    "regexp"
)

var ansiRegex = regexp.MustCompile(
    `\x1b\[[0-9;]*m`,
)

// Copyright (c) 2026 Senzdetta