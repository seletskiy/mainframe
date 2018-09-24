package text

import (
	"regexp"
	"strings"
)

var (
	reTag = regexp.MustCompile(`(?P<tag>[a-z_]+)`)

	reValueInt    = regexp.MustCompile(`(?P<int>\d+)`)
	reValueColor  = regexp.MustCompile(`#(?P<color>[\da-f]{3}|[\da-f]{6})`)
	reValueString = regexp.MustCompile(`(?P<string>"(?:\\.|[^\\"])*")`)

	reValue = regexp.MustCompile(
		strings.NewReplacer(
			`{int}`, reValueInt.String(),
			`{color}`, reValueColor.String(),
			`{string}`, reValueString.String(),
		).Replace(
			`(?:{int}|{color}|{string})`,
		),
	)

	reKey      = regexp.MustCompile(`(?P<key>[a-z_]+)`)
	reKeyValue = regexp.MustCompile(
		strings.NewReplacer(
			`{key}`, reKey.String(),
			`{value}`, reValue.String(),
		).Replace(
			`{key}(?::\s*{value})?`,
		),
	)

	reGarbage = regexp.MustCompile(`(?P<garbage>.+)`)

	reMessage = regexp.MustCompile(
		strings.NewReplacer(
			`{tag}`, reTag.String(),
			`{key_value}`, reKeyValue.String(),
			`{garbage}`, reGarbage.String(),
		).Replace(
			`(?:^{tag}|{key_value})(?:\s+|$)|{garbage}`,
		),
	)
)
