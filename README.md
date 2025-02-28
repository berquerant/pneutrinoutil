# pneutrionoutil

My [NEUTRIONO](https://studio-neutrino.com/) utilities.

# Usage

## CLI

``` shell
./dist/pneutrinoutil --neutrinoDir ./dist/NEUTRINO --workDir ./tmp --score /path/to/some.musicxml
```

## HTTP server

``` shell
./dist/pneutrinoutil-server --neutrinoDir ./dist/NEUTRINO --workDir ./tmp --pneutrinoutil ./dist/pneutrinoutil
```

[Swagger is available](http://localhost:9101/v1/swagger/index.html).

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

Download NEUTRINO and singer libraries, and build pneutrinoutil.

``` shell
./task ansible
```
