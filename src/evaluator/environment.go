package evaluator

import (
	"fmt"
)

// Environment stores variable bindings
type Environment struct {
	values    map[string]string
	enclosing *Environment // Reference to the enclosing environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		values:    make(map[string]string),
		enclosing: nil,
	}
}

// NewLocalEnvironment creates a new environment with the given enclosing environment
func NewLocalEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]string),
		enclosing: enclosing,
	}
}

// Define defines a new variable in the environment
func (e *Environment) Define(name string, value string) {
	e.values[name] = value
}

// Get retrieves a variable's value from the environment
func (e *Environment) Get(name string) (string, error) {
	if value, ok := e.values[name]; ok {
		return value, nil
	}
	
	// If the variable isn't found in this environment, check the enclosing one
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	
	return "", NewRuntimeError(fmt.Sprintf("Undefined variable '%s'.", name), 1)
}

// Assign assigns a value to an existing variable
func (e *Environment) Assign(name string, value string, line int) (string, error) {
	// Check if the variable exists in this environment
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return value, nil
	}
	
	// If the variable isn't found in this environment, try to assign in the enclosing one
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value, line)
	}
	
	return "", NewRuntimeError(fmt.Sprintf("Undefined variable '%s'.", name), line)
} 