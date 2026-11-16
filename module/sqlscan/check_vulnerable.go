// https://github.com/Senzdetta/Comet

package sqlscan

import (
    "strings"
)

func checkVulnerable(body string, statusCode int) (bool, string) {
    for _, sig := range sqlSignatures {
        if strings.Contains(body, sig) {
            return true, "Error-Based: " + sig
        }
    }

    successKeywords := []string{
        "Welcome",
        "Login success",
        "Logged in as",
        "Dashboard",
    }

    if statusCode == 200 {
        for _, keyword := range successKeywords {
            if strings.Contains(body, keyword) {
                return true, "Login Bypass / Keyword Match: " + keyword
            }
        }
    }

    return false, ""
}

// Copyright (c) 2026 Senzdetta