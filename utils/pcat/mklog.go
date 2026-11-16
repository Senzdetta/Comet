// https://github.com/Senzdetta/Comet

package pcat

func Mklog(name, tag, message string) {
    log := getLogger(name)
    log.Log(tag, message)
}

// Copyright (c) 2026 Senzdetta