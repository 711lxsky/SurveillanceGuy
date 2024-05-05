package util

import (
	"fmt"
	"regexp"
)

// MatchTargetByRegexPattern
// 根据正则表达式匹配目标内容
func MatchTargetByRegexPattern(content []byte, pattern string) (string, error) {
	compileRes, _ := regexp.Compile(pattern)
	items := compileRes.FindSubmatch(content)
	if len(items) >= 2 {
		res := string(items[1])
		return res, nil
	} else {
		return "", fmt.Errorf("Can't match target")
	}
}
