package interpretercore

func isAlphaChar(ch byte) bool {
	return ch == '_' ||
		('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z')
}

func isAlphaNumChar(ch byte) bool {
	return ch == '_' ||
		('0' <= ch && ch <= '9') ||
		('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z')
}

func isWhitespace(ch byte) bool {
	return ch == ' ' ||
		ch == '\n' ||
		ch == '\t' ||
		ch == '\r'
}
