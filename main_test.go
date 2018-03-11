package main

import (
	"io/ioutil"
	"reflect"
	"sort"
	"testing"
)

func TestScan(t *testing.T) {
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

func TestExtract(t *testing.T) {

}
