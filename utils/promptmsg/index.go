package promptmsgutils

import (
	"regexp"
)

// TODO: is regexp object mutable on matches and replaces?
// If so then how to reset its state?

func ReplacePromptMsgPrefix(promptMsg string, prefix string) string {
	msgPrefixRegex := regexp.MustCompile("![a-z]+")
	
	return msgPrefixRegex.ReplaceAllString(
		promptMsg,
		prefix,
	)
}

func GetPromptMsgPrefix(promptMsg string) string {
	msgPrefixRegex := regexp.MustCompile("![a-z]+")

	res := msgPrefixRegex.FindString(
		promptMsg,
	)
	if len(res) == 0 {
		return ""
	}

	return res[1:]
}
