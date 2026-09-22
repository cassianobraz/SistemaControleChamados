// Package docs embeds the hand-written OpenAPI 3.0 specification so it
// ships inside the compiled binary — no extra file needs to be copied into
// the Docker image or read from disk at runtime.
package docs

import _ "embed"

//go:embed openapi.json
var OpenAPISpec []byte
