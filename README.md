# Caddy-zlog

## Overview

`zlog` is a log middleware for [Caddy](https://github.com/caddyserver/caddy), it's based on https://github.com/rs/zerolog.

## Installation

Rebuild caddy as follows:

1. `git clone https://github.com/liuzl/caddy-zlog`
2. copy `caddy-zlog` to `github.com/caddyserver/caddy/caddyhttp/zlog`
3. add `_ "github.com/caddyserver/caddy/caddyhttp/zlog"` to file `github.com/caddyserver/caddy/caddyhttp/caddyhttp.go`
4. add `zlog` to the variable `directives` in file `github.com/caddyserver/caddy/caddyhttp/httpserver/plugin.go`
5. `cd github.com/caddyserver/caddy/caddy && go build`

## Caddyfile syntax

```
127.0.0.1 {
    zlog {
        log_dir ./server_zerolog
    }
}
```
