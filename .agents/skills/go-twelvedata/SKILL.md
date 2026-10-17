---
name: go-twelvedata
description: "Use the unofficial go-twelvedata client to write Go code for Twelve Data market prices, quotes, time series, technical indicators, fundamentals, instrument metadata and batch queries."
homepage: https://github.com/psyb0t/go-twelvedata
user-invocable: true
metadata:
  openclaw:
    emoji: "💹"
    requires:
      bins:
        - go
      env:
        - TWELVEDATA_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the code you write talks HTTPS to api.twelvedata.com."
---

# go-twelvedata

Unofficial, typed Go client for the Twelve Data API, in package `twelvedata` at the module root.

Read [references/setup.md](references/setup.md) before writing calls. It covers installation, authentication, representative operations, responses and request limits.

## Security & safety

- Read the API key from `TWELVEDATA_API_KEY`. Set `Authorization: apikey <key>` through `WithRequestEditorFn`. Never hardcode, commit or log the key or an authenticated request.
- Use `https://api.twelvedata.com`. The spec also permits an `apikey` query parameter, which can expose the key in URL logs.
- Calls consume API credits. Respect the user's plan, data permissions and request budget.

## When to use

- Go code needs prices, quotes, historical candles, technical indicators, fundamentals or instrument lists from Twelve Data.
- Go code needs to combine market data requests with the batch endpoint.
- The project already imports `github.com/psyb0t/go-twelvedata`.

## When NOT to use

- The code is not Go. This module ships only a Go package.
- You need a client maintained by Twelve Data. This one is unofficial and not affiliated with Twelve Data.
- You want to change the generated client itself. Read `AGENTS.md` in the repository.

## Usage

The README example also lives in `example_test.go`:

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

Replace `SYMBOL_X` with the requested instrument. Check `StatusCode()` and the appropriate decoded status field before using a response. A nil Go error covers transport and decoding, not HTTP success. `JSON200` can be nil.

The `Price` response field is a string. Number schemas map to `float64`; many market price fields remain strings as defined by the provider.

Before finishing a change, compile the calling code and test its request and response handling against a local HTTP server. Tests must not spend provider credits.
