# CHANGELOG

### 2026-10-06
```text
- Change --version output.
```

### 2026-10-08
```text
- Removing utils/cursor package.
```

### 2026-11-09
```text
- Converted --portscan to native Go (using net.Dial and syscall).
  It now supports TCP, UDP, and SCTP, and includes many new customization parameters.

- Removed support for files in --corsscan.
  It must now be a URL.
```

### 2026-11-10
```text
- Improvements to detection in --sqlscan.

- Improvements to --dnslookup.
  The default is now the native Go implementation, though you can still use the API-based version by specifying --mode api along with --api endpoint.
  It now also supports customization of threads, timeouts, custom API endpoints, and much more.
```

### 2026-11-11
```text
- There are a few minor changes to --header.
  The end result remains the same—this is simply to tidy up the code and refine the color styling.

- The --subdomain option now supports a native, API-free approach using DNS fuzzing (with customizable wordlists) and DNS-based IP lookups.
  The legacy API-based mode remains available via the --mode api flag, with further customization possible using the --api <api> option.
  Additionally, numerous other sub-flags have been added to allow for flexible configuration.
```

### 2026-11-12
```text
- Addition to the custom --api endpoint for --hostsearch, even though this is a legacy feature that has been entirely superseded by --subdomain.

- Added a new flag to --osint.
  --check-body is used to analyze the returned HTML, so as not to rely solely on the returned status code.

- Optional addition of a new flag to --osint.
  --wordlist <file> is used to search for usernames with a custom domain.

- Added domain list to module/osint/sites/osint_sites.txt from 700+ domains to 1060+ domains.

- Wordlist changes for --dirfuzzer.

- Implementation of native Go in --paramscan, and, of course, the legacy --mode api.

- Fixed panic race condition in module/dirfuzzer/fuzz.go.
```

### 2026-11-13
```text
- Adding --service <file> to --portscan allows for customizing the services database.

- Renamed functions in several files.

- Replaced the --dirfuzzer flag with --pathfuzzer, and updated the associated package as well.

- The --osint flag has been changed to --userhunt, along with its corresponding package.
```

### 2026-11-15
```text
- Animation changes for --subdomain, and code cleanup.

- Removed the dependency on third-party tools for --subtakeover, which now utilizes subdomain.BruteForceDNS instead.
  Also added new options: --dnsresolv <dns>:[<port>], --wordlist <file>, and --fingerprint <file>.
```