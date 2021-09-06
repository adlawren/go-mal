package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstParse(t *testing.T) {
	testCases := map[string]Node{
		"\"this is a test\"": &symbolNode{str: "\"this is a test\""},
		"(+ \"this is another\" \"test\")": &listNode{
			children: []Node{
				&symbolNode{str: "+"},
				&symbolNode{str: "\"this is another\""},
				&symbolNode{str: "\"test\""},
			},
		},
		"1":    &symbolNode{str: "1"},
		"1.23": &symbolNode{str: "1.23"},
		"(+ 1 2)": &listNode{
			children: []Node{
				&symbolNode{str: "+"},
				&symbolNode{str: "1"},
				&symbolNode{str: "2"},
			},
		},
		"(+ (* 1 (- 2 3)) 4 (* 5 (+ 6 7)))": &listNode{
			children: []Node{
				&symbolNode{str: "+"},
				&listNode{
					children: []Node{
						&symbolNode{str: "*"},
						&symbolNode{str: "1"},
						&listNode{
							children: []Node{
								&symbolNode{str: "-"},
								&symbolNode{str: "2"},
								&symbolNode{str: "3"},
							},
						},
					},
				},
				&symbolNode{str: "4"},
				&listNode{
					children: []Node{
						&symbolNode{str: "*"},
						&symbolNode{str: "5"},
						&listNode{
							children: []Node{
								&symbolNode{str: "+"},
								&symbolNode{str: "6"},
								&symbolNode{str: "7"},
							},
						},
					},
				},
			},
		},
	}

	for form, expectedNode := range testCases {
		reader, err := NewReader(form)
		assert.NoError(t, err)

		node, err := reader.ParseAst()
		assert.NoError(t, err)

		cmpNodes(t, expectedNode, node)
	}
}

func cmpNodes(t *testing.T, n1, n2 Node) {
	assert.Equal(t, len(n1.Children()), len(n2.Children()))

	if len(n1.Children()) > 0 {
		for idx := 0; idx < len(n1.Children()); idx++ {
			cmpNodes(t, n1.Children()[idx], n2.Children()[idx])
		}
	} else {
		assert.Equal(t, n1.String(), n2.String())
	}
}
