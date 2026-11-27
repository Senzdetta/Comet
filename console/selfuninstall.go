// https://github.com/Senzdetta/Comet

package console

import (
    "github.com/Senzdetta/Comet/module/selfuninstall"
)

type Selfuninstall struct{}
func (c Selfuninstall) Execute(args []string) {
    selfuninstall.Execute()
}

// Copyright (c) 2026 Senzdetta