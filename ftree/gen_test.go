package ftree_test

import (
	"testing"
	"time"

	"github.com/rhermes/object-storage-bench/ftree"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

// This is just a dumb test of the generator for now, until I can figure out a
// way to test it with purpose.
func TestGenDumb(t *testing.T) {
	// TODO(rHermes): Make these test actually useful
	c := ftree.Config{
		Src:      rand.NewSource(uint64(time.Now().UnixNano())),
		NumFiles: 10000,
		AvgDepth: 3,
		NewRatio: 0.01,
	}

	node := ftree.Generate(c)

	require.Len(t, node.Files(), int(c.NumFiles))
	require.InDelta(t, 3, node.AvgDepth(), 0.1)
}
