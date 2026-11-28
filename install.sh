#!/usr/bin/env bash
# https://github.com/SenzdettaCore/blsys

blsys_url="github.com/SenzdettaCore/blsys@latest"
msgblsys="Updating"

if ! command -v go > /dev/null 2>&1; then
    echo -e "\x1b[1;31m[!] \x1b[0mGolang is not installed!" >&2
    echo -e "\x1b[1;31m[!] \x1b[0mPlease installing first!" >&2
    exit 1
fi

if ! command -v blsys > /dev/null 2>&1; then
    msgblsys="Installing"
    echo -e "\x1b[0;33m[!] \x1b[0mBuild system: \x1b[0;32mblsys \x1b[0mNot found!" >&2
fi

if [[ -f "${PREFIX:-}/etc/os-release" && ! -d '/storage/emulated/0/Android/data' ]]; then
    if [[ ${EUID} -ne 0 ]]; then
        echo -e "\x1b[1;31m[!] \x1b[0mEUID: \x1b[0;32m${EUID} \x1b[0mroot required!" >&2
        exit 1
    fi
fi

echo -ne "\x1b[1;34m[*] \x1b[0m${msgblsys}: \x1b[0;32m${blsys_url}\x1b[0m"
export GOPROXY='direct'
export GOBIN="${PREFIX:-/usr}/bin" 
command go install "${blsys_url}" > /dev/null 2>&1
if [[ ${?} -eq 0 ]]; then
    echo -e " \x1b[1;90m- \x1b[0mexit: \x1b[0;32m${?}\x1b[0m"
else
    echo -e " \x1b[1;90m- \x1b[0merror: \x1b[1;31m${?}\x1b[0m"
    exit 1
fi

unset GOPROXY
unset GOBIN
unset blsys_url
unset msgblsys

export __root__="$(
    command blsys \
        --lookpath "${BASH_SOURCE[0]}" \
        --rtt 0
)"; readonly __root__

export manifest="${__root__}/manifest.meta"; readonly manifest

user_platform="$(command blsys --platform)"
echo -ne "\x1b[1;34m[*] \x1b[0mChecking platform support: \x1b[0;32m${user_platform}\x1b[0m"
if command blsys --platform --support "${manifest}"; then
    echo -e " \x1b[1;90m- \x1b[0mexit: \x1b[0;32m${?}\x1b[0m"
else
    echo -e " \x1b[1;90m- \x1b[0merror: \x1b[1;31m${?}\x1b[0m"
    echo -e "\x1b[1;31m[!] \x1b[0mUnsupported platform: \x1b[0;32m${user_platform}\x1b[0m" >&2
    exit 1
fi

euid="$(command blsys --geteuid)"
if [[ "${user_platform}" != 'android' ]]; then
    if [[ ${euid} -ne 0 ]]; then
        echo -e "\x1b[1;31m[!] \x1b[0mEUID: \x1b[0;32m${euid} \x1b[0mroot required!" >&2
        exit 1
    fi
fi

unset user_platform
unset euid

if ((
    BASH_VERSINFO[0] < 4 ||
    (BASH_VERSINFO[0] == 4 && BASH_VERSINFO[1] < 4)
)); then
    echo -e "\x1b[0;33m[!] \x1b[0mWarning bash version: \x1b[0;31m${BASH_VERSION} \x1b[1;90m< \x1b[0;32mv4.4\x1b[0m" >&2
fi

source <(command blsys --codegen 'import')

if [[ ! -f "${__root__}/.blsys/install.sh" ]]; then
    echo -e "\x1b[1;31m[!] \x1b[0mFile: \x1b[0;32m${__root__}/.blsys/install.sh \x1b[0mnot found!" >&2
    exit 1
fi

echo -e "\x1b[1;34m[*] \x1b[0mCall: \x1b[0;32mblsys::install\x1b[0m"
source "${__root__}/.blsys/install.sh"
blsys::install "${@}"
exit 0

# Copyright (c) 2026 Senzdetta