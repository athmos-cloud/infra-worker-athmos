package resource

import (
	"strings"
)

func formatResourceName(name string) string {
	toRemove := []string{" ", "_", ".", "@", "$", "%", "^", "&", "*", "(", ")", "[", "]", "{", "}", "|", "\\", "/", "?", "<", ">", "!", "`", "~", "+", "=", ",", ";", ":"}
	res := name
	for _, char := range toRemove {
		res = strings.ReplaceAll(res, char, "")
	}
	return strings.ToLower(res)
}
