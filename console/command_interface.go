// https://github.com/Senzdetta/Comet

package console

import (
    "embed"
)

//go:embed wordlist/*
var WordlistFS embed.FS

type Command interface {
    Execute(args []string)
}

// Copyright (c) 2026 Senzdetta