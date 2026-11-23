// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "fmt"
    "strings"
    "sync"
    "net/http"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func fuzz(
    client *http.Client,
    target, directory string,
    idx int,
    recursive bool,
    total int,
    results *[]FuzzResult,
    dirs *[]string,
    mu *sync.Mutex,
) {
    directory = strings.Trim(
        strings.TrimSpace(directory),
        "/",
    )

    if directory == "" {
        return
    }

    url := fmt.Sprintf(
        "%s/%s",
        target, directory,
    )

    status, length, contentType, location, bodySample, ok := fetchResponseInfo(
        client, url,
    )

    if !ok {
        status = 0
        length = 0
    }

    mu.Lock()
    pcat.Printf(
        "%s[*] %sFuzzing %s(%s%d%s/%s%d%s)%s: %s%s%s/%s%s%s",
        color.B, color.N, color.DG, color.GG, idx, color.DG,
        color.GG, total, color.DG, color.N,
        color.GG, target, color.DG, color.BB, directory, color.N,
    )
    mu.Unlock()

    if status == 200 || status == 301 || status == 302 {
        entryType := detectTypeFromResponse(
            client, url, status,
            contentType, location, bodySample,
        )

        mu.Lock()
        *results = append(
            *results, FuzzResult{
                URL: url,
                Status: status,
                Type: entryType,
                Length: length,
            },
        )

        if recursive && entryType == "dir" {
            *dirs = append(*dirs, directory)
        }
        mu.Unlock()
    }
}

// Copyright (c) 2026 Senzdetta