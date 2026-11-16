// https://github.com/Senzdetta/Comet

package portscan

import (
    "strconv"
)

func lookupService(
    port int,
    proto string,
    servicePath string,
) Service {
    serviceOnce.Do(
        func() {
            loadServices(servicePath)
        },
    )

    key := strconv.Itoa(port) + "/" + proto
    if s, ok := serviceMap[key]; ok {
        return s
    }

    return Service{
        Name: "unknown",
        Prob: "unknown",
        Desc: "unknown",
    }
}

// Copyright (c) 2026 Senzdetta