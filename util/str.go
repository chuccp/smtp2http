package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func DeduplicateIds(ids string) string {
	vs := strings.Split(ids, ",")
	vvs := RemoveDuplicates(vs)

	return strings.Join(vvs, ",")

}

func RemoveDuplicates(nums []string) []string {
	m := make(map[string]bool)
	var result []string
	for _, num := range nums {
		if _, ok := m[num]; !ok {
			m[num] = true
			result = append(result, num)
		}
	}
	return result
}

func StringToUintIds(ids string) []uint {
	uIds := make([]uint, 0)
	receiveEmailIds := strings.Split(ids, ",")
	for _, i2 := range receiveEmailIds {
		atoi, err := strconv.Atoi(i2)
		if err != nil {
			continue
		}
		uIds = append(uIds, uint(atoi))
	}
	return uIds
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
func GetCachePath(rootPath, filename string) string {
	id := uuid.New().String()
	name := MD5([]byte(id))
	ext := path.Ext(filename)
	if len(ext) > 0 {
		name = name + ext
	}
	return path.Join(rootPath, FormatDate(time.Now()), name)
}
