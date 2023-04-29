[![Go Reference](https://pkg.go.dev/badge/github.com/dcilke/hz.svg)](https://pkg.go.dev/github.com/dcilke/hz)
![Build Status](https://github.com/dcilke/hz/actions/workflows/release.yml/badge.svg)

# hz

Human readable streaming formatter based on `zerolog.ConsoleWriter{}`

## Install

### Go

```zsh
go get github.com/dcilke/hz
```

### Homebrew

```zsh
brew tap dcilke/taps
brew install dcilke/taps/hz
```

## Help

```zsh
hz --help
Usage:
  hz [FILE]

Application Options:
  -l, --level=    only output lines at this level
  -s, --strict    exclude non JSON output
  -f, --flat      flatten objects
  -v, --vertical  vertical output
  -r, --raw       raw output
  -n, --no-pin    exclude pinning of fields

Help Options:
  -h, --help      Show this help message
```

## Config

Default command options can be specified in a config file located at `$HOME/.config/hz/config.yml`.

```yaml
level:
  - trace
  - debug
  - info
  - warn
  - error
  - fatal
  - panic
strict: false
flat: false
vertical: false
plain: false
noPin: false
```

## Why?

I use [zerolog](https://github.com/rs/zerolog) for structured logging and want to be able to quickly tap into the log streams.

## Looking for something similar?

- [pino-pretty](https://github.com/pinojs/pino-pretty)
- [jq](https://github.com/stedolan/jq)
