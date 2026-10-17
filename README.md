# go-twelvedata

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-twelvedata.svg)](https://pkg.go.dev/github.com/psyb0t/go-twelvedata)
[![CI](https://github.com/psyb0t/go-twelvedata/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-twelvedata/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-twelvedata/badges/version.svg)](https://github.com/psyb0t/go-twelvedata/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-twelvedata/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-twelvedata/badges/importers.svg)](https://github.com/psyb0t/go-twelvedata/blob/badges/importers.md)

Unofficial Go client for the [Twelve Data API](https://api.twelvedata.com/doc/swagger/openapi.json), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Twelve Data publishes. I am not affiliated with Twelve Data.

The package is `twelvedata`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-twelvedata
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-twelvedata"
)

func main() {
	auth := twelvedata.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"apikey "+os.Getenv("TWELVEDATA_API_KEY"),
			)

			return nil
		},
	)

	client, err := twelvedata.NewClientWithResponses(
		"https://api.twelvedata.com",
		auth,
	)
	if err != nil {
		panic(err)
	}

	const symbol = "SYMBOL_X"

	res, err := client.GetPriceWithResponse(
		context.Background(),
		&twelvedata.GetPriceParams{Symbol: new(symbol)},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

Replace `SYMBOL_X` with the instrument you need. Check `StatusCode()` and `JSON200` before using the result. A nil Go error covers transport and decoding, not HTTP success. Prices such as `JSON200.Price` are strings; number schemas map to `float64`.

The client covers all 187 paths, including prices, quotes, time series, technical indicators, fundamentals, instrument lists and `POST /batch`. Method names follow the spec's operation IDs.

## Auth

Twelve Data accepts `Authorization: apikey <key>`. Read the key from `TWELVEDATA_API_KEY` and set it with the request editor shown above. The spec also permits an `apikey` query parameter. Prefer the header to keep credentials out of URL logs.

The base URL is `https://api.twelvedata.com`. Endpoint access, available data and API credits depend on your plan and permissions. Handle HTTP 429 with bounded backoff. The generated client has no automatic retries or throttling.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://api.twelvedata.com/doc/swagger/openapi.json` to `api/openapi.json`.
- `scripts/patch-spec.sh` writes `api/openapi.patched.json`. It converts OpenAPI 3.1 schema `examples` arrays to a single `example` and declares OpenAPI 3.0.3 so kin-openapi can read the document. It removes component-level `x-go-name` overrides that reuse names such as `Meta` for different schemas, letting the unique component keys name the types. Property and parameter overrides and wire schemas remain intact. It removes no paths or operations.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64` and resolves type and enum naming collisions.

`.github/workflows/regen.yml` does the same every day in CI. When the spec changes the client, it runs `make lint`, `make test` and `make sec`. A failing check can trigger up to three repair attempts and a review. Passing changes get a CHANGELOG entry and a release tag.

## Agent integrations

The [skill](.agents/skills/go-twelvedata) teaches a coding agent how to use go-twelvedata when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-twelvedata@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-twelvedata@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-twelvedata
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Twelve Data's published API specification. The specification itself is not redistributed here. The downloaded document has no `info.license` or `info.termsOfService` field. API access and data use remain subject to the provider's terms.
