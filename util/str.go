package util

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
)

func MD5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	md5Hash := hex.EncodeToString(hashValue)
	return md5Hash
}

func SplitPath(path string) []string {
	path = strings.ReplaceAll(path, "\\", "/")
	vs := strings.Split(path, "/")
	ps := make([]string, 0)
	for _, i2 := range vs {
		i2 = strings.TrimSpace(i2)
		if len(i2) == 0 {
			continue
		}
		ps = append(ps, i2)
	}
	return ps
}

func IsMatchPath(path, math string) bool {
	math = ReplaceAllRegex(math, "\\*[a-zA-Z]+", ".*")
	re := regexp.MustCompile(math)
	return re.MatchString(path)

}
func ReplaceAllRegex(path, regex, math string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(path, math)
}

func ContainsAny(s string, strs ...string) bool {
	for _, str := range strs {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}

func ContainsAnyIgnoreCase(s string, strs ...string) bool {
	sLower := strings.ToLower(s)
	for _, str := range strs {
		if strings.Contains(sLower, strings.ToLower(str)) {
			return true
		}
	}
	return false
}
func EqualsAnyIgnoreCase(s string, strs ...string) bool {
	sLower := strings.ToLower(s)
	for _, str := range strs {
		if sLower == strings.ToLower(str) {
			return true
		}
	}
	return false
}
func BoolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
