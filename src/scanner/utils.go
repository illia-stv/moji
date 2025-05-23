package scanner

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_' ||
		c > 127 // Allow UTF-8 characters (including emojis)
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
} 