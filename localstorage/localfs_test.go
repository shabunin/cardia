package localstorage

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"testing"
)

func tmpDir() (string, int, error) {
	p, err := os.MkdirTemp("", "cardia*")
	if err != nil {
		return "", 0, err
	}
	err = os.MkdirAll(path.Join(p, "subfolder1/dir11"), 0750)
	if err != nil {
		return p, 0, err
	}
	err = os.MkdirAll(path.Join(p, "subfolder1/dir12"), 0750)
	if err != nil {
		return p, 0, err
	}
	err = os.MkdirAll(path.Join(p, "subfolder2/dir21"), 0750)
	if err != nil {
		return p, 0, err
	}

	err = os.WriteFile(path.Join(p, "subfolder1/hello.txt"), []byte("hello"), 0750)
	if err != nil {
		return p, 0, err
	}
	err = os.WriteFile(path.Join(p, "subfolder1/dir12/friend.txt"), []byte("friend"), 0750)
	if err != nil {
		return p, 0, err
	}
	err = os.WriteFile(path.Join(p, "subfolder2/goodbye.txt"), []byte("friend"), 0750)
	if err != nil {
		return p, 0, err
	}

	cnt := 0
	fs.WalkDir(os.DirFS(p), ".",
		func(path string, d fs.DirEntry, err error) error {
			cnt += 1
			return nil
		})

	return p, cnt, nil
}

func TestLocalStorage(t *testing.T) {

	p, c, err := tmpDir()
	if p != "" {
		defer os.RemoveAll(p)
	}
	if err != nil {
		t.Error(err)
		return
	}

	lfs := NewLocalFs(p, &Config{
		CacheSize:     10 * 1024 * 1024,
		CacheDuration: 0,
	})

	cnt := 0
	fs.WalkDir(lfs, ".",
		func(path string, d fs.DirEntry, err error) error {
			fmt.Println(path)
			cnt += 1
			return nil
		})
	if cnt != c {
		t.Error("cnt = ", cnt, "; c = ", c)
	}

	f, err := lfs.Open("subfolder1/hello.txt")
	if err != nil {
		t.Error(err)
	}
	all, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(all))
	if !bytes.Equal(all, []byte("hello")) {
		t.Error("wrong read")
	}

	sfs, err := fs.Sub(lfs, "subfolder1")
	if err != nil {
		t.Error(err)
	}

	f, err = sfs.Open("dir12/friend.txt")
	if err != nil {
		t.Error(err)
	}
	all, err = io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(all))
	if !bytes.Equal(all, []byte("friend")) {
		t.Error("wrong read")
	}

	// path traversal
	_, err = sfs.Open("../subfolder2/goodbye.txt")
	if err == nil {
		t.Error("should return 'invalid name' error")
	}

	_, err = sfs.Open("new.txt")
	if err == nil {
		t.Error("should return error because there is no such file new.txt")
	}

	if wfs, ok := sfs.(WriteFS); ok {
		f, err = wfs.Create("new.txt")
		if err != nil {
			t.Error(err)
		}
		if w, ok := f.(io.Writer); ok {
			_, err = w.Write([]byte("new sensation"))
			if err != nil {
				t.Error(err)
			}
		} else {
			t.Error("file should support io.Writer interface")
		}

	} else {
		t.Error("should support localstorage.WriteFS interface")
	}

	f, err = sfs.Open("new.txt")
	if err != nil {
		t.Error(err)
	}
	all, err = io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(all))
	if !bytes.Equal(all, []byte("new sensation")) {
		t.Error("wrong read")
	}

	// symlink to target in trusted root dir
	target := path.Join(p, "subfolder1/symtarget.txt")
	err = os.WriteFile(target, []byte("target content"), 0750)
	if err != nil {
		t.Error(err)
	}
	symlink := path.Join(p, "subfolder1/symlink1")
	err = os.Symlink(target, symlink)
	if err != nil {
		t.Error(err)
	}

	f, err = sfs.Open("symlink1")
	if err != nil {
		t.Error(err)
	}
	all, err = io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(all, []byte("target content")) {
		t.Error("wrong read")
	}

	// symlink outside of trusted root
	target = "/etc/passwd"
	symlink = path.Join(p, "subfolder1/symlink2")
	err = os.Symlink(target, symlink)
	if err != nil {
		t.Error(err)
	}

	_, err = sfs.Open("symlink2")
	if err == nil {
		t.Error("there should be error!")
	} else {
		fmt.Println(err)
	}

}
