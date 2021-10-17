package ftree

import (
	"fmt"
	"sort"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Config struct {
	// The seed to the random generator.
	Seed uint64

	// Number of files to generate
	NumFiles uint64

	// The average depth of the folders
	AvgDepth uint32

	// Chance of new folders, must be between 0 and 1
	NewRatio float64
}

// Generate a tree relationship.
func Generate(c Config) Node {
	//This algorithm was thought up by a young genius I know, I can claim only
	// the implementation :)
	src := rand.NewSource(c.Seed)

	// Will be used the generation of file depths
	pos := distuv.Poisson{
		Lambda: float64(c.AvgDepth),
		Src:    src,
	}

	// We start from zero, so it's one less
	padding := numDigits(c.NumFiles - 1)

	rootNode := Node{
		Children: make(map[string]Node),
	}

	// Pool is the pool of files still left to be placed. They are grouped by
	// the level they are going to use.
	pool := make(map[string]uint64, c.NumFiles)
	for i := uint64(0); i < c.NumFiles; i++ {
		// Filenames are simply numbers with zero names to be all of the same length
		fname := fmt.Sprintf("f%0*d", padding, i)
		depth := uint64(pos.Rand())

		pool[fname] = depth
		rootNode.Children[fname] = Node{}

		fmt.Printf("name: %s, depth: %d\n", fname, depth)
	}

	fmt.Printf("=== Going to generating the tree ===\n\n")

	// Now we must process the tree
	genTreeStruct(pool, src, c.NewRatio, rootNode, 0)

	return rootNode
}

func genTreeStruct(pool map[string]uint64, src rand.Source, p float64, n Node, level uint64) {
	// This function is called on each folder.

	// we know each node is a file, but we don't know how many of them are
	// to be redistributed. We loop through and check to see.
	filesToDist := make([]string, 0)
	for name, _ := range n.Children {
		l := pool[name]
		if l < level {
			panic("this should never happen!")
		}
		if l == level {
			continue
		}

		filesToDist = append(filesToDist, name)
	}
	fmt.Printf("We are redistributing %d of the %d files.\n", len(filesToDist), len(n.Children))

	if len(filesToDist) == 0 {
		return
	}

	// we sort the file names, to make sure we get a deterministic output
	sort.Strings(filesToDist)

	// we are using -1 in the N and +1 on the outside, so we always generate at least one folder.
	numNewFolders := uint64(distuv.Binomial{N: float64(len(filesToDist) - 1), P: p, Src: src}.Rand()) + 1
	// fmt.Printf("Number of new folders: %d\n", numNewFolders)

	// We start from zero, so it's one less
	padding := numDigits(numNewFolders - 1)

	newFolders := []string{}

	// now we must distribute the files over. We create an identical one for each possibility.
	cat := distuv.NewCategorical(repeatFloat(1.0, numNewFolders), src)
	for _, name := range filesToDist {
		folderName := fmt.Sprintf("d%0*.0f", padding, cat.Rand())
		nb, ok := n.Children[folderName]
		if !ok {
			nb.Children = make(map[string]Node)
			n.Children[folderName] = nb
			newFolders = append(newFolders, folderName)
		}

		nb.Children[name] = n.Children[name]
		delete(n.Children, name)
		// fmt.Printf("for file %s we are putting it into the new folder %s\n", name, folderName)
	}

	// Now we simply need to run over the directories we created
	for _, newFolder := range newFolders {
		genTreeStruct(pool, src, p, n.Children[newFolder], level+1)
	}

}

// utilitty function to create a float array
func repeatFloat(f float64, n uint64) []float64 {
	xs := make([]float64, n)
	for i, _ := range xs {
		xs[i] = f
	}
	return xs
}
