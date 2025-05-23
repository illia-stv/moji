package evaluator

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (e *Evaluator) evalMultiply(expr string) (string, error) {
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
	
	// Check for booleans or nil, which are invalid for multiplication
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and multiply
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		// The operands must be numbers - runtime error
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	result := leftNum * rightNum
	// Format the result
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

func (e *Evaluator) evalDivide(expr string) (string, error) {
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
	
	// Check for booleans or nil, which are invalid for division
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and divide
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		// The operands must be numbers - runtime error
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	if rightNum == 0 {
		// Division by zero, throw a runtime error
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Division by zero.", line)
	}
	result := leftNum / rightNum
	// Format the result
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

func (e *Evaluator) evalAdd(expr string) (string, error) {
	// Special case for empty additions like (+ ())
	innerExpr := strings.TrimPrefix(expr, "(+ ")
	innerExpr = strings.TrimSuffix(innerExpr, ")")
	innerExpr = strings.TrimSpace(innerExpr)
	
	if innerExpr == "()" {
		return "()", nil
	}
	
	// Extract the operands
	left, right, ok := splitOperands(expr)
	if !ok {
		return expr, NewEvaluationError(ErrInvalidExpression, expr)
	}
	//
	// Check for empty parentheses in operands
	if left == "()" || right == "()" {
		return "()", nil
	}
	
	// Check if operands are explicitly marked as strings
	leftIsString := strings.HasPrefix(left, "(string ") && strings.HasSuffix(left, ")")
	rightIsString := strings.HasPrefix(right, "(string ") && strings.HasSuffix(right, ")")
	
	// Check for quoted string literals too
	if !leftIsString {
		leftIsString = strings.HasPrefix(strings.TrimSpace(left), "\"") && strings.HasSuffix(strings.TrimSpace(left), "\"")
	}
	if !rightIsString {
		rightIsString = strings.HasPrefix(strings.TrimSpace(right), "\"") && strings.HasSuffix(strings.TrimSpace(right), "\"")
	}
	
	// If both are string literals, do string concatenation
	if leftIsString && rightIsString {
		// Evaluate both operands recursively to get the string content
		leftValue, err := e.evaluateExpression(left)
		if err != nil {
			return expr, err
		}
		rightValue, err := e.evaluateExpression(right)
		if err != nil {
			return expr, err
		}
		
		// If the values are quoted strings, remove quotes before concatenation
		if strings.HasPrefix(leftValue, "\"") && strings.HasSuffix(leftValue, "\"") {
			leftValue = leftValue[1:len(leftValue)-1]
		}
		if strings.HasPrefix(rightValue, "\"") && strings.HasSuffix(rightValue, "\"") {
			rightValue = rightValue[1:len(rightValue)-1]
		}
		
		// Return the concatenated string with quotes
		return "\"" + leftValue + rightValue + "\"", nil
	}
	
	// If one is a string literal and one is a number literal, it's an error
	leftIsNumber := isNumeric(left) && !leftIsString
	rightIsNumber := isNumeric(right) && !rightIsString
	
	if (leftIsString && rightIsNumber) || (leftIsNumber && rightIsString) {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be two numbers or two strings.", line)
	}
	
	// For non-string literals, proceed with normal evaluation
	leftValue, err := e.evaluateExpression(left)
	if err != nil {
		return expr, err
	}
	rightValue, err := e.evaluateExpression(right)
	if err != nil {
		return expr, err
	}
	
	// Special handling for empty parentheses results
	if leftValue == "()" || rightValue == "()" {
		return "()", nil
	}
	
	// Check for booleans or nil, which are invalid for addition
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be two numbers or two strings.", line)
	}
	
	// Check if both values are numeric for addition
	leftIsNumberValue := isNumeric(leftValue)
	rightIsNumberValue := isNumeric(rightValue)
	
	// If both are numbers, perform numeric addition
	if leftIsNumberValue && rightIsNumberValue {
		leftNum, _ := strconv.ParseFloat(leftValue, 64)
		rightNum, _ := strconv.ParseFloat(rightValue, 64)
		result := leftNum + rightNum
		// Format the result
		if result == float64(int64(result)) {
			return strconv.FormatInt(int64(result), 10), nil
		}
		return strconv.FormatFloat(result, 'f', -1, 64), nil
	}
	
	// If neither is a number, treat as string concatenation
	if !leftIsNumberValue && !rightIsNumberValue {
		// If the values are quoted strings, remove quotes before concatenation
		if strings.HasPrefix(leftValue, "\"") && strings.HasSuffix(leftValue, "\"") {
			leftValue = leftValue[1:len(leftValue)-1]
		}
		if strings.HasPrefix(rightValue, "\"") && strings.HasSuffix(rightValue, "\"") {
			rightValue = rightValue[1:len(rightValue)-1]
		}
		
		// Return the concatenated string with quotes
		return "\"" + leftValue + rightValue + "\"", nil
	}
	
	// If one is a number and one is a string, it's a mixed type error
	line := 1 // Default to line 1
	return expr, NewRuntimeError("Operands must be two numbers or two strings.", line)
}

// isNumeric checks if a value is a numeric value
// Returns true if the value is a number, false otherwise
func isNumeric(value string) bool {
	// Try to parse as a number
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// Helper function to check if a value is a boolean
func isBoolean(value string) bool {
	return value == "true" || value == "false"
}

func (e *Evaluator) evalSubtract(expr string) (string, error) {
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
	
	// Check for booleans or nil, which are invalid for subtraction
	if isBoolean(leftValue) || isBoolean(rightValue) || leftValue == "nil" || rightValue == "nil" {
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	
	// Convert to numbers and subtract
	leftNum, err1 := strconv.ParseFloat(leftValue, 64)
	rightNum, err2 := strconv.ParseFloat(rightValue, 64)
	if err1 != nil || err2 != nil {
		// The operands must be numbers - runtime error
		line := 1 // Default to line 1
		return expr, NewRuntimeError("Operands must be numbers.", line)
	}
	result := leftNum - rightNum
	// Format the result
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

func (e *Evaluator) evalUnaryMinus(expr string) (string, error) {
	// Extract the operand by removing "(- " and ")"
	operand := strings.TrimPrefix(expr, "(- ")
	operand = strings.TrimSuffix(operand, ")")
	
	// Find the token associated with this expression to get the line number
	line := 1 // Default to line 1 if we can't find the token info
	
	// Evaluate the operand recursively
	value, err := e.evaluateExpression(operand)
	if err != nil {
		return expr, err
	}
	
	// Convert to number and negate
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		// The operand is not a number, throw a runtime error
		fmt.Fprintf(os.Stderr, "Runtime error: operand %q is not a number\n", value)
		// Use exact message "Operand must be a number." as per the specification
		return expr, NewRuntimeError("Operand must be a number.", line)
	}
	
	result := -num
	// Format the result
	if result == float64(int64(result)) {
		return strconv.FormatInt(int64(result), 10), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil
}
