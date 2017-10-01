package utils

import "errors"

var (
	MDMissingSections = errors.New("Missing required sections for article")
	MDEmpty           = errors.New("Cannot process empty markdown")
)
