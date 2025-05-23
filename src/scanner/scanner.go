package scanner

import (
	"fmt"
	"os"
	"strings"

	"moji/src/scanner/constants"
	"moji/src/scanner/types"
)

type scanner struct {
	source    string
	current   int
	start     int
	line      int
	tokens    []types.Token
	hadError  bool
}

func NewScanner(source string) *scanner {
	return &scanner{
		source:   source,
		current:  0,
		start:    0,
		line:     1,
		tokens:   []types.Token{},
		hadError: false,
	}
}

func (s *scanner) ScanTokens() []types.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.addToken(constants.EOF, nil)
	return s.tokens
}

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	case '(':
		s.addToken(constants.LEFT_PAREN, nil)
	case ')':
		s.addToken(constants.RIGHT_PAREN, nil)
	case '{':
		s.addToken(constants.LEFT_BRACE, nil)
	case '}':
		s.addToken(constants.RIGHT_BRACE, nil)
	case ',':
		s.addToken(constants.COMMA, nil)
	case '.':
		s.addToken(constants.DOT, nil)
	case '-':
		s.addToken(constants.MINUS, nil)
	case '+':
		s.addToken(constants.PLUS, nil)
	case ';':
		s.addToken(constants.SEMICOLON, nil)
	case '*':
		s.addToken(constants.STAR, nil)
	case '=':
		if s.match('=') {
			s.addToken(constants.EQUAL_EQUAL, nil)
		} else {
			s.addToken(constants.EQUAL, nil)
		}
	case '!':
		if s.match('=') {
			s.addToken(constants.BANG_EQUAL, nil)
		} else {
			s.addToken(constants.BANG, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(constants.LESS_EQUAL, nil)
		} else {
			s.addToken(constants.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(constants.GREATER_EQUAL, nil)
		} else {
			s.addToken(constants.GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(constants.SLASH, nil)
		}
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) || c > 127 {
			// Check for emoji operators
			if c == 240 && s.current < len(s.source) && len(s.source) >= s.current+3 {
				// Check for 👉 emoji (F0 9F 91 89)
				if s.source[s.current] == 159 && s.source[s.current+1] == 145 && s.source[s.current+2] == 137 {
					// This is the 👉 emoji (U+1F449)
					s.current += 3 // Skip the remaining bytes of the emoji
					if s.peek() == '=' && s.current < len(s.source) {
						s.advance() // Consume the '='
						s.addToken(constants.EQUAL_EQUAL, nil)
					} else {
						s.addToken(constants.EQUAL, nil)
					}
				// Check for 📝 emoji (U+1F4DD)
				} else if s.source[s.current] == 159 && s.source[s.current+1] == 147 && s.source[s.current+2] == 157 {
					// This is the 📝 emoji (U+1F4DD)
					s.current += 3 // Skip the remaining bytes of the emoji
					if s.peek() == '=' && s.current < len(s.source) {
						s.advance() // Consume the '='
						s.addToken(constants.EQUAL_EQUAL, nil)
					} else {
						s.addToken(constants.EQUAL, nil)
					}
				} else {
					s.identifier()
				}
			} else if c == 226 && s.current < len(s.source) && len(s.source) >= s.current+5 {
				// Check for ⚖️ emoji (E2 9A 96 EF B8 8F)
				if s.source[s.current] == 154 && s.source[s.current+1] == 150 && 
				   s.source[s.current+2] == 239 && s.source[s.current+3] == 184 && 
				   s.source[s.current+4] == 143 {
					// This is the ⚖️ emoji (U+2696 U+FE0F)
					s.current += 5 // Skip the remaining bytes of the emoji
					s.addToken(constants.EQUAL_EQUAL, nil)
				// Check for ▶️ emoji (E2 96 B6 EF B8 8F)
				} else if s.source[s.current] == 150 && s.source[s.current+1] == 182 && 
				   s.source[s.current+2] == 239 && s.source[s.current+3] == 184 && 
				   s.source[s.current+4] == 143 {
					// This is the ▶️ emoji (U+25B6 U+FE0F)
					s.current += 5 // Skip the remaining bytes of the emoji
					s.addToken(constants.GREATER, nil)
				// Check for ◀️ emoji (E2 97 80 EF B8 8F)
				} else if s.source[s.current] == 151 && s.source[s.current+1] == 128 && 
				   s.source[s.current+2] == 239 && s.source[s.current+3] == 184 && 
				   s.source[s.current+4] == 143 {
					// This is the ◀️ emoji (U+25C0 U+FE0F)
					s.current += 5 // Skip the remaining bytes of the emoji
					s.addToken(constants.LESS, nil)
				// Check for ↩️ emoji (E2 86 A9 EF B8 8F)
				} else if s.source[s.current] == 134 && s.source[s.current+1] == 169 && 
				   s.source[s.current+2] == 239 && s.source[s.current+3] == 184 && 
				   s.source[s.current+4] == 143 {
					// This is the ↩️ emoji (U+21A9 U+FE0F)
					s.current += 5 // Skip the remaining bytes of the emoji
					s.addToken(constants.ELSE, nil)
				} else {
					s.identifier()
				}
			} else {
				s.identifier()
			}
		} else {
			s.error(fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() byte {
	ch := s.source[s.current]
	s.current++
	return ch
}

func (s *scanner) addToken(tokenType types.TokenType, literal interface{}) {
	if tokenType == constants.EOF {
		s.tokens = append(s.tokens, types.Token{TokenType: tokenType, Lexeme: "", Literal: literal, Line: s.line})
	} else {
		s.tokens = append(s.tokens, types.Token{TokenType: tokenType, Lexeme: string(s.source[s.start:s.current]), Literal: literal, Line: s.line})
	}
}

func (s *scanner) error(message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", s.line, message)
	s.hadError = true
}

func (s *scanner) HasError() bool {
	return s.hadError
}

func Scan(fileContents []byte) {
	scanner := NewScanner(string(fileContents))
	scanner.ScanTokens()
	scanner.PrintTokens()
	
	if scanner.HasError() {
		os.Exit(65)
	}
}

func (s *scanner) PrintTokens() {
	for _, token := range s.tokens {
		fmt.Println(token.String())
	}
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *scanner) string() {
	// The opening quotation mark is already consumed. Now consume until we hit a closing quotation mark
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error("Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	// Trim the surrounding quotes
	str := s.source[s.start+1 : s.current-1]
	s.addToken(constants.STRING, str)
}

func (s *scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	lexeme := s.source[s.start:s.current]
	// Parse the number and format it
	value := lexeme
	if !strings.Contains(value, ".") {
		value = value + ".0"
	} else {
		// Remove trailing zeros after decimal point
		parts := strings.Split(value, ".")
		if len(parts) == 2 {
			// Remove trailing zeros from the decimal part
			decimal := strings.TrimRight(parts[1], "0")
			if decimal == "" {
				decimal = "0"
			}
			value = parts[0] + "." + decimal
		}
	}
	s.addToken(constants.NUMBER, value)
}

func (s *scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *scanner) identifier() {
	// For multi-byte characters like emojis
	for !s.isAtEnd() && (isAlphaNumeric(s.peek()) || s.peek() > 127) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, exists := constants.Keywords[text]
	if !exists {
		tokenType = constants.IDENTIFIER
	}
	s.addToken(tokenType, nil)
}
