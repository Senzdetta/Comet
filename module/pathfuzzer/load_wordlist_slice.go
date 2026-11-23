// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "strings"
)

func loadWordlistSlice(openWords string) ([]string, error) {
    scanner, closeFn, _, err := openWordlist(WordlistFS, openWords)
    if err != nil {
        return nil, err
    }
    defer closeFn()

    var words []string
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        if idx := strings.IndexAny(line, " \t"); idx != -1 {
            rest := strings.TrimSpace(line[idx:])
            if strings.HasPrefix(rest, "#") {
                line = strings.TrimSpace(line[:idx])
            }
        }

        w := strings.Trim(line, "/")
        if w != "" {
            words = append(words, w)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return words, nil
}

// Copyright (c) 2026 Senzdetta