<div align="center">
    <img src="https://img.shields.io/badge/Comet-v0.1.20261119-blue?style=square&logo=go&v=1" />
    <img src="https://img.shields.io/badge/Platform-Linux-blue?style=square&logo=linux&v=1" />
    <a href="LICENSE">
        <img src="https://img.shields.io/badge/License-AGPL--3.0-blue?style=square&logo=github&v=1" />
    </a>
</div>

# Comet
Comet is a lightweight, high-performance CLI tool built for efficient reconnaissance and vulnerability scanning.

## Features
- Multi-protocol port scanner (TCP, UDP, SCTP).
- HTML link crawling (hyperlink extraction).
- Username enumeration across social media platforms.
- CORS misconfiguration and subdomain takeover scanner.
- DNS bruteforcing, SQLi scanning, and fuzzing.
- And more features.

## Disclaimer
Please read [disclaimer](.docs/disclaimer.md) before using this tool. <br>
Use this software at your own risk. <br>
The author is not responsible for any damage, data loss, or issues that may result from its use.

## Installation
Quick install:
```bash
git clone https://github.com/Senzdetta/Comet
bash Comet/install.sh
```
For more detailed installation instructions, see [install](.docs/install.md).

## Usage Example
```bash
comet --portscan 127.0.0.1 --protocol tcp,udp,sctp --only open
comet --pathfuzzer http://192.168.1.1/ --recursive
comet --corsscan http://192.168.1.1 --origin http://evil.com
comet --dnslookup google.com --dnsresolv 8.8.8.8:53
comet --userhunt someone --threads 10 --timeout 50 --check-body true
```
And more commands.

## More
For further information, visit: [more](.docs/__more__.md)

<code>Made with ❤ ~love.</code>