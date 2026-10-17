// Package twelvedata is a Go client for the Twelve Data market data API,
// generated with oapi-codegen from the OpenAPI document Twelve Data publishes.
package twelvedata

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.json
