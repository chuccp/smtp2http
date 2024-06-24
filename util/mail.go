package util

import (
	"fmt"
	"regexp"
)

func FormatMail(name, mail string) string {
	if len(name) == 0 {
		return mail
	}
	return fmt.Sprintf(`"%s" <%s>`, name, mail)
}

func ExtractEmails(text string) []string {
	// 定义正则表达式
	emailRegex := regexp.MustCompile(
		`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)

	// 使用FindAllString方法查找所有匹配的电子邮件地址
	emails := emailRegex.FindAllString(text, -1)

	// 去除重复的电子邮件地址
	uniqueEmails := removeDuplicates(emails)

	return uniqueEmails
}

// removeDuplicates 函数用于去除电子邮件地址列表中的重复项
func removeDuplicates(emails []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range emails {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
