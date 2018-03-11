package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
)

// Scan list the files from directory.
func Scan(dir string, c chan string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			c <- path.Join(dir, info.Name())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Extract get mp3 audio for video file.
func Extract(c chan<- string, dir string) error {
	return nil
}

func main() {
	s := flag.String("s", "", "Source directory.")
	// o := flag.String("o", "", "Output directory.")

	flag.Parse()

	go func() {
		if err := Scan(*s, c); err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}
		close(c)
	}()
}
