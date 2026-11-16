// https://github.com/Senzdetta/Comet

package pcat

import (
    "github.com/Senzdetta/Comet/utils/logger"
)

func getLogger(name string) *logger.Logger {
    loggerMu.Lock()
    defer loggerMu.Unlock()

    if l, exists := loggers[name]; exists {
        return l
    }

    l := logger.NewLogger(name)
    loggers[name] = l

    return l
}

// Copyright (c) 2026 Senzdetta