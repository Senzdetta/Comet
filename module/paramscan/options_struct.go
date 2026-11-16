// https://github.com/Senzdetta/Comet

package paramscan

type Options struct {
    Target string
    Mode string
    Threads int
    Timeout int
    Retry int
    Fuzz bool
    CustomAPI string
    WordlistPath string
    ValuePath string
}

// Copyright (c) 2026 Senzdetta