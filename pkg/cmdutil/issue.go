package cmdutil

import (
	"regexp"
	"strings"
)

const (
	IssueFieldKeyServiceList = "customfield_12300"
)

var (
	nameAlphaPattern = regexp.MustCompile(`^[a-zA-Z0-9_].*$`)
)

func MatchServiceName(srcName, target string) bool {
	if !strings.HasPrefix(srcName, target) {
		return false
	}
	if len(srcName) > len(target) && nameAlphaPattern.MatchString(srcName[len(target):]) {
		return false
	}
	return true
}
