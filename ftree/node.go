package ftree

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

// FromFiles returns a tree built on the files in the files tree.
func FromFiles(files []string) Node {
	root := Node{Children: make(map[string]Node)}

	for _, file := range files {
		// Clean the file to make sure it's nice.
		file = path.Clean(file)
		parts := strings.Split(file, "/")

		// Just guard against the file ending with a /
		if strings.HasSuffix(file, "/") {
			panic("files should not end in a slash")
		}

		cNode := root
		for i, part := range parts {
			if _, ok := cNode.Children[part]; !ok {
				// The last part should only be created
				if i != len(parts)-1 {
					cNode.Children[part] = Node{Children: make(map[string]Node)}
				} else {
					cNode.Children[part] = Node{}
				}
			}

			cNode = cNode.Children[part]
		}
	}

	return root
}

// IsFile returns if the node is a file.
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
	for key := range n.Children {
		keys = append(keys, key)
	}

	indent := strings.Repeat(" ", level)

	sort.Strings(keys)

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

	keys := make([]string, 0, len(n.Children))
	for key := range n.Children {
		keys = append(keys, key)
	}

	// We sort the keys, to ensure stable output.
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

// AvgDepth returns the average depth of the files, as measured from this node.
// It is important to notice that files in the top level directory have a depth
// of 1.
func (n Node) AvgDepth() float64 {
	sum, count := n.avgDepth(0)

	return float64(sum) / float64(count)
}

// returns the number of files under this this node. I returns
// the sum and the count of files under it.
func (n Node) avgDepth(level uint64) (uint64, uint64) {
	if n.IsFile() {
		return level, 1
	}

	sum, cnt := uint64(0), uint64(0)

	for _, node := range n.Children {
		ks, kc := node.avgDepth(level + 1)
		sum += ks
		cnt += kc
	}

	return sum, cnt
}

// MaxDepth returns the MaxDepth for node.
func (n Node) MaxDepth() uint64 {
	if n.IsFile() {
		return 0
	}

	ret := uint64(0)

	for _, nd := range n.Children {
		nm := nd.MaxDepth()
		if nm > ret {
			ret = nm
		}
	}

	// We add one for the current node
	return ret + 1
}
