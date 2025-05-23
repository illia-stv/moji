package evaluator

import (
	"fmt"
	"os"
	"strconv"
)

func (e *Evaluator) evalGreater(expr string) (string, error) {
	// Extract the operands
	left, right, ok := splitOperands(expr)
	if !ok {
		return expr, NewEvaluationError(ErrInvalidExpression, expr)
	}
	
	// Evaluate both operands recursively
	leftValue, err := e.evaluateExpression(left)
	if err != nil {
		return expr, err
	}
	rightValue, err := e.evaluateExpression(right)
	if err != nil {
		return expr, err
	}
	
	// Check for booleans or nil, which are invalid for comparison
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and compare
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse operands as numbers for > comparison: %s, %s\n", leftValue, rightValue)
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Compare the values
	return strconv.FormatBool(leftNum > rightNum), nil
}

func (e *Evaluator) evalGreaterEqual(expr string) (string, error) {
	// Extract the operands
	left, right, ok := splitOperands(expr)
	if !ok {
		return expr, NewEvaluationError(ErrInvalidExpression, expr)
	}
	
	// Evaluate both operands recursively
	leftValue, err := e.evaluateExpression(left)
	if err != nil {
		return expr, err
	}
	rightValue, err := e.evaluateExpression(right)
	if err != nil {
		return expr, err
	}
	
	// Check for booleans or nil, which are invalid for comparison
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and compare
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse operands as numbers for >= comparison: %s, %s\n", leftValue, rightValue)
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Compare the values
	return strconv.FormatBool(leftNum >= rightNum), nil
}

func (e *Evaluator) evalLess(expr string) (string, error) {
	// Extract the operands
	left, right, ok := splitOperands(expr)
	if !ok {
		return expr, NewEvaluationError(ErrInvalidExpression, expr)
	}
	
	// Evaluate both operands recursively
	leftValue, err := e.evaluateExpression(left)
	if err != nil {
		return expr, err
	}
	rightValue, err := e.evaluateExpression(right)
	if err != nil {
		return expr, err
	}
	
	// Check for booleans or nil, which are invalid for comparison
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and compare
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse operands as numbers for < comparison: %s, %s\n", leftValue, rightValue)
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Compare the values
	return strconv.FormatBool(leftNum < rightNum), nil
}

func (e *Evaluator) evalLessEqual(expr string) (string, error) {
	// Extract the operands
	left, right, ok := splitOperands(expr)
	if !ok {
		return expr, NewEvaluationError(ErrInvalidExpression, expr)
	}
	
	// Evaluate both operands recursively
	leftValue, err := e.evaluateExpression(left)
	if err != nil {
		return expr, err
	}
	rightValue, err := e.evaluateExpression(right)
	if err != nil {
		return expr, err
	}
	
	// Check for booleans or nil, which are invalid for comparison
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and compare
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse operands as numbers for <= comparison: %s, %s\n", leftValue, rightValue)
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Compare the values
	return strconv.FormatBool(leftNum <= rightNum), nil
} 