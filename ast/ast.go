package ast

import (
	"fmt"
	"regexp"
	"strings"
)

type Node interface {
	String() string
	Children() []Node
	Eval() (interface{}, error)
}

type symbolNode struct {
	str string
}

func newSymbolNode(form string) Node {
	return &symbolNode{str: form}
}

func (s *symbolNode) String() string {
	return s.str
}

func (s *symbolNode) Children() []Node { return []Node{} }

// TODO
func (s *symbolNode) Eval() (interface{}, error) {
	return nil, nil
}

type listNode struct {
	children []Node
}

func newListNode(children []Node) Node {
	return &listNode{children: children}
}

func (l *listNode) String() string {
	childStrings := []string{}
	for _, child := range l.Children() {
		childStrings = append(childStrings, child.String())
	}

	return "(" + strings.Join(childStrings, " ") + ")"
}

func (l *listNode) Children() []Node { return l.children }

// TODO
func (l *listNode) Eval() (interface{}, error) {
	return nil, nil
}

type Reader struct {
	idx    int
	tokens []string
}

func NewReader(form string) (*Reader, error) {
	tokens, err := Tokenize(form)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize form: %w", err)
	}

	r := &Reader{
		idx:    -1,
		tokens: tokens,
	}

	return r, nil
}

func Tokenize(form string) ([]string, error) {
	var tokens []string
	regexp, err := regexp.Compile("[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\\\.|[^\\\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)")
	if err != nil {
		return tokens, fmt.Errorf("failed to compile regular expression: %w", err)
	}

	tokens = regexp.FindAllString(form, -1)
	for idx := 0; idx < len(tokens); idx++ {
		tokens[idx] = strings.TrimSpace(tokens[idx])
	}

	return tokens, nil
}

func (r *Reader) ParseAst() (Node, error) {
	nextToken := r.nextToken()
	if nextToken == "(" {
		lst, err := r.parseList()
		if err != nil {
			return nil, fmt.Errorf("failed to parse list: %w", err)
		}

		return lst, nil
	} else {
		sym, err := r.parseSymbol()
		if err != nil {
			return nil, fmt.Errorf("failed to parse symbol: %w", err)
		}

		return sym, nil
	}

	return nil, nil
}

func (r *Reader) parseList() (Node, error) {
	if !r.next() || r.token() != "(" {
		return nil, fmt.Errorf("expected '(' token")
	}

	var children []Node
	for {
		nextToken := r.nextToken()
		if nextToken == ")" {
			r.next()
			break
		}

		rootNode, err := r.ParseAst()
		if err != nil {
			return nil, fmt.Errorf("failed to parse AST: %w", err)
		}

		children = append(children, rootNode)
	}

	return newListNode(children), nil
}

func (r *Reader) parseSymbol() (Node, error) {
	if !r.next() {
		return nil, fmt.Errorf("expected symbol token")
	}

	return newSymbolNode(r.token()), nil
}

func (r *Reader) next() bool {
	r.idx += 1

	return r.idx < len(r.tokens)
}

func (r *Reader) token() string {
	return r.tokens[r.idx]
}

func (r *Reader) nextToken() string {
	if r.idx < len(r.tokens)-1 {
		return r.tokens[r.idx+1]
	}

	return ""
}
