package filetree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var nodeTestTree = Node{
	Children: map[string]Node{
		"bus": Node{
			Children: map[string]Node{
				"two": Node{},
				"abc": Node{
					Children: map[string]Node{
						"peter": Node{},
					},
				},
			},
		},
		"apple": Node{},
	},
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
