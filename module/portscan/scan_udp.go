// https://github.com/Senzdetta/Comet

package portscan

import (
    "syscall"
    "net"
    "strconv"
    "time"
    "errors"
)

func scanUDP(
    ip string,
    port int,
    timeout time.Duration,
) (string, string) {
    addr := net.JoinHostPort(ip, strconv.Itoa(port))
    conn, err := net.Dial("udp", addr)
    if err != nil {
        return "filtered", ""
    }
    defer conn.Close()

    probe, ok := udpProbes[port]
    if !ok {
        probe = []byte{0x00}
    }

    conn.SetDeadline(time.Now().Add(timeout))
    conn.Write(probe)

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err == nil {
        return "open", udpBanner(port, buf[:n])
    }

    if errors.Is(err, syscall.ECONNREFUSED) {
        return "closed", ""
    }

    return "open|filtered", ""
}

// Copyright (c) 2026 Senzdetta