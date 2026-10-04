// https://github.com/Senzdetta/Comet

package main

import (
    "os"
    "strings"
    "github.com/Senzdetta/Comet/console"
)

func main() {
    args := os.Args[1:]
    input := strings.Join(args, " ")
    console.CometConsole(input)
}

// Copyright (c) 2026 Senzdetta