// https://github.com/Senzdetta/Comet

package portscan

import (
    "net"
    "strconv"
    "time"
)

func scanTCP(
    ip string,
    port int,
    timeout time.Duration,
) (string, string) {
    addr := net.JoinHostPort(ip, strconv.Itoa(port))
    conn, err := net.DialTimeout("tcp", addr, timeout)
    if err != nil {
        return classify(err), ""
    }
    defer conn.Close()

    buf := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(timeout))
    n, _ := conn.Read(buf)

    if n == 0 {
        conn.SetDeadline(time.Now().Add(timeout))
        conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
        n, _ = conn.Read(buf)
    }

    return "open", cleanBanner(buf[:n])
}

// Copyright (c) 2026 Senzdetta