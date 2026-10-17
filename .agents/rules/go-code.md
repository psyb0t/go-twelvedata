# Go code

Refines `~/.claude/rule-details/go/wrapping-errors.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` is generated. Never edit it by hand. Change `config.yaml` or the spec patch in `scripts/`, then run `make generate` and commit the result.
- **GO2** Spec fixes live in `scripts/`, applied to the downloaded spec before generation. `scripts/patch-spec.sh` converts schema `examples` arrays to one `example`, declares OpenAPI 3.0.3 for kin-openapi and removes colliding component-level `x-go-name` overrides. Keep property and parameter naming overrides. Each patch says in its header which spec problem it fixes. Never remove paths or operations.
- **GO3** The downloaded spec in `api/` is gitignored and never committed.
- **GO4** Number fields stay `float64` through the `type-mapping` in `config.yaml`. Money in `float32` loses precision.
- **GO5** Hand-written Go (`gen.go` and the tests) wraps errors with `ctxerrors.Wrap` against `ctxerrors/commerr` sentinels. No `fmt.Errorf`, no `errors.New` outside a sentinel declaration.
- **GO6** The exported API is whatever oapi-codegen emits. A config change that renames or removes an exported name is a breaking release (see `releases.md`).
- **GO7** `make lint` passes with zero issues. Hand-written lines are at most 80 columns, the `.golangci.yml` limit.
