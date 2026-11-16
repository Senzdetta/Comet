// https://github.com/Senzdetta/Comet

package subdomain

import (
    "embed"
)

//go:embed wordlist/subdomains-top1million-20000.txt
var defaultWordlistFS embed.FS

// Copyright (c) 2026 Senzdetta