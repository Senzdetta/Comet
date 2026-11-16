// https://github.com/Senzdetta/Comet

package paramscan

func Execute(opts Options) {
    if opts.Mode == "api" {
        runAPIMode(opts)
    } else {
        runNativeMode(opts)
    }
}

// Copyright (c) 2026 Senzdetta