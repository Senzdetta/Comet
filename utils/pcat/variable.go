// https://github.com/Senzdetta/Comet

package pcat

import (
    "sync"
    "github.com/Senzdetta/Comet/utils/logger"
)

var (
    loggers = make(map[string]*logger.Logger)
    loggerMu sync.Mutex
)

// Copyright (c) 2026 Senzdetta