package parser

import (
	"fmt"
	"math"
	"math/big"

	"github.com/sarumaj/edu-taschenrechner/pkg/calc"
)

// make sure that the node type implements the Node interface
var _ Node = &node{}

// Node is an interface for nodes in the parse tree
type Node = NodeInterface[*node]

// NodeInterface is a generic interface for nodes in the parse tree
// It is used to define the methods that are common to all nodes
type NodeInterface[n any] interface {
	Evaluate(p *parser) (*big.Float, error)
	Float() (*big.Float, bool)
	IsLeaf() bool
	Left() n
	Right() n
	SetLeft(left any) n
	SetRight(right any) n
	SetValue(value string) n
	Value() string
}

// node implements the Node interface
type node struct {
	value string
	left  *node
	right *node
}

// Evaluate evaluates the node and returns the result
func (node *node) Evaluate(p *parser) (*big.Float, error) {
	if node.IsLeaf() { // Leaf node, check if it is a variable or a number
		if val, ok := p.Variables[node.value]; ok {
			return val, nil
		}

		if val, ok := node.Float(); ok {
			return val, nil
		}

		return nil, fmt.Errorf("undefined variable or function: %s", node.value)
	}

	// Handle function calls
	if fn, ok := p.Functions[node.value]; ok {
		// Collect all arguments
		var args []*big.Float
		// Extract the arguments from the nodes in the left subtree, from left to right
		for currentNode := node.Left(); currentNode != nil; currentNode = currentNode.Right() {
			arg, err := currentNode.Left().Evaluate(p)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}

		// Call the function with the evaluated arguments
		return fn(args...)
	}

	if node.Left() == nil {
		return nil, fmt.Errorf("missing left operand for operator %s", node.value)
	}

	// Evaluate the left subtree
	left, err := node.Left().Evaluate(p)
	if err != nil {
		return nil, err
	}

	// Handle unary operators
	switch node.Value() {
	case "!": // Factorial
		left, err = calc.Factorial(left, 1)
		if err != nil {
			return nil, err
		}

	case "°": // Convert the result from degrees to radians
		left = big.NewFloat(0).Mul(left, big.NewFloat(0).Quo(big.NewFloat(math.Pi), big.NewFloat(180)))

	case "√": // Square root
		if left.Cmp(big.NewFloat(0)) < 0 {
			return nil, fmt.Errorf("square root of a negative number")
		}
		left = big.NewFloat(0).Sqrt(left)

	case "-": // Unary minus
		if node.Right() == nil {
			return big.NewFloat(0).Neg(left), nil
		}

	}

	// If there is no right node, return the result
	if node.Right() == nil {
		return left, nil
	}

	// Evaluate the right subtree
	right, err := node.Right().Evaluate(p)
	if err != nil {
		return nil, err
	}

	// Handle binary operators
	switch node.Value() {
	case "+": // Addition
		return big.NewFloat(0).Add(left, right), nil

	case "-": // Subtraction
		return big.NewFloat(0).Sub(left, right), nil

	case "*": // Multiplication
		return big.NewFloat(0).Mul(left, right), nil

	case "/": // Division
		if right.Cmp(big.NewFloat(0)) == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return big.NewFloat(0).Quo(left, right), nil

	case "^": // Exponentiation
		return calc.Pow(left, right)

	default:
		return nil, fmt.Errorf("unsupported operator: %s", node.value)
	}
}

// Float converts the node value to a big.Float
func (node *node) Float() (*big.Float, bool) {
	return big.NewFloat(0).SetString(node.value)
}

// IsLeaf returns true if the node is a leaf node
func (node *node) IsLeaf() bool {
	return node.left == nil && node.right == nil
}

// Left returns the left node
func (node *node) Left() *node { return node.left }

// Right returns the right node
func (node *node) Right() *node { return node.right }

// SetLeft sets the left node
func (n *node) SetLeft(left any) *node {
	n.left, _ = left.(*node)
	return n
}

// SetRight sets the right node
func (n *node) SetRight(right any) *node {
	n.right, _ = right.(*node)
	return n
}

// SetValue sets the node value
func (node *node) SetValue(value string) *node {
	node.value = value
	return node
}

// Value returns the node value
func (node *node) Value() string { return node.value }

// NewNode creates a new node
func NewNode(value string) Node {
	return &node{value: value}
}
