package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestScanBuffered(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}
	c := make(chan string, 5)
	// defer close(c)

	expected := []string{}
	actual := []string{}

	for i := 0; i < 5; i++ {
		f, _ := ioutil.TempFile(dir, "")
		expected = append(expected, f.Name())
	}

	go func() {
		if err := Scan(dir, c); err != nil {
			t.Fatalf("error scanning directory: %v", err)
		}
		close(c)
	}()

	for s := range c {
		if s != "" {
			actual = append(actual, s)
		}
	}

	sort.Strings(expected)
	sort.Strings(actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected to be equal, got %v and %v", expected, actual)
	}
}

func TestScanUnbuffered(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}
	c := make(chan string)

	expected := []string{}
	actual := []string{}

	for i := 0; i < 5; i++ {
		f, _ := ioutil.TempFile(dir, "")
		expected = append(expected, f.Name())
		f.Close()
	}

	go func() {
		if err := Scan(dir, c); err != nil {
			t.Fatalf("error scanning directory: %v", err)
		}
		close(c)
	}()

	for s := range c {
		if s != "" {
			actual = append(actual, s)
		}
	}

	sort.Strings(expected)
	sort.Strings(actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected to be equal, got %v and %v", expected, actual)
	}
}

func TestExtract(t *testing.T) {
	t.Parallel()

	c := make(chan string)
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}

	outDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}

	src := []string{}

	for i := 0; i < 5; i++ {
		f, _ := ioutil.TempFile(dir, "")
		f.Close()
		src = append(src, f.Name())
	}

	go func() {
		if err := Scan(dir, c); err != nil {
			t.Fatalf("error scanning directory: %v", err)
		}
		close(c)
	}()

	// we use `cp` command instead of `ffmpeg`
	if err := Extract(c, "cp", outDir); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dst := []string{}

	filepath.Walk(outDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			dst = append(dst, path.Join(outDir, info.Name()))
		}

		return nil
	})

	sort.Strings(src)
	sort.Strings(dst)

	if len(src) != len(dst) {
		t.Fatalf("expected to be equal, got %v and %v", len(src), len(dst))
	}
}
