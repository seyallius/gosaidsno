// Package aspect - advice defines the advice types and execution chain for AOP
package aspect

import (
	"sort"
	"sync"
)

// -------------------------------------------- Constants & Variables --------------------------------------------

const (
	Before         AdviceType = iota // Before advice executes before the target function.
	After                            // After advice executes after the target function (always runs, even on panic).
	Around                           // Around advice wraps the target function execution (can skip it).
	AfterReturning                   // AfterReturning advice executes only if the function returns successfully (no panic/error).
	AfterThrowing                    // AfterThrowing advice executes only if the function panics.
)

// -------------------------------------------- Public Functions --------------------------------------------

// AdviceType represents the type of advice to apply.
type AdviceType int

// AdviceFunc is the signature for advice functions.
// It receives the execution context and can modify it.
// The context.Context inside the Context struct can be used for cancellation and deadlines.
type AdviceFunc func(c *Context) error

// Advice represents a single piece of advice attached to a function.
type Advice struct {
	Type     AdviceType
	Handler  AdviceFunc
	Priority int // Higher priority executes first (for same type).
}

// AdviceChain manages a collection of advice for a single function.
type AdviceChain struct {
	before         []Advice
	after          []Advice
	around         []Advice
	afterReturning []Advice
	afterThrowing  []Advice
	mu             sync.RWMutex
}

// NewAdviceChain creates a new empty advice chain.
func NewAdviceChain() *AdviceChain {
	return &AdviceChain{
		before:         make([]Advice, 0),
		after:          make([]Advice, 0),
		around:         make([]Advice, 0),
		afterReturning: make([]Advice, 0),
		afterThrowing:  make([]Advice, 0),
	}
}

// -------------------------------------------- Public Functions --------------------------------------------

// Add adds advice to the chain based on its type.
func (ac *AdviceChain) Add(advice Advice) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	switch advice.Type {
	case Before:
		ac.before = append(ac.before, advice)
	case After:
		ac.after = append(ac.after, advice)
	case Around:
		ac.around = append(ac.around, advice)
	case AfterReturning:
		ac.afterReturning = append(ac.afterReturning, advice)
	case AfterThrowing:
		ac.afterThrowing = append(ac.afterThrowing, advice)
	}
}

// ExecuteBefore runs all Before advice in order of priority.
func (ac *AdviceChain) ExecuteBefore(c *Context) error {
	ac.mu.RLock()
	advice := append([]Advice(nil), ac.before...)
	ac.mu.RUnlock()

	return ac.executeAdviceList(advice, c)
}

// ExecuteAfter runs all After advice in order of priority.
func (ac *AdviceChain) ExecuteAfter(c *Context) error {
	ac.mu.RLock()
	advice := append([]Advice(nil), ac.after...)
	ac.mu.RUnlock()

	return ac.executeAdviceList(advice, c)
}

// ExecuteAround runs all Around advice in order of priority.
func (ac *AdviceChain) ExecuteAround(c *Context) error {
	ac.mu.RLock()
	advice := append([]Advice(nil), ac.around...)
	ac.mu.RUnlock()

	return ac.executeAdviceList(advice, c)
}

// ExecuteAfterReturning runs all AfterReturning advice in order of priority.
func (ac *AdviceChain) ExecuteAfterReturning(c *Context) error {
	ac.mu.RLock()
	advice := append([]Advice(nil), ac.afterReturning...)
	ac.mu.RUnlock()

	return ac.executeAdviceList(advice, c)
}

// ExecuteAfterThrowing runs all AfterThrowing advice in order of priority.
func (ac *AdviceChain) ExecuteAfterThrowing(c *Context) error {
	ac.mu.RLock()
	advice := append([]Advice(nil), ac.afterThrowing...)
	ac.mu.RUnlock()

	return ac.executeAdviceList(advice, c)
}

// HasAround returns true if the chain has Around advice.
func (ac *AdviceChain) HasAround() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	return len(ac.around) > 0
}

// Count returns the total number of advice in the chain.
func (ac *AdviceChain) Count() int {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	return len(ac.before) +
		len(ac.after) +
		len(ac.around) +
		len(ac.afterReturning) +
		len(ac.afterThrowing)
}

// -------------------------------------------- Private Helper Functions --------------------------------------------

// executeAdviceList runs a list of advice in priority order.
func (ac *AdviceChain) executeAdviceList(adviceList []Advice, c *Context) error {
	if len(adviceList) == 0 {
		return nil
	}

	// Sort by priority (highest first)
	sortedAdviceList := make([]Advice, len(adviceList))
	copy(sortedAdviceList, adviceList)

	sort.Slice(sortedAdviceList, func(i, j int) bool {
		return sortedAdviceList[i].Priority > sortedAdviceList[j].Priority
	})

	// Execute in order
	for _, advice := range sortedAdviceList {
		// Check if context is cancelled before executing advice
		select {
		case <-c.Context().Done():
			return c.Context().Err()
		default:
			// Context not cancelled, continue execution
		}

		if err := advice.Handler(c); err != nil {
			return err
		}
	}
	return nil
}
