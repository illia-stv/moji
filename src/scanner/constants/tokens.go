package constants

import "moji/src/scanner/types"

const (
	// Single-character tokens.
	LEFT_PAREN  types.TokenType = "LEFT_PAREN"
	RIGHT_PAREN types.TokenType = "RIGHT_PAREN"
	LEFT_BRACE  types.TokenType = "LEFT_BRACE"
	RIGHT_BRACE types.TokenType = "RIGHT_BRACE"
	COMMA       types.TokenType = "COMMA"
	DOT         types.TokenType = "DOT"
	MINUS       types.TokenType = "MINUS"
	PLUS        types.TokenType = "PLUS"
	SEMICOLON   types.TokenType = "SEMICOLON"
	SLASH       types.TokenType = "SLASH"
	STAR        types.TokenType = "STAR"

	// One or two character tokens.
	BANG          types.TokenType = "BANG"
	BANG_EQUAL    types.TokenType = "BANG_EQUAL"
	EQUAL         types.TokenType = "EQUAL"
	EQUAL_EQUAL   types.TokenType = "EQUAL_EQUAL"
	GREATER       types.TokenType = "GREATER"
	GREATER_EQUAL types.TokenType = "GREATER_EQUAL"
	LESS          types.TokenType = "LESS"
	LESS_EQUAL    types.TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER types.TokenType = "IDENTIFIER"
	STRING     types.TokenType = "STRING"
	NUMBER     types.TokenType = "NUMBER"

	// Keywords.
	AND    types.TokenType = "AND"
	CLASS  types.TokenType = "CLASS"
	ELSE   types.TokenType = "ELSE"
	FALSE  types.TokenType = "FALSE"
	FUN    types.TokenType = "FUN"
	FOR    types.TokenType = "FOR"
	IF     types.TokenType = "IF"
	NIL    types.TokenType = "NIL"
	OR     types.TokenType = "OR"
	PRINT  types.TokenType = "PRINT"
	RETURN types.TokenType = "RETURN"
	SUPER  types.TokenType = "SUPER"
	THIS   types.TokenType = "THIS"
	TRUE   types.TokenType = "TRUE"
	VAR    types.TokenType = "VAR"
	WHILE  types.TokenType = "WHILE"
	EOF    types.TokenType = "EOF"
)

var Keywords = map[string]types.TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"↩️":     ELSE,
	"false":  FALSE,
	"⛔️":     FALSE,
	"for":    FOR,
	"fun":    FUN,
	"🔀":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"📢":     PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"✅":     TRUE,
	"🎁":     VAR,
	"🔄":     WHILE,
} 