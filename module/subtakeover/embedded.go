// https://github.com/Senzdetta/Comet

package subtakeover

import (
    _ "embed"
)

//go:embed wordlist/seclists/subdomains-top1million-20000.txt
var dataWordlist string

//go:embed wordlist/fingerprint.hcl
var dataFingerprint string

// Copyright (c) 2026 Senzdetta