package api

import _ "embed"

// OpenAPIYAML is the OpenAPI 3.0 spec for the anon-skrzynka API (embedded).
//go:embed openapi.yaml
var OpenAPIYAML []byte
