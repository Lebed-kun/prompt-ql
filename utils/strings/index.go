package stringsutils

import "strings"

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func TrimWhitespace(s string) string {
	res := strings.Builder{}
	sRune := []rune(s)

	ptr := 0
	for ptr < len(sRune) && isWhitespace(sRune[ptr]) {
		ptr++
	}

	for ptr < len(sRune) {
		if sRune[ptr] == '\n' {
			res.WriteRune(sRune[ptr])
			ptr++

			for ptr < len(sRune) && isWhitespace(sRune[ptr]) {
				ptr++
			}
		} else if isWhitespace(sRune[ptr]) {
			res.WriteRune(sRune[ptr])
			ptr++

			for ptr < len(sRune) && isWhitespace(sRune[ptr]) && sRune[ptr] != '\n' {
				ptr++
			}
		} else {
			res.WriteRune(sRune[ptr])
			ptr++
		}
	}

	for ptr < len(sRune) && isWhitespace(sRune[ptr]) {
		ptr++
	}

	return res.String()
}
