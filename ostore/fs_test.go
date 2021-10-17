package ostore_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/rhermes/object-storage-bench/ftree"
	"github.com/rhermes/object-storage-bench/ostore"
)

func TestFsStore(t *testing.T) {
	// We create a mem fs for testing.
	fs := afero.NewMemMapFs()

	store := ostore.NewFS(fs)

	files := []string{
		"abc.txt",
		"d/e/geg.txt",
		"d/h/test.txt",
		"d/lol.txt",
	}

	tree := ftree.FromFiles(files)

	require.NoError(t, store.Create(tree))

	for _, file := range files {
		fi, err := fs.Stat(file)
		require.NoError(t, err)
		require.False(t, fi.IsDir())
	}
}
