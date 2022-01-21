package lib

import (
	"errors"
)

type WildcardType int

const (
	exactMatch WildcardType = iota
	matchOne
	matchRest
)

type Node struct {
	pattern string

	part       string
	children   []*Node
	wildcard   WildcardType // is a wildcard node
	isCritical bool         // critical node
}

func (node *Node) matchPattern(part string) bool {
	return node.wildcard != exactMatch || part == node.part
}

func (node *Node) findMatchedChild(part string) (*Node, bool) {
	for _, child := range node.children {
		if child.matchPattern(part) {
			return child, true
		}
	}
	return nil, false
}

func newNode(pattern, part string, wildcard WildcardType) *Node {
	return &Node{
		pattern:  pattern,
		part:     part,
		wildcard: wildcard,
		children: make([]*Node, 0),
	}
}

type TrieTree struct {
	root *Node
}

func (tree *TrieTree) addPath(pattern string, parts []string) error {
	ptr := tree.root

	// match parts
	for _, part := range parts {
		if child, ok := ptr.findMatchedChild(part); ok {
			// follow the path
			ptr = child
		} else {
			// insert node and continue

			wildcard := exactMatch
			switch part[0] {
			case ':':
				wildcard = matchOne
			case '*':
				wildcard = matchRest
			}

			node := newNode("", part, wildcard)
			ptr.children = append(ptr.children, node)
			ptr = node
		}
	}

	if ptr.isCritical {
		// error: cannot add 2 nodes
		return errors.New("cannot insert two identical paths")
	}
	ptr.isCritical = true
	ptr.pattern = pattern
	return nil
}

func (tree *TrieTree) matchPath(pathParts []string) (*Node, bool) {
	ptr := tree.root

	// split path
	matchPosition := 0

	// match parts
	for idx, part := range pathParts {
		if child, ok := ptr.findMatchedChild(part); ok {
			ptr = child
			if child.wildcard != exactMatch {
				matchPosition = idx
				if child.wildcard == matchRest {
					break
				}
			}
		} else {
			return nil, false
		}
	}
	result := ptr.isCritical

	if ptr.wildcard == matchRest && matchPosition == len(pathParts) {
		// matches nothing
		result = false
	}
	return ptr, result
}

func newTrieTree() *TrieTree {
	return &TrieTree{
		root: newNode("/", "", exactMatch),
	}
}
