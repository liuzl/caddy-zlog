# Caddy-zlog

## Overview

`zlog` is a log middleware for [Caddy](https://github.com/mholt/caddy), it's based on https://github.com/rs/zerolog.

## Installation

Rebuild caddy as follows:

1. `git clone https://github.com/liuzl/caddy-zlog`
2. copy `caddy-zlog` to `github.com/mholt/caddy/caddyhttp/zlog`
3. add `_ "github.com/mholt/caddy/caddyhttp/zlog"` to file `github.com/mholt/caddy/caddyhttp/caddyhttp.go`
4. add `zlog` to the variable `directives` in file `github.com/mholt/caddy/caddyhttp/httpserver/plugin.go`
5. `cd github.com/mholt/caddy/caddy && go build`

## Caddyfile syntax

```
127.0.0.1 {
    zlog
}
```
