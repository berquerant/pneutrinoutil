# pneutrionoutil

My [NEUTRIONO](https://studio-neutrino.com/) utilities.

# Usage

``` shell
./dist/pneutrinoutil --neutrinoDir ./dist/NEUTRINO --workDir ./tmp --score /path/to/some.musicxml
```

# Requirements

- macOS
- Go
- [uv](https://github.com/astral-sh/uv)

# Installation

Prepare libraries.

``` shell
uv sync
go mod download
```

Download NEUTRINO and singer libraries, and build pneutrinoutil

``` shell
./task ansible
```
