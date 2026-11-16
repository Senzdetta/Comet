// https://github.com/Senzdetta/Comet

package subtakeover

import (
    "bufio"
    "strings"
)

func parseWordlist(content string) []string {
    var words []string
    scanner := bufio.NewScanner(strings.NewReader(content))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        words = append(words, line)
    }
    return words
}

// Copyright (c) 2026 Senzdetta