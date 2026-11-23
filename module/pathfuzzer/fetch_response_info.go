// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "io"
    "strings"
    "net/http"
)

func fetchResponseInfo(client *http.Client, rawURL string) (
    status int,
    length int64,
    contentType string,
    location string,
    bodySample string,
    ok bool,
) {
    req, err := http.NewRequest("GET", rawURL, nil)
    if err != nil {
        return 0, 0, "", "", "", false
    }

    resp, err := client.Do(req)
    if err != nil {
        return 0, 0, "", "", "", false
    }
    defer resp.Body.Close()

    status = resp.StatusCode
    contentType = strings.ToLower(
        resp.Header.Get("Content-Type"),
    )
    location = resp.Header.Get("Location")

    sampleBuf := make([]byte, 2048)
    nSample, _ := io.ReadFull(resp.Body, sampleBuf)

    if nSample < 0 {
        nSample = 0
    }

    bodySample = strings.ToLower(string(sampleBuf[:nSample]))

    var extraBytes int64
    if resp.ContentLength >= 0 {
        length = resp.ContentLength
    } else {
        extraBytes, _ = io.Copy(io.Discard, resp.Body)
        length = int64(nSample) + extraBytes
    }

    return status,
        length,
        contentType,
        location,
        bodySample,
        true
}

// Copyright (c) 2026 Senzdetta