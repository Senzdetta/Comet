// https://github.com/Senzdetta/Comet

package console

import (
    "github.com/Senzdetta/Comet/module/help"
)

type Help struct{}
func (c Help) Execute(args []string) {
    help.Execute()
}

// Copyright (c) 2026 Senzdetta