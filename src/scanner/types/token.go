package types

import "fmt"

type TokenType string

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (t *Token) String() string {
	if t.Literal == nil {
		return fmt.Sprintf("%s %s %s", t.TokenType, t.Lexeme, "null")
	}
	return fmt.Sprintf("%s %s %v", t.TokenType, t.Lexeme, t.Literal)
} 