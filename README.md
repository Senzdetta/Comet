# Comet
Comet is a lightweight CLI tool built for high-efficiency security reconnaissance.

## Features
- Massive username enumeration across 700+ domains.
- High-performance directory fuzzing.
- HTML link crawling (hyperlink extraction).
- Parameter discovery via Wayback Machine.
- And more features.

## Disclaimer
Please read [disclaimer](.docs/disclaimer.md) before using this tool. </br>
Use this software at your own risk. </br>
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

<code>Made by ❤~love</code>