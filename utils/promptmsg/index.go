package promptmsgutils

import (
	"regexp"
)

var msgPrefixRegex = regexp.MustCompile("![a-z]+")

func ReplacePromptMsgPrefix(promptMsg string, prefix string) string {
	return msgPrefixRegex.ReplaceAllString(
		promptMsg,
		prefix,
	)
}
