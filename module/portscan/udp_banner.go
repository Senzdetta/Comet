// https://github.com/Senzdetta/Comet

package portscan

import (
    "strconv"
    "fmt"
)

func udpBanner(port int, b []byte) string {
    switch port {
        case 53:
            if len(b) >= 12 && b[0] == 0x12 && b[1] == 0x34 {
                rc := int(b[3] & 0x0f)
                name := strconv.Itoa(rc)

                if rc < len(dnsRcodes) {
                    name = dnsRcodes[rc]
                }

                an := int(b[6])<<8 | int(b[7])
                ns := int(b[8])<<8 | int(b[9])

                ra := "no recursion"
                if b[3]&0x80 != 0 {
                    ra = "recursion available"
                }

                return fmt.Sprintf(
                    "DNS response (%s, answers=%d, authority=%d, %s)",
                    name, an, ns, ra,
                )
            }
        case 123:
            if len(b) >= 48 {
                return fmt.Sprintf(
                    "NTP v%d mode %d stratum %d",
                    (b[0]>>3)&7, b[0]&7, b[1],
                )
            }
        case 161:
            if len(b) > 0 && b[0] == 0x30 {
                return "SNMP response"
            }
    }

    printable := 0
    for _, c := range b {
        if c >= 32 && c <= 126 {
            printable++
        }
    }

    if len(b) == 0 || printable*2 < len(b) {
        return fmt.Sprintf(
            "binary data (%d bytes)",
            len(b),
        )
    }

    return cleanBanner(b)
}

// Copyright (c) 2026 Senzdetta