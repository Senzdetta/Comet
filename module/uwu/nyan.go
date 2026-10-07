// https://github.com/Senzdetta/Comet

package uwu

import (
    "fmt"
    "time"
)

func Nyan() {
    faces := []string{
        "(｡◕‿◕｡)",
        "(≧◡≦)",
        "ʕ•ᴥ•ʔ",
        "(・ω・)",
        "(๑˃ᴗ˂)ﻭ",
        "(ง'̀-'́)ง",
        "(=^･ω･^=)",
    }

    fixface := "(・ω・)"
    delay := 200 * time.Millisecond
    end := time.After(5 * time.Second)
    nyaa := 0

    fmt.Print("\x1b[?25l")
    for {
        select {
            case <-end:
                fmt.Printf(
                    "\r%s\x1b[K\x1b[?25h\n",
                    fixface,
                )
                return
            default:
                fmt.Printf(
                    "\r%s\x1b[K",
                    faces[nyaa%len(faces)],
                )
            time.Sleep(delay)
            nyaa++
        }
    }
}

// Copyright (c) 2026 Senzdetta