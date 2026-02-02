# golink

A tiny HTTP redirect server for “go links”. Configure a list of short names (like `hr`) and their long targets, start the server, and then visit `http://localhost:<port>/<short>` to get redirected.

## Config

`golink` reads a YAML config file specified via `-config`.

Example:

```yaml
port: 8080
default_scheme: https
links:
  - short: example
    long: https://example.com
```

- `port` is required.
- `default_scheme` is required and is used when a link target does not include `http://` or `https://`.
- `links` is a list of objects with:
  - `short`: the path segment (without a leading `/`)
  - `long`: the redirect target. If it does not include a scheme, `default_scheme` will be used.

## Run

```sh
go run ./cmd/golink -config ./config.yaml
```

## Usage

- Requesting `/` returns a plain-text message.
- Requesting `/<short>` redirects to the configured target.
- Unknown shorts return `404`.

Examples:

```sh
curl -i http://localhost:8080/
curl -i http://localhost:8080/example
curl -i http://localhost:8080/unknown
```

## Notes

- If you configure `long: example.com` and `default_scheme: https`, the redirect target will be `https://example.com`.
