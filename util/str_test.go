package util

import (
	"regexp"
	"testing"
)

func TestName(t *testing.T) {
	re := regexp.MustCompile("^/assets")
	fa := re.MatchString("/assets/log-BwvnxVjg.js")
	t.Log(fa)
}
