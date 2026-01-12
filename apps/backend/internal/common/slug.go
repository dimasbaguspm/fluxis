package common

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	nonWordRegexp    = regexp.MustCompile(`[^a-z0-9_]+`)
	multiUnderRegexp = regexp.MustCompile(`_+`)
)

func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = nonWordRegexp.ReplaceAllString(s, "_")
	s = multiUnderRegexp.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		return "n_a"
	}
	return s
}

func SlugifyUnique(s string) string {
	base := Slugify(s)
	ts := time.Now().Unix()
	return base + "_" + strconv.FormatInt(ts, 10)
}
