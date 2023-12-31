package cmdutil

import (
	"regexp"
	"strings"
)

const (
	IssueFieldKeyServiceList = "customfield_12300"
	ServicePublishedMark     = "（已发布）"
)

var (
	publishedPattern   = regexp.MustCompile(`^.+已发布.*$`)
	serviceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+`)
)

func MatchServiceName(srcName, target string) bool {
	if !strings.HasPrefix(srcName, target) {
		return false
	}
	return GetPureServiceName(srcName) == target
}

func GetPureServiceName(name string) string {
	return serviceNamePattern.FindString(name)
}

func IsPublishedService(service string) bool {
	return publishedPattern.MatchString(service)
}
