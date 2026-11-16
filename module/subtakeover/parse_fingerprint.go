// https://github.com/Senzdetta/Comet

package subtakeover

import (
    "github.com/hashicorp/hcl/v2/hclsimple"
)

func parseFingerprint(content string) map[string]string {
    var config FingerprintConfig
    err := hclsimple.Decode(
        "fingerprint.hcl",
        []byte(content), nil, &config,
    )

    if err != nil {
        return make(map[string]string)
    }

    return config.Signatures
}

// Copyright (c) 2026 Senzdetta