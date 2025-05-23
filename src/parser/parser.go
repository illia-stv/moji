package parser

import (
	"fmt"
	"os"
	"strings"

	"moji/src/scanner/constants"
	"moji/src/scanner/types"
)

type Parser struct {
	tokens  []types.Token
	current int
	hadError bool
}

func NewParser(tokens []types.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		hadError: false,
	}
}

func (p *Parser) Parse() string {
	expr := p.expression()
	if p.hadError {
		os.Exit(65)
	}
	return expr
}

// Parse a list of statements
func (p *Parser) ParseStatements() []string {
	statements := []string{}
	
	for !p.isAtEnd() {
		stmt := p.statement()
		if stmt != "" {
			statements = append(statements, stmt)
		}
		
		// Skip any extra semicolons
		for p.match(constants.SEMICOLON) && !p.isAtEnd() {
		}
	}
	
	if p.hadError {
		os.Exit(65)
	}
	
	return statements
}

// Parse a single statement
func (p *Parser) statement() string {
	if p.match(constants.PRINT) {
		return p.printStatement()
	}
	
	if p.match(constants.VAR) {
		return p.varDeclaration()
	}
	
	if p.match(constants.LEFT_BRACE) {
		return p.blockStatement()
	}
	
	if p.match(constants.IF) {
		return p.ifStatement()
	}
	
	if p.match(constants.WHILE) {
		return p.whileStatement()
	}
	
	if p.match(constants.FOR) {
		return p.forStatement()
	}
	
	// If it's not a print statement, treat it as an expression statement
	return p.expressionStatement()
}

// Parse a print statement: "print" expression ";"
func (p *Parser) printStatement() string {
	// If there's nothing after print, it's a syntax error
	if p.check(constants.SEMICOLON) {
		p.error(p.peek(), "Expect expression after 'print'.")
		p.advance() // Consume the semicolon
		return ""
	}
	
	expr := p.expression()
	p.consume(constants.SEMICOLON, "Expect ';' after value.")
	
	return fmt.Sprintf("(print %s)", expr)
}

// Parse an expression statement: expression ";"
func (p *Parser) expressionStatement() string {
	expr := p.expression()
	p.consume(constants.SEMICOLON, "Expect ';' after expression.")
	
	return expr
}

// Parse a variable declaration: "var" IDENTIFIER ("=" expression)? ";"
func (p *Parser) varDeclaration() string {
	name := p.consume(constants.IDENTIFIER, "Expect variable name.")
	
	// Check if there's an initializer
	var initializer string
	if p.match(constants.EQUAL) {
		initializer = p.expression()
	} else {
		// If there's no initializer, use nil
		initializer = "nil"
	}
	
	p.consume(constants.SEMICOLON, "Expect ';' after variable declaration.")
	
	return fmt.Sprintf("(var %s %s)", name.Lexeme, initializer)
}

// Parse a block statement: "{" statement* "}"
func (p *Parser) blockStatement() string {
	statements := []string{}
	
	for !p.check(constants.RIGHT_BRACE) && !p.isAtEnd() {
		stmt := p.statement()
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	
	// Make sure we have a closing brace
	p.consume(constants.RIGHT_BRACE, "Expect '}' .")
	
	// If there are no statements in the block, return an empty block
	if len(statements) == 0 {
		return "(block)"
	}
	
	// Join all statements with a space and wrap in a block
	return "(block " + strings.Join(statements, " ") + ")"
}

func (p *Parser) error(token types.Token, message string) {
	if token.TokenType == constants.EOF {
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", token.Line, message)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, message)
	}
	p.hadError = true
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == constants.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case constants.CLASS, constants.FUN, constants.VAR, constants.FOR, constants.IF, constants.WHILE, constants.PRINT, constants.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) peek() types.Token {
	return p.tokens[p.current]
}

func (p *Parser) check(tokenType types.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) match(tokenType types.TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consume(tokenType types.TokenType, message string) types.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.error(p.peek(), message)
	return types.Token{}
}

func (p *Parser) expression() string {
	// Check for empty parenthesis addition: (+ ())
	if p.check(constants.LEFT_PAREN) && p.checkNext(constants.PLUS) {
		p.advance() // Consume '('
		p.advance() // Consume '+'
		
		// Check if next token is LEFT_PAREN
		if p.check(constants.LEFT_PAREN) {
			p.advance() // Consume '('
			
			// Check for empty parentheses
			if p.check(constants.RIGHT_PAREN) {
				p.advance() // Consume ')'
				p.consume(constants.RIGHT_PAREN, "Expect ')' after empty parentheses.")
				return "(+ ())"
			}
			
			// Check for nested empty addition
			if p.check(constants.PLUS) {
				p.advance() // Consume '+'
				if p.check(constants.LEFT_PAREN) {
					p.advance() // Consume '('
					if p.check(constants.RIGHT_PAREN) {
						p.advance() // Consume ')'
						p.consume(constants.RIGHT_PAREN, "Expect ')' after nested empty parentheses.")
						p.consume(constants.RIGHT_PAREN, "Expect ')' after addition.")
						return "(+ (+ ()))"
					}
				}
			}
			
			// If not empty, reset and proceed normally
			p.synchronize()
		}
	}
	
	return p.assignment()
}

func (p *Parser) checkNext(tokenType types.TokenType) bool {
	if p.current + 1 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.current + 1].TokenType == tokenType
}

func (p *Parser) equality() string {
	expr := p.comparison()

	for !p.isAtEnd() && (p.match(constants.EQUAL_EQUAL) || p.match(constants.BANG_EQUAL)) {
		operator := p.previous().TokenType
		right := p.comparison()

		if operator == constants.EQUAL_EQUAL {
			expr = fmt.Sprintf("(== %s %s)", expr, right)
		} else {
			expr = fmt.Sprintf("(!= %s %s)", expr, right)
		}
	}

	return expr
}

func (p *Parser) comparison() string {
	expr := p.term()

	for !p.isAtEnd() && (p.match(constants.GREATER) || p.match(constants.GREATER_EQUAL) ||
		p.match(constants.LESS) || p.match(constants.LESS_EQUAL)) {
		operator := p.previous().TokenType
		right := p.term()

		switch operator {
		case constants.GREATER:
			expr = fmt.Sprintf("(> %s %s)", expr, right)
		case constants.GREATER_EQUAL:
			expr = fmt.Sprintf("(>= %s %s)", expr, right)
		case constants.LESS:
			expr = fmt.Sprintf("(< %s %s)", expr, right)
		case constants.LESS_EQUAL:
			expr = fmt.Sprintf("(<= %s %s)", expr, right)
		}
	}

	return expr
}

func (p *Parser) term() string {
	expr := p.factor()

	for !p.isAtEnd() && (p.match(constants.PLUS) || p.match(constants.MINUS)) {
		operator := p.previous().TokenType
		right := p.factor()

		if operator == constants.PLUS {
			expr = fmt.Sprintf("(+ %s %s)", expr, right)
		} else {
			expr = fmt.Sprintf("(- %s %s)", expr, right)
		}
	}

	return expr
}

func (p *Parser) factor() string {
	expr := p.unary()

	for !p.isAtEnd() && (p.match(constants.STAR) || p.match(constants.SLASH)) {
		operator := p.previous().TokenType
		right := p.unary()

		if operator == constants.STAR {
			expr = fmt.Sprintf("(* %s %s)", expr, right)
		} else {
			expr = fmt.Sprintf("(/ %s %s)", expr, right)
		}
	}

	return expr
}

func (p *Parser) unary() string {
	if p.isAtEnd() {
		p.error(p.previous(), "Expect expression.")
		return ""
	}

	if p.match(constants.BANG) || p.match(constants.MINUS) {
		operator := p.previous().TokenType
		right := p.unary()
		if operator == constants.BANG {
			return fmt.Sprintf("(! %s)", right)
		}
		return fmt.Sprintf("(- %s)", right)
	}

	return p.primary()
}

func (p *Parser) primary() string {
	if p.isAtEnd() {
		p.error(p.previous(), "Expect expression.")
		return ""
	}

	// For error cases like blocks where expressions are expected
	if p.check(constants.LEFT_BRACE) {
		p.error(p.peek(), "Expect expression.")
		return ""
	}

	token := p.advance()
	switch token.TokenType {
	case constants.TRUE:
		return "true"
	case constants.FALSE:
		return "false"
	case constants.NIL:
		return "nil"
	case constants.NUMBER:
		return token.Literal.(string)
	case constants.STRING:
		// Mark string literals explicitly with a prefix to differentiate them from number literals
		return "(string " + token.Literal.(string) + ")"
	case constants.IDENTIFIER:
		// Include the line number with the variable reference
		return fmt.Sprintf("(var-ref %s %d)", token.Lexeme, token.Line)
	case constants.LEFT_PAREN:
		// Check for empty parentheses
		if p.check(constants.RIGHT_PAREN) {
			p.advance() // Consume the right paren
			return "()"
		}
		
		expr := p.expression()
		p.consume(constants.RIGHT_PAREN, "Expect ')' after expression.")
		return fmt.Sprintf("(group %s)", expr)
	default:
		p.error(token, "Expect expression.")
		return ""
	}
}

func (p *Parser) advance() types.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() types.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.tokens[p.current].TokenType == constants.EOF
}

func (p *Parser) assignment() string {
	expr := p.logicalOr()
	
	if p.match(constants.EQUAL) {
		equals := p.previous()
		// Since assignment is right-associative, we recursively call assignment
		// to handle chained assignments like a = b = c
		value := p.assignment()
		
		// Check if the left-hand side is a valid variable identifier
		if strings.HasPrefix(expr, "(var-ref ") && strings.HasSuffix(expr, ")") {
			// Extract the variable name and line
			content := strings.TrimPrefix(expr, "(var-ref ")
			content = strings.TrimSuffix(content, ")")
			
			parts := strings.Split(content, " ")
			if len(parts) >= 2 {
				varName := parts[0]
				line := parts[1]
				
				// Create an assignment expression
				return fmt.Sprintf("(assign %s %s %s)", varName, line, value)
			}
		}
		
		p.error(equals, "Invalid assignment target.")
	}
	
	return expr
}

func (p *Parser) logicalOr() string {
	expr := p.logicalAnd()
	
	for p.match(constants.OR) {
		right := p.logicalAnd()
		expr = fmt.Sprintf("(or %s %s)", expr, right)
	}
	
	return expr
}

func (p *Parser) logicalAnd() string {
	expr := p.equality()
	
	for p.match(constants.AND) {
		right := p.equality()
		expr = fmt.Sprintf("(and %s %s)", expr, right)
	}
	
	return expr
}

// Parse an if statement: "if" "(" expression ")" statement ("else" statement)?
func (p *Parser) ifStatement() string {
	p.consume(constants.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(constants.RIGHT_PAREN, "Expect ')' after if condition.")
	
	thenBranch := p.statement()
	
	var elseBranch string
	if p.match(constants.ELSE) {
		// In the else branch, we expect a statement, not a declaration
		if p.check(constants.VAR) {
			p.error(p.peek(), "Expect expression.")
			p.synchronize()
			return ""
		}
		
		elseBranch = p.statement()
		return fmt.Sprintf("(if %s %s %s)", condition, thenBranch, elseBranch)
	}
	
	return fmt.Sprintf("(if %s %s)", condition, thenBranch)
}

// Parse a while statement: "while" "(" expression ")" statement
func (p *Parser) whileStatement() string {
	p.consume(constants.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(constants.RIGHT_PAREN, "Expect ')' after while condition.")
	
	body := p.statement()
	
	return fmt.Sprintf("(while %s %s)", condition, body)
}

// Parse a for statement: "for" "(" (varDecl | exprStmt | ";") expression? ";" expression? ")" statement
func (p *Parser) forStatement() string {
	p.consume(constants.LEFT_PAREN, "Expect '(' after 'for'.")
	
	// Parse initializer
	var initializer string
	if p.match(constants.SEMICOLON) {
		// No initializer
		initializer = ""
	} else if p.match(constants.VAR) {
		initializer = p.varDeclaration()
	} else {
		// Handle expression or syntax error
		initializer = p.expressionStatement()
	}
	
	// Parse condition
	var condition string
	if !p.check(constants.SEMICOLON) {
		condition = p.expression()
		if p.hadError {
			// If there was an error parsing the condition, synchronize and continue
			p.synchronize()
		}
	} else {
		// If no condition is provided, use 'true'
		condition = "true"
	}
	p.consume(constants.SEMICOLON, "Expect ';' after loop condition.")
	
	// Parse increment
	var increment string
	if !p.check(constants.RIGHT_PAREN) {
		increment = p.expression()
		if p.hadError {
			// If there was an error parsing the increment, synchronize and continue
			p.synchronize()
		}
	} else {
		increment = ""
	}
	p.consume(constants.RIGHT_PAREN, "Expect ')' after for clauses.")
	
	// If there were any errors, don't proceed with desugaring
	if p.hadError {
		return ""
	}
	
	// Parse body
	body := p.statement()
	
	// Desugar for loop into a while loop with a block
	
	// If there's an increment, make the body a block containing the original body and the increment
	if increment != "" {
		body = fmt.Sprintf("(block %s %s)", body, increment)
	}
	
	// Create the while loop with the condition
	body = fmt.Sprintf("(while %s %s)", condition, body)
	
	// If there's an initializer, make the loop a block containing the initializer and the while loop
	if initializer != "" {
		body = fmt.Sprintf("(block %s %s)", initializer, body)
	}
	
	return body
}
