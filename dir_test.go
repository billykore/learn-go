package learning

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestFileIsExist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = replaceDir(fmt.Sprintf("%s/test", wd))
	if err != nil {
		t.Fatal(err)
	}
}

func replaceDir(path string) error {
	err := os.Mkdir(path, 0754)
	if errors.Is(err, os.ErrExist) {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
		err = os.Mkdir(path, 0754)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestRemove(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(fmt.Sprintf("%s/test", wd)); err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(fmt.Sprintf("%s/test", wd))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFiles(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dirname := wd + "/test"

	d, err := os.Open(dirname)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == "a.go" {
			err := os.Remove(fmt.Sprintf("%s/a.go", dirname))
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	t.Log("Reading " + dirname)
}
