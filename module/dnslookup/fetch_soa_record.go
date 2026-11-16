// https://github.com/Senzdetta/Comet

package dnslookup

import (
    "fmt"
    "time"
    "github.com/miekg/dns"
    "github.com/Senzdetta/Comet/utils/color"
    "github.com/Senzdetta/Comet/utils/pcat"
)

func fetchSOARecord(
    domain, dnsServer string,
    timeout int,
) string {
    m := new(dns.Msg)
    m.SetQuestion(
        dns.Fqdn(domain),
        dns.TypeSOA,
    )
    m.RecursionDesired = true

    c := new(dns.Client)
    c.Timeout = time.Duration(timeout) * time.Second

    r, _, err := c.Exchange(m, dnsServer)
    if err != nil || r == nil || r.Rcode != dns.RcodeSuccess {
        return ""
    }

    for _, ans := range r.Answer {
        if soa, ok := ans.(*dns.SOA); ok {
            pcat.Printf(
                "%s[*] %sLookup: %s%s %s(%sns: %s%s%s, %smbox: %s%s%s, %sserial: %s%d%s, %srefresh: %s%d%s, %sretry: %s%d%s, %sexpire: %s%r%s, %sminttl: %s%s%s)%s",
                color.B, color.N, color.GG, domain, color.DG,
                color.N, color.BB, soa.Ns, color.DG,
                color.N, color.BB, soa.Mbox, color.DG,
                color.N, color.BB, soa.Serial, color.DG,
                color.N, color.BB, soa.Refresh, color.DG,
                color.N, color.BB, soa.Retry, color.DG,
                color.N, color.BB, soa.Expire, color.DG,
                color.N, color.BB, soa.Minttl, color.DG,
                color.N,
            )

            return fmt.Sprintf(
                "%s %s %d %d %d %d %d",
                soa.Ns, soa.Mbox, soa.Serial,
                soa.Refresh, soa.Retry, soa.Expire,
                soa.Minttl,
            )
        }
    }

    return ""
}

// Copyright (c) 2026 Senzdetta