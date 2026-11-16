// https://github.com/Senzdetta/Comet

package pcat

func visualLength(str string) int {
    cleanStr := ansiRegex.ReplaceAllString(str, "")
    return len(cleanStr)
}

// Copyright (c) 2026 Senzdetta