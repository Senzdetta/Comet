// https://github.com/Senzdetta/Comet

package pcat

func truncate(str string, maxWidth int) string {
    if maxWidth <= 3 {
        return "..."
    }

    if visualLength(str) <= maxWidth {
        return str
    }

    runes := []rune(str)
    visibleCount := 0
    inAnsi := false

    var result []rune

    for i := 0; i < len(runes); i++ {
        if runes[i] == '\x1b' {
            inAnsi = true
        }

        result = append(result, runes[i])

        if inAnsi {
            if runes[i] == 'm' {
                inAnsi = false
            }
            continue
        }

        visibleCount++
        if visibleCount >= maxWidth-3 {
            result = append(result, []rune("...")...)
            break
        }
    }

    return string(result)
}

// Copyright (c) 2026 Senzdetta