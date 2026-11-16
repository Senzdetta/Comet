// https://github.com/Senzdetta/Comet

package subdomain

import (
    "os"
    "strings"
)

func loadWordlist(path string) ([]string, error) {
    var data []byte
    var err error

    if path == "" {
        data, err = defaultWordlistFS.ReadFile(
            "wordlist/subdomains-top1million-20000.txt",
        )
    } else {
        data, err = os.ReadFile(path)
    }

    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(data), "\n")
    var words []string

    seen := make(map[string]bool)
    for _, l := range lines {
        l = strings.TrimSpace(l)
        if l != "" && !strings.HasPrefix(l, "#") {
            if !seen[l] {
                words = append(words, l)
                seen[l] = true
            }
        }
    }

    return words, nil
}

// Copyright (c) 2026 Senzdetta