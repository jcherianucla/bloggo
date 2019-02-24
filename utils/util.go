package utils

import "strings"

func JoinStrings(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		_, err := sb.WriteString(str)
		if err != nil {
			continue
		}
	}
	return sb.String()
}
