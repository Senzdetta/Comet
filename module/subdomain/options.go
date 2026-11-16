// https://github.com/Senzdetta/Comet

package subdomain

type Options struct {
	Domain    string
	Mode      string
	Threads   int
	Timeout   int
	Retry     int
	Api       string
	LookupIP  bool
	Wordlist  string
	DNSResolv string
}

// Copyright (c) 2026 Senzdetta