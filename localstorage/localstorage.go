package localstorage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type WriteFS interface {
	fs.FS
	Create(name string) (fs.File, error)
}

// localfs should serve isolated directories
// with parse traversal prevention
// https://www.stackhawk.com/blog/golang-path-traversal-guide-examples-and-prevention/
type localfs struct {
	trustedRoot string
	subFs       fs.FS
}

func New(dir string) fs.FS {
	if !filepath.IsAbs(dir) {
		base, _ := os.Getwd()
		dir = path.Join(base, dir)
	}
	return &localfs{
		trustedRoot: dir,
		subFs:       os.DirFS(dir),
	}
}

func inTrustedRoot(path string, trustedRoot string) error {
	cleanRoot := filepath.Clean(trustedRoot)
	for path != "/" {
		if path == cleanRoot {
			return nil
		}
		path = filepath.Dir(path)
	}
	return errors.New("path is outside of trusted root")
}

func (t *localfs) verifyPath(path string) (string, error) {

	c := filepath.Clean(path)

	err := inTrustedRoot(c, t.trustedRoot)
	if err != nil {
		return c, fmt.Errorf("unsafe or invalid path specified: %w", err)
	}

	_, error := os.Stat(c)
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		// file does not exist
		return c, nil
	}

	// if file exist, check if it is symlink
	r, err := filepath.EvalSymlinks(c)
	if err != nil {
		return c, fmt.Errorf("eval symlink unsafe or invalid path specified: %w", err)
	}

	err = inTrustedRoot(r, t.trustedRoot)
	if err != nil {
		return r, fmt.Errorf("unsafe or invalid path specified: %w", err)
	}

	return r, nil
}

func (t *localfs) Open(name string) (fs.File, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}
	return t.subFs.Open(name)
}

func (t *localfs) Sub(name string) (fs.FS, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}
	return New(fullPath), nil
}

func (t *localfs) ReadFile(name string) ([]byte, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(t.subFs, name)
}

func (t *localfs) ReadDir(name string) ([]fs.DirEntry, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}

	return fs.ReadDir(t.subFs, name)
}

func (t *localfs) Stat(name string) (fs.FileInfo, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}
	return fs.Stat(t.subFs, name)
}

// Create extending a bit standard fs interfaces.
func (t *localfs) Create(name string) (fs.File, error) {
	fullPath := path.Join(t.trustedRoot, name)
	_, err := t.verifyPath(fullPath)
	if err != nil {
		return nil, err
	}
	return os.Create(fullPath)
}