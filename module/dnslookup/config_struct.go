// https://github.com/Senzdetta/Comet

package dnslookup

type Config struct {
    Domain string
    Timeout int
    Retries int
    Threads int
    Mode string
    APIURL string
    DNSResolv string
}

// Copyright (c) 2026 Senzdetta