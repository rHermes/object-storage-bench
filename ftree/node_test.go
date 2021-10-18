package ftree_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rhermes/object-storage-bench/ftree"
)

var nodeTestTree = ftree.Node{
	Children: map[string]ftree.Node{
		"bus": {
			Children: map[string]ftree.Node{
				"two": {},
				"abc": {
					Children: map[string]ftree.Node{
						"peter": {},
					},
				},
			},
		},
		"apple": {},
	},
}

func TestFromFiles(t *testing.T) {
	t.Parallel()

	files := []string{
		"apple",
		"bus/two",
		"bus/abc/peter",
	}

	n := ftree.FromFiles(files)
	require.Equal(t, nodeTestTree, n)
}

func TestNodePrint(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	nodeTestTree.Print(&buf)

	expt := strings.TrimSpace(`
apple
bus/
 abc/
  peter
 two
`) + "\n"

	require.Equal(t, expt, buf.String())
}

func TestNodeFiles(t *testing.T) {
	t.Parallel()

	expt := []string{
		"apple",
		"bus/abc/peter",
		"bus/two",
	}

	require.Equal(t, expt, nodeTestTree.Files())
}

func TestNodeAvgDepth(t *testing.T) {
	t.Parallel()

	x := nodeTestTree.AvgDepth()
	require.Equal(t, (1.0+3.0+2.0)/3.0, x)
}

func TestNodeMaxDepth(t *testing.T) {
	t.Parallel()
	require.Equal(t, uint64(3), nodeTestTree.MaxDepth())
}
