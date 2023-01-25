package files

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	//create a temporary directory
	dir, err := ioutil.TempDir("", "files")
	fmt.Println(dir)
	if err != nil {
		t.Fatal(err)
	}

	l, err := NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	return l, dir, func() {
		//os.RemoveAll(dir)
	}
}

func TestSaveContentsOfReader(t *testing.T) {
	savePath := "1/test.png"
	fileContents := "Hello World"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	//check file has been correctly written
	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)

	//check contents of file
	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(b))
}

func TestGetsContentAndWritesToWriter(t *testing.T) {
	savePath := "1/test.png"
	fileContents := "Hello World"
	l, _, cleanup := setupLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	//check file has been correctly written
	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()

	//check contents of file
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(b))
}
