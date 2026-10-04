// https://github.com/Senzdetta/Comet

package console

import (
    "github.com/Senzdetta/Comet/module/version"
)

type Version struct{}
func (c Version) Execute(args []string) {
    version.CometVersion()
}

// Copyright (c) 2026 Senzdetta