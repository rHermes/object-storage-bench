package ftree

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDumb(t *testing.T) {
	c := Config{
		Seed:     13,
		NumFiles: 1000,
		AvgDepth: 3,
		NewRatio: 0.1,
	}

	fmt.Println("=== ENTERING GENEREATE ===")
	node := Generate(c)
	fmt.Println("=== EXITING GENERATE ===")

	node.Print(os.Stdout)
	// for _, file := range node.Files() {
	// 	fmt.Println(file)
	// }

	g := make(map[int][]string)

	g[2] = append(g[2], "wow")
	g[2] = append(g[2], "wew")

	require.Equal(t, []string{"wow", "wew"}, g[2])
}
