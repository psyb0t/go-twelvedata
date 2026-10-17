#!/bin/bash
# Converts schema examples arrays to a single example for kin-openapi and
# declares OpenAPI 3.0.3. No paths, operations or schemas are removed.
# Component x-go-name overrides reuse names such as Meta and Indicator across
# distinct schemas. Remove those overrides so unique component keys name Go
# types. Property and parameter overrides remain unchanged.
# Normalize typographic dashes and promotional wording in prose annotations
# before they become generated Go comments. Schema and wire values stay intact.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

main() {
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.patched.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp before the local leaves scope
	trap "rm -f -- '${tmp}'" EXIT

	jq '
def prose:
  gsub("\u2013|\u2014"; ",")
  | gsub("comprehensive"; "detailed")
  | gsub("leverage"; "use");
  walk(
    if type == "object" and (.examples? | type == "array") then
      (if has("example") or (.examples | length == 0) then .
       else .example = .examples[0] end) | del(.examples)
    else . end
  )
  | .openapi = "3.0.3"
  | .components.schemas |= with_entries(.value |= del(."x-go-name"))
  | walk(
      if type == "object" then
        with_entries(
          if (.key == "description" or .key == "summary" or .key == "example")
             and (.value | type == "string") then .value |= prose
          else . end
        )
      else . end
    )
' "${SPEC_DIR}/openapi.json" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/openapi.patched.json"
	log INFO "converted schema examples to OpenAPI 3.0 form"
}

main "$@"
