package common

import (
	"regexp"
)

const HEADER_ACCESS_TOKEN = "Access-Token"
const HEADER_CLIENT_ID = "Client-Id"
const HEADER_CLIENT_SECRET = "Client-Secret"

var HEADER_AUTH_REGEX = regexp.MustCompile(`^Bearer (.*)$`)

var SUPPORTED_DATA_TYPES = map[string]bool{
	"text":        true,
	"timestamptz": true,
	"uuid":        true,
	"smallint":    true,
	"integer":     true,
	"bigint":      true,
	"decimal":     true,
	"numeric":     true,
}
