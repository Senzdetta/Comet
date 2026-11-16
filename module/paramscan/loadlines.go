// https://github.com/Senzdetta/Comet

package paramscan

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "github.com/Senzdetta/Comet/utils/color"
)

func loadLines(customPath string, embeddedPath string) []string {
    var scanner *bufio.Scanner
    if customPath != "" {
        if _, err := os.Stat(customPath); err == nil {
            file, err := os.Open(customPath)
            if err != nil {
                fmt.Printf(
                    "%s[!] %sFailed to open file %s: %v\n",
                    color.R, color.N, customPath, err,
                )
                return nil
            }
            defer file.Close()
            scanner = bufio.NewScanner(file)
        } else {
            scanner = bufio.NewScanner(
                strings.NewReader(customPath),
            )
        }
    } else {
        content, err := embeddedFiles.ReadFile(embeddedPath)
        if err != nil {
            return nil
        }

        scanner = bufio.NewScanner(
            strings.NewReader(string(content)),
        )
    }

    var lines []string
    seen := make(map[string]bool)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        if !seen[line] {
            lines = append(lines, line)
            seen[line] = true
        }
    }

    return lines
}

// Copyright (c) 2026 Senzdetta