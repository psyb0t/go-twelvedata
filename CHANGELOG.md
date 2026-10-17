# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Generated doc comments keep the API's own description text. The spec patch no longer rewrites dashes or wording in descriptions.
- `make sec` uses Semgrep 1.181.0. Its pinned requirements now hash every published file, so the development image also builds on arm64.

## v0.1.0, 2026-10-17

- Added `twelvedata`, a REST client generated with oapi-codegen from the Twelve Data API specification, version 0.0.1, OpenAPI 3.1.0. It covers all 187 paths: prices, quotes, time series, technical indicators, fundamentals, instrument metadata, calendars and batch requests.
- Number schemas map to `float64`. String price fields retain their string representation.
- `make generate` downloads [the public spec](https://api.twelvedata.com/doc/swagger/openapi.json), converts schema `examples` arrays to a single `example`, declares OpenAPI 3.0.3 for the generator and removes component-level `x-go-name` overrides that cause duplicate Go type names. Property and parameter names stay intact. No operations are removed. The spec is not committed.
- Spec prose annotations have typographic dashes and promotional wording normalized before generation. Wire schemas stay intact.
