package filetree

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
)

// Node represents a node in a file. If a node does not have any children,
// it is considered to be a file otherwise it's a directory. This matches
// the way object stores see the world as their folders are simply constructs.
type Node struct {
	Children map[string]Node
}

func (n Node) IsFile() bool {
	return len(n.Children) == 0
}

// Print prints the tree for debugging. It's in lexograpihc order.
func (n Node) Print(w io.Writer) {
	n.printTree(w, 0)
}

func (n Node) printTree(w io.Writer, level int) {
	// This should really never happen, but we can guard against it.
	if len(n.Children) == 0 {
		return
	}

	// sort the keys before printing
	keys := make([]string, 0, len(n.Children))
	for key, _ := range n.Children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	indent := strings.Repeat(" ", level)
	for _, key := range keys {
		nd := n.Children[key]
		if nd.IsFile() {
			fmt.Fprintf(w, "%s%s\n", indent, key)
		} else {
			fmt.Fprintf(w, "%s%s/\n", indent, key)
			nd.printTree(w, level+1)
		}
	}
}

// Files returns the list of full files in the tree
// The files are returned in sortede order.
func (n Node) Files() []string {
	return n.getFiles("")
}

func (n Node) getFiles(prefix string) []string {
	if n.IsFile() {
		return nil
	}

	// sort the keys before printing
	keys := make([]string, 0, len(n.Children))
	for key, _ := range n.Children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	files := make([]string, 0)
	for _, key := range keys {
		nd := n.Children[key]
		nprefix := path.Join(prefix, key)

		if nd.IsFile() {
			files = append(files, nprefix)
		} else {
			files = append(files, nd.getFiles(nprefix)...)
		}
	}

	return files
}
