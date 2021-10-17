package ostore

import (
	"fmt"
	"path"

	"github.com/spf13/afero"

	"github.com/rhermes/object-storage-bench/ftree"
)

var _ Store = &FsStore{}

// A FsStore implements Store, for an afero.FS file store.
type FsStore struct {
	fs afero.Fs
}

func NewFS(fs afero.Fs) *FsStore {
	return &FsStore{
		fs: fs,
	}
}

// Create creates the node structure from the filetree in the underlying
// fs.
func (fs *FsStore) Create(n ftree.Node) error {
	return fs.createTree("", n)
}

func (fs *FsStore) createTree(fname string, n ftree.Node) error {
	if n.IsFile() {
		f, err := fs.fs.Create(fname)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		// we don't need to write anything to this file, so we just close it
		if err := f.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}

		return nil
	}

	for name, node := range n.Children {
		if err := fs.createTree(path.Join(fname, name), node); err != nil {
			return err
		}
	}

	return nil
}
