// https://github.com/Senzdetta/Comet

package portscan

import (
    "os"
    "strings"
)

func loadServices(servicePath string) {
    serviceMap = make(map[string]Service)
    var content string

    if servicePath != "" {
        if fileBytes, err := os.ReadFile(servicePath); err == nil {
            content = string(fileBytes)
        } else {
            content = nmapServices
        }
    } else {
        content = nmapServices
    }

    for _, line := range strings.Split(content, "\n") {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        desc := ""
        if i := strings.Index(line, "#"); i >= 0 {
            desc = strings.TrimSpace(line[i+1:])
            line = line[:i]
        }

        f := strings.Fields(line)
        if len(f) < 2 {
            continue
        }

        key := strings.ToLower(f[1])
        if _, ok := serviceMap[key]; ok {
            continue
        }

        s := Service{
            Name: f[0],
            Prob: "unknown",
            Desc: desc,
        }

        if len(f) >= 3 {
            s.Prob = f[2]
        }

        if s.Desc == "" {
            s.Desc = "unknown"
        }

        serviceMap[key] = s
    }
}

// Copyright (c) 2026 Senzdetta