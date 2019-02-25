package utils

import "strings"

// JoinStrings takes a variadic number of strings and uses
// a string builder to join them together
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

// HandleErr will panic on a fatal error
func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
