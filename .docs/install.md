# INSTALL

### install.sh
Usage:
```bash
git clone https://github.com/Senzdetta/Comet
bash Comet/install.sh # <options...>
```

`install.sh` optional options (can be used together):
```text
* --home
└── override $HOME value.
```

### go install
Usage:
```bash
go install github.com/Senzdetta/Comet@latest
# optional:
# ln -sfv "$(go env GOPATH)/bin/Comet" ${PREFIX:-/usr}/bin/comet
```