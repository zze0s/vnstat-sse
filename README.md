# go-vnstat

Expose vnstat live data over a SSE api.

## Setup

Requires Go 1.22.

Download dependencies

```bash
go mod download
```

Build binary. This outputs a binary to `dist/go-vnstat`. Make it executable `chmod +x` then you can run it with `./dist/go-vnstat`.

```bash
make build
```

## Usage

Setup a client that points to `http://ip:8200/events?stream=live-data` for json feed. `/events?stream=messages` contains html useful for HTMX SSE.

