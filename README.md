[![Go Report Card](https://goreportcard.com/badge/github.com/berquerant/pneutrinoutil)](https://goreportcard.com/report/github.com/berquerant/pneutrinoutil)

# pneutrionoutil

My [NEUTRIONO](https://studio-neutrino.com/) utilities.

# Usage

## CLI

``` shell
./dist/pneutrinoutil --score /path/to/some.musicxml
```

## HTTP server

``` shell
./task start
```

[Swagger is available](http://localhost:9101/v1/swagger/index.html).

# Requirements

- macOS
- Go
- [uv](https://github.com/astral-sh/uv)
- [direnv](https://github.com/direnv/direnv)
- Docker
- AWS CLI

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
