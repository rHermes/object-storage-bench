package ftree

import (
	"fmt"
	"sort"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// A Config describes the parameters used to generate a filetree with the
// Generate call.
type Config struct {
	// The random source fo the generation
	Src rand.Source

	// Number of files to generate
	NumFiles uint64

	// The average depth of the folders.
	AvgDepth uint32

	// Chance of new folders, must be between 0 and 1
	NewRatio float64
}

// Generate a tree relationship.
func Generate(c Config) Node {
	// This algorithm was thought up by a young genius I know, I can claim only
	// the implementation :)
	if c.AvgDepth == 0 {
		panic("avg depth of zero is not possible")
	}

	// Will be used the generation of file depths
	pos := distuv.Poisson{
		Lambda: float64(c.AvgDepth - 1),
		Src:    c.Src,
	}

	// We start from zero, so it's one less
	padding := numDigits(c.NumFiles - 1)
	rootNode := Node{Children: make(map[string]Node)}

	// Pool is the pool of files still left to be placed. They are grouped by
	// the level they are going to use.
	pool := make(map[string]uint64, c.NumFiles)

	for i := uint64(0); i < c.NumFiles; i++ {
		// Filenames are simply numbers with zero names to be all of the same length
		fname := fmt.Sprintf("f%0*d", padding, i)
		depth := uint64(pos.Rand())

		pool[fname] = depth
		rootNode.Children[fname] = Node{}
	}

	// Now we must process the tree
	genTreeStruct(pool, c, rootNode, 0)

	return rootNode
}

func genTreeStruct(pool map[string]uint64, c Config, n Node, level uint64) {
	// we know each node is a file, but we don't know how many of them are
	// to be redistributed. We loop through and check to see.
	filesToDist := make([]string, 0)

	for name := range n.Children {
		l := pool[name]
		if l < level {
			panic("this should never happen!")
		} else if l == level {
			continue
		}

		filesToDist = append(filesToDist, name)
	}
	// fmt.Printf("We are redistributing %d of the %d files.\n", len(filesToDist), len(n.Children))

	if len(filesToDist) == 0 {
		return
	}

	// we sort the file names, to make sure we get a deterministic output
	sort.Strings(filesToDist)

	// we are using -1 in the N and +1 on the outside, so we always generate at least one folder.
	numNewFolders := uint64(distuv.Binomial{
		N: float64(len(filesToDist) - 1),
		P: c.NewRatio, Src: c.Src,
	}.Rand()) + 1

	// We start from zero, so it's one less
	padding := numDigits(numNewFolders - 1)

	newFolders := []string{}

	// now we must distribute the files over. We create an identical one for each possibility.
	cat := distuv.NewCategorical(repeatFloat(1.0, numNewFolders), c.Src)
	for _, name := range filesToDist {
		folderName := fmt.Sprintf("d%0*.0f", padding, cat.Rand())

		nb, ok := n.Children[folderName]
		if !ok {
			nb.Children = make(map[string]Node)
			n.Children[folderName] = nb

			newFolders = append(newFolders, folderName)
		}

		// we move the child down and remove it from the current directory.
		nb.Children[name] = n.Children[name]
		delete(n.Children, name)
	}

	// Now we simply need to run over the directories we created
	for _, newFolder := range newFolders {
		genTreeStruct(pool, c, n.Children[newFolder], level+1)
	}
}

// utilitty function to create a float array.
func repeatFloat(f float64, n uint64) []float64 {
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = f
	}

	return xs
}
