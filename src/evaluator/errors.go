package evaluator

import (
	"fmt"
)

// Error types
const (
	ErrInvalidExpression = "invalid expression format"
)

// EvaluationError represents an error during expression evaluation
type EvaluationError struct {
	Type string
	Expr string
}

// NewEvaluationError creates a new EvaluationError
func NewEvaluationError(errType, expr string) *EvaluationError {
	return &EvaluationError{
		Type: errType,
		Expr: expr,
	}
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("evaluation error: %s (expression: %s)", e.Type, e.Expr)
}

// RuntimeError represents a runtime error during program execution
type RuntimeError struct {
	Message string
	Line    int
}

// NewRuntimeError creates a new RuntimeError
func NewRuntimeError(message string, line int) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Line:    line,
	}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.Message, e.Line)
}

// GetFormattedMessage returns the message in the expected format
func (e *RuntimeError) GetFormattedMessage() string {
	return e.Message + "\n[line " + fmt.Sprint(e.Line) + "]"
}

// Common error messages
const (
	ErrDivisionByZero    = "division by zero"
	ErrInvalidOperands   = "invalid operands for operation"
	ErrInvalidOperator   = "invalid operator"
	ErrInvalidGroup      = "invalid group expression"
	ErrOperandNotNumber  = "Operand must be a number."
)
