package transformer

import (
	"regexp"
	"strings"
)

func CreateSlug(s string) string {
	s = strings.ToLower(s)

	reg, _ := regexp.Compile("[^a-z0-9\\s-]+")
	s = reg.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, " ", "-")

	reg2, _ := regexp.Compile("-+")
	s = reg2.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}
