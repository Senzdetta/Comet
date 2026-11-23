// https://github.com/Senzdetta/Comet

package pathfuzzer

import (
    "strings"
    "net/http"
)

func detectTypeFromResponse(
    client *http.Client,
    baseURL string,
    status int,
    contentType string,
    location string,
    bodySample string,
) string {
    if status == 301 || status == 302 || strings.HasSuffix(
        location, "/",
    ) {
        return "dir"
    }

    if isDirectoryListing(bodySample) {
        return "dir"
    }

    slashURL := baseURL + "/"
    slashStatus, _, _, _, slashBody, ok := fetchResponseInfo(
        client, slashURL,
    )

    if ok {
        if isDirectoryListing(slashBody) {
            return "dir"
        }

        if (slashStatus == 200 ||
            slashStatus == 301 ||
            slashStatus == 302 ||
            slashStatus == 403) {
                return "dir"
        }

        if slashStatus == 404 {
            return "file"
        }
    }

    if isAssetContentType(contentType) {
        return "file"
    }

    return "dir"
}

// Copyright (c) 2026 Senzdetta