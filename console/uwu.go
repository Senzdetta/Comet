// https://github.com/Senzdetta/Comet

package console

import (
    "time"
    "fmt"
    "github.com/Senzdetta/Comet/utils/cursor"
    "github.com/Senzdetta/Comet/module/uwu"
)

type Uwu struct{}
func (c Uwu) Execute(args []string) {
    cursor.Hide()
    uwu.Nyan(5 * time.Second)
    cursor.Visible()

    fmt.Println()
}

// Copyright (c) 2026 Senzdetta