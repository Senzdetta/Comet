// https://github.com/Senzdetta/Comet

package console

import (
    "github.com/Senzdetta/Comet/module/helper"
)

type Helper struct{}
func (c Helper) Execute(args []string) {
    helper.CometHelper()
}

// Copyright (c) 2026 Senzdetta