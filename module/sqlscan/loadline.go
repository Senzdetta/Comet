// https://github.com/Senzdetta/Comet

package sqlscan

import (
    "bufio"
    "strings"
    "io"
)

func loadLine(reader io.Reader) ([]string, error) {
    unique := make(map[string]bool)
    var items []string

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        if !unique[line] {
            unique[line] = true
            items = append(items, line)
        }
    }

    return items, scanner.Err()
}

// Copyright (c) 2026 Senzdetta