// Package filters provides the default host filter list for Rhine and
// a small utility function to generate a regexp expression from a slice of strings.
package filters

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// HostFilter is the default host filter for Rhine.
	// Requests to hosts in the hostFilter slice will be blocked.
	HostFilter = GenerateFilter([]string{
		`android\.bugly\.qq\.com`,
		`sessions\.bugsnag\.com`,
		`app\.adjust\.com`,
	})
)

// GenerateFilter compiles a regexp expression for a given list of URLs
func GenerateFilter(list []string) *regexp.Regexp {
	ret := ""
	for _, v := range list {
		ret += fmt.Sprintf(`(^.*%s.*$)|`, v)
	}
	ret = strings.TrimRight(ret, "|")
	return regexp.MustCompile(ret)
}
