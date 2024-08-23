package parser

import (
	"fmt"
	"strings"

	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Tokens is an interface for tokenizing an expression and parsing it into a parse tree
type Tokens interface {
	Compare(others ...string) bool
	Tree() (Node, error)
}

// tokens is a list of tokens which implements the Tokens interface
type tokens []string

// append appends a new token to the token list
func (tokens *tokens) append(token string) {
	*tokens = append(*tokens, token)
}

// consume consumes the next token from the token list
// and returns it. If the token list is empty, it returns an empty string.
func (tokens *tokens) consume() string {
	if len(*tokens) > 0 {
		token := (*tokens)[0]
		*tokens = (*tokens)[1:]
		return token
	}

	return ""
}

// len returns the number of tokens in the list.
func (tokens *tokens) len() int {
	return len(*tokens)
}

// parseExpr parses an expression and returns the root node of the parse tree
func (tokens *tokens) parseExpr() (Node, error) {
	node, err := tokens.parseAddSub()
	if err != nil {
		return nil, err
	}

	return node, nil
}

// parseAddSub parses addition and subtraction
func (tokens *tokens) parseAddSub() (Node, error) {
	// parse multiplication and division first
	node, err := tokens.parseMulDiv()
	if err != nil {
		return nil, err
	}

	for tokens.len() > 0 {
		if tokens.peek() != "+" && tokens.peek() != "-" {
			break // Not an addition or subtraction operator
		}

		// consume the operator
		operator := tokens.consume()

		// parse the right side of the expression
		right, err := tokens.parseMulDiv()
		if err != nil {
			return nil, err
		}

		// create a new node with the operator and the left and right nodes
		node = NewNode(operator).SetLeft(node).SetRight(right)
	}

	return node, nil
}

// parseMulDiv parses multiplication and division
func (tokens *tokens) parseMulDiv() (Node, error) {
	// parse factors first
	node, err := tokens.parseFactor()
	if err != nil {
		return nil, err
	}

	for tokens.len() > 0 {
		if tokens.peek() != "*" && tokens.peek() != "/" {
			break // Not a multiplication or division operator
		}

		// consume the operator
		operator := tokens.consume()

		// parse the right side of the expression
		right, err := tokens.parseFactor()
		if err != nil {
			return nil, err
		}

		// create a new node with the operator and the left and right nodes
		node = NewNode(operator).SetLeft(node).SetRight(right)
	}

	return node, nil
}

// parseFactor parses a factor (number, variable, function call, or sub-expression)
func (tokens *tokens) parseFactor() (Node, error) {
	if tokens.len() == 0 {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	// consume the token
	switch token := tokens.consume(); {
	case token == "(": // Handle sub-expression
		node, err := tokens.parseExpr()
		if err != nil {
			return nil, err
		}

		if tokens.len() == 0 || tokens.peek() != ")" {
			return nil, fmt.Errorf("missing closing parenthesis")
		}

		_ = tokens.consume() // consume the ')'
		return node, nil

	case token == "-": // Handle unary minus
		node, err := tokens.parseFactor()
		if err != nil {
			return nil, err
		}

		return NewNode("-").SetLeft(NewNode("0")).SetRight(node), nil

	case tokens.len() > 0 && tokens.peek() == "(": // Handle function call
		_ = tokens.consume() // consume the '('

		var args []Node // arguments to the function
		for tokens.len() > 0 && tokens.peek() != ")" {
			arg, err := tokens.parseExpr()
			if err != nil {
				return nil, err
			}

			args = append(args, arg)

			// consume the ',' if there are more arguments
			if tokens.len() > 0 && tokens.peek() == "," {
				_ = tokens.consume()
			}
		}

		if tokens.len() > 0 && tokens.peek() != ")" {
			return nil, fmt.Errorf("missing closing parenthesis in function call")
		}

		_ = tokens.consume() // consume the ')'

		// Create a function node where the arguments are linked as a list
		var argsNode Node
		for i := len(args) - 1; i >= 0; i-- { // Link the arguments as a list, right to left
			argsNode = NewNode("").SetLeft(args[i]).SetRight(argsNode)
		}

		// token is the function name
		// arguments are linked as a list in the left child of the function node
		return NewNode(token).SetLeft(argsNode), nil

	default: // Handle any other token
		return NewNode(token), nil

	}
}

// peek returns the next token in the list without consuming it.
// If the list is empty, it returns an empty string.
func (tokens *tokens) peek() string {
	if len(*tokens) > 0 {
		return (*tokens)[0]
	}

	return ""
}

// Compare compares the tokens with the given strings
func (tokens tokens) Compare(others ...string) bool {
	if tokens.len() != len(others) {
		return false
	}

	for i, token := range tokens {
		if token != others[i] {
			return false
		}
	}

	return true
}

// Tree parses the expression and returns the root node of the parse tree
func (tokens *tokens) Tree() (Node, error) {
	return tokens.parseExpr()
}

// Tokenize splits the expression into tokens
func Tokenize(expr string) (Tokens, error) {
	var tokens tokens
	var token strings.Builder

	for i := 0; i < len([]rune(expr)); i++ {
		switch ch := []rune(expr)[i]; {
		case ch == ' ': // Skip whitespace

		case
			runes.IsDigit(ch), ch == '.', // Handle numbers (including floating point)
			// Handle letters (for variable names and function names)
			runes.InRange(ch, 'a', 'z'), runes.InRange(ch, 'A', 'Z'), ch == '_', i > 0 && runes.IsDigit(ch):

			_, _ = token.WriteRune(ch)

		case runes.IsAnyOf(ch, "(),+-*/"): // Handle operators and parentheses
			if token.Len() > 0 {
				tokens.append(token.String())
				token.Reset()
			}
			tokens.append(string(ch))

		default:
			if token.Len() > 0 {
				tokens.append(token.String())
				token.Reset()
			}

		}
	}
	if token.Len() > 0 {
		tokens.append(token.String())
	}

	return &tokens, nil
}
